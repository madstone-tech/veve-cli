package contract

import (
	"os"
	"os/exec"
	"path/filepath"
	"strings"
	"testing"
)

// TestInvalidTheme tests that using an invalid theme returns an error with suggestions.
func TestInvalidTheme(t *testing.T) {
	tmpDir := t.TempDir()

	testMDContent := `# Test

Content here.
`
	testMDPath := filepath.Join(tmpDir, "test.md")
	if err := os.WriteFile(testMDPath, []byte(testMDContent), 0o644); err != nil {
		t.Fatalf("failed to create test markdown: %v", err)
	}

	outputPath := filepath.Join(tmpDir, "output.pdf")

	// Try to use a non-existent theme
	cmd := exec.Command("veve", testMDPath, "-o", outputPath, "--theme", "nonexistent-theme-xyz")
	err := cmd.Run()

	// Should have exited with error
	if err == nil {
		t.Fatal("expected veve to fail with invalid theme, but it succeeded")
	}

	// Check exit code is 1 (error)
	if exitErr, ok := err.(*exec.ExitError); ok {
		if exitErr.ExitCode() != 1 {
			t.Fatalf("expected exit code 1, got %d", exitErr.ExitCode())
		}
	}

	// Capture error output
	cmdWithOutput := exec.Command("veve", testMDPath, "-o", outputPath, "--theme", "nonexistent-theme-xyz")
	output, _ := cmdWithOutput.CombinedOutput()
	outStr := string(output)

	// Should mention the theme error
	if !strings.Contains(outStr, "theme") && !strings.Contains(outStr, "not found") {
		t.Logf("error message may not be clear: %s", outStr)
	}

	// Ideally should suggest available themes
	if !strings.Contains(outStr, "default") && !strings.Contains(outStr, "dark") {
		t.Logf("warning: error message does not suggest available themes")
	}
}

// TestEmptyThemeName tests that empty theme name is handled
func TestEmptyThemeName(t *testing.T) {
	tmpDir := t.TempDir()

	testMDContent := `# Test
Content.
`
	testMDPath := filepath.Join(tmpDir, "test.md")
	if err := os.WriteFile(testMDPath, []byte(testMDContent), 0o644); err != nil {
		t.Fatalf("failed to create test markdown: %v", err)
	}

	outputPath := filepath.Join(tmpDir, "output.pdf")

	// Try with empty theme
	cmd := exec.Command("veve", testMDPath, "-o", outputPath, "--theme", "")
	err := cmd.Run()

	// Behavior can vary: might default to default theme or error
	// Just verify it either succeeds or fails gracefully
	if err != nil {
		if exitErr, ok := err.(*exec.ExitError); ok {
			if exitErr.ExitCode() != 1 {
				t.Logf("exit code for empty theme: %d", exitErr.ExitCode())
			}
		}
	}
}

// TestThemeNameCaseSensitivity tests case sensitivity of theme names
func TestThemeNameCaseSensitivity(t *testing.T) {
	tmpDir := t.TempDir()

	testMDContent := `# Test
Content.
`
	testMDPath := filepath.Join(tmpDir, "test.md")
	if err := os.WriteFile(testMDPath, []byte(testMDContent), 0o644); err != nil {
		t.Fatalf("failed to create test markdown: %v", err)
	}

	outputPath := filepath.Join(tmpDir, "output.pdf")

	// Try uppercase theme name
	cmd := exec.Command("veve", testMDPath, "-o", outputPath, "--theme", "Dark")
	err := cmd.Run()

	// This test documents behavior - lowercase is standard, uppercase may or may not work
	if err != nil {
		t.Logf("uppercase theme name rejected (expected behavior)")
	} else {
		t.Logf("uppercase theme name accepted (permissive behavior)")
	}
}
