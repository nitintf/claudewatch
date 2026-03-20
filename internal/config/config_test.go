package config

import "testing"

func TestDefaultConfig(t *testing.T) {
	cfg := DefaultConfig()
	if cfg.Theme != "dracula" {
		t.Errorf("expected theme dracula, got %q", cfg.Theme)
	}
}

func TestLoadReturnsDefaultWhenNoFile(t *testing.T) {
	cfg, err := Load()
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if cfg.Theme == "" {
		t.Error("expected non-empty theme")
	}
}
