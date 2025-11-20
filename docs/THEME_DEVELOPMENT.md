# Theme Development Guide

This guide explains how to create custom themes for veve-cli.

## Theme Format

A veve theme is a CSS file with optional YAML metadata at the top. The metadata allows you to document your theme with name, author, description, and version information.

### Basic Theme Structure

```css
---
name: my-custom-theme
author: Your Name
description: A custom theme for veve-cli
version: 1.0.0
---
/* Your CSS rules here */
body {
  font-family: Arial, sans-serif;
  color: #333;
  background: white;
}

h1 {
  font-size: 28pt;
  font-weight: bold;
  color: #000;
}
```

### Metadata Fields

The YAML front matter (between `---` markers) supports the following fields:

| Field | Required | Type | Description |
|-------|----------|------|-------------|
| `name` | No | string | Theme name (defaults to filename if not provided) |
| `author` | No | string | Theme author name (defaults to "Unknown") |
| `description` | No | string | Theme description (defaults to "Custom theme") |
| `version` | No | string | Theme version (defaults to "1.0.0") |

**Note:** All metadata fields are optional. If omitted, sensible defaults will be applied.

## Creating a Theme

### Step 1: Choose a Location

You have two options for storing custom themes:

1. **User themes directory**: `~/.config/veve/themes/` (Linux/macOS)
   - Themes here are discovered automatically on startup
   - Recommended for personal theme collections

2. **Local file**: Specify the full path with `--theme /path/to/theme.css`
   - Useful for project-specific themes
   - No installation needed

### Step 2: Write Your CSS

Create a `.css` file with your styling rules. The CSS will be applied to the PDF output via Pandoc's styling mechanism.

Example: `mytheme.css`

```css
---
name: mytheme
author: Jane Developer
description: A minimalist theme
version: 1.0.0
---
body {
  font-family: 'Georgia', serif;
  font-size: 12pt;
  line-height: 1.5;
  color: #222;
  background: #fafafa;
  margin: 2em;
}

h1 {
  font-size: 32pt;
  font-weight: bold;
  color: #1a5490;
  margin-top: 1em;
  margin-bottom: 0.5em;
}

h2 {
  font-size: 24pt;
  color: #2a6fb0;
  margin-top: 0.8em;
  margin-bottom: 0.4em;
  border-bottom: 2px solid #e0e0e0;
}

code {
  font-family: 'Courier New', monospace;
  background: #f0f0f0;
  padding: 2px 6px;
  border-radius: 3px;
}

blockquote {
  margin-left: 2em;
  padding-left: 1em;
  border-left: 3px solid #ccc;
  color: #666;
  font-style: italic;
}

a {
  color: #1a5490;
  text-decoration: underline;
}

table {
  border-collapse: collapse;
  margin: 1em 0;
}

table th {
  background: #f0f0f0;
  padding: 0.5em;
  border: 1px solid #ddd;
  text-align: left;
}

table td {
  padding: 0.5em;
  border: 1px solid #ddd;
}
```

### Step 3: Install the Theme

#### Option A: User Themes Directory

```bash
# Create the themes directory if it doesn't exist
mkdir -p ~/.config/veve/themes/

# Copy your theme
cp mytheme.css ~/.config/veve/themes/

# List available themes
veve theme list
```

#### Option B: Local File Path

```bash
# Use the theme from anywhere
veve input.md --theme /path/to/mytheme.css -o output.pdf
```

## Using Your Theme

Once installed, use your theme with the `--theme` flag:

```bash
# Use an installed theme
veve input.md --theme mytheme -o output.pdf

# Use a local file
veve input.md --theme ./mytheme.css -o output.pdf

# Use a built-in theme
veve input.md --theme dark -o output.pdf
```

## Theme Validation

Your theme CSS will be validated for:

- **Non-empty content**: The CSS must not be blank
- **Balanced braces**: All `{` must have corresponding `}`
- **Valid syntax**: Basic CSS syntax checking

If validation fails, veve-cli will report the error and exit.

## CSS Best Practices

### 1. Use Semantic HTML Elements

The Pandoc conversion produces semantic HTML that maps to CSS selectors:

- `h1`, `h2`, `h3`, etc. for headings
- `p` for paragraphs
- `code`, `pre` for code blocks
- `blockquote` for block quotes
- `ul`, `ol`, `li` for lists
- `table`, `tr`, `td`, `th` for tables
- `img` for images
- `a` for links

### 2. Define Font Families

```css
body {
  font-family: 'Segoe UI', Arial, sans-serif;
}

code {
  font-family: 'Courier New', 'Monaco', monospace;
}
```

### 3. Use Relative Units

For better scaling and responsive design:

```css
/* Good */
h1 { font-size: 28pt; }
body { margin: 2em; }
code { padding: 2px 6px; }

/* Less ideal */
h1 { font-size: 50px; }
body { margin: 100px; }
```

### 4. Avoid Complex Selectors

Stick to simple, direct selectors for maximum compatibility:

```css
/* Good */
h1 { color: #333; }
a { text-decoration: underline; }

/* May not work with PDF conversion */
h1 + p { margin-top: 0; }
li:first-child { font-weight: bold; }
```

### 5. Test Your Theme

Create a test markdown file and convert it with your theme:

```bash
veve test.md --theme mytheme -o test.pdf
```

Then open the PDF and verify:
- Text is readable
- Colors are correct
- Spacing looks good
- Code blocks are formatted properly

## Built-in Themes

veve-cli includes three built-in themes you can use as reference:

### default
A clean, simple theme suitable for most documents.

### dark
A dark theme with light text, useful for reduced-eye-strain reading.

### academic
A formal theme suitable for academic papers and formal documents.

View these themes in the `themes/` directory to see examples of well-structured themes.

## Troubleshooting

### Theme Not Found

If veve-cli can't find your theme:

```bash
# Check installed themes
veve theme list

# Verify file exists
ls ~/.config/veve/themes/mytheme.css

# Try with full path
veve input.md --theme ~/.config/veve/themes/mytheme.css -o output.pdf
```

### CSS Not Applied

If your CSS isn't being applied:

1. Check for CSS syntax errors
2. Verify the CSS file is readable
3. Try with a built-in theme to confirm the pipeline works
4. Check that selectors match HTML elements in the output

### PDF Generation Fails

If PDF generation fails with your theme:

1. Validate CSS has balanced braces: `{ }` count must match
2. Test with the default theme
3. Simplify your CSS to identify the problematic rule
4. Check Pandoc version compatibility

## Advanced Features

### Font Embedding (Future)

Future versions of veve-cli will support custom fonts. For now, stick to system fonts or web-safe fonts.

### Theme Inheritance (Future)

Future versions may support theme inheritance, allowing you to extend built-in themes.

## Contributing Themes

If you create a great theme, consider contributing it to the veve-cli project! Open a pull request with your theme and we may include it as a built-in option.

## See Also

- [Pandoc CSS Documentation](https://pandoc.org/MANUAL.html#css)
- [CSS Selectors Reference](https://developer.mozilla.org/en-US/docs/Web/CSS/CSS_Selectors)
- Built-in themes: `veve theme list`
