package auth

import (
	"context"
	"crypto/sha256"
	"encoding/json"
	"fmt"
	"os"
	"os/exec"
	"path/filepath"
	"runtime"
	"strings"
	"time"
)

// Credentials holds OAuth credentials from Claude Code.
type Credentials struct {
	AccessToken      string
	SubscriptionType string
}

type rawCredentials struct {
	ClaudeAiOauth struct {
		AccessToken      string `json:"accessToken"`
		SubscriptionType string `json:"subscriptionType"`
	} `json:"claudeAiOauth"`
}

// Read loads Claude Code OAuth credentials from keychain or file.
func Read() (Credentials, error) {
	if runtime.GOOS == "darwin" {
		if creds, err := readKeychain(); err == nil {
			return creds, nil
		}
	}
	return readFile()
}

// PlanName maps a subscription type to a display name.
func PlanName(subType string) string {
	lower := strings.ToLower(subType)
	switch {
	case strings.Contains(lower, "max"):
		return "Max"
	case strings.Contains(lower, "pro"):
		return "Pro"
	case strings.Contains(lower, "team"):
		return "Team"
	case strings.Contains(lower, "enterprise"):
		return "Enterprise"
	default:
		return ""
	}
}

func readKeychain() (Credentials, error) {
	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()
	out, err := exec.CommandContext(ctx,
		"/usr/bin/security", "find-generic-password",
		"-s", keychainService(), "-w",
	).Output()
	if err != nil {
		return Credentials{}, fmt.Errorf("keychain: %w", err)
	}
	return parseRaw(out)
}

func readFile() (Credentials, error) {
	dir := os.Getenv("CLAUDE_CONFIG_DIR")
	if dir == "" {
		home, err := os.UserHomeDir()
		if err != nil {
			return Credentials{}, err
		}
		dir = filepath.Join(home, ".claude")
	}
	data, err := os.ReadFile(filepath.Join(dir, ".credentials.json"))
	if err != nil {
		return Credentials{}, err
	}
	return parseRaw(data)
}

func parseRaw(data []byte) (Credentials, error) {
	var raw rawCredentials
	if err := json.Unmarshal(data, &raw); err != nil {
		return Credentials{}, err
	}
	return Credentials{
		AccessToken:      raw.ClaudeAiOauth.AccessToken,
		SubscriptionType: raw.ClaudeAiOauth.SubscriptionType,
	}, nil
}

func keychainService() string {
	name := "Claude Code-credentials"
	if dir := os.Getenv("CLAUDE_CONFIG_DIR"); dir != "" {
		h := sha256.Sum256([]byte(dir))
		name += fmt.Sprintf("-%x", h[:4])
	}
	return name
}
