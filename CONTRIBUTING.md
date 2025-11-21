# Contributing to veve-cli

Thank you for your interest in contributing to veve-cli! This document provides guidelines and instructions for contributing to the project.

## Code of Conduct

Please be respectful and constructive in all interactions. We're building a welcoming community.

## How to Contribute

### Reporting Bugs

If you find a bug, please create an issue on GitHub with:

1. **Title**: Clear, concise description of the bug
2. **Description**: 
   - Steps to reproduce
   - Expected behavior
   - Actual behavior
   - System information (OS, Go version, Pandoc version)
3. **Example**: If possible, provide a minimal example to reproduce

Example:
```
**Bug**: PDF generation fails with special characters in filename

**Steps to reproduce**:
1. Create file named "test-ñ.md"
2. Run: veve test-ñ.md
3. Error occurs

**Expected**: PDF should be created
**Actual**: File not found error

**System**: macOS 12.5, Go 1.20, Pandoc 2.19
```

### Suggesting Features

For feature requests:

1. **Title**: Feature name
2. **Description**: 
   - Problem it solves
   - Use cases
   - Proposed solution (optional)
3. **Context**: Any relevant information

### Submitting Code

#### Setup Development Environment

```bash
# Clone the repository
git clone https://github.com/madstone-tech/veve-cli.git
cd veve-cli

# Install dependencies
go mod download

# Install Pandoc (required for testing)
brew install pandoc  # macOS
# or
apt-get install pandoc  # Linux
# or
choco install pandoc  # Windows

# Verify setup
go test ./...
```

#### Development Workflow

1. **Create a branch** for your feature/fix:
   ```bash
   git checkout -b feature/description
   # or
   git checkout -b fix/issue-number
   ```

2. **Make changes** following [Code Style](#code-style):
   ```bash
   # Edit files
   # Test locally
   go test ./...
   
   # Format code
   gofmt -w ./...
   
   # Lint code
   golangci-lint run ./...
   ```

3. **Write tests**:
   - Unit tests in `internal/`
   - Contract tests in `tests/contract/`
   - Aim for 80%+ coverage

4. **Commit with clear messages**:
   ```bash
   git commit -m "feat: Add feature name
   
   Description of changes and why.
   
   Fixes #123
   "
   ```

5. **Push and create a Pull Request**:
   ```bash
   git push origin feature/description
   ```

#### Code Style

Follow standard Go conventions:

- **Naming**: `camelCase` for functions/variables, `PascalCase` for exported
- **Formatting**: `gofmt` (automatic via pre-commit)
- **Linting**: `golangci-lint` (all checks must pass)
- **Comments**: 
  - Exported functions/types must have comments
  - Comments should be clear and concise
  - Avoid redundant comments

Example:
```go
// processMarkdown converts markdown to PDF using Pandoc.
// It returns an error if Pandoc is not installed or conversion fails.
func processMarkdown(inputPath, outputPath, theme string) error {
    // implementation
}
```

#### Testing Requirements

All code changes must include tests:

1. **Unit Tests** - Test functions in isolation
   ```bash
   go test ./internal/...
   ```

2. **Contract Tests** - Test CLI interface
   ```bash
   go test ./tests/contract/...
   ```

3. **All Tests** - Run full suite
   ```bash
   go test ./...
   ```

Minimum coverage: 80%

Check coverage:
```bash
go test -cover ./...
```

#### Documentation

Update documentation for:

- New features → Update README.md
- Theme features → Update docs/THEME_DEVELOPMENT.md
- Integration examples → Update docs/INTEGRATION.md
- Release process changes → Update docs/RELEASE.md

## Project Structure

```
veve-cli/
├── cmd/veve/              # CLI commands
│   ├── main.go           # Entry point
│   ├── root.go           # Root command
│   ├── convert.go        # Convert command
│   └── theme.go          # Theme commands
├── internal/             # Internal packages
│   ├── config/          # Configuration loading
│   ├── converter/       # PDF conversion logic
│   ├── logging/         # Logging utilities
│   ├── theme/           # Theme management
│   └── errors.go        # Error handling
├── tests/               # Test suites
│   └── contract/        # CLI contract tests
├── scripts/             # Build and utility scripts
├── docs/                # Documentation
├── themes/              # Built-in themes
└── .github/workflows/   # CI/CD workflows
```

## Theme Development

To develop custom themes:

1. See [docs/THEME_DEVELOPMENT.md](docs/THEME_DEVELOPMENT.md)
2. Create theme in `themes/` directory
3. Add tests in `tests/contract/`
4. Update README examples

### Adding Built-in Themes

1. Create CSS file in `themes/` directory
2. Add metadata (optional YAML comment)
3. Update `internal/theme/metadata.go` if adding options
4. Add contract tests for the theme
5. Document in README.md

## Review Process

When you submit a PR:

1. **Automated checks**:
   - Tests must pass
   - Code coverage must be 80%+
   - Linting must pass
   - Formatting must be correct

2. **Code review**:
   - Project maintainer reviews changes
   - May request modifications
   - Discusses design decisions

3. **Approval and merge**:
   - Once approved, PR will be merged
   - Contributor will be credited
   - Feature may be included in next release

## Building and Testing

### Run All Tests
```bash
go test ./...
```

### Run Specific Tests
```bash
go test ./internal/theme
go test ./tests/contract -run TestConvertBasic
```

### Build Binary
```bash
go build -o veve ./cmd/veve
./veve --version
```

### Test Locally
```bash
# Build
go build -o veve ./cmd/veve

# Test basic conversion
echo "# Hello" > test.md
./veve test.md -o test.pdf

# Test with theme
./veve test.md --theme dark -o test.pdf

# Test CLI help
./veve --help
./veve theme --help
```

## Release Process

For maintainers releasing new versions:

1. Update CHANGELOG.md
2. Create git tag: `git tag -a v0.x.x -m "Release v0.x.x"`
3. Push tag: `git push origin v0.x.x`
4. GitHub Actions automatically builds and releases

See [docs/RELEASE.md](docs/RELEASE.md) for details.

## Questions?

- Open a GitHub issue for bug reports
- Open a GitHub discussion for questions
- Check existing issues and documentation first

## Recognition

Contributors will be:
- Credited in release notes
- Added to contributor list
- Thanked in documentation

## License

By contributing, you agree that your contributions will be licensed under the MIT License.

---

Thank you for making veve-cli better! 🎉
