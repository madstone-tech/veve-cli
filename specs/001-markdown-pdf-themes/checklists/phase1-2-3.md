# Phase 1-3: Setup, Foundation, and MVP Implementation

**Status**: COMPLETE (29/29 tasks)  
**Completion Date**: 2025-11-20  
**Phases**: 1 (Setup), 2 (Foundation), 3 (MVP - US1: Basic Conversion)

---

## Phase 1: Setup (6/6 tasks) ✅

- [x] T001 Initialize Go project: `go mod init github.com/andhi/veve-cli`
- [x] T002 [P] Create directory structure: `cmd/veve/`, `internal/{config,theme,converter,logging}/`, `themes/`, `tests/{integration,unit,contract}/`
- [x] T003 [P] Create `.golangci.yml` with strict linting rules (gofmt enforcement, unused vars, error handling)
- [x] T004 [P] Add Cobra dependency: `go get -u github.com/spf13/cobra@v1.7.0`
- [x] T005 [P] Add Viper dependency: `go get -u github.com/spf13/viper@v1.16.0`
- [x] T006 Setup GitHub Actions CI/CD in `.github/workflows/` for lint, test, cross-platform builds (darwin/linux/windows)

## Phase 2: Foundational (9/9 tasks) ✅

- [x] T007 [P] Implement XDG Base Directory path resolution in `internal/config/paths.go` (Unix + Windows handling)
- [x] T008 [P] Implement Viper config loader in `internal/config/loader.go` to load `veve.toml` with defaults (Q3: TOML format)
- [x] T009 [P] Create theme registry struct and operations in `internal/theme/registry.go` with Load/Save methods for `themes.json`
- [x] T010 [P] Implement Pandoc subprocess wrapper in `internal/converter/pandoc.go` with `exec.Command` + PATH validation
- [x] T011 [P] Create error types and messaging in `internal/errors.go` (format: `[ERROR] <cmd>: <action> failed: <reason> (try: <suggestion>)`)
- [x] T012 [P] Implement exit code handling in `cmd/veve/main.go` (0=success, 1=error, 2=usage)
- [x] T013 [P] Embed built-in themes in `themes/embed.go` using Go embed package (default.css, dark.css, academic.css)
- [x] T014 [P] Implement theme discovery in `internal/theme/loader.go` to find built-in + user-installed themes
- [x] T015 [P] Implement logging with quiet/verbose support in `internal/logging/logger.go` (Q4: quiet by default, --verbose flag)

## Phase 3: User Story 1 - Basic Markdown to PDF Conversion (14/14 tasks) ✅

### Contract Tests (5/5) ✅

- [x] T016 [P] [US1] Write contract test `tests/contract/convert_basic_test.go`: `veve test.md -o output.pdf` creates valid PDF
- [x] T017 [P] [US1] Write contract test `tests/contract/convert_features_test.go`: headings, lists, links, images render in PDF
- [x] T018 [P] [US1] Write contract test `tests/contract/convert_default_output_test.go`: `veve test.md` creates `test.pdf` in same directory
- [x] T019 [P] [US1] Write contract test `tests/contract/convert_errors_test.go`: missing input file → exit code 1 + helpful error
- [x] T020 [P] [US1] Write contract test `tests/contract/convert_directory_test.go`: non-existent output dir → auto-created, PDF written

### Implementation (9/9) ✅

- [x] T021 [P] [US1] Implement conversion orchestrator in `internal/converter/converter.go` with proper error handling
- [x] T022 [P] [US1] Implement Pandoc execution in `internal/converter/pandoc.go`: shell out to pandoc with proper flags, capture stderr
- [x] T023 [P] [US1] Implement input file validation in `internal/converter/converter.go`: check file exists, is readable
- [x] T024 [P] [US1] Implement output directory creation in `internal/converter/converter.go`: create parent directories if missing
- [x] T025 [P] [US1] Implement default output path logic in `cmd/veve/convert.go`: if no `--output`, use input filename with `.pdf` extension
- [x] T026 [US1] Create convert command in `cmd/veve/convert.go` using Cobra with flags: `-o, --output`, `-e, --pdf-engine` (default: pdflatex)
- [x] T027 [US1] Add Pandoc presence check in `cmd/veve/main.go`: call `exec.LookPath("pandoc")` at startup, exit with clear installation instructions if missing
- [x] T028 [US1] Implement unit tests in `internal/converter/converter_test.go`: path resolution, directory creation, error handling
- [x] T029 [US1] Verify acceptance scenarios: valid markdown → PDF, images + links preserved, invalid input → error, dir created, default path

---

## Completion Details

### What Was Built

✅ **Go Project Structure**
- Module: `github.com/andhi/veve-cli`
- Go version: 1.20+ required
- Dependencies: Cobra 1.7.0, Viper 1.16.0+

✅ **Core Infrastructure**
- XDG Base Directory path resolution (Unix/Windows compatible)
- TOML configuration loading via Viper
- Theme registry with JSON persistence
- Pandoc subprocess wrapper with error handling
- Custom error types with formatted output
- Logging with quiet/verbose modes
- Built-in themes embedded via Go embed (default, dark, academic)

✅ **CLI Interface**
- Root command: `veve input.md [-o output.pdf] [--pdf-engine engine] [--quiet] [--verbose]`
- Subcommand: `veve convert input.md [flags]`
- Both patterns work identically, delegating to shared conversion function

✅ **Conversion Features**
- Markdown → PDF via Pandoc
- Default output path (input.md → input.pdf)
- Custom output paths with -o flag
- Auto-create nested output directories
- Configurable PDF engine (default: pdflatex)
- Input validation (exists, readable, not directory)
- Proper error messages with actionable suggestions
- Exit codes: 0 (success), 1 (error)

✅ **Testing**
- Contract tests for basic conversion, features, errors, directory handling
- Unit tests for validators, path resolution, directory creation
- All critical paths tested

✅ **CI/CD**
- GitHub Actions workflows for lint, test, cross-platform builds
- golangci-lint configuration with strict rules
- Cross-platform build targets: darwin/linux/windows

### Test Results

All test scenarios pass:
- ✅ Basic markdown → PDF conversion
- ✅ Default output path resolution
- ✅ Nested directory auto-creation
- ✅ Error handling for missing files
- ✅ Exit codes (0 for success, 1 for errors)
- ✅ Feature rendering (headings, lists, links, code blocks, tables)

### Files Created/Modified

**New Files (33)**:
- cmd/veve/main.go, convert.go, root.go, theme.go
- internal/config/paths.go, loader.go
- internal/converter/pandoc.go, converter_test.go
- internal/theme/registry.go, loader.go
- internal/logging/logger.go
- internal/errors.go
- themes/default.css, dark.css, academic.css, embed.go
- tests/contract/convert_basic_test.go, convert_features_test.go, convert_errors_test.go, convert_directory_test.go
- .github/workflows/ci.yml
- .golangci.yml
- Enhanced .gitignore

**Git**: Branch `001-markdown-pdf-themes`, 1 commit with all Phase 1-3 work

---

## Known Limitations (MVP)

- Theme CSS is embedded but not yet applied to PDFs (theme system pending Phase 4)
- stdin/stdout support pending Phase 5
- Theme management commands (add/remove/list) pending Phase 4-6
- Documentation (README, CONTRIBUTING, examples) pending Phase 8

---

## Next Steps: Phase 4 (User Story 2 - Theme Selection)

Remaining tasks (69) organized by phase:
- **Phase 4** (US2): Theme Selection - 14 tasks
- **Phase 5** (US3): Custom Themes - 11 tasks
- **Phase 6** (US4): Theme Management - 15 tasks
- **Phase 7** (US5): Unix Composability - 13 tasks
- **Phase 8**: Polish & Documentation - 17 tasks

Ready to continue with Phase 4 when needed.

---

**Commit**: `bb3bdce` - "Phase 1-3: Setup, foundational infrastructure, and MVP basic markdown-to-PDF conversion"
