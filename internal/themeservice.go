package internal

import (
	"encoding/json"
	"fmt"
	"os"
	"path/filepath"
	"sync"
)

const themesFileName = "themes.json"

// LayoutTheme is a named, reusable entity-detail layout preset. The built-in
// "standard" theme has Builtin=true and cannot be deleted.
type LayoutTheme struct {
	ID      string         `json:"id"`
	Name    string         `json:"name"`
	Builtin bool           `json:"builtin"`
	Modules []ModuleLayout `json:"modules"`
}

// ThemeStore is the root structure written to themes.json.
type ThemeStore struct {
	Version        int           `json:"version"`
	DefaultThemeID string        `json:"defaultThemeId"`
	Themes         []LayoutTheme `json:"themes"`
}

// ThemeService handles reading and writing themes.json for a world folder.
// It is stateless — all theme state lives on disk.
type ThemeService struct {
	mu sync.RWMutex
}

// GetThemes reads themes.json from folderPath. Returns nil (no error) if the
// file doesn't exist — the frontend synthesizes the built-in Standard theme in
// that case.
func (s *ThemeService) GetThemes(folderPath string) (*ThemeStore, error) {
	s.mu.RLock()
	defer s.mu.RUnlock()
	if folderPath == "" {
		return nil, fmt.Errorf("folder path cannot be empty")
	}
	path := filepath.Join(folderPath, themesFileName)
	data, err := os.ReadFile(path)
	if os.IsNotExist(err) {
		return nil, nil
	}
	if err != nil {
		return nil, fmt.Errorf("read themes.json: %w", err)
	}
	var store ThemeStore
	if err := json.Unmarshal(data, &store); err != nil {
		return nil, fmt.Errorf("parse themes.json: %w", err)
	}
	return &store, nil
}

// SaveThemes atomically writes store to themes.json in folderPath.
func (s *ThemeService) SaveThemes(folderPath string, store ThemeStore) error {
	s.mu.Lock()
	defer s.mu.Unlock()
	if folderPath == "" {
		return fmt.Errorf("folder path cannot be empty")
	}
	path := filepath.Join(folderPath, themesFileName)
	if err := writeJSONAtomic(path, store); err != nil {
		return fmt.Errorf("save themes.json: %w", err)
	}
	return nil
}
