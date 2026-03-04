package main

import "time"

// WorldInfo describes a recently opened world entry stored in app config.
type WorldInfo struct {
	Name       string    `json:"name"`
	Path       string    `json:"path"`
	LastOpened time.Time `json:"lastOpened"`
}

// AppConfig is the top-level app-level config persisted to the OS config dir.
// It is separate from world data — it lives at %APPDATA%/LitGuardian/config.json.
type AppConfig struct {
	RecentWorlds []WorldInfo `json:"recentWorlds"`
}
