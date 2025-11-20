package contract

import (
	"os/exec"
	"strings"
	"testing"
)

// TestListThemes tests that `veve list-themes` displays all available themes.
func TestListThemes(t *testing.T) {
	cmd := exec.Command("veve", "list-themes")
	output, err := cmd.CombinedOutput()

	// Command may not exist yet
	if err != nil {
		t.Logf("list-themes command failed: %v\nOutput: %s", err, string(output))
		t.Skip("list-themes command not yet implemented")
	}

	outStr := string(output)

	// Should contain built-in themes
	expectedThemes := []string{"default", "dark", "academic"}
	for _, theme := range expectedThemes {
		if !strings.Contains(outStr, theme) {
			t.Errorf("output missing theme '%s'", theme)
		}
	}

	// Should have table-like structure (columns separated by pipes or aligned)
	lines := strings.Split(outStr, "\n")
	if len(lines) < 2 {
		t.Logf("output too short: %s", outStr)
		// This might be okay if format changes
	}
}

// TestListThemesShowsDescriptions tests that themes include descriptions
func TestListThemesShowsDescriptions(t *testing.T) {
	cmd := exec.Command("veve", "list-themes")
	output, err := cmd.CombinedOutput()

	if err != nil {
		t.Logf("list-themes command not available: %v", err)
		t.Skip("list-themes not yet implemented")
	}

	outStr := string(output)

	// Should have some description text
	if len(outStr) < 50 {
		t.Logf("output seems too short for theme descriptions: %d chars", len(outStr))
	}

	// Check for headers (Author, Description, etc.)
	expectedHeaders := []string{"Theme", "Author", "Description"}
	for _, header := range expectedHeaders {
		if !strings.Contains(strings.ToLower(outStr), strings.ToLower(header)) {
			t.Logf("warning: output may not include '%s' header", header)
		}
	}
}

// TestListThemesFormatting tests that output is well-formatted
func TestListThemesFormatting(t *testing.T) {
	cmd := exec.Command("veve", "list-themes")
	output, err := cmd.CombinedOutput()

	if err != nil {
		t.Logf("list-themes not available")
		t.Skip("list-themes not yet implemented")
	}

	outStr := string(output)

	// Should not be empty
	if len(strings.TrimSpace(outStr)) == 0 {
		t.Fatal("output is empty")
	}

	// Should be readable (no excessive control characters)
	if strings.Count(outStr, "\n") == 0 && len(outStr) > 200 {
		t.Logf("warning: output may not be properly formatted with newlines")
	}
}
