---
name: Bug Report
about: Report a bug to help us improve veve-cli
title: "[BUG] "
labels: ["bug"]
assignees: ""
---

## Description

A clear and concise description of what the bug is.

## Steps to Reproduce

Steps to reproduce the behavior:

1. Run command: `...`
2. Input file: `...`
3. Expected output: `...`
4. See error

## Expected Behavior

What you expected to happen.

## Actual Behavior

What actually happened instead.

## System Information

**Environment:**
- **OS**: macOS / Linux / Windows
- **Version**: [e.g., macOS 12.5, Ubuntu 22.04, Windows 11]
- **veve-cli version**: `veve --version`
- **Go version**: `go version` (if building from source)
- **Pandoc version**: `pandoc --version`

**Theme (if applicable):**
- Built-in theme used: default / dark / academic / custom
- Theme file path: (if using custom theme)

## Error Message

If applicable, provide the full error message:

```
[paste error output here]
```

## Additional Context

Add any other context about the problem here, such as:
- Markdown file complexity (links, images, code blocks, etc.)
- PDF engine used (pdflatex, xelatex, etc.)
- Any recent changes to your system
- Relevant configuration settings

## Minimal Example

If possible, provide a minimal example that reproduces the bug:

**input.md:**
```markdown
# Test

Your minimal example here
```

**Command:**
```bash
veve input.md --theme default -o output.pdf
```

**Expected vs Actual:**
[Describe the difference]
