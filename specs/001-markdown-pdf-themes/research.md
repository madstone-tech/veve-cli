# Research: Markdown-to-PDF Converter with Theme Management

**Branch**: `001-markdown-pdf-themes` | **Date**: 2025-11-20

This document resolves all technical clarifications identified in the implementation plan. Findings inform the Phase 1 design decisions.

---

## 1. Pandoc Subprocess Execution in Go

### Decision
Use `os/exec.Command` to shell out to Pandoc binary. Do NOT embed Pandoc library or attempt cgo bindings.

### Rationale
- **Simplicity**: `os/exec` is Go standard library, no external dependencies
- **Reliability**: Pandoc has been battle-tested for 15+ years; leverages upstream updates automatically
- **Portability**: Single Pandoc binary on PATH; cross-platform (brew, apt, choco)
- **Maintenance**: Zero maintenance cost; no Haskell runtime, no version compatibility tracking
- **User control**: Users can customize Pandoc installation (LaTeX engine choice, filters, plugins)

### Pattern: Basic Wrapper

```go
// converter/pandoc.go
type PandocCmd struct {
    Engine   string // pdflatex, weasyprint, lualatex
    Args     []string
    CSSPath  string
}

func (p *PandocCmd) Convert(inputPath, outputPath string) error {
    cmd := exec.Command("pandoc",
        inputPath,
        "-o", outputPath,
        "--pdf-engine=" + p.Engine,
    )
    if p.CSSPath != "" {
        cmd.Args = append(cmd.Args, "--css="+p.CSSPath)
    }
    if err := cmd.Run(); err != nil {
        return fmt.Errorf("pandoc conversion failed: %w", err)
    }
    return nil
}
```

### Error Handling: Pandoc Not Found

```go
// At startup or first use
pandocPath, err := exec.LookPath("pandoc")
if err != nil {
    return fmt.Errorf("pandoc not found in PATH\n"+
        "Install with: brew install pandoc (macOS)\n"+
        "           or apt-get install pandoc (Linux)\n"+
        "           or choco install pandoc (Windows)\n"+
        "           or https://pandoc.org/installing.html")
}
```

### Alternatives Considered

| Alternative | Pros | Cons | Decision |
|---|---|---|---|
| Embed Pandoc library (Haskell) | No external tool | Large binary (50MB+), Haskell runtime, maintenance burden | ❌ Rejected |
| cgo binding to C library | Better perf | Hard to maintain, breaks portability, version sync nightmare | ❌ Rejected |
| Shell out to pandoc (chosen) | Simple, reliable, user-controlled | Requires Pandoc install | ✅ Selected |

---

## 2. XDG Base Directory Implementation

### Decision
Follow XDG Base Directory Specification (`~/.config/veve/` on Unix systems, `%APPDATA%\veve\` on Windows).

### Directory Structure

```
~/.config/veve/                      (Unix: Linux, macOS)
├── veve.json                        (Config file: PDF engine, defaults)
├── themes/                          (User-installed themes)
│   ├── dark.css
│   ├── academic.css
│   └── custom-brand.css
└── themes.json                      (Theme registry: metadata)

%APPDATA%\veve\                      (Windows)
├── veve.json
├── themes\
└── themes.json
```

### Rationale

- **Unix Philosophy**: Follows Linux/macOS conventions; users expect `~/.config/`
- **Portability**: Works on Linux, macOS, Windows with minimal conditional logic
- **XDG Compliance**: Industry standard for application configuration
- **User Control**: Themes are discoverable, editable, and shareable
- **No system pollution**: Config isolated from system directories

### Implementation Pattern

```go
// internal/config/paths.go
import "os/user"

func ConfigDir() (string, error) {
    homeDir, err := os.UserHomeDir()
    if err != nil {
        return "", err
    }
    
    configDir := filepath.Join(homeDir, ".config", "veve")
    
    // On Windows, use %APPDATA%\veve
    if runtime.GOOS == "windows" {
        appData := os.Getenv("APPDATA")
        if appData != "" {
            configDir = filepath.Join(appData, "veve")
        }
    }
    
    return configDir, nil
}

func EnsureConfigDir() error {
    dir, err := ConfigDir()
    if err != nil {
        return err
    }
    return os.MkdirAll(filepath.Join(dir, "themes"), 0755)
}
```

### Alternatives Considered

| Alternative | Pros | Cons | Decision |
|---|---|---|---|
| Fixed relative paths (e.g., `./config/`) | Simple | Breaks portability, clutters repo | ❌ Rejected |
| Env var only (`VEVE_CONFIG_DIR`) | Flexible | Requires setup, discoverability poor | ❌ Rejected |
| XDG Base Directory (chosen) | Standard, portable | Slightly more code | ✅ Selected |

---

## 3. Built-in Theme Embedding Strategy

### Decision
Embed theme CSS files into the veve binary using Go's `embed` package (Go 1.16+).

### Rationale

- **Zero external dependencies**: Themes shipped with binary; no separate downloads
- **Fast startup**: No filesystem access needed for built-in themes
- **Reliability**: Can't lose theme files; always available
- **Cross-platform**: Single binary works everywhere
- **User override**: Still load from `~/.config/veve/themes/` for custom themes

### Implementation Pattern

```go
// themes/embed.go
//go:embed *.css
var builtInThemes embed.FS

func LoadBuiltInTheme(name string) ([]byte, error) {
    return fs.ReadFile(builtInThemes, name+".css")
}
```

### Built-in Themes to Ship (v1.0)

1. **default.css** - Clean, professional, neutral colors (black text, white background)
2. **dark.css** - Dark mode (light text, dark background; suitable for terminal users)
3. **academic.css** - Paper-like (serif fonts, justified text, professional margins)

Future themes (v1.1+):
- modern.css (colorful, sans-serif, contemporary)
- minimal.css (extreme simplicity, tiny margins)
- print.css (optimized for printing, minimal colors)

### Directory Structure

```
themes/
├── default.css          # Default theme
├── dark.css
├── academic.css
└── embed.go             # Embed directive
```

### Alternatives Considered

| Alternative | Pros | Cons | Decision |
|---|---|---|---|
| Embed themes in binary (chosen) | Single file, portable | Larger binary (~+50KB) | ✅ Selected |
| Ship themes separately | Smaller binary | Breaks portability, extra files | ❌ Rejected |
| Generate at runtime | Flexible | Unnecessary complexity | ❌ Rejected |

---

## 4. Theme CSS → PDF Application Mechanism

### Decision
Use Pandoc's `--css` flag for HTML-based engines, or inject into LaTeX template for LaTeX engines.

### Implementation Pattern: CSS for weasyprint

```go
// converter/converter.go
func applyTheme(cmd *exec.Cmd, themePath string) {
    // For HTML→PDF engines (weasyprint)
    cmd.Args = append(cmd.Args, "--css="+themePath)
}
```

### Implementation Pattern: LaTeX Template

```go
// converter/converter.go
func applyLaTeXTheme(cmd *exec.Cmd, templatePath string) {
    // For LaTeX-based engines (pdflatex)
    cmd.Args = append(cmd.Args, "--template="+templatePath)
}
```

### Rationale

- **Native support**: Pandoc handles CSS/template injection natively
- **No custom rendering**: Avoid reimplementing PDF features
- **User flexibility**: Themes leverage Pandoc's full feature set
- **Simplicity**: Just pass file paths; Pandoc does the work

### Theme Development Guide

**CSS Theme Format** (for weasyprint/HTML engines):
```css
/* themes/dark.css */
body {
    background-color: #1e1e1e;
    color: #e0e0e0;
    font-family: "Source Sans Pro", sans-serif;
    font-size: 11pt;
    line-height: 1.6;
}
h1, h2, h3 { color: #4a9eff; }
code { background: #333; padding: 2px 4px; }
```

**LaTeX Template** (for pdflatex/xelatex):
```latex
% themes/academic.latex
\documentclass{article}
\usepackage[margin=1in]{geometry}
\usepackage{xcolor}
\definecolor{accentcolor}{HTML}{2E5090}
% ... Pandoc template directives ...
```

### Alternatives Considered

| Alternative | Pros | Cons | Decision |
|---|---|---|---|
| CSS injection (chosen) | Native, simple | Limited to HTML engines | ✅ Selected |
| LaTeX template (chosen) | Powerful, professional | Steeper learning curve | ✅ Selected |
| Custom CSS→PDF library | Full control | Massive effort, reinvent wheel | ❌ Rejected |

---

## 5. Theme Registry Format

### Decision
Use simple JSON manifest at `~/.config/veve/themes.json` (optional; directory scanning is fallback).

### Schema

```json
{
  "version": "1.0",
  "themes": [
    {
      "name": "dark",
      "path": "~/.config/veve/themes/dark.css",
      "type": "css",
      "author": "veve maintainers",
      "description": "Dark mode theme with light text on dark background",
      "source_url": "https://github.com/veve-cli/themes/dark"
    },
    {
      "name": "custom-brand",
      "path": "~/.config/veve/themes/custom-brand.css",
      "type": "css",
      "author": "Company Brand Team",
      "description": "Corporate branding with company colors"
    }
  ]
}
```

### Implementation Pattern

```go
// internal/theme/registry.go
type ThemeRegistry struct {
    Version string     `json:"version"`
    Themes  []Theme    `json:"themes"`
}

type Theme struct {
    Name        string `json:"name"`
    Path        string `json:"path"`
    Type        string `json:"type"` // "css" or "latex"
    Author      string `json:"author"`
    Description string `json:"description"`
    SourceURL   string `json:"source_url"`
}

func Load() (*ThemeRegistry, error) {
    registryPath := filepath.Join(configDir, "themes.json")
    data, err := os.ReadFile(registryPath)
    if err != nil && !os.IsNotExist(err) {
        return nil, err
    }
    
    var reg ThemeRegistry
    if len(data) > 0 {
        if err := json.Unmarshal(data, &reg); err != nil {
            return nil, fmt.Errorf("invalid themes.json: %w", err)
        }
    }
    return &reg, nil
}

func (r *ThemeRegistry) Add(t Theme) error {
    // Check for duplicates
    for _, existing := range r.Themes {
        if existing.Name == t.Name {
            return fmt.Errorf("theme %q already exists", t.Name)
        }
    }
    r.Themes = append(r.Themes, t)
    return r.Save()
}

func (r *ThemeRegistry) Save() error {
    registryPath := filepath.Join(configDir, "themes.json")
    data, _ := json.MarshalIndent(r, "", "  ")
    return os.WriteFile(registryPath, data, 0644)
}
```

### Fallback: Directory Scanning

If `themes.json` doesn't exist, scan `~/.config/veve/themes/` for `.css` and `.latex` files:

```go
func DiscoverThemes() ([]Theme, error) {
    themesDir := filepath.Join(configDir, "themes")
    entries, err := os.ReadDir(themesDir)
    // Iterate and find *.css, *.latex files
    // Return as auto-discovered themes
}
```

### Rationale

- **Human-readable**: JSON is editable, debuggable
- **Version-able**: Can evolve schema without breaking
- **Optional**: Works with or without registry file
- **Lightweight**: Single file; no database needed

### Alternatives Considered

| Alternative | Pros | Cons | Decision |
|---|---|---|---|
| JSON registry (chosen) | Readable, structured | Manual editing required | ✅ Selected |
| YAML registry | Simpler syntax | Extra dependency (go-yaml) | ❌ Rejected |
| Directory scanning only | No config file | No metadata (author, URL) | ⚠ Fallback |
| Database (SQLite) | Powerful | Overkill for CLI tool | ❌ Rejected |

---

## 6. PDF Engine Selection

### Decision
Default to **pdflatex** for professional output. Users can override with `--pdf-engine` flag or set default in config.

### Engine Comparison

| Engine | Input Format | Quality | Speed | LaTeX Required | Notes |
|--------|---|---|---|---|---|
| **pdflatex** | LaTeX | ⭐⭐⭐⭐⭐ Professional | Medium | Yes | Best for professional docs; default |
| **xelatex** | LaTeX | ⭐⭐⭐⭐⭐ Professional | Slow | Yes | Unicode support; advanced fonts |
| **lualatex** | LaTeX | ⭐⭐⭐⭐⭐ Professional | Slow | Yes | Lua scripting; modern |
| **weasyprint** | HTML/CSS | ⭐⭐⭐⭐ Good | Fast | No | Lightweight; CSS-based themes |
| **prince** | HTML/CSS | ⭐⭐⭐⭐⭐ Professional | Fast | No | Commercial; tagged PDFs |
| **typst** | Custom | ⭐⭐⭐⭐ Good | Very Fast | No | Modern, experimental |

### Strategy

1. **Default (MVP v1.0)**: `pdflatex` with LaTeX-based themes
   - Rationale: Professional output, proven, widely available
   - Assumption: Users willing to install LaTeX (~500MB)

2. **Alternative (v1.1)**: `--pdf-engine weasyprint` for lighter setup
   - Rationale: No LaTeX required; CSS-based themes more accessible
   - Path: `veve input.md --pdf-engine weasyprint -o output.pdf`

3. **Future (v2.0)**: Automatic engine detection based on system
   - Rationale: Ease of use for users without LaTeX

### Implementation Pattern

```go
// cmd/veve/convert.go
var pdfEngine string // Flag: --pdf-engine

func init() {
    convertCmd.Flags().StringVarP(&pdfEngine, "pdf-engine", "e", "pdflatex",
        "PDF engine: pdflatex, xelatex, weasyprint, typst")
}

func convertMarkdown(inputPath, outputPath string) error {
    cmd := exec.Command("pandoc",
        inputPath,
        "-o", outputPath,
        "--pdf-engine=" + pdfEngine,
    )
    return cmd.Run()
}
```

### Rationale for pdflatex as Default

- Most widely installed (included in most TeX distributions)
- Stable and mature (30+ years)
- Professional typography (proper kerning, ligatures, hyphenation)
- Best support for complex documents (citations, cross-refs)
- Matches Marked 2's output quality

### Alternatives Considered

| Alternative | Pros | Cons | Decision |
|---|---|---|---|
| pdflatex default (chosen) | Professional, stable | Requires LaTeX | ✅ Selected (v1.0) |
| weasyprint default | Lightweight, no LaTeX | Less professional output | ⚠ Future default (v1.1) |
| Auto-detect | Best UX | Complex logic | ⚠ Future (v2.0) |

---

## 7. stdin/stdout Handling

### Decision
Use `-` as special filename to indicate stdin (input) or stdout (output).

### Unix Convention Pattern

```bash
# Read from stdin, write to file
cat input.md | veve - -o output.pdf

# Read from file, write to stdout
veve input.md -o -

# Read from stdin, write to stdout
cat input.md | veve - -o - > output.pdf

# Pipe through other tools
veve input.md -o - | gs -sDEVICE=jpeg -o output.jpg -

# Batch processing
for file in *.md; do veve "$file" -o "${file%.md}.pdf"; done
```

### Implementation Pattern

```go
// internal/converter/converter.go
func Convert(inputPath, outputPath string, opts Options) error {
    input := inputPath
    if input == "-" {
        // Read from stdin
        input = "/dev/stdin" // Unix
        // Windows: Use special handling or skip
    }
    
    output := outputPath
    if output == "-" {
        // Write to stdout; use temp file then pipe
        // Or: Use Pandoc's native stdout support
    }
    
    cmd := exec.Command("pandoc", input, "-o", output, ...)
    return cmd.Run()
}
```

### Edge Cases

- **Windows stdin/stdout**: May need special handling; test cross-platform
- **Binary data**: PDF is binary; ensure no text encoding issues
- **Piping**: Verify Pandoc writes cleanly to stdout without extra output

### Rationale

- **Unix tradition**: Established pattern (grep, cat, sed, awk)
- **Composability**: Enables piping and shell integration
- **Flexibility**: Same tool works for interactive + scripted usage

### Alternatives Considered

| Alternative | Pros | Cons | Decision |
|---|---|---|---|
| `-` special (chosen) | Standard convention | Less explicit | ✅ Selected |
| `--stdin` / `--stdout` flags | Explicit | Verbose, non-standard | ❌ Rejected |
| Both `-` and flags | Most flexible | Confusing | ❌ Rejected |

---

## 8. Error Handling & Messaging

### Decision
Consistent error format: Exit codes (0, 1, 2) + actionable messages to stderr.

### Exit Code Scheme

```
0: Success
1: Conversion/runtime error (file I/O, Pandoc failure, theme not found)
2: Usage error (invalid flags, missing required args)
```

### Error Message Format

```
[ERROR] <command>: <action> failed: <reason>
(try: <suggestion>)
```

### Examples

```
[ERROR] convert: input file not found
(try: veve README.md -o README.pdf)

[ERROR] convert: pandoc not found in PATH
(try: brew install pandoc)

[ERROR] convert: theme "invalid" not found
(try: veve --list-themes to see available themes)

[ERROR] theme add: download failed (network error)
(try: check your internet connection and try again)

[ERROR] theme remove: theme "dark" is built-in; cannot remove
(try: use custom theme instead)
```

### Implementation Pattern

```go
// internal/errors.go
type VeveError struct {
    Command    string // e.g., "convert", "theme add"
    Action     string // e.g., "conversion", "file read"
    Reason     string // e.g., "file not found"
    Suggestion string // e.g., "check file path"
    ExitCode   int    // 1 or 2
}

func (e *VeveError) Error() string {
    return fmt.Sprintf("[ERROR] %s: %s failed: %s\n(try: %s)",
        e.Command, e.Action, e.Reason, e.Suggestion)
}

// Usage in commands
if err := pandoc.Convert(...); err != nil {
    log.SetFlags(0)
    log.Fatalln(&VeveError{
        Command:    "convert",
        Action:     "conversion",
        Reason:     err.Error(),
        Suggestion: "ensure pandoc is installed (brew install pandoc)",
        ExitCode:   1,
    })
}
```

### Rationale

- **Clarity**: Users know exactly what went wrong
- **Actionability**: Each error includes next steps
- **Scriptability**: Exit codes allow conditional logic
- **Consistency**: All errors follow same format

### Alternatives Considered

| Alternative | Pros | Cons | Decision |
|---|---|---|---|
| Structured errors (chosen) | Clear, actionable | Slightly verbose | ✅ Selected |
| Generic errors | Simple | Not helpful to users | ❌ Rejected |
| JSON error output | Machine-readable | Overkill for CLI | ❌ Rejected |

---

## Summary Table

| Clarification | Decision | Confidence | Risk |
|---|---|---|---|
| Pandoc subprocess | Shell out with `os/exec` | High | Low (upstream-maintained) |
| XDG Base Directory | `~/.config/veve/` with Windows override | High | Low (standard) |
| Built-in themes | Embed in binary with Go `embed` | High | Low (Go 1.16+ stable) |
| Theme CSS application | Pandoc `--css` + `--template` flags | High | Low (native support) |
| Theme registry | JSON manifest + directory fallback | Medium | Low (simple format) |
| PDF engine default | `pdflatex` with override flag | High | Medium (requires LaTeX install) |
| stdin/stdout | `-` special filename | High | Low (Unix convention) |
| Error handling | Exit codes + structured messages | High | Low (clear semantics) |

---

## Next Steps

1. **Phase 1 Design**: Use these decisions to fill `data-model.md`, `contracts/`, `quickstart.md`
2. **Implementation**: Reference this document during coding for patterns + edge cases
3. **Testing**: Validate research assumptions (especially cross-platform stdin/stdout, theme embedding)
4. **Documentation**: Include decision rationale in README + developer guide
