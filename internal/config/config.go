// Package config handles JSON persistence of user settings.
package config

import (
	"encoding/json"
	"os"
	"path/filepath"
)

// Config holds all persisted user settings.
type Config struct {
	ModDir      string         `json:"mod_dir"`
	MO2Dir      string         `json:"mo2_dir"`      // MO2 instance root directory
	MO2Profile  string         `json:"mo2_profile"`  // last-used MO2 profile name
	Priorities  map[string]int `json:"priorities"`   // archive name → priority (1-99, 0 = unset)
	BackupLimit int            `json:"backup_limit"` // max modlist backups to keep; 0 = use default (20)
}

// EffectiveBackupLimit returns the configured backup limit, or 20 when unset.
func (c *Config) EffectiveBackupLimit() int {
	if c.BackupLimit <= 0 {
		return 20
	}

	return c.BackupLimit
}

// DefaultPath returns the path to the config file in the user's config dir.
func DefaultPath() string {
	dir, _ := os.UserConfigDir()

	return filepath.Join(dir, "cp77-modorder", "config.json")
}

// Load reads and deserialises the config file at path.
// Returns an empty Config (not an error) when the file does not exist yet.
func Load(path string) (*Config, error) {
	cfg := &Config{Priorities: make(map[string]int)}
	data, err := os.ReadFile(path)
	if os.IsNotExist(err) {
		return cfg, nil
	}
	if err != nil {
		return nil, err
	}

	if err := json.Unmarshal(data, cfg); err != nil {
		return nil, err
	}

	if cfg.Priorities == nil {
		cfg.Priorities = make(map[string]int)
	}

	return cfg, nil
}

// Save serialises and writes the config to path, creating parent dirs as needed.
func (c *Config) Save(path string) error {
	if err := os.MkdirAll(filepath.Dir(path), 0o755); err != nil {
		return err
	}

	data, err := json.MarshalIndent(c, "", "  ")
	if err != nil {
		return err
	}

	return os.WriteFile(path, data, 0o644)
}
