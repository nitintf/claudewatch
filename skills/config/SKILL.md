---
name: config
description: Configure claudewatch status line settings interactively
allowed-tools: Bash, Read, Write, AskUserQuestion
---

Configure the claudewatch status line. Read the current config from `~/.config/claudewatch/config.toml` first to show current values, then ask the user what they want to change.

## Questions to ask

Ask the question, showing the current value as context:

1. **Theme** — Which theme?
   Options: dracula, catppuccin-mocha, catppuccin-latte, nord, tokyo-night, gruvbox, solarized-dark, solarized-light, one-dark, rosepine
   Show current value.

## After collecting answers

Write the updated config to `~/.config/claudewatch/config.toml`:

```toml
theme = "<theme>"
```

Tell the user the config has been updated. If they changed the theme, tell them to restart Claude Code to see the new theme.

Note: Plan type and usage limits are auto-detected from your Claude Code credentials — no configuration needed.
