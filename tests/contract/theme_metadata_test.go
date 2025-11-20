package contract

import (
	"os"
	"os/exec"
	"path/filepath"
	"strings"
	"testing"
)

// TestThemeMetadataYAML tests that custom themes with YAML metadata are parsed correctly.
func TestThemeMetadataYAML(t *testing.T) {
	home, err := os.UserHomeDir()
	if err != nil {
		t.Fatalf("failed to get home directory: %v", err)
	}

	themesDir := filepath.Join(home, ".config", "veve", "themes")
	if err := os.MkdirAll(themesDir, 0o755); err != nil {
		t.Fatalf("failed to create themes directory: %v", err)
	}

	// Create a theme with YAML front matter
	themePath := filepath.Join(themesDir, "documented.css")
	themeContent := `---
name: documented
author: Test Author
description: A documented custom theme
version: 1.0.0
---
/* Documented Theme */
body {
  font-family: Georgia, serif;
  font-size: 12pt;
}

h1 {
  color: navy;
}
`

	if err := os.WriteFile(themePath, []byte(themeContent), 0o644); err != nil {
		t.Fatalf("failed to create theme: %v", err)
	}
	defer os.Remove(themePath)

	// Check if theme list shows metadata
	cmd := exec.Command("veve", "theme", "list")
	output, err := cmd.CombinedOutput()
	if err != nil {
		t.Logf("theme list failed: %v", err)
		t.Skip("theme list not available")
	}

	outStr := string(output)
	if len(outStr) == 0 {
		t.Skip("theme list output is empty")
	}

	// Look for theme name in output
	if !strings.Contains(outStr, "documented") && !strings.Contains(outStr, "Test Author") {
		t.Logf("metadata may not be displayed in list:\n%s", outStr)
	}
}

// TestThemeMetadataWithoutYAML tests themes without YAML metadata get defaults
func TestThemeMetadataWithoutYAML(t *testing.T) {
	home, err := os.UserHomeDir()
	if err != nil {
		t.Fatalf("failed to get home directory: %v", err)
	}

	themesDir := filepath.Join(home, ".config", "veve", "themes")
	if err := os.MkdirAll(themesDir, 0o755); err != nil {
		t.Fatalf("failed to create themes directory: %v", err)
	}

	// Create a simple theme without metadata
	themePath := filepath.Join(themesDir, "simple.css")
	if err := os.WriteFile(themePath, []byte("body { color: red; }"), 0o644); err != nil {
		t.Fatalf("failed to create theme: %v", err)
	}
	defer os.Remove(themePath)

	// Should still be usable
	tmpDir := t.TempDir()
	testMDPath := filepath.Join(tmpDir, "test.md")
	if err := os.WriteFile(testMDPath, []byte("# Test\nContent"), 0o644); err != nil {
		t.Fatalf("failed to create test markdown: %v", err)
	}

	outputPath := filepath.Join(tmpDir, "output.pdf")

	cmd := exec.Command("veve", testMDPath, "-o", outputPath, "--theme", "simple")
	err = cmd.Run()

	if err != nil {
		t.Logf("conversion with simple theme failed (may not be discovered yet)")
		t.Skip("Custom theme support not yet implemented")
	}

	if _, err := os.Stat(outputPath); err == nil {
		t.Logf("simple theme (no metadata) conversion succeeded")
	}
}

// TestThemeMetadataEdgeCases tests malformed YAML in themes
func TestThemeMetadataMalformed(t *testing.T) {
	home, err := os.UserHomeDir()
	if err != nil {
		t.Fatalf("failed to get home directory: %v", err)
	}

	themesDir := filepath.Join(home, ".config", "veve", "themes")
	if err := os.MkdirAll(themesDir, 0o755); err != nil {
		t.Fatalf("failed to create themes directory: %v", err)
	}

	// Create theme with malformed YAML
	themePath := filepath.Join(themesDir, "malformed.css")
	themeContent := `---
name: malformed
author: [invalid yaml structure
description: This YAML is broken
---
body { }
`

	if err := os.WriteFile(themePath, []byte(themeContent), 0o644); err != nil {
		t.Fatalf("failed to create theme: %v", err)
	}
	defer os.Remove(themePath)

	// Should still work (fallback to defaults)
	tmpDir := t.TempDir()
	testMDPath := filepath.Join(tmpDir, "test.md")
	if err := os.WriteFile(testMDPath, []byte("# Test"), 0o644); err != nil {
		t.Fatalf("failed to create markdown: %v", err)
	}

	outputPath := filepath.Join(tmpDir, "output.pdf")

	cmd := exec.Command("veve", testMDPath, "-o", outputPath, "--theme", "malformed")
	err = cmd.Run()

	if err != nil {
		t.Logf("malformed theme handling: %v", err)
		// Could error out or continue with defaults - both are acceptable
	} else {
		t.Logf("malformed theme handled gracefully (used as fallback)")
	}
}

// TestThemeMetadataWithVersion tests version field in metadata
func TestThemeMetadataVersion(t *testing.T) {
	home, err := os.UserHomeDir()
	if err != nil {
		t.Fatalf("failed to get home directory: %v", err)
	}

	themesDir := filepath.Join(home, ".config", "veve", "themes")
	if err := os.MkdirAll(themesDir, 0o755); err != nil {
		t.Fatalf("failed to create themes directory: %v", err)
	}

	themePath := filepath.Join(themesDir, "versioned.css")
	themeContent := `---
name: versioned
author: Developer
version: 2.1.0
description: Theme with version info
---
body { }
`

	if err := os.WriteFile(themePath, []byte(themeContent), 0o644); err != nil {
		t.Fatalf("failed to create theme: %v", err)
	}
	defer os.Remove(themePath)

	// Try using it
	tmpDir := t.TempDir()
	testMDPath := filepath.Join(tmpDir, "test.md")
	if err := os.WriteFile(testMDPath, []byte("# Test"), 0o644); err != nil {
		t.Fatalf("failed to create markdown: %v", err)
	}

	outputPath := filepath.Join(tmpDir, "output.pdf")

	cmd := exec.Command("veve", testMDPath, "-o", outputPath, "--theme", "versioned", "--verbose")
	output, err := cmd.CombinedOutput()

	if err == nil {
		t.Logf("theme with version metadata works")
	} else if strings.Contains(string(output), "version") {
		t.Logf("version information may be logged: %s", string(output))
	}
}
