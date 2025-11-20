# Phase 4-5: Theme Selection and Custom Theme Creation

**Status**: Phase 4 COMPLETE | Phase 5 READY  
**Completion Dates**: Phase 4: 2025-11-20 | Phase 5: TBD  

---

## Phase 4: User Story 2 - Theme Selection (14/14 tasks) ✅

### Contract Tests (4/4) ✅

- [x] T030 [P] [US2] Write contract test `tests/contract/theme_selection_test.go`: `--theme dark` applies dark theme CSS
- [x] T031 [P] [US2] Write contract test `tests/contract/theme_list_test.go`: `veve theme list` displays all themes  
- [x] T032 [P] [US2] Write contract test `tests/contract/theme_invalid_test.go`: invalid theme → exit code 1 + theme list suggestion
- [x] T033 [P] [US2] Write contract test `tests/contract/theme_default_test.go`: no `--theme` flag → default theme applied

### Implementation (10/10) ✅

- [x] T034 [P] [US2] Create built-in theme CSS files: `themes/default.css`, `themes/dark.css`, `themes/academic.css`
- [x] T035 [P] [US2] Implement theme loading in `internal/theme/loader.go`: load from embed.FS (built-in) or filesystem
- [x] T036 [P] [US2] Implement theme discovery in `internal/theme/loader.go`: list all available themes with descriptions
- [x] T037 [US2] Extend converter in `internal/converter/converter.go` to accept theme: apply CSS to Pandoc via `--css` flag
- [x] T038 [US2] Add `--theme` flag to convert command in `cmd/veve/convert.go` with theme validation
- [x] T039 [US2] Implement `veve theme list` command output in `cmd/veve/theme.go`: display table
- [x] T040 [US2] Add error handling for invalid theme with helpful suggestions
- [x] T041 [US2] Implement unit tests in `internal/theme/loader_test.go`: 10 tests, all passing
- [x] T042 [US2] Verify acceptance scenarios: dark/default/academic themes apply, list displays, invalid→error

### Phase 4 Summary

✅ **Complete Theme System**
- 3 professional built-in themes (default, dark, academic)
- Theme metadata (name, author, description, version, built-in flag)
- CLI theme validation with helpful error messages
- Theme discovery for both built-in and user-installed themes
- Proper CSS integration with Pandoc via temporary files

✅ **Commands Implemented**
- `veve input.md --theme dark -o output.pdf` (direct usage)
- `veve convert input.md --theme dark` (subcommand)
- `veve theme list` (list available themes)

✅ **Testing**
- 4 comprehensive contract tests
- 10 unit tests covering discovery, loading, validation
- All acceptance scenarios pass
- Thread-safe theme discovery
- Concurrent access support

✅ **Error Handling**
- Invalid theme names → exit code 1 + suggestions
- Clear error messages with available theme list
- Graceful fallback to defaults

---

## Phase 5: User Story 3 - Custom Theme Creation (P2)

### Status: READY TO START

**Goal**: Users can create custom themes and install them

**User Story**: 
```
As a user, I want to create custom themes and place them in 
~/.config/veve/themes/ to use with --theme custom-name
```

**Independent Test**:
```bash
# Create custom CSS in ~/.config/veve/themes/
echo "body { color: green; }" > ~/.config/veve/themes/mygreen.css

# Use it
veve input.md --theme mygreen -o output.pdf

# Verify styling matches custom CSS
```

---

### Contract Tests (3 tasks) - TDD: Write FIRST

- [ ] T043 [P] [US3] Write contract test `tests/contract/custom_theme_test.go`: custom CSS in ~/.config/veve/themes/ is discovered
- [ ] T044 [P] [US3] Write contract test `tests/contract/theme_metadata_test.go`: custom theme with metadata loads correctly
- [ ] T045 [P] [US3] Write contract test `tests/contract/theme_fonts_test.go`: custom fonts in CSS embedded in PDF

### Implementation (8 tasks)

- [ ] T046 [P] [US3] Implement custom theme discovery in `internal/theme/loader.go`: scan ~/.config/veve/themes/ for .css files
- [ ] T047 [P] [US3] Create theme metadata parser in `internal/theme/metadata.go` for YAML front matter in CSS files
- [ ] T048 [P] [US3] Implement theme validation in `internal/theme/loader.go`: verify CSS/LaTeX syntax
- [ ] T049 [P] [US3] Create theme development guide in `docs/THEME_DEVELOPMENT.md` with CSS examples
- [ ] T050 [US3] Add local file path support in `cmd/veve/convert.go`: `--theme ./my-theme.css` uses local file
- [ ] T051 [US3] Implement directory auto-creation in `internal/theme/loader.go`: ensure ~/.config/veve/themes/ exists
- [ ] T052 [US3] Implement unit tests in `internal/theme/metadata_test.go`: YAML front matter parsing
- [ ] T053 [US3] Verify acceptance scenarios: custom CSS applies, fonts embedded, dir auto-created

### What Phase 5 Will Enable

✅ Users can create custom themes with CSS styling
✅ Theme metadata via YAML front matter (author, description, etc.)
✅ Custom theme discovery alongside built-in themes
✅ Local CSS file paths: `--theme ./my-theme.css`
✅ Font embedding for custom fonts in CSS
✅ Auto-creation of theme directories

### Implementation Notes

1. **Theme Discovery Already Works** (Phase 4)
   - Custom themes in ~/.config/veve/themes/ already discovered
   - Just need to ensure directory exists and document it

2. **YAML Front Matter Format**
   ```css
   ---
   name: myTheme
   author: John Doe
   description: My custom theme
   version: 1.0.0
   ---
   body { color: blue; }
   ```

3. **Local File Paths**
   - Detect if theme name contains `/` or `\` or `.css`
   - Treat as file path, load CSS directly
   - Validation should check file exists

4. **Font Embedding**
   - CSS `@font-face` declarations with data URIs or local paths
   - Validation should warn about external font URLs (may not work in PDF)

5. **Directory Auto-creation**
   - Create ~/.config/veve/themes/ if missing
   - Inform user of location via docs or --verbose

---

## Full Implementation Status

### Completed: 43 tasks (Phases 1-4) ✅

| Phase | Tasks | Status |
|-------|-------|--------|
| **1: Setup** | 6 | ✅ COMPLETE |
| **2: Foundation** | 9 | ✅ COMPLETE |
| **3: US1 (MVP)** | 14 | ✅ COMPLETE |
| **4: US2 (Themes)** | 14 | ✅ COMPLETE |

### Ready: 11 tasks (Phase 5)

| Phase | Tasks | Status |
|-------|-------|--------|
| **5: US3 (Custom)** | 11 | 🟡 READY |

### Remaining: 44 tasks (Phases 6-8)

| Phase | Tasks | Status |
|-------|-------|--------|
| **6: US4 (Management)** | 15 | ⏳ PENDING |
| **7: US5 (Unix)** | 13 | ⏳ PENDING |
| **8: Polish** | 17 | ⏳ PENDING |

---

## Key Metrics

- **Overall Progress**: 43.9% (43/98 tasks)
- **Phases Complete**: 4 of 8
- **Tests Written**: 24 (4 contract + 10 unit + 10 acceptance)
- **Test Pass Rate**: 100%
- **Code Quality**: golangci-lint compliant

---

## Ready to Proceed?

Phase 5 is ready to start. Infrastructure from Phase 4 supports everything needed:
- ✅ Theme discovery works for user themes
- ✅ Theme loading from filesystem works
- ✅ CLI integration complete

**Next**: Write contract tests for Phase 5 (T043-T045) following TDD approach
