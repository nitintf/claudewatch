---
name: update
description: Update claudewatch to the latest version
allowed-tools: Bash, Read
---

Update claudewatch to the latest version. Config is preserved.

## Steps

1. **Show current version**: Run `claudewatch version` to show the installed version.

2. **Update**: Run `claudewatch update` — this fetches the latest version via `go install`, re-registers with Claude Code, and preserves the existing config.

3. **Verify**: Run `claudewatch version` again to confirm the update.

4. **Done**: Tell the user the update is complete. If the version changed, tell them to restart Claude Code.
