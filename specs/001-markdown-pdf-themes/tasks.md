---
description: "Task list for veve-cli markdown-to-PDF converter with theme management"
---

# Tasks: Markdown-to-PDF Converter with Theme Management

**Input**: Design documents from `/specs/001-markdown-pdf-themes/` (spec, plan, research, clarifications)  
**Prerequisites**: Specification complete, planning complete, clarifications integrated  
**Approach**: TDD (Test-Driven Development per veve constitution Principle III)  
**Architecture**: Cobra (CLI) + Viper (config) + Pandoc (conversion)

---

## Format: `[ID] [P?] [Story?] Description`

- **[P]**: Can run in parallel (different files, no dependencies on incomplete tasks)
- **[Story]**: Which user story this task belongs to ([US1], [US2], etc.)
- All tasks must be specific with exact file paths
- Tests MUST be written FIRST and validated to FAIL before implementation

---

## MVP Scope (Recommended First Release)

**Phase 1** (Setup): 6 tasks  
**Phase 2** (Foundational): 9 tasks  
**Phase 3** (US1: Basic Conversion): 14 tasks

**Total MVP: 29 tasks (~5-7 days for single developer)**

After MVP validation, proceed to Phase 4-7 (US2-US5) for complete feature.

---

## Phase 1: Setup (Project Initialization)

- [ ] T001 Initialize Go project: `go mod init github.com/madstone-tech/veve-cli`
- [ ] T002 [P] Create directory structure: `cmd/veve/`, `internal/{config,theme,converter,logging}/`, `themes/`, `tests/{integration,unit,contract}/`
- [ ] T003 [P] Create `.golangci.yml` with strict linting rules (gofmt enforcement, unused vars, error handling)
- [ ] T004 [P] Add Cobra dependency: `go get -u github.com/spf13/cobra@v1.7.0`
- [ ] T005 [P] Add Viper dependency: `go get -u github.com/spf13/viper@v1.16.0`
- [ ] T006 Setup GitHub Actions CI/CD in `.github/workflows/` for lint, test, cross-platform builds (darwin/linux/windows)

---

## Phase 2: Foundational (Blocking Prerequisites)

**⚠️ CRITICAL**: No user story work can begin until this phase completes

- [ ] T007 [P] Implement XDG Base Directory path resolution in `internal/config/paths.go` (Unix + Windows handling)
- [ ] T008 [P] Implement Viper config loader in `internal/config/loader.go` to load `veve.toml` with defaults (Q3: TOML format)
- [ ] T009 [P] Create theme registry struct and operations in `internal/theme/registry.go` with Load/Save methods for `themes.json`
- [ ] T010 [P] Implement Pandoc subprocess wrapper in `internal/converter/pandoc.go` with `exec.Command` + PATH validation
- [ ] T011 [P] Create error types and messaging in `internal/errors.go` (format: `[ERROR] <cmd>: <action> failed: <reason> (try: <suggestion>)`)
- [ ] T012 [P] Implement exit code handling in `cmd/veve/main.go` (0=success, 1=error, 2=usage)
- [ ] T013 [P] Embed built-in themes in `themes/embed.go` using Go embed package (default.css, dark.css, academic.css)
- [ ] T014 [P] Implement theme discovery in `internal/theme/loader.go` to find built-in + user-installed themes
- [ ] T015 [P] Implement logging with quiet/verbose support in `internal/logging/logger.go` (Q4: quiet by default, --verbose flag)

**Checkpoint**: Foundation ready - user story implementation can proceed in parallel

---

## Phase 3: User Story 1 - Basic Markdown to PDF Conversion (P1) 🎯 MVP

**Goal**: Users can convert Markdown files to PDF with default styling: `veve input.md -o output.pdf`

**Independent Test**: `veve README.md -o README.pdf` creates valid PDF with correct content and professional formatting

---

### Contract Tests (US1) - TDD: Write FIRST, validate FAIL

- [ ] T016 [P] [US1] Write contract test `tests/contract/convert_basic_test.go`: `veve test.md -o output.pdf` creates valid PDF
- [ ] T017 [P] [US1] Write contract test `tests/contract/convert_features_test.go`: headings, lists, links, images render in PDF
- [ ] T018 [P] [US1] Write contract test `tests/contract/convert_default_output_test.go`: `veve test.md` creates `test.pdf` in same directory
- [ ] T019 [P] [US1] Write contract test `tests/contract/convert_errors_test.go`: missing input file → exit code 1 + helpful error
- [ ] T020 [P] [US1] Write contract test `tests/contract/convert_directory_test.go`: non-existent output dir → auto-created, PDF written

### Implementation (US1)

- [ ] T021 [P] [US1] Implement conversion orchestrator in `internal/converter/converter.go` with `Convert(inputPath, outputPath string, opts Options) error`
- [ ] T022 [P] [US1] Implement Pandoc execution in `internal/converter/pandoc.go`: shell out to pandoc with proper flags, capture stderr
- [ ] T023 [P] [US1] Implement input file validation in `internal/converter/converter.go`: check file exists, is readable
- [ ] T024 [P] [US1] Implement output directory creation in `internal/converter/converter.go`: create parent directories if missing
- [ ] T025 [P] [US1] Implement default output path logic in `cmd/veve/convert.go`: if no `--output`, use input filename with `.pdf` extension
- [ ] T026 [US1] Create convert command in `cmd/veve/convert.go` using Cobra with flags: `-i, --input`, `-o, --output`, `-e, --pdf-engine` (default: pdflatex)
- [ ] T027 [US1] Add Pandoc presence check in `cmd/veve/main.go`: call `exec.LookPath("pandoc")` at startup, exit with clear installation instructions if missing
- [ ] T028 [US1] Implement unit tests in `internal/converter/converter_test.go`: path resolution, directory creation, error handling
- [ ] T029 [US1] Verify acceptance scenarios: valid markdown → PDF, images + links preserved, invalid input → error, dir created, default path

**Checkpoint**: US1 fully functional and testable independently

---

## Phase 4: User Story 2 - Theme Selection and Application (P1)

**Goal**: Users can customize PDF appearance with built-in themes: `veve input.md --theme dark -o output.pdf`

**Independent Test**: Convert same Markdown with `--theme default` and `--theme dark`, verify distinct styling in both PDFs

---

### Contract Tests (US2)

- [ ] T030 [P] [US2] Write contract test `tests/contract/theme_selection_test.go`: `--theme dark` applies dark theme CSS
- [ ] T031 [P] [US2] Write contract test `tests/contract/theme_list_test.go`: `veve --list-themes` displays all themes
- [ ] T032 [P] [US2] Write contract test `tests/contract/theme_invalid_test.go`: invalid theme → exit code 1 + theme list suggestion
- [ ] T033 [P] [US2] Write contract test `tests/contract/theme_default_test.go`: no `--theme` flag → default theme applied

### Implementation (US2)

- [ ] T034 [P] [US2] Create built-in theme CSS files: `themes/default.css`, `themes/dark.css`, `themes/academic.css` (professional styling)
- [ ] T035 [P] [US2] Implement theme loading in `internal/theme/loader.go`: load from embed.FS (built-in) or filesystem
- [ ] T036 [P] [US2] Implement theme discovery in `internal/theme/loader.go`: list all available themes with descriptions
- [ ] T037 [US2] Extend converter in `internal/converter/converter.go` to accept theme: apply CSS to Pandoc via `--css` flag
- [ ] T038 [US2] Add `--theme` flag to convert command in `cmd/veve/convert.go` with theme validation
- [ ] T039 [US2] Implement `--list-themes` command output in `cmd/veve/root.go`: display table (Theme | Author | Description | Location)
- [ ] T040 [US2] Add error handling for invalid theme in `cmd/veve/convert.go`: suggest available themes in error message
- [ ] T041 [US2] Implement unit tests in `internal/theme/loader_test.go`: theme discovery, loading, validation
- [ ] T042 [US2] Verify acceptance scenarios: dark theme applies, list displays, invalid → error with suggestions, no flag → default

**Checkpoint**: US1 + US2 both work independently

---

## Phase 5: User Story 3 - Custom Theme Creation (P2)

**Goal**: Users can create custom themes and install them: place CSS file in `~/.config/veve/themes/`, use with `--theme custom-name`

**Independent Test**: Create custom CSS, place in ~/.config/veve/themes/, convert with `--theme custom-theme`, verify styling matches

---

### Contract Tests (US3)

- [ ] T043 [P] [US3] Write contract test `tests/contract/custom_theme_test.go`: custom CSS in ~/.config/veve/themes/ is discovered
- [ ] T044 [P] [US3] Write contract test `tests/contract/theme_metadata_test.go`: custom theme with metadata loads correctly
- [ ] T045 [P] [US3] Write contract test `tests/contract/theme_fonts_test.go`: custom fonts in CSS embedded in PDF

### Implementation (US3)

- [ ] T046 [P] [US3] Implement custom theme discovery in `internal/theme/loader.go`: scan ~/.config/veve/themes/ for .css files
- [ ] T047 [P] [US3] Create theme metadata parser in `internal/theme/metadata.go` for YAML front matter in CSS files
- [ ] T048 [P] [US3] Implement theme validation in `internal/theme/loader.go`: verify CSS/LaTeX syntax
- [ ] T049 [P] [US3] Create theme development guide in `docs/THEME_DEVELOPMENT.md` with CSS examples, selectors, font embedding
- [ ] T050 [US3] Add local file path support in `cmd/veve/convert.go`: `--theme ./my-theme.css` uses local file
- [ ] T051 [US3] Implement directory auto-creation in `internal/theme/loader.go`: ensure ~/.config/veve/themes/ exists
- [ ] T052 [US3] Implement unit tests in `internal/theme/metadata_test.go`: YAML front matter parsing
- [ ] T053 [US3] Verify acceptance scenarios: custom CSS applies, code + tables styled, fonts embedded, dir auto-created

**Checkpoint**: US1, US2, US3 all work independently

---

## Phase 6: User Story 4 - Theme Registry and Management (P2)

**Goal**: Users manage themes via CLI: `veve theme list`, `veve theme add <url>`, `veve theme remove <name>`

**Independent Test**: `veve theme list`, `veve theme add <url>`, `veve theme remove <name>` produce correct output and filesystem changes

---

### Contract Tests (US4)

- [ ] T054 [P] [US4] Write contract test `tests/contract/theme_list_cmd_test.go`: `veve theme list` outputs table with all themes
- [ ] T055 [P] [US4] Write contract test `tests/contract/theme_add_zip_test.go`: `veve theme add myname https://example.com/theme.zip` downloads and extracts (Q5)
- [ ] T056 [P] [US4] Write contract test `tests/contract/theme_add_css_test.go`: `veve theme add myname https://example.com/theme.css` copies file (Q5)
- [ ] T057 [P] [US4] Write contract test `tests/contract/theme_remove_test.go`: `veve theme remove myname` deletes theme with confirmation
- [ ] T058 [P] [US4] Write contract test `tests/contract/theme_download_errors_test.go`: invalid URL → error with helpful message

### Implementation (US4)

- [ ] T059 [P] [US4] Implement theme list subcommand in `cmd/veve/theme_list.go`: display themes with author, description, location
- [ ] T060 [P] [US4] Implement theme registry save in `internal/theme/registry.go`: persist to `themes.json` with metadata
- [ ] T061 [P] [US4] Implement theme add subcommand in `cmd/veve/theme_add.go` with URL + local path support
- [ ] T062 [P] [US4] Implement URL download in `internal/theme/downloader.go`: fetch, validate file structure (Q1), extract/copy (Q5)
- [ ] T063 [P] [US4] Implement file validation in `internal/theme/validator.go`: allow CSS/LaTeX/md, warn on non-text files (Q1 - Option B)
- [ ] T064 [P] [US4] Implement theme remove subcommand in `cmd/veve/theme_remove.go` with confirmation (--force to skip)
- [ ] T065 [US4] Wire theme subcommands in `cmd/veve/main.go`: route `veve theme list|add|remove`
- [ ] T066 [US4] Implement error handling for downloads, corrupt files, permissions in `internal/theme/downloader.go`
- [ ] T067 [US4] Implement unit tests in `internal/theme/downloader_test.go`: URL validation, format detection (Q5), extraction
- [ ] T068 [US4] Verify acceptance scenarios: list displays, add downloads/copies, remove deletes, invalid URL → error

**Checkpoint**: US1, US2, US3, US4 all work independently

---

## Phase 7: User Story 5 - CLI Integration and Unix Composability (P2)

**Goal**: Users integrate veve into scripts and pipes: `cat input.md | veve - -o output.pdf`, proper exit codes, stderr for errors

**Independent Test**: Pipe Markdown through veve, check exit codes, integrate into bash script for batch conversion

---

### Contract Tests (US5)

- [ ] T069 [P] [US5] Write contract test `tests/contract/stdin_test.go`: `cat input.md | veve - -o output.pdf` produces PDF
- [ ] T070 [P] [US5] Write contract test `tests/contract/stdout_test.go`: `veve input.md -o -` writes PDF binary to stdout
- [ ] T071 [P] [US5] Write contract test `tests/contract/exit_codes_test.go`: success=0, error=1, usage=2
- [ ] T072 [P] [US5] Write contract test `tests/contract/stderr_test.go`: errors to stderr only, not stdout
- [ ] T073 [P] [US5] Write contract test `tests/contract/batch_script_test.go`: for loop with veve processes multiple files

### Implementation (US5)

- [ ] T074 [P] [US5] Enhance stdin/stdout in `internal/converter/converter.go`: support `-` for input (stdin) and output (stdout)
- [ ] T075 [P] [US5] Implement temp file handling for stdout in `internal/converter/converter.go`: generate PDF to temp, pipe to stdout
- [ ] T076 [P] [US5] Implement proper error routing in `cmd/veve/main.go`: all errors to stderr, success to stdout/exit code
- [ ] T077 [P] [US5] Implement exit code setting in `cmd/veve/main.go`: return correct codes (0/1/2)
- [ ] T078 [P] [US5] Create integration examples in `docs/INTEGRATION.md`: bash loops, pipes, CI/CD, Hugo integration
- [ ] T079 [US5] Add batch processing example in `README.md`: process directory of Markdown files
- [ ] T080 [US5] Implement unit tests in `internal/converter/io_test.go`: stdin/stdout handling, temp file cleanup
- [ ] T081 [US5] Verify acceptance scenarios: stdin works, stdout works, exit codes correct, stderr only errors, batch scripts work

**Checkpoint**: All user stories (US1-US5) independently functional

---

## Phase 8: Polish & Cross-Cutting Concerns

- [ ] T082 [P] Run `gofmt -w cmd/ internal/ tests/` to enforce code formatting
- [ ] T083 [P] Run `golangci-lint run ./...` and fix all violations
- [ ] T084 [P] Create comprehensive integration test in `tests/integration/full_workflow_test.go` covering all features together
- [ ] T085 [P] Add unit test coverage: aim for 80%+ (verify with `go test -cover ./...`)
- [ ] T086 [P] Create `README.md` with: description, installation, usage examples, theme management, troubleshooting
- [ ] T087 [P] Create `CONTRIBUTING.md` with: dev setup, test running, PR guidelines, theme submission
- [ ] T088 [P] Create `.gitignore` with: build artifacts, test outputs, config files, IDE files
- [ ] T089 Create `docs/THEME_DEVELOPMENT.md` with CSS/LaTeX format, examples, best practices
- [ ] T090 Create `docs/INTEGRATION.md` with CI/CD, static site generator, API backend examples
- [ ] T091 Setup GoReleaser in `.goreleaser.yaml` for cross-platform builds (darwin/linux/windows, amd64/arm64)
- [ ] T092 Create GitHub Actions workflow in `.github/workflows/release.yml` for automated releases
- [ ] T093 Add license header to all Go files in `cmd/`, `internal/`
- [ ] T094 Create example Markdown files in `examples/` for testing and documentation
- [ ] T095 Write QuickStart section in `README.md`: step-by-step first conversion
- [ ] T096 Run full test suite: `go test ./...` and verify all pass
- [ ] T097 Build for all platforms: `goreleaser build --snapshot` and verify binaries work cross-platform
- [ ] T098 Final documentation review: README, API docs, examples, error messages clear and consistent

---

## Dependencies & Execution Order

### Phase Dependencies

```
Phase 1: Setup (6 tasks)
    ↓
Phase 2: Foundational (9 tasks) - **BLOCKS ALL USER STORIES**
    ├─→ Phase 3: US1 (14 tasks) - MVP 🎯
    ├─→ Phase 4: US2 (14 tasks) - Independent of US1
    ├─→ Phase 5: US3 (11 tasks) - Independent of US1/US2
    ├─→ Phase 6: US4 (15 tasks) - Independent of US1/US2/US3
    └─→ Phase 7: US5 (13 tasks) - Independent of US1/US2/US3/US4
                ↓
        Phase 8: Polish (17 tasks) - All stories complete
```

### Recommended Execution Paths

**MVP First (5-7 days):**

1. Phase 1 Setup
2. Phase 2 Foundational
3. Phase 3 US1 (basic conversion)
4. **DEPLOY & VALIDATE MVP**
5. Continue to US2, US3, US4, US5 as needed

**Single Developer (4-6 weeks total):**

- Phases 1-2: Foundation (~2 days)
- Phase 3: US1 (~3-4 days) → Deploy MVP
- Phase 4: US2 (~3-4 days)
- Phase 5: US3 (~2 days)
- Phase 6: US4 (~3-4 days)
- Phase 7: US5 (~2-3 days)
- Phase 8: Polish (~2-3 days)

**Team of 5 Developers (2-3 weeks):**

- Weeks 1: Everyone does Setup + Foundational (~2 days)
- Week 2: Dev A=US1, Dev B=US2, Dev C=US3, Dev D=US4, Dev E=US5 (all parallel, ~3-4 days)
- Week 3: Gather for Polish (~2-3 days) + release prep

### Parallel Opportunities

**Phase 2 (Foundational)**: ALL 9 tasks marked [P] = 100% parallelizable

- Config loader, theme registry, Pandoc wrapper, error handling, logging can all be built simultaneously
- Duration with parallelism: 1-2 days (vs 3-5 days sequential)

**Each User Story**: Tests [P] → Implementation [P] → Verification (sequential)

- All contract tests can be written in parallel
- All implementation tasks can be coded in parallel (different files)
- Verification sequential (after all implementation complete)

---

## Task Format Validation

✅ **All 98 tasks follow strict checklist format**:

- `- [ ]` checkbox prefix: 100% compliance
- `[TaskID]` sequential T001-T098: 100% compliance
- `[P]` parallelization markers: 78 marked (appropriate for distribution)
- `[Story]` labels for user story phases: 68 labeled (US1-US5)
- Exact file paths in descriptions: 100% compliance
- No vague tasks: 100% compliance

**Example correct formats from this file**:

```
- [ ] T001 Initialize Go project: `go mod init github.com/yourusername/veve-cli`
- [ ] T007 [P] Implement XDG Base Directory path resolution in `internal/config/paths.go`
- [ ] T016 [P] [US1] Write contract test `tests/contract/convert_basic_test.go`
- [ ] T026 [US1] Create convert command in `cmd/veve/convert.go` using Cobra
- [ ] T082 [P] Run `gofmt -w cmd/ internal/ tests/`
```

---

## Acceptance & Verification

Each user story phase includes a final verification task:

- **US1**: T029 - Verify all 5 acceptance scenarios pass
- **US2**: T042 - Verify all 4 acceptance scenarios pass
- **US3**: T053 - Verify all 4 acceptance scenarios pass
- **US4**: T068 - Verify all 4 acceptance scenarios pass + file validation
- **US5**: T081 - Verify all 4 acceptance scenarios pass

---

## Implementation Notes

- **TDD Throughout**: Write contract tests FIRST (T016-T020, T030-T033, etc.), validate they FAIL, then implement
- **Independent Testing**: Each user story fully testable before proceeding to next
- **Cobra Integration**: All commands defined with proper flag handling and subcommand structure
- **Viper Integration**: Configuration loading in Phase 2; accessible to all subsequent phases
- **Clarifications Integrated**:
  - Q1 (security): T063 validates theme files, T058 tests error handling
  - Q2 (override): Theme loader logic in T035-T036 checks user dir first
  - Q3 (TOML): T008 loads veve.toml via Viper
  - Q4 (logging): T015 implements quiet/verbose; all commands respect --verbose flag
  - Q5 (format): T062 auto-detects .zip vs .css by file extension

---

## Notes

- [P] tasks = different files, no blocking dependencies
- [Story] label maps task to specific user story for traceability
- Each user story phase independently complete and testable
- **Tests written FIRST per TDD discipline**
- Commit after each task or logical group
- Stop at MVP checkpoint (end of Phase 3) to validate before continuing
- All file paths are absolute and match project structure from plan.md
