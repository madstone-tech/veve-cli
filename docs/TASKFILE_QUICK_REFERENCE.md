# Task Quick Reference Card

## Essential Commands

### Development
```bash
task dev        # Download dependencies
task build      # Build binary
task run        # Run from source
task clean      # Remove artifacts
```

### Testing
```bash
task test-quick    # Unit tests only (1-2 sec) ⚡
task precommit     # Pre-commit checks (1-2 min) ⚡
task test          # All tests (5-10 min) 🐌
task test-contract # Contract tests only (5+ min) 🐌
```

### Code Quality
```bash
task fmt     # Format code
task lint    # Run linter
task vet     # Go vet
task all     # Format + lint + test + build
```

### Installation
```bash
task install        # Build and install
task uninstall      # Remove binary
task uninstall-full # Remove + config
```

### Themes
```bash
task config-init           # Create config dir
task themes-list           # List available themes
task themes-install FILE=. # Install theme
task config-clean          # Remove all config
```

## Workflow Recipes

### Quick Feedback (< 1 min)
```bash
task fmt
task test-quick
task build
```

### Before Committing (1-2 min)
```bash
task precommit
```

### Before Pushing (5-10 min)
```bash
task all
```

### Full CI Check (10-15 min)
```bash
task test     # With all contract tests
```

## Task Performance

| Command | Time | When to Use |
|---------|------|-------------|
| `test-quick` | 1-2s | Daily development |
| `precommit` | 1-2m | Before git commit |
| `test` | 5-10m | Before git push |
| `test-contract` | 5+ m | Full integration test |
| `all` | 5-10m | Before push |
| `ci` | 5-10m | CI pipeline |

## Useful Flags

```bash
task -l          # List main tasks
task --list-all  # List ALL tasks (35 total)
task -h          # Help
```

## Troubleshooting

**Tests timing out?**
- Use `task test-quick` instead (unit tests only)
- Full `task test` takes 5-10 minutes (normal)

**golangci-lint not found?**
- Run `task lint` (installs automatically)

**Permission denied installing?**
- Use `sudo task install` or edit INSTALL_PATH

**Can't find theme command?**
- Run `task config-init` first
- Check `task themes-list`

## All Available Tasks

Development: `dev` `run` `build` `build-release` `clean`
Testing: `test` `test-quick` `test-unit` `test-contract` `test-coverage` `test-theme` `test-verbose`
Quality: `fmt` `fmt-check` `lint` `vet` `precommit`
Install: `install` `install-verify` `uninstall` `uninstall-full`
Config: `config-init` `config-show` `config-clean`
Themes: `themes-list` `themes-install`
Docs: `docs` `help` `version`
Workflows: `all` `ci` `dev-setup` `distclean` `reset`

---

**For more details:** `task --list-all` or read `TASKFILE_GUIDE.md`
