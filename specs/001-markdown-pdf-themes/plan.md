# Implementation Plan: Markdown-to-PDF Converter with Theme Management

**Branch**: `001-markdown-pdf-themes` | **Date**: 2025-11-20 | **Spec**: `/specs/001-markdown-pdf-themes/spec.md`  
**Input**: Feature specification + clarifications (5 ambiguities resolved via structured process)  
**Status**: Ready for Phase 1 design execution

## Summary

veve-cli is a cross-platform CLI markdown-to-PDF converter inspired by Marked 2 (macOS app). It leverages Pandoc for professional document conversion and provides an extensible theme system for customizing PDF appearance. Core MVP enables simple conversion (`veve input.md -o output.pdf`) with theme selection (`--theme dark`). Built-in themes + theme management subcommand (`veve theme list|add|remove`) enable both casual users and content teams. Unix philosophy: composable, scriptable, stdin/stdout-aware. **Architecture leverages Cobra for CLI structure and Viper for configuration management.**

## Technical Context

**Language/Version**: Go 1.20+ (idiomatic, modern stdlib)

**Primary Dependencies**:
- Cobra 1.7+ (CLI framework: subcommands, flags, help)
- Viper 1.16+ (configuration management: TOML parsing, env var merging)
- Pandoc 2.18+ (external: markdown→PDF conversion engine)
- Go standard library (os, exec, filepath, io, encoding/json, archive/zip)

**Storage**: XDG Base Directory compliant
- Config: `~/.config/veve/veve.toml` (TOML format - resolved via Q3)
- Themes: `~/.config/veve/themes/` (CSS/LaTeX files)
- Registry: `~/.config/veve/themes.json` (JSON manifest)

**Testing**: `testing` package (Go standard) + table-driven tests
- Unit: command logic, theme loading, config parsing, validation
- Integration: end-to-end conversions, theme conflicts (user override), file validation
- Contract: CLI interface (args, exit codes, stderr/stdout, --verbose flag)

**Target Platform**: Linux (amd64, arm64), macOS (amd64, arm64), Windows (amd64)

**Project Type**: Single CLI binary (no web/mobile)

**Performance Goals**:
- Simple conversion: <5 seconds per file (excluding Pandoc startup)
- Theme switching: zero overhead (just path selection)
- Batch processing: 100s of files via shell loops or pipes

**Constraints**:
- Pandoc MUST be in PATH (required dependency)
- Memory: reasonable for 10MB+ files
- Error messages: clear, actionable, <200 chars
- Exit codes: 0 (success), 1 (error), 2 (usage)
- Logging: quiet by default; `--verbose` enables details (Q4 clarification)

**Scale/Scope**:
- ~2000-3000 lines of Go code (MVP)
- 3 built-in themes (embedded via Go embed)
- Unlimited custom themes (user-managed)
- Theme distribution: .zip (collections) + .css (single files) - Q5 clarification

---

## Constitution Check

*GATE: Must pass before Phase 0 research. Re-check after Phase 1 design.*

### Veve Constitution Alignment (v1.0.0)

**Principle I: Language & Tooling (Go-First)** ✅
- MUST use modern, idiomatic Go → Go 1.20+ confirmed
- Go Modules for dependency management → Cobra, Viper via go.mod
- GoReleaser for cross-platform builds → To configure in Phase 2
- **Status**: PASS

**Principle II: CLI-First Architecture** ✅
- MUST expose all features via CLI → Cobra framework (primary dep)
- Cobra framework mandatory → Confirmed
- Text + JSON output formats → Required in spec (FR-016+)
- stdin/stdout/stderr separation → Quiet by default + --verbose (Q4)
- **Status**: PASS

**Principle III: Test-Driven Development** ✅
- MUST have unit + integration tests → 27 test tasks in spec
- 80%+ coverage for new code → Measured post-implementation
- Contract tests for CLI interfaces → Required for Cobra commands
- **Status**: PASS

**Principle IV: Theme Management via Registry** ✅
- Registry: themes.json in ~/.config/veve/ → Confirmed (Q3)
- Storage: ~/.config/veve/themes/ → XDG-compliant
- Theme conflicts: user override built-in → Q2 clarification (Option A)
- Remote + local sources: .zip + .css → Q5 clarification (Option B)
- Theme selection via --theme flag → FR-003
- **Status**: PASS

**Principle V: Code Quality & Documentation Standards** ✅
- gofmt enforcement → CI requirement
- golangci-lint configuration → To provide
- GoDoc comments on public APIs → Required
- **Status**: PASS

### Gate Evaluation

**RESULT**: ✅ **PASS** - No violations. All principles satisfied or enhanced by clarifications.

**Key Alignment with Clarifications**:
- Q3 (TOML config): Viper natively supports TOML → seamless integration
- Q2 (theme override): File path precedence simpler than namespacing → aligned with Principle IV
- Q4 (quiet by default): Unix philosophy standard → aligns Principle II
- Q5 (.zip + .css): Flexible distribution → enables Principle IV extensibility

---

## Project Structure

### Documentation (this feature)

```text
specs/001-markdown-pdf-themes/
├── spec.md              # Feature spec (195 lines, clarifications integrated)
├── plan.md              # This file (implementation plan + constitution check)
├── research.md          # Phase 0: technical decisions (already completed)
├── data-model.md        # Phase 1: entities, relationships, validation
├── contracts/           # Phase 1: Cobra command contracts
│   └── cli-interface.md # CLI structure, flags, exit codes
├── quickstart.md        # Phase 1: developer guide
└── tasks.md             # Phase 2: 97 tasks (to be updated with clarifications)
```

### Source Code (repository root)

```text
cmd/veve/
├── main.go              # Entry point, Cobra root init, Viper config setup
├── root.go              # Root command setup, global --verbose flag
├── convert.go           # Primary convert command (input file → PDF)
├── theme.go             # `veve theme` command group setup
├── theme_list.go        # List themes command
├── theme_add.go         # Add theme command (with .zip/.css detection - Q5)
├── theme_remove.go      # Remove theme command
└── errors.go            # Error types, exit code handling

internal/
├── config/
│   ├── config.go        # Config struct with fields from veve.toml
│   ├── loader.go        # Viper integration for TOML loading (Q3)
│   └── paths.go         # XDG path resolution (Unix + Windows)
├── theme/
│   ├── registry.go      # Theme discovery, list, add, remove
│   ├── loader.go        # Load theme from embed or filesystem
│   ├── validator.go     # File validation for downloads (FR-017 - Q1)
│   └── metadata.go      # Theme manifest parsing
├── converter/
│   ├── converter.go     # Orchestrate Markdown→PDF conversion
│   ├── pandoc.go        # Pandoc subprocess wrapper
│   └── output.go        # PDF output and directory handling
└── logging/
    └── logger.go        # Logging with quiet/verbose support (Q4)

themes/                  # Built-in themes (embedded)
├── default.css
├── dark.css
├── academic.css
└── embed.go             # Go embed directive

tests/
├── integration/         # End-to-end tests
├── unit/                # Unit tests
└── contract/            # CLI contract tests (Cobra commands)

.golangci.yml           # Linting config
go.mod                  # Cobra, Viper, pinned versions
go.sum                  # Dependency checksums
```

**Structure Rationale**:
- `cmd/veve/`: Cobra command handlers (CLI entry points)
- `internal/`: Core logic packages (config, theme, converter, logging)
- `themes/`: Embedded built-in themes (compiled into binary)
- `tests/`: Comprehensive test suites

---

## Clarifications Integration Map

| # | Clarification | Decision | Affected Areas |
|---|---------------|----------|----------------|
| Q1 | Theme security | Validate + warn (Option B) | FR-017, theme_add.go, validator.go |
| Q2 | Theme naming | User override (Option A) | registry.go, loader.go, registry structure |
| Q3 | Config format | TOML (Option B) | internal/config/loader.go, veve.toml storage |
| Q4 | Logging | Quiet + --verbose (Option B) | root.go, logging/logger.go, all commands |
| Q5 | Distribution | .zip + .css (Option B) | theme_add.go, auto-detection logic |

---

## Phase 0: Research (Completed)

All clarifications resolved via structured ambiguity process (see Clarifications section in spec.md).

**Research Areas Addressed**:
- Cobra command structure and flag handling
- Viper configuration loading (TOML format)
- Theme file validation patterns
- Quiet/verbose logging integration
- Theme distribution format detection
- File system operations (XDG paths, directory creation)

**No remaining "NEEDS CLARIFICATION" items** - all technical decisions made.

---

## Phase 1: Design & Contracts (Next Steps)

**Deliverables to generate**:

1. **data-model.md**: Entity definitions (Config, Theme, ThemeRegistry, ConversionRequest, PandocConfig) with:
   - Fields and types
   - Validation rules (from FR requirements)
   - Relationships and override behavior (Q2)
   - State transitions if applicable

2. **contracts/cli-interface.md**: Cobra command interface specifications:
   - Command hierarchy: `veve`, `veve theme list|add|remove`
   - Global flags: `--verbose`, `--help`
   - Convert flags: `-o`, `-t`, `-e`, `--template`, `--css`
   - Theme flags: URL/path handling for add, confirmation for remove
   - Exit codes: 0/1/2 semantics
   - Output formats: human-readable (quiet by default) + verbose mode (Q4)

3. **quickstart.md**: Developer guide for:
   - Installation and setup
   - First conversion example
   - Theme selection walkthrough
   - Configuration (veve.toml editing, env var overrides)
   - Logging/troubleshooting with --verbose
   - Integration examples

4. **Agent context update**: Record Cobra/Viper technology stack decisions

---

## Phase 2: Tasks (via /speckit.tasks)

Will generate 97+ tasks organized by:
- Phase 1: Setup (Cobra scaffold, Viper init)
- Phase 2: Foundational (config loader, theme registry, Pandoc wrapper)
- Phase 3-7: User stories (conversion, themes, theme mgmt, composability)
- Phase 8: Polish (docs, linting, testing, release)

**Updated with clarifications**: File validation (Q1), theme override (Q2), TOML loading (Q3), logging (Q4), format detection (Q5)

---

## Next Steps

**Immediate**:
1. ✅ Specification complete + clarifications integrated
2. ✅ Planning complete + constitution verified
3. ⏭ **Execute Phase 1 design**:
   ```bash
   /speckit.plan 001  # Continue to generate data-model, contracts, quickstart
   ```

**Then**:
4. Generate task breakdown: `/speckit.tasks 001`
5. Begin implementation with TDD approach (tests first, per Principle III)

---

**Status**: ✅ **READY FOR PHASE 1 DESIGN EXECUTION**

All clarifications integrated. Constitution verified. Technical architecture defined. Ready to proceed with data modeling, contract generation, and task breakdown.
