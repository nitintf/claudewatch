---
name: config
description: Configure claudewatch status line settings interactively
allowed-tools: Bash, Read, Write, AskUserQuestion
---

Configure the claudewatch status line. Read the current config from `~/.config/claudewatch/config.toml` first to show ALL current values to the user.

## Step 1: Show current config

Read and display the full current config as a formatted summary, e.g.:

```
Current claudewatch config:
  theme       = dracula
  show_plan   = true
  show_5h     = true
  show_7d     = true
  show_extra  = true
  show_cost   = false
  show_cwd    = false
  show_branch = false
```

## Step 2: Ask what to change

Ask the user what they want to configure:
- **Theme** — Change the color theme
- **Segments** — Toggle which segments are shown/hidden
- **Both** — Change theme and segments

## Step 3a: Theme (if selected)

Ask which theme. Show ALL 10 options:
dracula, catppuccin-mocha, catppuccin-latte, nord, tokyo-night, gruvbox, solarized-dark, solarized-light, one-dark, rosepine

Since AskUserQuestion only supports 4 options, split into two questions:
1. First ask with 4 popular options + indicate "Other" for more
2. If they pick "Other", show the remaining themes

## Step 3b: Segments (if selected)

Show all segment toggles with their current values. Ask which segments to toggle (multi-select):
- show_plan (Plan name in model bracket)
- show_5h (5-hour usage quota)
- show_7d (7-day usage quota)
- show_extra (Pay-as-you-go extra usage)
- show_cost (Session cost)
- show_cwd (Working directory name)
- show_branch (Git branch)

Note: model name and context window are always shown and cannot be disabled.

Selected segments will be TOGGLED (true→false, false→true).

## Step 4: Write config

Write the FULL updated config to `~/.config/claudewatch/config.toml` with ALL keys explicitly set. Example:

```toml
theme = "dracula"
show_plan = true
show_5h = true
show_7d = true
show_extra = true
show_cost = false
show_cwd = false
show_branch = false
```

## Step 5: Confirm

Show the updated config summary. If the theme changed, tell them to restart Claude Code. Segment changes take effect on the next status line refresh (automatic).
