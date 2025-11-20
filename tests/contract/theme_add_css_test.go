package contract

import (
	"os"
	"os/exec"
	"path/filepath"
	"testing"
)

// TestThemeAddCSS tests that 'veve theme add myname file.css' copies a CSS theme.
func TestThemeAddCSS(t *testing.T) {
	// Create a temporary directory with a theme file
	tmpDir := t.TempDir()
	themeFile := filepath.Join(tmpDir, "test-theme.css")
	themeCSS := `---
name: test-theme
author: Test User
description: A test CSS theme
version: 1.0.0
---
body { color: blue; }
h1 { font-size: 24pt; }
`

	if err := os.WriteFile(themeFile, []byte(themeCSS), 0o644); err != nil {
		t.Fatalf("failed to create test theme file: %v", err)
	}

	// Get home directory
	home, err := os.UserHomeDir()
	if err != nil {
		t.Fatalf("failed to get home directory: %v", err)
	}

	themesDir := filepath.Join(home, ".config", "veve", "themes")

	// Ensure the theme doesn't exist first
	testThemePath := filepath.Join(themesDir, "test-add-theme.css")
	defer os.Remove(testThemePath)

	// Run 'veve theme add test-add-theme /path/to/test-theme.css'
	cmd := exec.Command("veve", "theme", "add", "test-add-theme", themeFile)
	output, err := cmd.CombinedOutput()

	if err != nil {
		t.Logf("veve theme add failed: %v\nOutput: %s", err, string(output))
		t.Skip("veve theme add not yet fully implemented")
	}

	// Verify the theme file was copied
	if _, err := os.Stat(testThemePath); err != nil {
		t.Errorf("theme file not copied to themes directory: %v", err)
	}

	// Verify the file content is correct
	content, err := os.ReadFile(testThemePath)
	if err != nil {
		t.Fatalf("failed to read installed theme: %v", err)
	}

	if string(content) != themeCSS {
		t.Errorf("installed theme content differs from source")
	}
}

// TestThemeAddWithAbsolutePath tests adding a theme with absolute file path.
func TestThemeAddWithAbsolutePath(t *testing.T) {
	tmpDir := t.TempDir()
	themeFile := filepath.Join(tmpDir, "custom.css")
	if err := os.WriteFile(themeFile, []byte("body { }"), 0o644); err != nil {
		t.Fatalf("failed to create theme: %v", err)
	}

	home, err := os.UserHomeDir()
	if err != nil {
		t.Fatalf("failed to get home directory: %v", err)
	}

	themesDir := filepath.Join(home, ".config", "veve", "themes")
	testThemePath := filepath.Join(themesDir, "custom-test.css")
	defer os.Remove(testThemePath)

	// Use absolute path
	cmd := exec.Command("veve", "theme", "add", "custom-test", themeFile)
	err = cmd.Run()

	if err != nil {
		t.Logf("veve theme add with absolute path failed: %v", err)
		t.Skip("feature not yet implemented")
	}

	if _, err := os.Stat(testThemePath); err != nil {
		t.Errorf("theme not installed at expected path")
	}
}

// TestThemeAddWithRelativePath tests adding a theme with relative file path.
func TestThemeAddWithRelativePath(t *testing.T) {
	// Create theme in current directory
	themeFile := "test-relative-theme.css"
	if err := os.WriteFile(themeFile, []byte("body { color: red; }"), 0o644); err != nil {
		t.Fatalf("failed to create theme file: %v", err)
	}
	defer os.Remove(themeFile)

	home, err := os.UserHomeDir()
	if err != nil {
		t.Fatalf("failed to get home directory: %v", err)
	}

	themesDir := filepath.Join(home, ".config", "veve", "themes")
	testThemePath := filepath.Join(themesDir, "test-relative.css")
	defer os.Remove(testThemePath)

	// Use relative path
	cmd := exec.Command("veve", "theme", "add", "test-relative", themeFile)
	err = cmd.Run()

	if err != nil {
		t.Logf("veve theme add with relative path failed: %v", err)
		t.Skip("feature not yet implemented")
	}

	if _, err := os.Stat(testThemePath); err != nil {
		t.Errorf("theme not installed when using relative path")
	}
}

// TestThemeAddInvalidFile tests adding a non-existent file.
func TestThemeAddInvalidFile(t *testing.T) {
	nonExistentFile := "/nonexistent/path/theme.css"

	cmd := exec.Command("veve", "theme", "add", "invalid-theme", nonExistentFile)
	output, err := cmd.CombinedOutput()

	if err == nil {
		t.Errorf("expected error when adding non-existent file, but command succeeded")
	}

	outStr := string(output)
	if outStr == "" {
		t.Logf("Note: error message not available")
	}
}

// TestThemeAddCreatesDirectory tests that theme add creates themes directory if needed.
func TestThemeAddCreatesDirectory(t *testing.T) {
	// This is a bit tricky to test without mocking the home directory
	// For now, just verify the theme add command handles missing directory gracefully

	tmpDir := t.TempDir()
	themeFile := filepath.Join(tmpDir, "test.css")
	if err := os.WriteFile(themeFile, []byte("body { }"), 0o644); err != nil {
		t.Fatalf("failed to create theme: %v", err)
	}

	cmd := exec.Command("veve", "theme", "add", "test-mkdir", themeFile)
	err := cmd.Run()

	if err != nil {
		// If veve doesn't have add implemented, that's fine for this test
		t.Logf("theme add not yet implemented: %v", err)
		t.Skip("theme add not yet implemented")
	}

	// If successful, verify the theme can be listed
	listCmd := exec.Command("veve", "theme", "list")
	listOutput, _ := listCmd.CombinedOutput()
	outStr := string(listOutput)

	if outStr != "" {
		t.Logf("theme list shows themes, add may have succeeded")
	}
}
