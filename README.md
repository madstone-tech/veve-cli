# veve-cli - Markdown to PDF Converter with Theme Support

[![CI](https://github.com/andhi/veve-cli/actions/workflows/ci.yml/badge.svg)](https://github.com/andhi/veve-cli/actions/workflows/ci.yml)
[![Release](https://github.com/andhi/veve-cli/actions/workflows/release.yml/badge.svg)](https://github.com/andhi/veve-cli/actions/workflows/release.yml)
[![Go Version](https://img.shields.io/badge/go-1.20+-blue)](https://golang.org)
[![License](https://img.shields.io/badge/license-MIT-green)](LICENSE)

A fast, themeable markdown-to-PDF converter built with Go. Convert your markdown files to beautiful PDFs with built-in and custom themes. Perfect for documentation, reports, and technical writing.

## Features

- 📄 **Markdown to PDF** - Fast, reliable conversion via Pandoc
- 🎨 **Theme Support** - 3 built-in themes (default, dark, academic) + unlimited custom themes
- ⚙️ **Theme Management** - List, add, and remove themes via CLI
- 📝 **Custom Themes** - Create themes in `~/.config/veve/themes/` with YAML metadata
- 🔀 **Unix Composability** - Full stdin/stdout support for piping and scripting
- 🛠️ **Configuration** - TOML-based config with XDG Base Directory support
- 🚀 **Cross-Platform** - macOS, Linux, Windows support
- 🔧 **Pandoc Flexibility** - Configurable PDF engines and Pandoc options
- 📦 **Shell Completions** - bash, zsh, and fish autocompletion support

## Installation

### macOS (Homebrew - Recommended)

```bash
brew tap andhi/tap
brew install veve
```

Shell completions are automatically installed with Homebrew.

### Linux

#### Download Pre-built Binary

```bash
# AMD64
curl -sL https://github.com/andhi/veve-cli/releases/latest/download/veve_Linux_x86_64.tar.gz | tar xz
sudo mv veve /usr/local/bin/

# ARM64
curl -sL https://github.com/andhi/veve-cli/releases/latest/download/veve_Linux_arm64.tar.gz | tar xz
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

Download the latest `.zip` file from [releases](https://github.com/andhi/veve-cli/releases/latest), extract it, and add the folder to your PATH.

### From Go

```bash
# Latest version
go install github.com/andhi/veve-cli/cmd/veve@latest

# Specific version
go install github.com/andhi/veve-cli/cmd/veve@v0.1.0
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

Flags:
- `-o, --output string` - Output PDF file path (default: input filename with .pdf extension)
- `-t, --theme string` - Theme to use for PDF styling (default: "default")
- `-e, --pdf-engine string` - Pandoc PDF engine to use (default: "pdflatex")
- `--quiet` - Suppress non-error output
- `--verbose` - Enable verbose output
- `-h, --help` - Show help message
- `-v, --version` - Show version

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

- Go 1.20 or later
- Pandoc 2.18 or later
- Git

### Build

```bash
# Clone repository
git clone https://github.com/andhi/veve-cli.git
cd veve-cli

# Build
go build -o veve ./cmd/veve

# Install to system
sudo mv veve /usr/local/bin/

# Verify
veve --version
```

### Development

```bash
# Run tests
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
git clone https://github.com/andhi/veve-cli.git
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

- [Theme Development Guide](docs/THEME_DEVELOPMENT.md) - Create custom themes
- [Integration Examples](docs/INTEGRATION.md) - Use veve in your workflow
- [Release Guide](docs/RELEASE.md) - Create releases
- [Contributing Guide](CONTRIBUTING.md) - Development guidelines

## Performance

Typical performance metrics:

| Task | Time |
|------|------|
| Simple document (< 10 pages) | < 2 seconds |
| Complex document (< 50 pages) | 2-5 seconds |
| Large document (> 50 pages) | 5-10 seconds |

*Times vary based on Pandoc, PDF engine, and system performance*

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

- **Issues**: [GitHub Issues](https://github.com/andhi/veve-cli/issues)
- **Discussions**: [GitHub Discussions](https://github.com/andhi/veve-cli/discussions)
- **Documentation**: See `/docs` directory

## Roadmap

- [ ] Web UI for theme preview
- [ ] Built-in theme marketplace
- [ ] Docker image
- [ ] Package managers (apt, rpm, brew, etc.)
- [ ] VS Code extension for theme development
- [ ] Template variables (author, date, etc.)
- [ ] Batch processing UI
- [ ] PDF merge capability

---

Made with ❤️ by the veve-cli team

**Transform your documentation with veve-cli** 📄✨
