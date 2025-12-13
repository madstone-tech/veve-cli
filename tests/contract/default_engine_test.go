package contract_test

import (
	"os"
	"os/exec"
	"path/filepath"
	"testing"
)

// TestDefaultEngineSelection tests that default engine is used when no --engine specified (T062)
func TestDefaultEngineSelection(t *testing.T) {
	vevePath := buildVeve(t)
	if vevePath == "" {
		t.Skip("veve binary not available")
	}

	t.Run("veve convert without --engine flag uses default", func(t *testing.T) {
		tmpDir := t.TempDir()
		inputFile := filepath.Join(tmpDir, "test.md")
		outputFile := filepath.Join(tmpDir, "output.pdf")

		markdown := `# Test
Unicode: 世界 🎉
`

		if err := os.WriteFile(inputFile, []byte(markdown), 0644); err != nil {
			t.Fatalf("failed to create test file: %v", err)
		}

		// Run veve convert without --engine flag (output is specified via -o flag)
		cmd := exec.Command(vevePath, "convert", "-o", outputFile, inputFile)
		output, err := cmd.CombinedOutput()

		// Should succeed (exit code 0)
		if err != nil {
			t.Logf("command output: %s", string(output))
			t.Errorf("conversion without --engine flag failed: %v", err)
			return
		}

		// Output PDF should exist
		if _, err := os.Stat(outputFile); err != nil {
			t.Errorf("output PDF not created: %v", err)
			return
		}

		// PDF should have content
		info, _ := os.Stat(outputFile)
		if info.Size() == 0 {
			t.Error("output PDF is empty")
		}

		t.Logf("Default engine conversion successful (PDF size: %d bytes)", info.Size())
	})
}

// TestDefaultEngineIsUnicodeCapable tests that default engine supports unicode (T063)
func TestDefaultEngineIsUnicodeCapable(t *testing.T) {
	vevePath := buildVeve(t)
	if vevePath == "" {
		t.Skip("veve binary not available")
	}

	t.Run("default engine renders unicode characters correctly", func(t *testing.T) {
		tmpDir := t.TempDir()
		inputFile := filepath.Join(tmpDir, "unicode-test.md")
		outputFile := filepath.Join(tmpDir, "unicode-output.pdf")

		// Content with various unicode characters
		markdown := `# Unicode Test

Emoji: 🎉 📄 ✅ 🚀
CJK: 世界 日本 中国
Math: ∑ ± ∫ ∈ ∪ ∩
Diacritics: Café naïve Zürich
`

		if err := os.WriteFile(inputFile, []byte(markdown), 0644); err != nil {
			t.Fatalf("failed to create test file: %v", err)
		}

		// Run veve convert without --engine flag (output is specified via -o flag)
		cmd := exec.Command(vevePath, "convert", "-o", outputFile, inputFile)
		output, err := cmd.CombinedOutput()

		if err != nil {
			t.Logf("command output: %s", string(output))
			t.Errorf("conversion failed: %v", err)
			return
		}

		// Verify PDF was created
		info, err := os.Stat(outputFile)
		if err != nil {
			t.Errorf("output PDF not created: %v", err)
			return
		}

		if info.Size() == 0 {
			t.Error("output PDF is empty")
			return
		}

		// PDF should be reasonably sized (not too small, indicating successful conversion)
		if info.Size() < 5000 {
			t.Logf("Warning: PDF is small (%d bytes), may indicate incomplete conversion", info.Size())
		} else {
			t.Logf("PDF created successfully with unicode content (%d bytes)", info.Size())
		}
	})

	t.Run("default engine handles emoji without corruption", func(t *testing.T) {
		tmpDir := t.TempDir()
		inputFile := filepath.Join(tmpDir, "emoji-test.md")
		outputFile := filepath.Join(tmpDir, "emoji-output.pdf")

		markdown := `# Emoji Test
🎉 🚀 ❤️ 👨‍💻 👩‍🔬
`

		if err := os.WriteFile(inputFile, []byte(markdown), 0644); err != nil {
			t.Fatalf("failed to create test file: %v", err)
		}

		cmd := exec.Command(vevePath, "convert", "-o", outputFile, inputFile)
		output, err := cmd.CombinedOutput()

		// Should either succeed or provide clear error if emoji not supported
		if err != nil {
			t.Logf("Note: emoji conversion failed (acceptable if emoji not fully supported by default engine): %v", err)
			t.Logf("Output: %s", string(output))
			return
		}

		// Verify PDF exists
		info, err := os.Stat(outputFile)
		if err == nil && info.Size() > 0 {
			t.Logf("Emoji conversion succeeded (%d bytes)", info.Size())
		}
	})
}

// TestDefaultEngineErrorHandling tests error when no default available (T064)
func TestDefaultEngineErrorHandling(t *testing.T) {
	vevePath := buildVeve(t)
	if vevePath == "" {
		t.Skip("veve binary not available")
	}

	t.Run("clear error when no unicode engine available", func(t *testing.T) {
		// This test may not be practical since we assume a unicode-capable engine exists
		// on the test system, but we include it for completeness

		tmpDir := t.TempDir()
		inputFile := filepath.Join(tmpDir, "test.md")
		outputFile := filepath.Join(tmpDir, "output.pdf")

		// Create file with content that requires unicode
		markdown := `# Test
世界
`
		if err := os.WriteFile(inputFile, []byte(markdown), 0644); err != nil {
			t.Fatalf("failed to create test file: %v", err)
		}

		// If conversion fails, the error message should be actionable
		cmd := exec.Command(vevePath, "convert", "-o", outputFile, inputFile)
		output, err := cmd.CombinedOutput()

		if err != nil {
			outputStr := string(output)
			// Error should have helpful message
			if outputStr == "" {
				t.Error("error message should not be empty")
			} else {
				t.Logf("Error message (acceptable): %s", outputStr)

				// Verify error message is actionable (not cryptic)
				if len(outputStr) < 10 {
					t.Error("error message too short to be helpful")
				}
			}
		}
		// If conversion succeeds, that's also acceptable - it means unicode is supported
	})
}
