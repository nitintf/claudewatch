---
name: setup
description: Install and configure claudewatch status line for Claude Code
allowed-tools: Bash, Read, Write, AskUserQuestion
---

Install claudewatch — a themed status line for Claude Code.

## Steps

1. **Check Go**: Verify Go is installed (`go version`). If not installed, tell the user to install Go first from https://go.dev/dl/ and stop.

2. **Install binary**: Run `go install github.com/nitintf/claudewatch@latest`

3. **Register with Claude Code**: Run `claudewatch install`

4. **Verify**: Read `~/.claude/settings.json` and confirm statusLine is set.

5. **Choose theme**: Ask which theme they want:
   Options: dracula, catppuccin-mocha, catppuccin-latte, nord, tokyo-night, gruvbox, solarized-dark, solarized-light, one-dark, rosepine

6. **Write config**: Write to `~/.config/claudewatch/config.toml`:
   ```toml
   theme = "<chosen-theme>"
   ```

7. **Done**: Tell the user to restart Claude Code. Plan type and usage limits are auto-detected from credentials — no configuration needed.
