---
name: uninstall
description: Uninstall claudewatch status line from Claude Code
allowed-tools: Bash, Read, AskUserQuestion
---

Uninstall claudewatch from Claude Code.

## Steps

1. **Confirm**: Ask the user if they are sure they want to uninstall claudewatch. This will remove the status line, config, and cache.

2. **Run uninstall**: Run `claudewatch uninstall`

3. **Remove binary**: Run `rm -f $(which claudewatch)`

4. **Done**: Tell the user claudewatch has been fully removed. Restart Claude Code to apply.
