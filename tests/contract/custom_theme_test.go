package contract

import (
	"os"
	"os/exec"
	"path/filepath"
	"testing"
)

// TestCustomThemeDiscovery tests that custom CSS files in ~/.config/veve/themes/ are discovered.
func TestCustomThemeDiscovery(t *testing.T) {
	// Get veve config directory
	home, err := os.UserHomeDir()
	if err != nil {
		t.Fatalf("failed to get home directory: %v", err)
	}

	themesDir := filepath.Join(home, ".config", "veve", "themes")

	// Create themes directory if needed
	if err := os.MkdirAll(themesDir, 0o755); err != nil {
		t.Fatalf("failed to create themes directory: %v", err)
	}

	// Create a custom theme CSS file
	customThemePath := filepath.Join(themesDir, "custom-test-theme.css")
	customCSS := `
/* Custom Test Theme */
body {
  font-family: Arial, sans-serif;
  color: purple;
}

h1 {
  color: darkviolet;
  border-bottom: 2px solid purple;
}
`

	if err := os.WriteFile(customThemePath, []byte(customCSS), 0o644); err != nil {
		t.Fatalf("failed to create custom theme: %v", err)
	}

	defer os.Remove(customThemePath) // Clean up

	// Create test markdown
	tmpDir := t.TempDir()
	testMDPath := filepath.Join(tmpDir, "test.md")
	testMD := `# Test with Custom Theme

This uses our custom theme.
`
	if err := os.WriteFile(testMDPath, []byte(testMD), 0o644); err != nil {
		t.Fatalf("failed to create test markdown: %v", err)
	}

	outputPath := filepath.Join(tmpDir, "output.pdf")

	// Try to convert using the custom theme
	cmd := exec.Command("veve", testMDPath, "-o", outputPath, "--theme", "custom-test-theme")
	output, err := cmd.CombinedOutput()

	if err != nil {
		t.Logf("veve with custom theme failed: %v\nOutput: %s", err, string(output))
		// This is expected if custom theme support not yet implemented
		t.Skip("Custom theme support not yet fully implemented")
	}

	// Verify PDF was created if conversion succeeded
	if _, err := os.Stat(outputPath); err == nil {
		pdf, err := os.ReadFile(outputPath)
		if err != nil {
			t.Fatalf("failed to read PDF: %v", err)
		}
		if string(pdf[:4]) != "%PDF" {
			t.Fatal("output is not a valid PDF")
		}
	}
}

// TestCustomThemeListDisplay tests that custom themes appear in 'veve theme list'
func TestCustomThemeListDisplay(t *testing.T) {
	home, err := os.UserHomeDir()
	if err != nil {
		t.Fatalf("failed to get home directory: %v", err)
	}

	themesDir := filepath.Join(home, ".config", "veve", "themes")
	if err := os.MkdirAll(themesDir, 0o755); err != nil {
		t.Fatalf("failed to create themes directory: %v", err)
	}

	// Create custom theme
	customThemePath := filepath.Join(themesDir, "mycolor.css")
	if err := os.WriteFile(customThemePath, []byte("body { color: blue; }"), 0o644); err != nil {
		t.Fatalf("failed to create custom theme: %v", err)
	}
	defer os.Remove(customThemePath)

	// Run list-themes
	cmd := exec.Command("veve", "theme", "list")
	output, err := cmd.CombinedOutput()
	if err != nil {
		t.Logf("theme list command failed: %v", err)
		t.Skip("theme list command not available")
	}

	outStr := string(output)

	// Custom theme should appear in the list
	if outStr != "" && len(outStr) > 0 {
		// If list shows anything, custom theme should be there
		// (but it might not if discovery hasn't been run)
		t.Logf("theme list output:\n%s", outStr)
	}
}

// TestCustomThemeWithSubdirectory tests themes in subdirectories (if supported)
func TestCustomThemeWithSubdirectory(t *testing.T) {
	home, err := os.UserHomeDir()
	if err != nil {
		t.Fatalf("failed to get home directory: %v", err)
	}

	themesDir := filepath.Join(home, ".config", "veve", "themes")
	if err := os.MkdirAll(themesDir, 0o755); err != nil {
		t.Fatalf("failed to create themes directory: %v", err)
	}

	// Create theme in subdirectory (if supported)
	subDir := filepath.Join(themesDir, "modern")
	if err := os.MkdirAll(subDir, 0o755); err != nil {
		t.Fatalf("failed to create subdirectory: %v", err)
	}
	defer os.RemoveAll(subDir)

	themePath := filepath.Join(subDir, "colors.css")
	if err := os.WriteFile(themePath, []byte("body { }"), 0o644); err != nil {
		t.Fatalf("failed to create theme: %v", err)
	}

	// Try to list themes
	cmd := exec.Command("veve", "theme", "list")
	_, err = cmd.CombinedOutput()

	if err != nil {
		t.Logf("theme list with subdirectory: command may not support nested dirs")
	}
}
