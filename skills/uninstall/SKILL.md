---
name: uninstall
description: Uninstall claudewatch status line from Claude Code
allowed-tools: Bash, Read, AskUserQuestion
---

Uninstall claudewatch from Claude Code.

## Steps

1. **Confirm**: Ask the user if they are sure they want to uninstall claudewatch. This will remove the status line, cache, and binary. Config is preserved at `~/.config/claudewatch/config.toml` for future reinstalls.

2. **Run uninstall**: Run `claudewatch uninstall` — this removes the statusLine from settings, clears the cache, and removes the binary. Config is kept.

3. **Done**: Tell the user claudewatch has been removed. Config was preserved at `~/.config/claudewatch/` in case they reinstall. Restart Claude Code to apply.
