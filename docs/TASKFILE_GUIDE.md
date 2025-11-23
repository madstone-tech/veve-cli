# veve-cli Taskfile Guide

This project uses [Task](https://taskfile.dev/) to streamline development workflows. Task is a simple, portable alternative to Make.

## Installation

If you don't have Task installed, install it:

```bash
# macOS
brew install go-task/tap/go-task

# Linux
sh -c "$(curl --location https://taskfile.dev/install.sh)" -- -d -b /usr/local/bin

# Or use Go
go install github.com/go-task/task/v3/cmd/task@latest
```

## Available Tasks

### Development Tasks

```bash
task dev              # Setup development environment (download modules)
task run              # Run veve-cli from source code
task build            # Build veve-cli binary
task build-release    # Build optimized release binary (stripped)
task clean            # Remove build artifacts and temp files
```

### Testing Tasks

```bash
task test             # Run all tests with race detector (slow - 5-10 minutes, includes contract tests)
task test-unit        # Run unit tests only (internal/) - same as test-quick
task test-quick       # Run quick tests only (internal/ only, no contract tests) - 1-2 seconds
task test-contract    # Run contract tests only (tests/contract/) - slow, 5+ minutes
task test-coverage    # Run tests with coverage report (HTML) - slow
task test-theme       # Run theme-specific tests (metadata, fonts, parsing) - fast
task test-verbose     # Run all tests with verbose output - slow
```

**Note on test performance:**
- **Fast** (< 5 seconds): `test-quick`, `test-unit`, `test-theme`
- **Slow** (5-10 minutes): `test`, `test-contract`, `test-coverage`

The contract tests are slow because they invoke the actual veve command which performs PDF conversion. For quick feedback during development, use `task test-quick` or `task test-unit`.

### Code Quality Tasks

```bash
task fmt              # Format code with gofmt
task fmt-check        # Check if code needs formatting
task lint             # Run golangci-lint linter
task vet              # Run go vet analyzer
task precommit        # Run all pre-commit checks (fmt, lint, test)
```

### Installation Tasks

```bash
task install          # Build and install veve-cli to /usr/local/bin
task install-verify   # Verify veve-cli is installed and working
task uninstall        # Remove veve-cli from system
task uninstall-full   # Remove veve-cli and all local config
```

### Configuration Tasks

```bash
task config-init      # Create ~/.config/veve/themes directory
task config-show      # Show veve config directory contents
task config-clean     # Remove all veve config and themes
task themes-list      # List available themes
task themes-install   # Install theme file (FILE=path/to/theme.css)
```

### Documentation Tasks

```bash
task docs             # Show theme development guide
task help             # Show veve-cli help
task version          # Show veve-cli version
```

### Workflow Tasks

```bash
task all              # Format, lint, test, and build everything
task ci               # Run CI pipeline checks (like pre-commit)
task dev-setup        # Setup complete development environment
task distclean        # Deep clean (remove vendor, cache, artifacts)
task reset            # Reset to pristine state (download fresh dependencies)
```

## Quick Start

### First Time Setup

```bash
# Setup development environment
task dev-setup

# Run tests
task test

# Build the binary
task build
```

### Before Committing

```bash
# Quick pre-commit checks (fast - 1-2 minutes)
task precommit

# Or for a more thorough check including contract tests (slow - 5-10 minutes):
task all

# Individual steps:
task fmt
task lint
task test-quick        # Unit tests only (fast)
task test-contract     # Full contract tests (slow)
```

### Installation

```bash
# Build and install to /usr/local/bin
task install

# Verify installation
task install-verify

# Later, to remove
task uninstall
```

### Working with Themes

```bash
# Initialize config
task config-init

# List themes
task themes-list

# Install a custom theme
task themes-install FILE=~/my-theme.css

# View config
task config-show

# Clean up
task config-clean
```

## Common Workflows

### Daily Development

```bash
# Format code
task fmt

# Run quick unit tests
task test-unit

# Build locally
task build

# Run the app
./veve --help
```

### Before Push

```bash
# Run all checks
task precommit

# Or more comprehensive
task all
```

### Full Test Suite

```bash
# Unit + Contract tests with race detector
task test

# With coverage report
task test-coverage

# Theme-specific only
task test-theme
```

### CI/CD

```bash
# Runs format check, vet, lint, and test
task ci
```

## Advanced Usage

### Custom Build Version

```bash
# Pass VERSION variable
task -s build VERSION=1.0.0
```

### Installing to Custom Path

Edit `Taskfile.yml` and change the `INSTALL_PATH` variable:

```yaml
vars:
  INSTALL_PATH: /custom/path
```

### Running Tasks in Parallel

Task supports `--parallel` flag:

```bash
task --parallel test-unit test-contract
```

### Listing Tasks with Details

```bash
task --list-all        # Show all tasks with descriptions
task --list            # Show main tasks
```

## Troubleshooting

### Permission Denied on `task install`

The `/usr/local/bin` directory may not be writable. Use sudo:

```bash
sudo task install
```

Or install to a user directory:

Edit Taskfile.yml and set `INSTALL_PATH: ~/.local/bin`

### golangci-lint Not Found

Task will automatically install it:

```bash
task lint  # Installs golangci-lint if needed
```

### Config Directory Issues

Manually create config:

```bash
mkdir -p ~/.config/veve/themes
task config-show
```

### Contract Tests Timing Out

Contract tests are slow (5+ minutes) because they run the actual `veve` command for PDF conversion. They have a 300-second timeout.

If tests timeout:

```bash
# Run quick tests instead
task test-quick         # Fast unit tests

# Or run just theme tests
task test-theme         # Just metadata/fonts/parsing tests

# For CI, allow more time
go test -v -timeout 600s ./...
```

The full `task test` command includes contract tests and can take 5-10 minutes depending on system performance.

## Task Variables

The Taskfile defines these variables:

- `BINARY_NAME`: veve (the output binary name)
- `MAIN_PKG`: ./cmd/veve (main package path)
- `VERSION`: dev (version string for builds)
- `INSTALL_PATH`: /usr/local/bin (installation directory)

Override them when running tasks:

```bash
task -s build VERSION=1.0.0
```

## Integration with IDEs

### VS Code

Install the [Task Explorer extension](https://marketplace.visualstudio.com/items?itemName=actano.vscode-taskexplorer)

### GoLand / IntelliJ IDEA

Task is recognized as an external tool. Configure it under:
Settings → Tools → Task

## See Also

- [Task Documentation](https://taskfile.dev/)
- [veve-cli README](README.md)
- [Theme Development Guide](docs/THEME_DEVELOPMENT.md)
