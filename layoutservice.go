package main

import (
	"encoding/json"
	"fmt"
	"os"
	"path/filepath"
)

const layoutFileName = "layout.json"

// ModuleLayout stores the persisted layout config for a single module.
type ModuleLayout struct {
	ID      string `json:"id"`
	ColSpan int    `json:"colSpan"`
	Visible bool   `json:"visible"`
}

// ViewLayout stores the ordered list of modules for a single view.
type ViewLayout struct {
	Modules []ModuleLayout `json:"modules"`
}

// DashboardLayout is the root structure written to layout.json.
type DashboardLayout struct {
	Version int                     `json:"version"`
	Views   map[string]ViewLayout   `json:"views"`
}

// LayoutService handles reading and writing layout.json for a world folder.
// It is stateless — all layout state lives on disk.
type LayoutService struct{}

// GetLayout reads layout.json from folderPath. Returns nil (no error) if the
// file doesn't exist — the frontend applies registry defaults in that case.
func (s *LayoutService) GetLayout(folderPath string) (*DashboardLayout, error) {
	if folderPath == "" {
		return nil, fmt.Errorf("folder path cannot be empty")
	}
	path := filepath.Join(folderPath, layoutFileName)
	data, err := os.ReadFile(path)
	if os.IsNotExist(err) {
		return nil, nil
	}
	if err != nil {
		return nil, fmt.Errorf("read layout.json: %w", err)
	}
	var layout DashboardLayout
	if err := json.Unmarshal(data, &layout); err != nil {
		return nil, fmt.Errorf("parse layout.json: %w", err)
	}
	return &layout, nil
}

// SaveLayout atomically writes layout to layout.json in folderPath.
func (s *LayoutService) SaveLayout(folderPath string, layout DashboardLayout) error {
	if folderPath == "" {
		return fmt.Errorf("folder path cannot be empty")
	}
	path := filepath.Join(folderPath, layoutFileName)
	if err := writeJSONAtomic(path, layout); err != nil {
		return fmt.Errorf("save layout.json: %w", err)
	}
	return nil
}
