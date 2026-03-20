package theme

import (
	"embed"
	"fmt"
	"os"
	"path/filepath"
	"strings"

	toml "github.com/pelletier/go-toml/v2"
)

//go:embed builtin/*.toml
var builtinFS embed.FS

// Colors defines the color palette for a theme.
type Colors struct {
	Bg      string `toml:"bg"`
	Fg      string `toml:"fg"`
	Accent  string `toml:"accent"`
	Success string `toml:"success"`
	Warning string `toml:"warning"`
	Error   string `toml:"error"`
	Info    string `toml:"info"`
	Muted   string `toml:"muted"`
}

// Separators defines the powerline glyph separators.
type Separators struct {
	Left      string `toml:"left"`
	LeftThin  string `toml:"left_thin"`
	Right     string `toml:"right"`
	RightThin string `toml:"right_thin"`
}

// Icons defines the nerd font icons used in the status line.
type Icons struct {
	Model   string `toml:"model"`
	Context string `toml:"context"`
	Cost    string `toml:"cost"`
	Tools   string `toml:"tools"`
	Warn    string `toml:"warn"`
	Ok      string `toml:"ok"`
}

// Theme represents a complete claudewatch theme.
type Theme struct {
	Name       string     `toml:"name"`
	Author     string     `toml:"author"`
	Colors     Colors     `toml:"colors"`
	Separators Separators `toml:"separators"`
	Icons      Icons      `toml:"icons"`
}

// LoadBuiltin loads a built-in theme by name.
func LoadBuiltin(name string) (Theme, error) {
	filename := name + ".toml"
	data, err := builtinFS.ReadFile("builtin/" + filename)
	if err != nil {
		return Theme{}, fmt.Errorf("built-in theme %q not found: %w", name, err)
	}

	var t Theme
	if err := toml.Unmarshal(data, &t); err != nil {
		return Theme{}, fmt.Errorf("parsing theme %q: %w", name, err)
	}

	return t, nil
}

// LoadFromFile loads a theme from a filesystem path.
func LoadFromFile(path string) (Theme, error) {
	data, err := os.ReadFile(path)
	if err != nil {
		return Theme{}, err
	}

	var t Theme
	if err := toml.Unmarshal(data, &t); err != nil {
		return Theme{}, fmt.Errorf("parsing theme %q: %w", path, err)
	}

	return t, nil
}

// Load loads a theme by name — checks user themes dir first, then built-in.
func Load(name string) (Theme, error) {
	// Check user themes directory first
	home, err := os.UserHomeDir()
	if err == nil {
		userPath := filepath.Join(home, ".config", "claudewatch", "themes", name+".toml")
		if _, statErr := os.Stat(userPath); statErr == nil {
			return LoadFromFile(userPath)
		}
	}

	return LoadBuiltin(name)
}

// ListBuiltin returns the names of all built-in themes.
func ListBuiltin() ([]string, error) {
	entries, err := builtinFS.ReadDir("builtin")
	if err != nil {
		return nil, err
	}

	var names []string
	for _, e := range entries {
		if !e.IsDir() && strings.HasSuffix(e.Name(), ".toml") {
			names = append(names, strings.TrimSuffix(e.Name(), ".toml"))
		}
	}

	return names, nil
}

// ListUser returns the names of all user-installed themes.
func ListUser() ([]string, error) {
	home, err := os.UserHomeDir()
	if err != nil {
		return nil, err
	}

	dir := filepath.Join(home, ".config", "claudewatch", "themes")
	entries, err := os.ReadDir(dir)
	if err != nil {
		if os.IsNotExist(err) {
			return nil, nil
		}
		return nil, err
	}

	var names []string
	for _, e := range entries {
		if !e.IsDir() && strings.HasSuffix(e.Name(), ".toml") {
			names = append(names, strings.TrimSuffix(e.Name(), ".toml"))
		}
	}

	return names, nil
}
