package main

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
	Name      string    `json:"name"`
	CreatedAt time.Time `json:"createdAt"`
	UpdatedAt time.Time `json:"updatedAt"`
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
// It creates: world.json, entities/, links.json, tags.json, categories.json.
// Returns a WorldInfo suitable for adding to the recents list.
func (s *WorldService) CreateWorld(name, folderPath string) (*WorldInfo, error) {
	if name == "" {
		return nil, fmt.Errorf("world name cannot be empty")
	}
	if folderPath == "" {
		return nil, fmt.Errorf("folder path cannot be empty")
	}

	// Guard against creating a world on top of an existing one.
	metaPath := filepath.Join(folderPath, worldMetaFileName)
	if _, err := os.Stat(metaPath); err == nil {
		return nil, fmt.Errorf("folder already contains a world: %s", folderPath)
	}

	now := time.Now().UTC()

	// Write world.json atomically.
	meta := WorldMeta{Name: name, CreatedAt: now, UpdatedAt: now}
	if err := writeJSONAtomic(metaPath, meta); err != nil {
		return nil, fmt.Errorf("create world.json: %w", err)
	}

	// Create entities/ subdirectory.
	if err := os.MkdirAll(filepath.Join(folderPath, "entities"), 0755); err != nil {
		return nil, fmt.Errorf("create entities dir: %w", err)
	}

	// Create empty collection files.
	for _, filename := range []string{"links.json", "tags.json", "categories.json"} {
		p := filepath.Join(folderPath, filename)
		if err := writeJSONAtomic(p, []any{}); err != nil {
			return nil, fmt.Errorf("create %s: %w", filename, err)
		}
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

// writeJSONAtomic marshals v to indented JSON and writes it to path using
// a temp file + rename to prevent corruption on crash.
func writeJSONAtomic(path string, v any) error {
	if err := os.MkdirAll(filepath.Dir(path), 0755); err != nil {
		return err
	}
	data, err := json.MarshalIndent(v, "", "  ")
	if err != nil {
		return err
	}
	tmp := path + ".tmp"
	if err := os.WriteFile(tmp, data, 0644); err != nil {
		return err
	}
	return os.Rename(tmp, path)
}
