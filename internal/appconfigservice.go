package internal

import (
	"context"
	"encoding/json"
	"os"
	"path/filepath"
	"sort"
	"strings"
	"sync"
	"time"

	"github.com/wailsapp/wails/v3/pkg/application"
)

const (
	appConfigDirName  = "LitGuardian"
	appConfigFileName = "config.json"
	maxRecentWorlds   = 10
)

// AppConfigService manages app-level configuration (recent worlds list).
// Config is persisted at the OS user config directory:
//   - Windows: %APPDATA%/LitGuardian/config.json
//   - macOS:   ~/Library/Application Support/LitGuardian/config.json
//   - Linux:   ~/.config/LitGuardian/config.json
type AppConfigService struct {
	mu       sync.RWMutex
	config   AppConfig
	filePath string
}

// ServiceStartup implements the Wails 3 service lifecycle interface.
// Called automatically on app start — loads config from disk.
func (s *AppConfigService) ServiceStartup(ctx context.Context, options application.ServiceOptions) error {
	configDir, err := os.UserConfigDir()
	if err != nil {
		return err
	}
	s.filePath = filepath.Join(configDir, appConfigDirName, appConfigFileName)
	return s.load()
}

// GetRecentWorlds returns the recent worlds list sorted by LastOpened descending.
func (s *AppConfigService) GetRecentWorlds() []WorldInfo {
	s.mu.RLock()
	defer s.mu.RUnlock()

	result := make([]WorldInfo, len(s.config.RecentWorlds))
	copy(result, s.config.RecentWorlds)
	sort.Slice(result, func(i, j int) bool {
		return result[i].LastOpened.After(result[j].LastOpened)
	})
	return result
}

// canonPath returns a normalised, case-folded representation of p for
// duplicate detection. The stored path is always the original; only
// comparisons use the canonical form.
func canonPath(p string) string {
	return strings.ToLower(filepath.Clean(p))
}

// AddRecentWorld upserts a world entry by path, updates LastOpened, caps the
// list at maxRecentWorlds, and persists to disk.
func (s *AppConfigService) AddRecentWorld(name, path string) error {
	s.mu.Lock()
	defer s.mu.Unlock()

	now := time.Now().UTC()
	canon := canonPath(path)

	// Upsert: update existing entry if canonical path matches.
	for i, w := range s.config.RecentWorlds {
		if canonPath(w.Path) == canon {
			s.config.RecentWorlds[i].Name = name
			s.config.RecentWorlds[i].LastOpened = now
			return s.save()
		}
	}

	// Prepend new entry.
	s.config.RecentWorlds = append([]WorldInfo{{Name: name, Path: path, LastOpened: now}}, s.config.RecentWorlds...)

	// Cap the list.
	if len(s.config.RecentWorlds) > maxRecentWorlds {
		s.config.RecentWorlds = s.config.RecentWorlds[:maxRecentWorlds]
	}

	return s.save()
}

// RemoveRecentWorld removes a world entry by path and persists to disk.
func (s *AppConfigService) RemoveRecentWorld(path string) error {
	s.mu.Lock()
	defer s.mu.Unlock()

	canon := canonPath(path)
	filtered := s.config.RecentWorlds[:0]
	for _, w := range s.config.RecentWorlds {
		if canonPath(w.Path) != canon {
			filtered = append(filtered, w)
		}
	}
	s.config.RecentWorlds = filtered
	return s.save()
}

// load reads config from disk. A missing file is not an error — it just means
// no worlds have been opened yet.
func (s *AppConfigService) load() error {
	data, err := os.ReadFile(s.filePath)
	if os.IsNotExist(err) {
		s.config = AppConfig{}
		return nil
	}
	if err != nil {
		return err
	}
	return json.Unmarshal(data, &s.config)
}

// save atomically writes the config to disk using a uniquely-named temp file +
// fsync + rename to prevent corruption and races with concurrent saves.
// Must be called with s.mu held.
func (s *AppConfigService) save() error {
	dir := filepath.Dir(s.filePath)
	if err := os.MkdirAll(dir, 0755); err != nil {
		return err
	}

	data, err := json.MarshalIndent(s.config, "", "  ")
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
	return os.Rename(tmpName, s.filePath)
}
