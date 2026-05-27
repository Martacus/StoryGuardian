package internal

import (
	"encoding/json"
	"fmt"
	"os"
	"path/filepath"
	"time"

	"github.com/wailsapp/wails/v3/pkg/application"
)

const worldMetaFileName = "world.json"

// WorldMeta is persisted as world.json inside the world folder.
// It is the ground truth for a world's metadata on disk.
// Separate from WorldInfo, which is the lightweight recents-list entry in app config.
type WorldMeta struct {
	Name        string    `json:"name"`
	Description string    `json:"description"`
	CreatedAt   time.Time `json:"createdAt"`
	UpdatedAt   time.Time `json:"updatedAt"`
}

// WorldService handles world folder operations and the native OS folder picker.
// It is stateless — all data lives on disk, no in-memory cache.
type WorldService struct{}

// SelectFolder opens the native OS directory picker and returns the selected path,
// or an empty string if the user cancelled without selecting.
func (s *WorldService) SelectFolder() (string, error) {
	path, err := application.Get().Dialog.OpenFile().
		SetTitle("Select World Folder").
		CanChooseDirectories(true).
		CanChooseFiles(false).
		CanCreateDirectories(true).
		PromptForSingleSelection()
	if err != nil {
		return "", fmt.Errorf("folder picker: %w", err)
	}
	return path, nil
}

// CreateWorld scaffolds a new world at folderPath with the given name.
// If folderPath already contains files, the world is created in a child folder
// named after the world instead.
// It creates: world.json, entities/, links.json, tags.json, categories.json.
// Returns a WorldInfo suitable for adding to the recents list.
func (s *WorldService) CreateWorld(name, folderPath string) (*WorldInfo, error) {
	if name == "" {
		return nil, fmt.Errorf("world name cannot be empty")
	}
	if folderPath == "" {
		return nil, fmt.Errorf("folder path cannot be empty")
	}

	entries, err := os.ReadDir(folderPath)
	if err != nil && !os.IsNotExist(err) {
		return nil, fmt.Errorf("read folder: %w", err)
	}
	if len(entries) > 0 {
		folderPath = filepath.Join(folderPath, name)
		if _, err := os.Stat(folderPath); err == nil {
			return nil, fmt.Errorf("world folder already exists: %s", folderPath)
		} else if !os.IsNotExist(err) {
			return nil, fmt.Errorf("check world folder: %w", err)
		}
	}

	now := time.Now().UTC()

	// Track created paths so we can roll back on partial failure.
	var written []string
	rollback := func() {
		for i := len(written) - 1; i >= 0; i-- {
			os.Remove(written[i])
		}
	}

	// Write world.json atomically.
	metaPath := filepath.Join(folderPath, worldMetaFileName)
	meta := WorldMeta{Name: name, CreatedAt: now, UpdatedAt: now}
	if err := writeJSONAtomic(metaPath, meta); err != nil {
		rollback()
		return nil, fmt.Errorf("create world.json: %w", err)
	}
	written = append(written, metaPath)

	// Create entities/ subdirectory.
	entitiesDir := filepath.Join(folderPath, "entities")
	if err := os.MkdirAll(entitiesDir, 0755); err != nil {
		rollback()
		return nil, fmt.Errorf("create entities dir: %w", err)
	}
	written = append(written, entitiesDir)

	// Create images/ subdirectory.
	imagesDir := filepath.Join(folderPath, "images")
	if err := os.MkdirAll(imagesDir, 0755); err != nil {
		rollback()
		return nil, fmt.Errorf("create images dir: %w", err)
	}
	written = append(written, imagesDir)

	// Create empty collection files.
	for _, filename := range []string{"links.json", "tags.json", "categories.json"} {
		p := filepath.Join(folderPath, filename)
		if err := writeJSONAtomic(p, []any{}); err != nil {
			rollback()
			return nil, fmt.Errorf("create %s: %w", filename, err)
		}
		written = append(written, p)
	}

	return &WorldInfo{Name: name, Path: folderPath, LastOpened: now}, nil
}

// OpenWorld validates that folderPath is a valid world (contains world.json),
// reads the world name from disk, and returns a WorldInfo suitable for adding
// to the recents list. Returns a user-friendly error if world.json is absent.
func (s *WorldService) OpenWorld(folderPath string) (*WorldInfo, error) {
	if folderPath == "" {
		return nil, fmt.Errorf("folder path cannot be empty")
	}

	metaPath := filepath.Join(folderPath, worldMetaFileName)
	data, err := os.ReadFile(metaPath)
	if os.IsNotExist(err) {
		return nil, fmt.Errorf("no LitGuardian world found at %q — missing world.json", folderPath)
	}
	if err != nil {
		return nil, fmt.Errorf("read world.json: %w", err)
	}

	var meta WorldMeta
	if err := json.Unmarshal(data, &meta); err != nil {
		return nil, fmt.Errorf("parse world.json: %w", err)
	}

	return &WorldInfo{
		Name:       meta.Name,
		Path:       folderPath,
		LastOpened: time.Now().UTC(),
	}, nil
}

// GetWorldMeta reads world.json from folderPath and returns the full metadata.
func (s *WorldService) GetWorldMeta(folderPath string) (*WorldMeta, error) {
	if folderPath == "" {
		return nil, fmt.Errorf("folder path cannot be empty")
	}
	metaPath := filepath.Join(folderPath, worldMetaFileName)
	data, err := os.ReadFile(metaPath)
	if os.IsNotExist(err) {
		return nil, fmt.Errorf("no LitGuardian world found at %q — missing world.json", folderPath)
	}
	if err != nil {
		return nil, fmt.Errorf("read world.json: %w", err)
	}
	var meta WorldMeta
	if err := json.Unmarshal(data, &meta); err != nil {
		return nil, fmt.Errorf("parse world.json: %w", err)
	}
	return &meta, nil
}

// UpdateWorldMeta updates the name and description in world.json, bumps UpdatedAt,
// and writes atomically. Returns the updated meta so the caller doesn't need a second round-trip.
func (s *WorldService) UpdateWorldMeta(folderPath, name, description string) (*WorldMeta, error) {
	if folderPath == "" {
		return nil, fmt.Errorf("folder path cannot be empty")
	}
	if name == "" {
		return nil, fmt.Errorf("world name cannot be empty")
	}

	metaPath := filepath.Join(folderPath, worldMetaFileName)
	data, err := os.ReadFile(metaPath)
	if os.IsNotExist(err) {
		return nil, fmt.Errorf("no LitGuardian world found at %q — missing world.json", folderPath)
	}
	if err != nil {
		return nil, fmt.Errorf("read world.json: %w", err)
	}

	var meta WorldMeta
	if err := json.Unmarshal(data, &meta); err != nil {
		return nil, fmt.Errorf("parse world.json: %w", err)
	}

	meta.Name = name
	meta.Description = description
	meta.UpdatedAt = time.Now().UTC()

	if err := writeJSONAtomic(metaPath, meta); err != nil {
		return nil, fmt.Errorf("update world.json: %w", err)
	}
	return &meta, nil
}

// writeJSONAtomic marshals v to indented JSON and writes it to path using
// a uniquely-named temp file + fsync + rename to prevent corruption on crash
// and to avoid races when multiple goroutines write to the same directory.
func writeJSONAtomic(path string, v any) error {
	dir := filepath.Dir(path)
	if err := os.MkdirAll(dir, 0755); err != nil {
		return err
	}
	data, err := json.MarshalIndent(v, "", "  ")
	if err != nil {
		return err
	}
	tmp, err := os.CreateTemp(dir, ".tmp-*")
	if err != nil {
		return err
	}
	tmpName := tmp.Name()
	ok := false
	defer func() {
		if !ok {
			tmp.Close()
			os.Remove(tmpName)
		}
	}()
	if _, err := tmp.Write(data); err != nil {
		return err
	}
	if err := tmp.Sync(); err != nil {
		return err
	}
	if err := tmp.Close(); err != nil {
		return err
	}
	ok = true
	return os.Rename(tmpName, path)
}
