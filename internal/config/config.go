package config

import (
	"os"
	"path/filepath"

	toml "github.com/pelletier/go-toml/v2"
)

// Config holds the user's claudewatch configuration.
type Config struct {
	Theme      string `toml:"theme"`
	ShowPlan   *bool  `toml:"show_plan"`
	Show5h     *bool  `toml:"show_5h"`
	Show7d     *bool  `toml:"show_7d"`
	ShowExtra  *bool  `toml:"show_extra"`
	ShowCost   *bool  `toml:"show_cost"`
	ShowCwd    *bool  `toml:"show_cwd"`
	ShowBranch *bool  `toml:"show_branch"`
}

func boolPtr(b bool) *bool { return &b }

// DefaultConfig returns the default configuration.
func DefaultConfig() Config {
	return Config{
		Theme:      "dracula",
		ShowPlan:   boolPtr(true),
		Show5h:     boolPtr(true),
		Show7d:     boolPtr(true),
		ShowExtra:  boolPtr(true),
		ShowCost:   boolPtr(false),
		ShowCwd:    boolPtr(false),
		ShowBranch: boolPtr(false),
	}
}

// Enabled returns whether a segment is enabled (nil defaults to true).
func Enabled(b *bool) bool {
	if b == nil {
		return true
	}
	return *b
}

// Dir returns the config directory path (~/.config/claudewatch).
func Dir() (string, error) {
	home, err := os.UserHomeDir()
	if err != nil {
		return "", err
	}
	return filepath.Join(home, ".config", "claudewatch"), nil
}

// Path returns the config file path (~/.config/claudewatch/config.toml).
func Path() (string, error) {
	dir, err := Dir()
	if err != nil {
		return "", err
	}
	return filepath.Join(dir, "config.toml"), nil
}

// Load reads the config from disk. Returns default config if the file doesn't exist.
func Load() (Config, error) {
	cfg := DefaultConfig()

	p, err := Path()
	if err != nil {
		return cfg, err
	}

	data, err := os.ReadFile(p)
	if err != nil {
		if os.IsNotExist(err) {
			return cfg, nil
		}
		return cfg, err
	}

	if err := toml.Unmarshal(data, &cfg); err != nil {
		return DefaultConfig(), err
	}

	return cfg, nil
}

// Save writes the config to disk, creating the directory if needed.
func Save(cfg Config) error {
	p, err := Path()
	if err != nil {
		return err
	}

	if err := os.MkdirAll(filepath.Dir(p), 0o755); err != nil {
		return err
	}

	data, err := toml.Marshal(cfg)
	if err != nil {
		return err
	}

	return os.WriteFile(p, data, 0o644)
}
