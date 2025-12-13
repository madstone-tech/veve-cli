# veve-cli - Markdown to PDF Converter with Theme Support

[![CI](https://github.com/madstone-tech/veve-cli/actions/workflows/ci.yml/badge.svg)](https://github.com/madstone-tech/veve-cli/actions/workflows/ci.yml)
[![Release](https://github.com/madstone-tech/veve-cli/actions/workflows/release.yml/badge.svg)](https://github.com/madstone-tech/veve-cli/actions/workflows/release.yml)
[![Go Version](https://img.shields.io/badge/go-1.25+-blue)](https://golang.org)
[![License](https://img.shields.io/badge/license-MIT-green)](LICENSE)

A fast, themeable markdown-to-PDF converter built with Go. Convert your markdown files to beautiful PDFs with built-in and custom themes. Perfect for documentation, reports, and technical writing.

## Features

- 📄 **Markdown to PDF** - Fast, reliable conversion via Pandoc
- 🎨 **Theme Support** - 3 built-in themes (default, dark, academic) + unlimited custom themes
- ⚙️ **Theme Management** - List, add, and remove themes via CLI
- 📝 **Custom Themes** - Create themes in `~/.config/veve/themes/` with YAML metadata
- 🌐 **Remote Images** - Automatically download and embed remote images from HTTP/HTTPS URLs
- 🔀 **Unix Composability** - Full stdin/stdout support for piping and scripting
- 🛠️ **Configuration** - TOML-based config with XDG Base Directory support
- 🚀 **Cross-Platform** - macOS, Linux, Windows support
- 🔧 **Pandoc Flexibility** - Configurable PDF engines and Pandoc options
- 📦 **Shell Completions** - bash, zsh, and fish autocompletion support

## Installation

### macOS (Homebrew - Recommended)

```bash
brew tap madstone-tech/tap
brew install veve
```

Shell completions are automatically installed with Homebrew.

### Linux

#### Download Pre-built Binary

```bash
# AMD64
curl -sL https://github.com/madstone-tech/veve-cli/releases/latest/download/veve_Linux_x86_64.tar.gz | tar xz
sudo mv veve /usr/local/bin/

# ARM64
curl -sL https://github.com/madstone-tech/veve-cli/releases/latest/download/veve_Linux_arm64.tar.gz | tar xz
sudo mv veve /usr/local/bin/
```

#### From Package Managers

```bash
# Debian/Ubuntu (when available)
sudo apt install veve

# Fedora (when available)
sudo dnf install veve
```

### Windows

Download the latest `.zip` file from [releases](https://github.com/madstone-tech/veve-cli/releases/latest), extract it, and add the folder to your PATH.

### From Go

```bash
# Latest version
go install github.com/madstone-tech/veve-cli/cmd/veve@latest

# Specific version
go install github.com/madstone-tech/veve-cli/cmd/veve@v0.1.0
```

## Quick Start

### Basic Conversion

```bash
# Convert markdown to PDF
veve input.md -o output.pdf

# Use default output name (input.pdf)
veve input.md
```

### Theme Selection

```bash
# List all available themes
veve theme list

# Convert with a specific theme
veve input.md --theme dark -o output.pdf
veve input.md --theme academic -o output.pdf
```

### Remote Images

```bash
# Automatically download and embed remote images (enabled by default)
veve input.md -o output.pdf

# Explicit enable
veve input.md --enable-remote-images=true -o output.pdf

# Disable remote image downloading
veve input.md --enable-remote-images=false -o output.pdf

# Customize timeout and retries for unreliable networks
veve input.md \
  --remote-images-timeout=30 \
  --remote-images-max-retries=5 \
  -o output.pdf

# Use custom temp directory for downloads
veve input.md \
  --remote-images-temp-dir=/mnt/fast-storage \
  -o output.pdf
```

**Example Markdown:**

```markdown
# My Document

Here's a remote image from a CDN:

![Architecture Diagram](https://cdn.example.com/diagrams/architecture.png)

And another from an external source:

![Screenshot](https://docs.example.com/images/screenshot.png)

Local images still work too:

![Local Image](./local-image.png)
```

**Features:**
- 📥 Automatic download and embedding of HTTP/HTTPS image URLs
- ⚡ Concurrent downloads (5 images at a time by default)
- 🔄 Automatic retry with exponential backoff for transient failures
- 💾 Disk space limits (500MB per session, 100MB per image)
- 🧹 Automatic cleanup of temporary files
- ✅ Graceful degradation if some images fail to download
- 📊 Detailed error messages for troubleshooting

### Custom Themes

```bash
# Create custom theme directory
mkdir -p ~/.config/veve/themes

# Create a custom theme file
cat > ~/.config/veve/themes/mygreen.css << 'EOF'
---
name: mygreen
author: Your Name
description: My green theme
version: 1.0.0
---
body {
  font-family: Georgia, serif;
  color: darkgreen;
}
h1 {
  color: forestgreen;
  border-bottom: 2px solid darkgreen;
}
EOF

# Use your custom theme
veve input.md --theme mygreen -o output.pdf
```

### Theme Management

```bash
# List all themes (built-in + custom)
veve theme list

# Install a theme from file
veve theme add mytheme /path/to/mytheme.css

# Install a theme from URL
veve theme add mytheme https://example.com/themes/mytheme.css

# Remove a custom theme
veve theme remove mytheme
```

### Batch Processing

```bash
# Convert all markdown files in current directory
for file in *.md; do
  veve "$file" -o "${file%.md}.pdf"
done

# With specific theme
for file in *.md; do
  veve "$file" --theme dark -o "${file%.md}.pdf"
done
```

### Unix Piping

```bash
# Convert from stdin to stdout
cat input.md | veve - -o output.pdf

# Pipe to other commands
veve input.md -o - | curl -F "file=@-" https://api.example.com/upload

# Integration with other tools
pandoc-generated-md | veve - -o output.pdf
```

## Configuration

veve uses TOML for configuration. Config files are loaded from:

1. `~/.config/veve/veve.toml` (XDG Base Directory)
2. Environment variables (override config file)

### Example Configuration

```toml
# ~/.config/veve/veve.toml

# Default theme to use if not specified
default_theme = "dark"

# Default PDF engine
pdf_engine = "pdflatex"

# Quiet mode (suppress non-error output)
quiet = false

# Verbose mode (detailed output)
verbose = false
```

### Environment Variables

```bash
# Override configuration via environment
export VEVE_DEFAULT_THEME="dark"
export VEVE_PDF_ENGINE="xelatex"
export VEVE_QUIET="false"
export VEVE_VERBOSE="true"
```

## Command Reference

### Main Command

```bash
veve [input] [flags]
```

**Core Flags:**

- `-o, --output string` - Output PDF file path (default: input filename with .pdf extension)
- `-t, --theme string` - Theme to use for PDF styling (default: "default")
- `-e, --pdf-engine string` - Pandoc PDF engine to use (default: "pdflatex")
- `--quiet` - Suppress non-error output
- `--verbose` - Enable verbose output
- `-h, --help` - Show help message
- `-v, --version` - Show version

**Remote Images Flags:**

- `-r, --enable-remote-images` - Download and embed remote images (default: true)
- `--remote-images-timeout int` - Timeout in seconds per image download (default: 10)
- `--remote-images-max-retries int` - Maximum retry attempts for failed downloads (default: 3)
- `--remote-images-temp-dir string` - Custom temporary directory for downloads (default: system temp)

### Theme Commands

```bash
# List themes
veve theme list

# Add theme from file or URL
veve theme add <name> <path/url>

# Remove theme
veve theme remove <name>
veve theme remove <name> --force  # Skip confirmation
```

### Shell Completion

```bash
# Generate bash completion
veve completion bash

# Generate zsh completion
veve completion zsh

# Generate fish completion
veve completion fish

# Install completions
./scripts/install-completion.sh        # Auto-detect shell
./scripts/install-completion.sh bash   # Specific shell
./scripts/install-completion.sh all    # All supported shells
```

## Theme Development

Create custom themes with CSS styling. See [THEME_DEVELOPMENT.md](docs/THEME_DEVELOPMENT.md) for detailed guide.

### Basic Theme Structure

```css
---
name: mytheme
author: Your Name
description: A custom theme
version: 1.0.0
---

/* CSS styling */
body {
  font-family: Georgia, serif;
  color: #333;
}

h1 {
  color: #006699;
  border-bottom: 3px solid #006699;
}

code {
  background-color: #f5f5f5;
  padding: 2px 4px;
  border-radius: 3px;
}
```

### Theme Locations

- **Built-in themes**: Embedded in binary
- **User themes**: `~/.config/veve/themes/*.css`
- **Local themes**: Any path via `--theme /path/to/theme.css`

## Integration Examples

### Documentation Generation

```bash
# Generate PDF from markdown documentation
veve docs/guide.md -o guide.pdf --theme academic

# Batch convert documentation
for doc in docs/*.md; do
  veve "$doc" --theme academic -o "${doc%.md}.pdf"
done
```

### Static Site Generators

```bash
# Hugo integration in build script
for md in content/**/*.md; do
  veve "$md" -o "static/${md%.md}.pdf"
done

# Gatsby integration
npm run build && veve content/blog/*.md -o public/pdfs/
```

### CI/CD Pipelines

```bash
# GitHub Actions example
- name: Generate PDFs
  run: |
    for file in docs/*.md; do
      veve "$file" --theme default -o "build/${file%.md}.pdf"
    done

- name: Upload artifacts
  uses: actions/upload-artifact@v4
  with:
    path: build/*.pdf
```

## Remote Images Guide

### Quick Start

Remote images are **enabled by default**. Just use remote image URLs in your markdown:

```markdown
![diagram](https://cdn.example.com/diagram.png)
![photo](https://cdn.example.com/photo.jpg)
```

When you run `veve convert document.md`:
1. Images are automatically downloaded
2. Downloaded files are cached in temp directory
3. Markdown is rewritten with local paths
4. PDF is generated with embedded images
5. Temp files are cleaned up

### Performance Tips

**For Fast Networks:**
```bash
veve document.md  # Defaults are optimal for typical networks
```

**For Slow Networks:**
```bash
veve document.md \
  --remote-images-timeout=30 \
  --remote-images-max-retries=5
```

**For Large Image Batches:**
```bash
veve document.md \
  --remote-images-temp-dir=/mnt/fast-ssd  # Use faster storage
```

### Troubleshooting Remote Images

| Problem | Solution |
|---------|----------|
| Images not downloading | Feature is enabled by default. Check that URLs are correct |
| Timeout errors | Increase timeout: `--remote-images-timeout=30` |
| Rate limit errors (429) | Automatic retries handle this. Check image source |
| 404 errors | Verify image URLs in markdown are correct |
| Disk space exceeded | Reduce document size or split into multiple conversions |
| Cleanup warnings | Use custom temp dir: `--remote-images-temp-dir=./temp` |

For more details, see [specs/002-remote-images/quickstart.md](specs/002-remote-images/quickstart.md).

## Troubleshooting

### Pandoc not found

```
[ERROR] pandoc: Pandoc is required but not installed or not in PATH
```

**Solution**: Install Pandoc:

```bash
# macOS
brew install pandoc

# Linux (Debian/Ubuntu)
sudo apt-get install pandoc

# Linux (Fedora)
sudo dnf install pandoc

# Windows (Chocolatey)
choco install pandoc
```

### Theme not found

```
[ERROR] invalid theme 'mytheme': available themes are: [default dark academic]
```

**Solution**:

1. Check theme file exists: `ls ~/.config/veve/themes/mytheme.css`
2. Use correct theme name (without .css extension)
3. Use full path for local themes: `veve input.md --theme /path/to/mytheme.css`

### Encoding issues with special characters

**Solution**: Ensure your markdown file is UTF-8 encoded:

```bash
file -i input.md  # Check encoding
iconv -f utf-16 -t utf-8 input.md > input_utf8.md  # Convert if needed
```

### PDF generation slow

**Solution**:

1. Use faster PDF engine: `--pdf-engine xelatex` (faster than pdflatex)
2. Check Pandoc performance: `pandoc --version`
3. Simplify CSS in theme (reduce complexity)

## Building from Source

### Prerequisites

- Go 1.25 or later
- Pandoc 2.18 or later
- Git

### Build

```bash
# Clone repository
git clone https://github.com/madstone-tech/veve-cli.git
cd veve-cli

# Build
go build -o veve ./cmd/veve

# Install to system
sudo mv veve /usr/local/bin/

# Verify
veve --version
```

### Development

This project uses [Task](https://taskfile.dev/) to streamline common development workflows. Install Task if you haven't already:

```bash
# macOS
brew install go-task/tap/go-task

# Linux
sh -c "$(curl --location https://taskfile.dev/install.sh)" -- -d

# Or via Go
go install github.com/go-task/task/v3/cmd/task@latest
```

Then use Task for common operations:

```bash
# Setup development environment
task dev-setup

# Run tests
task test                 # All tests with race detector
task test-unit            # Unit tests only
task test-contract        # Contract tests only
task test-coverage        # Tests with coverage report

# Code quality
task fmt                  # Format with gofmt
task lint                 # Run linter
task precommit            # Run all pre-commit checks

# Build and install
task build                # Build binary
task install              # Build and install to /usr/local/bin
task uninstall            # Remove installation

# See all available tasks
task --list-all
```

For detailed information, see [TASKFILE_GUIDE.md](TASKFILE_GUIDE.md).

Traditional Go commands still work:

```bash
# Run tests directly
go test ./...

# Run with coverage
go test -cover ./...

# Generate completions
./scripts/generate-completions.sh

# Install completions
./scripts/install-completion.sh
```

## Contributing

Contributions welcome! Please see [CONTRIBUTING.md](CONTRIBUTING.md) for guidelines.

### Development Setup

```bash
# Clone and setup
git clone https://github.com/madstone-tech/veve-cli.git
cd veve-cli

# Install dependencies
go mod download

# Run tests
go test -v ./...

# Run linter
golangci-lint run ./...
```

## Release Process

See [docs/RELEASE.md](docs/RELEASE.md) for detailed release instructions.

Quick version:

```bash
# Tag a release
git tag -a v0.2.0 -m "Release v0.2.0"
git push origin v0.2.0

# GitHub Actions automatically builds and releases
```

## Documentation

- [Remote Images Guide](specs/002-remote-images/quickstart.md) - Automatic image downloading and embedding
- [Theme Development Guide](docs/THEME_DEVELOPMENT.md) - Create custom themes
- [Integration Examples](docs/INTEGRATION.md) - Use veve in your workflow
- [Release Guide](docs/RELEASE.md) - Create releases
- [Contributing Guide](CONTRIBUTING.md) - Development guidelines

## Performance

### Conversion Performance

Typical conversion times (Pandoc + PDF generation):

| Task                          | Time         |
| ----------------------------- | ------------ |
| Simple document (< 10 pages)  | < 2 seconds  |
| Complex document (< 50 pages) | 2-5 seconds  |
| Large document (> 50 pages)   | 5-10 seconds |

_Times vary based on Pandoc, PDF engine, and system performance_

### Remote Image Download Performance

When remote images are included:

| Scenario | Time |
|----------|------|
| Single image (2MB) | ~2 seconds |
| 5 images (10MB total) | ~2-3 seconds |
| 20 images (40MB total) | ~10 seconds |

_With 5 concurrent downloads (default). Times depend on network speed and image source responsiveness._

**Performance Tips:**
- Concurrent downloads (5 workers) provide ~2.5x speedup vs sequential
- Images are cached during conversion (no re-downloads for duplicates)
- Slow images don't block others (concurrent downloads continue)
- Timeouts prevent hanging on unresponsive sources (default 10s, configurable)

## Compatibility

### Operating Systems

- ✅ macOS (10.15+)
- ✅ Linux (Ubuntu 18.04+, Fedora 30+, etc.)
- ✅ Windows (10+)

### Go Versions

- ✅ Go 1.20+
- ✅ Go 1.21+

### Pandoc Versions

- ✅ Pandoc 2.18+
- ✅ Pandoc 2.19+
- ✅ Pandoc 3.x+

## License

MIT License - see [LICENSE](LICENSE) for details

## Acknowledgments

- Built with [Cobra](https://github.com/spf13/cobra) for CLI
- Configuration with [Viper](https://github.com/spf13/viper)
- Conversion powered by [Pandoc](https://pandoc.org)
- Inspiration from [Marked 2](https://marked2app.com)

## Support

- **Issues**: [GitHub Issues](https://github.com/madstone-tech/veve-cli/issues)
- **Discussions**: [GitHub Discussions](https://github.com/madstone-tech/veve-cli/discussions)
- **Documentation**: See `/docs` directory

## Changelog

### v0.2.0 (December 2025)
- ✨ **New Feature**: Automatic remote image downloading and embedding
  - Downloads HTTP/HTTPS images during conversion
  - Concurrent downloads (5 workers default)
  - Retry logic with exponential backoff
  - Disk space limits (500MB per session, 100MB per image)
  - Detailed error messages and logging
  - Graceful degradation on network failures
- 🔄 Improved error handling and reporting
- 📚 Enhanced documentation and examples

### v0.1.0 (November 2025)
- Initial release
- Markdown to PDF conversion
- Theme support (built-in and custom)
- Theme management CLI
- Configuration support

## Roadmap

- [ ] Web UI for theme preview
- [ ] Built-in theme marketplace
- [ ] Docker image with Pandoc
- [ ] Package managers (apt, rpm, brew, etc.)
- [ ] Template variables (author, date, etc.)
- [ ] PDF merge capability
- [ ] Batch image download progress indicator
- [ ] Image caching across multiple conversions

---

**Transform your documentation with veve-cli** 📄✨
