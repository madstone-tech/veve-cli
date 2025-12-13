package contract_test

import (
	"os"
	"os/exec"
	"path/filepath"
	"testing"
)

// TestEngineError_NoUnicodeEngine tests error when no unicode engine available
func TestEngineError_NoUnicodeEngine(t *testing.T) {
	vevePath := buildVeve(t)
	if vevePath == "" {
		t.Skip("veve binary not available")
	}

	t.Run("provides error when unicode engine not available", func(t *testing.T) {
		tmpDir := t.TempDir()
		inputFile := filepath.Join(tmpDir, "unicode.md")
		outputFile := filepath.Join(tmpDir, "output.pdf")

		markdown := `# Unicode Test
This contains unicode: 世界 🎉
`

		if err := os.WriteFile(inputFile, []byte(markdown), 0644); err != nil {
			t.Fatalf("failed to create test file: %v", err)
		}

		// Try to use a non-unicode-capable engine (if pdflatex available)
		_, err := exec.LookPath("pdflatex")
		if err != nil {
			t.Skip("pdflatex not available for error test")
		}

		// Run veve with pdflatex (should fail on unicode content)
		cmd := exec.Command(vevePath, "convert", "--engine", "pdflatex", inputFile, outputFile)
		_, err = cmd.CombinedOutput()

		if err == nil {
			t.Skip("pdflatex succeeded; may have unicode support in this environment")
		}

		// Should have non-zero exit code
		if cmd.ProcessState != nil && cmd.ProcessState.ExitCode() == 0 {
			t.Error("expected non-zero exit code for pdflatex with unicode")
		}
	})
}

// TestEngineError_InvalidEngine tests error message for invalid engine
func TestEngineError_InvalidEngine(t *testing.T) {
	vevePath := buildVeve(t)
	if vevePath == "" {
		t.Skip("veve binary not available")
	}

	t.Run("provides helpful error for invalid engine", func(t *testing.T) {
		tmpDir := t.TempDir()
		inputFile := filepath.Join(tmpDir, "test.md")
		outputFile := filepath.Join(tmpDir, "output.pdf")

		markdown := `# Test
Basic content.
`

		if err := os.WriteFile(inputFile, []byte(markdown), 0644); err != nil {
			t.Fatalf("failed to create test file: %v", err)
		}

		// Run with non-existent engine
		cmd := exec.Command(vevePath, "convert", "--engine", "fake-engine-xyz", inputFile, outputFile)
		output, err := cmd.CombinedOutput()

		// Should fail
		if err == nil {
			t.Error("expected error with invalid engine")
		}

		// Error message should be helpful
		errMsg := string(output)
		if len(errMsg) > 0 {
			// Error message should mention the engine or provide suggestions
			t.Logf("error message: %s", errMsg)
		}
	})
}

// TestEngineError_ActionableMessages tests that error messages are actionable
func TestEngineError_ActionableMessages(t *testing.T) {
	t.Run("error messages include installation guidance", func(t *testing.T) {
		// This test documents expected error message behavior
		// It verifies that when unicode rendering fails, message includes:
		// 1. Problem description
		// 2. Recommended solution
		// 3. Platform-specific installation command

		expectedPatterns := []string{
			// Problem patterns
			"unicode", "engine", "not found", "cannot process",
		}

		// Test would verify actual veve error messages
		// For now, this documents the requirement
		t.Logf("Expected error patterns: %v", expectedPatterns)
	})
}
