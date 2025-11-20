# Feature Specification: Markdown-to-PDF Converter with Theme Management

**Feature Branch**: `001-markdown-pdf-themes`  
**Created**: 2025-11-20  
**Status**: Draft  
**Input**: User description: "veve-cli is a markdown to pdf converter with theme management. it leverages pandoc. it is inspired by marked2 app for the mac. but should be cross platform and more with the unix philosophy. the goal is to be able to feed a markdown file and have it generate a pdf output with the choice of themes to apply."

## User Scenarios & Testing

### User Story 1 - Basic Markdown to PDF Conversion (Priority: P1)

A developer wants to convert a Markdown file to a PDF with default styling, without needing to install or configure additional tools beyond Pandoc. The conversion should be a simple one-command operation that works on any platform (macOS, Linux, Windows).

**Why this priority**: This is the MVP—the core functionality that defines veve-cli. Without this, there is no product. It must work reliably across all platforms to justify the cross-platform promise.

**Independent Test**: Can be fully tested by running `veve input.md -o output.pdf` on a simple Markdown file and verifying the output PDF is created with correct content and professional formatting.

**Acceptance Scenarios**:

1. **Given** a valid Markdown file with basic content (headings, paragraphs, lists), **When** user runs `veve input.md -o output.pdf`, **Then** a PDF file is created at the output path with properly formatted content
2. **Given** a Markdown file with images and links, **When** conversion completes, **Then** images are embedded and links are preserved in the PDF
3. **Given** invalid or missing input file, **When** user runs the command, **Then** a clear error message is displayed directing user to valid input
4. **Given** output path in a non-existent directory, **When** conversion is requested, **Then** directory is created and PDF is written successfully
5. **Given** no `--output` flag specified, **When** user runs `veve input.md`, **Then** PDF is written to `input.pdf` in the same directory

---

### User Story 2 - Theme Selection and Application (Priority: P1)

A developer wants to customize the PDF appearance by selecting from multiple built-in themes (default, dark, academic, etc.) without writing CSS or LaTeX templates. Theme selection should be a simple flag (`--theme <name>`), and switching themes should not require reinstallation or configuration.

**Why this priority**: Theme support is the primary differentiator from raw Pandoc. Marked 2 success was partly due to beautiful, easy-to-apply themes. This is core to veve's value proposition.

**Independent Test**: Can be fully tested by converting the same Markdown file with `--theme default` and `--theme dark`, verifying both PDFs are created with visually distinct styling applied correctly.

**Acceptance Scenarios**:

1. **Given** a Markdown file and user specifies `--theme dark`, **When** conversion completes, **Then** PDF is styled with dark theme colors and fonts
2. **Given** multiple built-in themes available, **When** user runs `veve --list-themes`, **Then** all available themes are listed with descriptions
3. **Given** an invalid theme name, **When** user specifies `--theme invalid-theme`, **Then** error message lists available themes and suggests corrections
4. **Given** no `--theme` flag specified, **When** conversion runs, **Then** default theme is automatically applied
5. **Given** custom theme CSS file path, **When** user runs `veve input.md --theme ./custom.css`, **Then** custom theme is applied to output PDF

---

### User Story 3 - Custom Theme Creation and Distribution (Priority: P2)

A content team wants to create custom themes matching their brand identity (company colors, fonts, logo placement). Themes should be shareable as standalone files and installable without recompiling veve-cli.

**Why this priority**: P2 because it enables extensibility. Core users (developers) don't need custom themes initially, but teams doing batch document generation will value this. It also builds community around the tool.

**Independent Test**: Can be fully tested by creating a custom theme CSS file, installing it, converting a Markdown file with `--theme custom-brand`, and verifying the output matches the custom styling.

**Acceptance Scenarios**:

1. **Given** a valid CSS file with Pandoc-compatible styling, **When** user places it in `~/.config/veve/themes/` and runs `veve input.md --theme custom-brand`, **Then** PDF is styled with custom theme
2. **Given** a Markdown file with code blocks and tables, **When** custom theme is applied, **Then** code highlighting and table styling respect custom CSS
3. **Given** user creates theme with custom fonts, **When** PDF is generated, **Then** fonts are properly embedded in PDF file
4. **Given** no theme directory exists, **When** user attempts to install custom theme, **Then** directory structure is created automatically

---

### User Story 4 - Theme Registry and Management (Priority: P2)

A user wants to list installed themes, preview theme samples, and manage theme installation without manual file manipulation. A built-in theme registry should allow discovering and installing community-contributed themes.

**Why this priority**: P2 because it enhances discoverability and user experience, but isn't strictly necessary for MVP. Users can manually manage themes initially.

**Independent Test**: Can be fully tested by running `veve theme list`, `veve theme add <theme-url>`, and `veve theme remove <theme-name>`, verifying each command produces expected output and filesystem changes.

**Acceptance Scenarios**:

1. **Given** user runs `veve theme list`, **When** command completes, **Then** all installed themes are displayed with descriptions and preview URLs
2. **Given** user runs `veve theme add https://example.com/themes/modern.zip`, **When** download completes, **Then** theme is extracted to `~/.config/veve/themes/` and available for use. Supports both `.zip` files (multiple themes) and single `.css` files (direct download).
3. **Given** user runs `veve theme remove dark`, **When** confirmation is provided, **Then** theme files are deleted and theme is no longer available
4. **Given** invalid theme URL, **When** download is attempted, **Then** clear error message explains what went wrong and suggests checking the URL
5. **Given** user downloads a theme with non-text files, **When** extraction completes, **Then** veve warns about suspicious files (executables, binaries) and allows user to proceed or cancel

---

### User Story 5 - CLI Integration and Unix Composability (Priority: P2)

A developer wants to integrate veve-cli into automation scripts, static site generators, and CI/CD pipelines. The tool should follow Unix conventions: accept stdin/stdout, produce predictable exit codes, and work seamlessly with pipes and shell redirections.

**Why this priority**: P2 because it's important for power users and automation, but not strictly required for interactive use cases. Core functionality can work without perfect Unix integration, though it's philosophically important to veve.

**Independent Test**: Can be fully tested by piping content through veve (`cat input.md | veve - -o output.pdf`), checking exit codes for various failure modes, and integrating veve into a bash script for batch conversion.

**Acceptance Scenarios**:

1. **Given** user pipes Markdown via stdin with `cat input.md | veve - -o output.pdf`, **When** command completes, **Then** PDF is written to output path
2. **Given** successful conversion, **When** veve completes, **Then** exit code is 0
3. **Given** invalid input or missing Pandoc, **When** veve encounters error, **Then** exit code is non-zero and error details are written to stderr
4. **Given** user specifies `--output -`, **When** conversion completes, **Then** PDF binary is written to stdout (useful for piping to other tools)

---

### Edge Cases

- What happens when Pandoc is not installed or not in PATH? → User receives clear installation instructions
- How does veve handle very large Markdown files (>100MB)? → Pandoc streams processing; memory usage should remain reasonable
- What happens if theme CSS references fonts not installed on system? → Pandoc falls back to system fonts; warning is displayed
- How does veve handle Markdown with unsupported syntax (e.g., custom HTML)? → Pandoc handles gracefully; result depends on Pandoc configuration
- What happens when writing to a read-only filesystem? → Clear error message indicating permission issue
- How does veve handle circular includes or recursive Markdown files? → Pandoc + user responsibility; document limitation

---

## Requirements

### Functional Requirements

- **FR-001**: System MUST accept a Markdown file path as input via command-line argument
- **FR-002**: System MUST generate a PDF file at a specified output path
- **FR-003**: System MUST support built-in theme selection via `--theme <name>` flag
- **FR-004**: System MUST provide clear error messages when Pandoc is not available in PATH
- **FR-005**: System MUST validate input Markdown file exists before attempting conversion
- **FR-006**: System MUST support reading Markdown from stdin when input is specified as `-`
- **FR-007**: System MUST write PDF to stdout when output is specified as `-`
- **FR-008**: System MUST use professional-grade PDF rendering via Pandoc with appropriate engine (pdflatex, weasyprint, or typst)
- **FR-009**: System MUST preserve Markdown features: headings, lists, code blocks, links, images, tables, blockquotes
- **FR-010**: System MUST embed images in PDF output
- **FR-011**: System MUST apply custom CSS themes to PDF output when specified via `--theme` or `--css` flags
- **FR-012**: System MUST store user-installed themes in `~/.config/veve/themes/` following XDG standards
- **FR-013**: System MUST support Pandoc template customization via `--template` flag for advanced users
- **FR-014**: System MUST provide a theme management subcommand: `veve theme {list|add|remove|preview}`
- **FR-015**: System MUST document all features clearly in README with installation, usage, and theme development guides
- **FR-016**: System MUST support quiet mode by default (only output errors to stderr); MUST accept `--verbose` or `-v` flag to enable detailed logging (theme resolution, Pandoc invocation, file operations)
- **FR-017**: System MUST validate downloaded theme files (from `veve theme add <url>`): accept only `.css`, `.latex`, `.md`, and metadata text files. MUST warn if non-text files (executables, binaries) are detected and allow user to cancel extraction.
- **FR-018**: System MUST support two theme distribution formats: `.zip` files (extract to theme subdirectory) and single `.css` files (copy directly). URL ending in `.css` triggers direct copy; other extensions trigger zip extraction.
- **FR-019**: System MUST load configuration from `~/.config/veve/veve.toml` (TOML format) for default settings. MUST use `~/.config/veve/themes.json` (JSON) for theme registry metadata. Both files are optional; hardcoded defaults apply if missing.

### Key Entities

- **Markdown Input**: Source file (or stdin) containing Markdown-formatted content
- **Theme**: CSS or Pandoc template file that defines visual styling for PDF output (colors, fonts, layout, spacing). User-installed themes with same name as built-in themes override the built-in version.
- **Theme Registry**: Directory structure (`~/.config/veve/themes/`) that stores installed themes. Themes are stored as individual CSS/LaTeX files or organized in subdirectories. User themes take precedence over built-in themes of the same name.
- **Theme Metadata**: Optional YAML/JSON file describing theme (name, author, description, preview URL). Stored as front matter in CSS files or in a separate `metadata.json` per theme.
- **PDF Output**: Generated PDF file with converted content and applied styling
- **Configuration**: TOML file (`~/.config/veve/veve.toml`) storing veve settings (default theme, PDF engine, logging level). Theme registry metadata stored separately in `~/.config/veve/themes.json`.

---

## Success Criteria

### Measurable Outcomes

- **SC-001**: Users can convert a simple Markdown file to PDF in under 5 seconds (excluding first Pandoc invocation)
- **SC-002**: Users can switch between 5+ built-in themes without any configuration or recompilation
- **SC-003**: 95% of test cases pass: Markdown features correctly render in PDF (headings, lists, code, tables, images)
- **SC-004**: Custom themes can be created and installed without requiring veve-cli modifications
- **SC-005**: veve-cli works identically on macOS, Linux, and Windows (verified with integration tests on all three platforms)
- **SC-006**: Documentation is clear enough that new users can complete a conversion task in under 5 minutes (measured via user testing)
- **SC-007**: Markdown input files up to 10MB convert successfully without memory errors
- **SC-008**: Tool provides helpful error messages in 95% of failure scenarios (e.g., missing Pandoc, invalid file paths, malformed theme CSS)

---

## Assumptions

- **Pandoc Installation**: Users will install Pandoc via package manager (brew, apt, choco, etc.) before using veve-cli. This is clearly documented in README.
- **PDF Engine**: Default PDF engine will be `pdflatex` for professional output; users can override via `--pdf-engine` flag.
- **Theme Format**: Built-in themes use CSS for HTML-based output or Pandoc's LaTeX template format. Users can contribute themes in either format.
- **XDG Compliance**: Theme registry follows XDG Base Directory spec (`~/.config/veve/` on Linux/macOS, `%APPDATA%\veve\` on Windows).
- **No GUI**: veve-cli is CLI-only; no GUI theme editor is planned (users edit theme files directly).
- **Marked 2 Feature Parity**: veve-cli is inspired by Marked 2's ease-of-use and theme system, but is not a full feature-for-feature replacement.
- **Configuration Format**: Main configuration stored in TOML format (`~/.config/veve/veve.toml`) for human readability. Theme registry metadata stored in JSON (`~/.config/veve/themes.json`) for structured data.
- **Theme Naming**: User-installed themes with same name as built-in themes override the built-in version (Unix config precedence pattern). No namespace separation or conflict errors.
- **Theme Distribution**: Themes may be distributed as `.zip` archives (for collections) or single `.css` files (for simple cases). veve detects format by file extension and handles extraction accordingly.
- **Security Model**: Theme file validation warns on suspicious files but does not block user. Users retain responsibility for trusting downloaded themes.

---

## Clarifications

### Session 2025-11-20

Ambiguities resolved via structured clarification process:

- Q: Theme file validation on download - how to prevent malicious themes? → A: Validate file structure (allow `.css`, `.latex`, `.md`, metadata files); warn on non-text files but allow user to proceed (Option B: Validate + warn model)
- Q: How should user themes handle naming conflicts with built-in themes? → A: User themes override built-in themes with same name (Unix config precedence, Option A)
- Q: Configuration file format - JSON, TOML, or YAML? → A: TOML for main config (`veve.toml`), JSON for theme registry (`themes.json`) (Option B)
- Q: Logging strategy - always verbose or on demand? → A: Quiet by default, `--verbose` flag enables detailed logs (Option B: Unix convention)
- Q: Theme distribution format - support only `.zip` or also single files? → A: Support both `.zip` (collections) and single `.css` files (simplicity) (Option B)

---

## Out of Scope (Future Enhancements)

- Live preview with auto-scroll (GUI feature, not CLI)
- Word document export (.docx) - Pandoc supports this; can be added in v1.1
- EPUB/ebook export - Pandoc supports; future enhancement
- Syntax highlighting for non-Markdown formats (Fountain, CriticMarkup, etc.) - Marked 2 has this; veve focuses on Markdown
- Web UI for theme management - CLI theme commands are sufficient for MVP
- Offline Pandoc fallback (pure Go PDF) - Scoped for v2.0 if needed
