package main

import (
	"encoding/json"
	"fmt"
	"io"
	"os"
	"os/exec"
	"path/filepath"
	"strings"

	"github.com/nitintf/claudewatch/internal/api"
	"github.com/nitintf/claudewatch/internal/auth"
	"github.com/nitintf/claudewatch/internal/config"
	"github.com/nitintf/claudewatch/internal/statusline"
	"github.com/nitintf/claudewatch/internal/theme"
)

func main() {
	if len(os.Args) > 1 {
		var err error
		switch os.Args[1] {
		case "install":
			err = install()
		case "uninstall":
			err = uninstall()
		}
		if err != nil {
			fmt.Fprintln(os.Stderr, err)
			os.Exit(1)
		}
		if os.Args[1] == "install" || os.Args[1] == "uninstall" {
			return
		}
	}

	if err := run(); err != nil {
		fmt.Fprintln(os.Stderr, err)
		os.Exit(1)
	}
}

func run() error {
	data, err := io.ReadAll(os.Stdin)
	if err != nil {
		return fmt.Errorf("reading stdin: %w", err)
	}
	if len(data) == 0 {
		return fmt.Errorf("no input on stdin — pipe Claude Code JSON or run: claudewatch install")
	}

	status, err := statusline.Parse(data)
	if err != nil {
		return err
	}

	cfg, err := config.Load()
	if err != nil {
		return err
	}

	t, err := theme.Load(cfg.Theme)
	if err != nil {
		t, _ = theme.LoadBuiltin("dracula")
	}

	// Read credentials for plan detection and API access.
	var plan string
	var usage *api.Usage
	creds, err := auth.Read()
	if err == nil {
		plan = auth.PlanName(creds.SubscriptionType)
		if creds.AccessToken != "" {
			usage, _ = api.FetchUsage(creds.AccessToken)
		}
	}

	output := statusline.Render(status, t, plan, usage)
	// Leading reset clears stale ANSI state from previous renders.
	// Non-breaking spaces prevent the terminal from collapsing whitespace.
	output = "\x1b[0m" + strings.ReplaceAll(output, " ", "\u00A0")
	// Extra newline adds visual padding below the status line.
	fmt.Println(output + "\n")
	return nil
}

func install() error {
	binaryPath, err := exec.LookPath("claudewatch")
	if err != nil {
		binaryPath, err = os.Executable()
		if err != nil {
			return fmt.Errorf("could not find claudewatch binary: %w", err)
		}
	}

	home, err := os.UserHomeDir()
	if err != nil {
		return err
	}

	settingsPath := filepath.Join(home, ".claude", "settings.json")

	settings := make(map[string]interface{})
	data, err := os.ReadFile(settingsPath)
	if err == nil {
		_ = json.Unmarshal(data, &settings)
	}

	settings["statusLine"] = map[string]interface{}{
		"type":    "command",
		"command": binaryPath,
	}

	_ = os.MkdirAll(filepath.Dir(settingsPath), 0o755)

	out, err := json.MarshalIndent(settings, "", "  ")
	if err != nil {
		return err
	}
	if err := os.WriteFile(settingsPath, out, 0o644); err != nil {
		return err
	}

	fmt.Printf("Installed claudewatch as Claude Code status line\n")
	fmt.Printf("  binary:   %s\n", binaryPath)
	fmt.Printf("  settings: %s\n", settingsPath)
	fmt.Printf("\nRestart Claude Code to see your status line.\n")
	return nil
}

func uninstall() error {
	home, err := os.UserHomeDir()
	if err != nil {
		return err
	}

	// Remove statusLine from Claude Code settings.
	settingsPath := filepath.Join(home, ".claude", "settings.json")
	data, err := os.ReadFile(settingsPath)
	if err == nil {
		settings := make(map[string]interface{})
		if err := json.Unmarshal(data, &settings); err == nil {
			delete(settings, "statusLine")
			out, err := json.MarshalIndent(settings, "", "  ")
			if err == nil {
				_ = os.WriteFile(settingsPath, out, 0o644)
			}
		}
	}

	// Remove config directory.
	cfgDir := filepath.Join(home, ".config", "claudewatch")
	_ = os.RemoveAll(cfgDir)

	// Remove usage cache.
	_ = os.Remove(filepath.Join(os.TempDir(), "claudewatch-usage.json"))

	fmt.Printf("Uninstalled claudewatch\n")
	fmt.Printf("  removed statusLine from %s\n", settingsPath)
	fmt.Printf("  removed config dir %s\n", cfgDir)
	fmt.Printf("\nRestart Claude Code to apply.\n")
	return nil
}
