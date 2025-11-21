# Changelog

All notable changes to veve-cli are documented in this file.

The format is based on [Keep a Changelog](https://keepachangelog.com/en/1.0.0/),
and this project adheres to [Semantic Versioning](https://semver.org/spec/v2.0.0.html).

## [0.2.0] - 2025-11-21

### Added

#### Core Features

- **Markdown to PDF Conversion** - Fast, reliable conversion using Pandoc
  - Support for complex markdown with code blocks, tables, and math
  - Multiple PDF engine support (pdflatex, xelatex, lualatex)
  - Configurable page settings (size, margins, font)
- **Built-in Themes** - Professional pre-configured themes
  - Default: Clean, minimal styling
  - Dark: Dark mode with high contrast
  - Academic: Professional academic paper style
- **Custom Theme Support** - Create and use custom themes
  - CSS-based styling system
  - YAML metadata for theme configuration
  - Local file path support with auto-directory creation
  - Automatic theme discovery in config directory
- **Theme Management** - Complete theme lifecycle management
  - List themes: `veve theme list`
  - Add themes from files: `veve theme add <name> <path>`
  - Download themes: `veve theme add <name> <url>`
  - Remove themes: `veve theme remove <name>`
  - View theme info: `veve theme info <name>`

#### Unix Composability

- **stdin/stdout Support** - Full Unix pipe compatibility
  - Read markdown from stdin: `cat input.md | veve`
  - Output PDF to stdout: `veve input.md | cat > output.pdf`
  - Batch processing: `ls *.md | xargs veve`
  - Piping to other tools: `veve input.md | pdfmerge`
- **Proper Exit Codes** - Meaningful status codes for scripting
  - 0: Success
  - 1: General error
  - 2: Invalid arguments
  - 64: Missing input file
  - 65: Pandoc not installed
  - 66: Theme not found
  - 67: Invalid theme
  - 70: PDF generation failed
  - 71: File creation failed

#### Shell Integration

- **Shell Completions** - First-class shell support
  - Bash completion script
  - Zsh completion script
  - Fish completion script
  - Auto-detection and installation: `veve completion [bash|zsh|fish]`
  - Completion for subcommands and flags
- **XDG Base Directory Support** - Standard configuration locations
  - Config: `~/.config/veve/`
  - State: `~/.local/state/veve/`
  - Cache: `~/.cache/veve/`
  - Windows: `%APPDATA%\veve\`

#### Professional Tooling

- **Release Infrastructure** - Automated builds and distribution
  - GoReleaser configuration for cross-platform builds
  - GitHub Actions CI/CD workflow
  - Pre-built binaries for macOS, Linux, Windows
  - Automatic checksums and archives
  - Semantic versioning support
  - Platform-specific installers
- **Logging** - Comprehensive debugging capabilities
  - Structured logging with levels (info, warn, error)
  - Debug mode for troubleshooting
  - Color-coded output
- **Code Quality**
  - Comprehensive test suite (100+ tests)
  - golangci-lint compliant
  - 80%+ code coverage
  - Proper error handling and validation

### Documentation

- **README.md** - 11KB comprehensive user guide
  - Installation instructions
  - Usage examples
  - Theme usage guide
  - Troubleshooting section
  - Performance metrics
  - Compatibility information
- **RELEASE.md** - 7.4KB release process guide
  - Step-by-step release instructions
  - GoReleaser configuration details
  - Local testing procedures
  - Installation methods
  - Verification instructions
  - Troubleshooting guide
- **THEME_DEVELOPMENT.md** - Theme creation guide
  - Theme structure and format
  - CSS styling system
  - Metadata configuration
  - Examples with source code
- **INTEGRATION.md** - Integration examples
  - Node.js/npm integration
  - Python integration
  - Bash script integration
  - GitHub Actions workflow examples
  - Docker usage examples

### Configuration

- **Default Configuration** - `~/.config/veve/config.yaml`
  - Default theme selection
  - PDF engine preference
  - Custom theme search paths
  - Output directory defaults

### Testing

- **Contract Tests** - 30+ CLI interface tests
  - Basic conversion functionality
  - Directory conversion
  - Error handling
  - Theme selection and validation
  - stdin/stdout piping
  - Batch processing
  - Exit code validation
  - Feature integration tests
- **Unit Tests** - 50+ unit tests
  - Theme loader functionality
  - Metadata parsing
  - Configuration loading
  - Registry management

## Version Information

- **Language**: Go 1.25+
- **Dependencies**: Pandoc 2.18+
- **Platforms**: macOS (10.15+), Linux (18.04+), Windows (10+)
- **Architecture**: amd64, arm64

## Release Artifacts

- `veve_Darwin_universal.tar.gz` - macOS (Intel + Apple Silicon)
- `veve_Darwin_x86_64.tar.gz` - macOS (Intel only)
- `veve_Darwin_arm64.tar.gz` - macOS (Apple Silicon only)
- `veve_Linux_x86_64.tar.gz` - Linux (amd64)
- `veve_Linux_arm64.tar.gz` - Linux (arm64)
- `veve_Windows_x86_64.zip` - Windows (amd64)
- `checksums.txt` - SHA256 checksums for all artifacts
- `veve_0.2.0_macos_all.tar.gz` - Homebrew formula artifacts

## Migration Guide

This is the first public release of veve-cli. No migration needed.

## Known Issues

None at this time. Please report issues at <https://github.com/madstone-tech/veve-cli/issues>

## Credits

Created by [Madstone Technology](https://github.com/madstone-tech)

## License

MIT License - See LICENSE file for details

---

## Future Roadmap

Potential features for future releases:

### v0.3.0 (Planned)

- Theme marketplace and online theme browser
- Built-in theme previewing
- Template support (cover pages, headers, footers)
- Output file optimization (compression)

### v0.4.0 (Planned)

- Web UI for drag-and-drop conversion
- Real-time preview of theme changes
- Batch conversion GUI
- Theme editor with live preview

### v1.0.0 (Planned)

- Plugin system for custom processors
- API server mode
- Docker image

---

For installation and usage instructions, see [README.md](README.md)

For release process details, see [docs/RELEASE.md](docs/RELEASE.md)
