package theme

import (
	"testing"
)

func TestLoadBuiltin(t *testing.T) {
	names := []string{
		"dracula", "catppuccin-mocha", "catppuccin-latte", "nord",
		"tokyo-night", "gruvbox", "solarized-dark", "solarized-light",
		"one-dark", "rosepine",
	}

	for _, name := range names {
		t.Run(name, func(t *testing.T) {
			th, err := LoadBuiltin(name)
			if err != nil {
				t.Fatalf("failed to load %q: %v", name, err)
			}
			if th.Name != name {
				t.Errorf("expected name %q, got %q", name, th.Name)
			}
			if th.Colors.Bg == "" {
				t.Error("expected non-empty bg color")
			}
			if th.Colors.Fg == "" {
				t.Error("expected non-empty fg color")
			}
			if th.Colors.Accent == "" {
				t.Error("expected non-empty accent color")
			}
			if th.Separators.Left == "" {
				t.Error("expected non-empty left separator")
			}
			if th.Icons.Model == "" {
				t.Error("expected non-empty model icon")
			}
		})
	}
}

func TestLoadBuiltinNotFound(t *testing.T) {
	_, err := LoadBuiltin("nonexistent-theme")
	if err == nil {
		t.Fatal("expected error for nonexistent theme")
	}
}

func TestListBuiltin(t *testing.T) {
	names, err := ListBuiltin()
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}

	if len(names) != 10 {
		t.Errorf("expected 10 built-in themes, got %d", len(names))
	}

	// Verify dracula is in the list
	found := false
	for _, n := range names {
		if n == "dracula" {
			found = true
			break
		}
	}
	if !found {
		t.Error("expected dracula in built-in themes")
	}
}

func TestLoadFallsBackToBuiltin(t *testing.T) {
	th, err := Load("dracula")
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if th.Name != "dracula" {
		t.Errorf("expected name %q, got %q", "dracula", th.Name)
	}
}
