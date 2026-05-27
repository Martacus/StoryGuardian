package internal

import (
	"encoding/base64"
	"fmt"
	"io"
	"os"
	"path/filepath"
	"sort"
	"strings"
	"sync"
	"time"

	"github.com/wailsapp/wails/v3/pkg/application"
)

// ImageInfo describes a single image file in a world's images/ directory.
type ImageInfo struct {
	FileName  string    `json:"fileName"`
	SizeBytes int64     `json:"sizeBytes"`
	ModTime   time.Time `json:"modTime"`
}

// ImageService handles image import, listing, and deletion for a world.
type ImageService struct {
	mu sync.RWMutex
}

// imagesDir returns the images/ subdirectory path for a world folder.
func (s *ImageService) imagesDir(folderPath string) string {
	return filepath.Join(folderPath, "images")
}

// ListImages returns all images in the world's images/ directory, sorted by
// modification time descending (newest first). Returns an empty slice (not
// an error) if the images/ directory doesn't exist yet.
func (s *ImageService) ListImages(folderPath string) ([]ImageInfo, error) {
	s.mu.RLock()
	defer s.mu.RUnlock()

	if folderPath == "" {
		return nil, fmt.Errorf("folder path cannot be empty")
	}

	dir := s.imagesDir(folderPath)
	entries, err := os.ReadDir(dir)
	if os.IsNotExist(err) {
		return []ImageInfo{}, nil
	}
	if err != nil {
		return nil, fmt.Errorf("read images dir: %w", err)
	}

	var images []ImageInfo
	for _, entry := range entries {
		if entry.IsDir() {
			continue
		}
		name := entry.Name()
		if !isImageFile(name) {
			continue
		}
		info, err := entry.Info()
		if err != nil {
			continue
		}
		images = append(images, ImageInfo{
			FileName:  name,
			SizeBytes: info.Size(),
			ModTime:   info.ModTime(),
		})
	}

	sort.Slice(images, func(i, j int) bool {
		return images[i].ModTime.After(images[j].ModTime)
	})

	return images, nil
}

// SelectImages opens a native file picker filtered to image types and returns
// the selected file paths. Returns an empty slice if the user cancels.
func (s *ImageService) SelectImages() ([]string, error) {
	paths, err := application.Get().Dialog.OpenFile().
		SetTitle("Select Images").
		AddFilter("Image Files", "*.png; *.jpg; *.jpeg; *.gif; *.webp; *.bmp").
		CanChooseFiles(true).
		CanChooseDirectories(false).
		PromptForMultipleSelection()
	if err != nil {
		return nil, fmt.Errorf("file picker: %w", err)
	}
	if paths == nil {
		return []string{}, nil
	}
	return paths, nil
}

// ImportImages copies source files into the world's images/ directory.
// Existing name conflicts get a numeric suffix (e.g. "map (1).png").
// Uses atomic write (temp file + rename) to prevent partial files on crash.
// Returns ImageInfo for each newly imported file.
func (s *ImageService) ImportImages(folderPath string, sourcePaths []string) ([]ImageInfo, error) {
	s.mu.Lock()
	defer s.mu.Unlock()

	if folderPath == "" {
		return nil, fmt.Errorf("folder path cannot be empty")
	}

	dir := s.imagesDir(folderPath)
	if err := os.MkdirAll(dir, 0755); err != nil {
		return nil, fmt.Errorf("create images dir: %w", err)
	}

	var imported []ImageInfo
	for _, src := range sourcePaths {
		info, err := importSingleImage(dir, src)
		if err != nil {
			return imported, fmt.Errorf("import %q: %w", filepath.Base(src), err)
		}
		imported = append(imported, info)
	}
	return imported, nil
}

// importSingleImage copies a single file into dir, resolving name conflicts.
func importSingleImage(dir, srcPath string) (ImageInfo, error) {
	base := filepath.Base(srcPath)
	destName := resolveImageNameConflict(dir, base)
	destPath := filepath.Join(dir, destName)

	srcFile, err := os.Open(srcPath)
	if err != nil {
		return ImageInfo{}, fmt.Errorf("open source: %w", err)
	}
	defer srcFile.Close()

	tmp, err := os.CreateTemp(dir, ".tmp-*")
	if err != nil {
		return ImageInfo{}, fmt.Errorf("create temp: %w", err)
	}
	tmpName := tmp.Name()
	ok := false
	defer func() {
		if !ok {
			tmp.Close()
			os.Remove(tmpName)
		}
	}()

	if _, err := io.Copy(tmp, srcFile); err != nil {
		return ImageInfo{}, fmt.Errorf("copy: %w", err)
	}
	if err := tmp.Sync(); err != nil {
		return ImageInfo{}, fmt.Errorf("sync: %w", err)
	}
	if err := tmp.Close(); err != nil {
		return ImageInfo{}, fmt.Errorf("close temp: %w", err)
	}
	ok = true

	if err := os.Rename(tmpName, destPath); err != nil {
		return ImageInfo{}, fmt.Errorf("rename: %w", err)
	}

	stat, _ := os.Stat(destPath)
	var size int64
	var modTime time.Time
	if stat != nil {
		size = stat.Size()
		modTime = stat.ModTime()
	}

	return ImageInfo{FileName: destName, SizeBytes: size, ModTime: modTime}, nil
}

// resolveImageNameConflict returns a unique file name in dir for the given base.
// If base doesn't exist, returns it unchanged. Otherwise appends "(1)", "(2)", ...
func resolveImageNameConflict(dir, base string) string {
	if _, err := os.Stat(filepath.Join(dir, base)); os.IsNotExist(err) {
		return base
	}
	ext := filepath.Ext(base)
	stem := strings.TrimSuffix(base, ext)
	for i := 1; ; i++ {
		candidate := fmt.Sprintf("%s (%d)%s", stem, i, ext)
		if _, err := os.Stat(filepath.Join(dir, candidate)); os.IsNotExist(err) {
			return candidate
		}
	}
}

// DeleteImage removes a single image from the world's images/ directory.
// Returns an error if fileName contains path separators (prevents path traversal).
func (s *ImageService) DeleteImage(folderPath, fileName string) error {
	s.mu.Lock()
	defer s.mu.Unlock()

	if folderPath == "" {
		return fmt.Errorf("folder path cannot be empty")
	}
	if fileName == "" {
		return fmt.Errorf("file name cannot be empty")
	}
	if strings.ContainsAny(fileName, "/\\") {
		return fmt.Errorf("invalid file name: %q", fileName)
	}

	path := filepath.Join(s.imagesDir(folderPath), fileName)
	if err := os.Remove(path); err != nil {
		if os.IsNotExist(err) {
			return fmt.Errorf("image not found: %q", fileName)
		}
		return fmt.Errorf("delete image: %w", err)
	}
	return nil
}

// GetImagePath returns the absolute path for an image file.
func (s *ImageService) GetImagePath(folderPath, fileName string) (string, error) {
	s.mu.RLock()
	defer s.mu.RUnlock()

	if folderPath == "" {
		return "", fmt.Errorf("folder path cannot be empty")
	}
	if strings.ContainsAny(fileName, "/\\") {
		return "", fmt.Errorf("invalid file name: %q", fileName)
	}

	path := filepath.Join(s.imagesDir(folderPath), fileName)
	if _, err := os.Stat(path); err != nil {
		return "", fmt.Errorf("image not found: %q", fileName)
	}
	return path, nil
}

// GetImageDataURL reads an image file and returns it as a base64-encoded data URL
// suitable for use in an HTML <img> src attribute.
func (s *ImageService) GetImageDataURL(folderPath, fileName string) (string, error) {
	s.mu.RLock()
	defer s.mu.RUnlock()

	if folderPath == "" {
		return "", fmt.Errorf("folder path cannot be empty")
	}
	if strings.ContainsAny(fileName, "/\\") {
		return "", fmt.Errorf("invalid file name: %q", fileName)
	}

	path := filepath.Join(s.imagesDir(folderPath), fileName)
	data, err := os.ReadFile(path)
	if err != nil {
		if os.IsNotExist(err) {
			return "", fmt.Errorf("image not found: %q", fileName)
		}
		return "", fmt.Errorf("read image: %w", err)
	}

	ext := strings.ToLower(filepath.Ext(fileName))
	mime := imageMimeType(ext)
	encoded := base64.StdEncoding.EncodeToString(data)
	return fmt.Sprintf("data:%s;base64,%s", mime, encoded), nil
}

// isImageFile returns true if the file extension is a supported image type.
func isImageFile(name string) bool {
	switch strings.ToLower(filepath.Ext(name)) {
	case ".png", ".jpg", ".jpeg", ".gif", ".webp", ".bmp", ".svg":
		return true
	}
	return false
}

// imageMimeType returns the MIME type for the given lowercase file extension.
func imageMimeType(ext string) string {
	switch ext {
	case ".jpg", ".jpeg":
		return "image/jpeg"
	case ".png":
		return "image/png"
	case ".gif":
		return "image/gif"
	case ".webp":
		return "image/webp"
	case ".bmp":
		return "image/bmp"
	case ".svg":
		return "image/svg+xml"
	default:
		return "image/octet-stream"
	}
}
