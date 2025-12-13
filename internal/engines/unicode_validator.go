// Package engines provides PDF engine detection, validation, and selection logic.
package engines

import (
	"context"
	"fmt"
	"os"
	"os/exec"
	"path/filepath"
	"time"
)

// ValidateUnicodeSupport tests if an engine can handle unicode/emoji content
// Uses test-based detection by attempting conversion of sample unicode document
// Returns TestResult with success/failure status
func ValidateUnicodeSupport(engine PDFEngine) *TestResult {
	result := &TestResult{
		Success: false,
	}

	// Create temporary directory for test
	tmpDir, err := os.MkdirTemp("", "veve-unicode-test-*")
	if err != nil {
		result.ErrorMessage = fmt.Sprintf("could not create temp directory: %v", err)
		return result
	}
	defer os.RemoveAll(tmpDir)

	// Create test markdown file
	testMDFile := filepath.Join(tmpDir, "unicode-test.md")
	testContent := getUnicodeTestContent()

	if err := os.WriteFile(testMDFile, []byte(testContent), 0644); err != nil {
		result.ErrorMessage = fmt.Sprintf("could not write test file: %v", err)
		return result
	}

	// Create output PDF path
	testPDFFile := filepath.Join(tmpDir, "output.pdf")

	// Execute conversion with timeout
	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()

	startTime := time.Now()

	// Use pandoc to convert with specified engine
	cmd := exec.CommandContext(ctx, "pandoc",
		"--from", "markdown",
		"--to", "pdf",
		"--pdf-engine", engine.Name,
		"--output", testPDFFile,
		testMDFile,
	)

	// Run conversion and capture output
	output, err := cmd.CombinedOutput()
	duration := time.Since(startTime)
	result.Duration = duration

	// Check if conversion succeeded
	if err != nil {
		result.ExitCode = getExitCode(err)
		result.Stderr = string(output)

		// Check for unicode-specific errors
		if isUnicodeError(string(output)) {
			result.ErrorMessage = fmt.Sprintf(
				"engine '%s' does not support unicode: %s",
				engine.Name, string(output),
			)
		} else {
			result.ErrorMessage = fmt.Sprintf(
				"engine '%s' failed test: %v",
				engine.Name, err,
			)
		}
		return result
	}

	// Verify output PDF was created
	if _, err := os.Stat(testPDFFile); err != nil {
		result.ErrorMessage = fmt.Sprintf("PDF output not created: %v", err)
		return result
	}

	// PDF was successfully created - test passed
	result.Success = true
	result.ExitCode = 0
	result.Stderr = string(output)

	return result
}

// isUnicodeError checks if error output indicates unicode-related failure
func isUnicodeError(output string) bool {
	// Common unicode error patterns
	unicodeErrors := []string{
		"Unicode character",
		"not set up for use with",
		"undefined control sequence",
		"missing character",
		"glyph",
		"encoding",
	}

	for _, pattern := range unicodeErrors {
		if isSubstring(output, pattern) {
			return true
		}
	}

	return false
}

// isSubstring checks if needle exists in haystack (case-insensitive)
func isSubstring(haystack, needle string) bool {
	// Simple substring check - can be enhanced for case-insensitive
	return len(haystack) > 0 && len(needle) > 0 &&
		(indexOf(haystack, needle) >= 0)
}

// indexOf finds index of substring (simple implementation)
func indexOf(s, substr string) int {
	for i := 0; i <= len(s)-len(substr); i++ {
		if s[i:i+len(substr)] == substr {
			return i
		}
	}
	return -1
}

// getExitCode extracts exit code from error
func getExitCode(err error) int {
	// Try to extract exit code from ExitError
	if ee, ok := err.(*exec.ExitError); ok {
		return ee.ExitCode()
	}
	return 1
}

// getUnicodeTestContent returns minimal but diverse unicode test content
func getUnicodeTestContent() string {
	return `# Unicode Test

Emoji: 🎉 📄 ✅ 🚀
CJK: 世界 日本 中国
Math: ∑ ± ∫ ∈
Diacritics: Café naïve Zürich
ZWJ: 👨‍💻 👩‍🔬

End test.
`
}

// ValidateEngineInstalled checks if engine binary exists and is executable
func ValidateEngineInstalled(engine PDFEngine) error {
	if !engine.IsInstalled {
		return fmt.Errorf("engine '%s' is not installed", engine.Name)
	}

	// Double-check that engine is actually in PATH
	_, err := exec.LookPath(engine.Name)
	return err
}
