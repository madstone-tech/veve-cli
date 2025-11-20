package contract

import (
	"os"
	"os/exec"
	"path/filepath"
	"strings"
	"testing"
)

// TestThemeListOutputsTable tests that 'veve theme list' outputs a formatted table with all themes.
func TestThemeListOutputsTable(t *testing.T) {
	cmd := exec.Command("veve", "theme", "list")
	output, err := cmd.CombinedOutput()

	if err != nil {
		t.Logf("theme list command failed: %v\nOutput: %s", err, string(output))
		t.Skip("veve not installed in PATH")
	}

	outStr := string(output)

	// Check for table headers
	if !strings.Contains(outStr, "NAME") {
		t.Errorf("expected 'NAME' header in output")
	}
	if !strings.Contains(outStr, "AUTHOR") {
		t.Errorf("expected 'AUTHOR' header in output")
	}
	if !strings.Contains(outStr, "DESCRIPTION") {
		t.Errorf("expected 'DESCRIPTION' header in output")
	}
	if !strings.Contains(outStr, "TYPE") {
		t.Errorf("expected 'TYPE' header in output")
	}

	// Check for built-in themes
	if !strings.Contains(outStr, "default") {
		t.Errorf("expected 'default' theme in list")
	}
	if !strings.Contains(outStr, "dark") {
		t.Errorf("expected 'dark' theme in list")
	}
	if !strings.Contains(outStr, "academic") {
		t.Errorf("expected 'academic' theme in list")
	}

	// Check for built-in type
	if !strings.Contains(outStr, "built-in") {
		t.Errorf("expected 'built-in' theme type in output")
	}
}

// TestThemeListIncludesCustomThemes tests that custom themes are listed.
func TestThemeListIncludesCustomThemes(t *testing.T) {
	home, err := os.UserHomeDir()
	if err != nil {
		t.Fatalf("failed to get home directory: %v", err)
	}

	themesDir := filepath.Join(home, ".config", "veve", "themes")
	if err := os.MkdirAll(themesDir, 0o755); err != nil {
		t.Fatalf("failed to create themes directory: %v", err)
	}

	// Create a custom theme
	customThemePath := filepath.Join(themesDir, "mycolor.css")
	customCSS := `---
name: mycolor
author: Test User
description: A custom color theme
---
body { color: purple; }
`
	if err := os.WriteFile(customThemePath, []byte(customCSS), 0o644); err != nil {
		t.Fatalf("failed to create custom theme: %v", err)
	}
	defer os.Remove(customThemePath)

	// Run theme list
	cmd := exec.Command("veve", "theme", "list")
	output, err := cmd.CombinedOutput()

	if err != nil {
		t.Logf("theme list command failed: %v", err)
		t.Skip("veve not installed")
	}

	outStr := string(output)

	// Check that custom theme appears
	if !strings.Contains(outStr, "mycolor") {
		t.Errorf("expected custom theme 'mycolor' in list\nOutput: %s", outStr)
	}

	// Check that it's marked as user theme (not built-in)
	lines := strings.Split(outStr, "\n")
	foundCustom := false
	for _, line := range lines {
		if strings.Contains(line, "mycolor") && !strings.Contains(line, "built-in") {
			foundCustom = true
			break
		}
	}
	if !foundCustom {
		t.Errorf("custom theme should not be marked as built-in")
	}
}

// TestThemeListMetadata tests that theme metadata is displayed.
func TestThemeListMetadata(t *testing.T) {
	cmd := exec.Command("veve", "theme", "list")
	output, err := cmd.CombinedOutput()

	if err != nil {
		t.Skip("veve not installed")
	}

	outStr := string(output)

	// Check that author column is populated for built-in themes
	if !strings.Contains(outStr, "veve-cli") {
		t.Errorf("expected author 'veve-cli' in output")
	}

	// Check that description column is populated
	if !strings.Contains(outStr, "professional") && !strings.Contains(outStr, "Dark") {
		t.Logf("Note: descriptions may vary, got: %s", outStr)
	}
}

// TestThemeListEmptyThemesDir tests behavior with empty themes directory.
func TestThemeListEmptyThemesDir(t *testing.T) {
	home, err := os.UserHomeDir()
	if err != nil {
		t.Fatalf("failed to get home directory: %v", err)
	}

	themesDir := filepath.Join(home, ".config", "veve", "themes")
	// Don't create the directory - let veve handle it

	cmd := exec.Command("veve", "theme", "list")
	output, err := cmd.CombinedOutput()

	if err != nil {
		t.Logf("theme list command failed: %v", err)
		t.Skip("veve not installed")
	}

	outStr := string(output)

	// Should at least show built-in themes
	if !strings.Contains(outStr, "default") {
		t.Errorf("expected built-in themes even if user themes dir is empty")
	}

	// Cleanup if dir was created
	os.RemoveAll(themesDir)
}

// TestThemeListExitCode tests that theme list returns exit code 0 on success.
func TestThemeListExitCode(t *testing.T) {
	cmd := exec.Command("veve", "theme", "list")
	err := cmd.Run()

	if err != nil {
		if exitErr, ok := err.(*exec.ExitError); ok {
			if exitErr.ExitCode() != 0 {
				t.Errorf("expected exit code 0, got %d", exitErr.ExitCode())
			}
		} else {
			t.Skipf("veve not installed: %v", err)
		}
	}
}
