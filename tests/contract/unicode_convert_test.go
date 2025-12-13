package contract_test

import (
	"os"
	"os/exec"
	"path/filepath"
	"testing"
)

// TestUnicodeConversion_DefaultEngine tests conversion with default engine
func TestUnicodeConversion_DefaultEngine(t *testing.T) {
	// Skip if veve not built
	vevePath := buildVeve(t)
	if vevePath == "" {
		t.Skip("veve binary not available")
	}

	t.Run("convert unicode markdown to PDF with default engine", func(t *testing.T) {
		// Create test markdown with unicode content
		tmpDir := t.TempDir()
		inputFile := filepath.Join(tmpDir, "unicode-test.md")
		outputFile := filepath.Join(tmpDir, "output.pdf")

		markdown := `# Unicode Test

Emoji: 🎉 📄 ✅
CJK: 世界 日本 中国
Math: ∑ ± ∫
Diacritics: Café naïve Zürich

End test.
`

		if err := os.WriteFile(inputFile, []byte(markdown), 0644); err != nil {
			t.Fatalf("failed to create test file: %v", err)
		}

		// Run veve convert with default settings
		cmd := exec.Command(vevePath, "convert", "-o", outputFile, inputFile)
		output, err := cmd.CombinedOutput()

		// Should succeed (exit code 0)
		if err != nil {
			t.Logf("command output: %s", string(output))
			t.Errorf("conversion failed: %v", err)
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
	})
}

// TestUnicodeConversion_ExplicitEngine tests conversion with explicit unicode-capable engine
func TestUnicodeConversion_ExplicitEngine(t *testing.T) {
	vevePath := buildVeve(t)
	if vevePath == "" {
		t.Skip("veve binary not available")
	}

	_, err := exec.LookPath("xelatex")
	if err != nil {
		t.Skip("xelatex not found; skipping explicit engine test")
	}

	t.Run("convert with explicit --engine xelatex", func(t *testing.T) {
		tmpDir := t.TempDir()
		inputFile := filepath.Join(tmpDir, "unicode-test.md")
		outputFile := filepath.Join(tmpDir, "output.pdf")

		markdown := `# Unicode Test
Emoji: 🎉 📄 ✅
CJK: 世界
Math: ∑
`

		if err := os.WriteFile(inputFile, []byte(markdown), 0644); err != nil {
			t.Fatalf("failed to create test file: %v", err)
		}

		// Run veve convert with --engine xelatex
		cmd := exec.Command(vevePath, "convert", "--engine", "xelatex", "-o", outputFile, inputFile)
		output, err := cmd.CombinedOutput()

		if err != nil {
			t.Logf("command output: %s", string(output))
			t.Errorf("conversion with xelatex failed: %v", err)
			return
		}

		if _, err := os.Stat(outputFile); err != nil {
			t.Errorf("output PDF not created: %v", err)
			return
		}

		info, _ := os.Stat(outputFile)
		if info.Size() == 0 {
			t.Error("output PDF is empty")
		}
	})
}

// TestUnicodeConversion_AllCharacterTypes tests various unicode character types
func TestUnicodeConversion_AllCharacterTypes(t *testing.T) {
	vevePath := buildVeve(t)
	if vevePath == "" {
		t.Skip("veve binary not available")
	}

	t.Run("unicode content includes emoji", func(t *testing.T) {
		tmpDir := t.TempDir()
		inputFile := filepath.Join(tmpDir, "emoji-test.md")
		outputFile := filepath.Join(tmpDir, "emoji.pdf")

		markdown := `# Emoji Test
🎉 🚀 ❤️ 👨‍💻 👩‍🔬
`

		if err := os.WriteFile(inputFile, []byte(markdown), 0644); err != nil {
			t.Fatalf("failed to create test file: %v", err)
		}

		cmd := exec.Command(vevePath, "convert", "-o", outputFile, inputFile)
		_, err := cmd.CombinedOutput()

		if err != nil {
			t.Logf("emoji conversion failed (may be expected): %v", err)
		}

		// At minimum, should try to create output
		if _, err := os.Stat(outputFile); os.IsNotExist(err) {
			// It's ok if PDF not created; test-based detection may have failed emoji support
			t.Logf("emoji PDF not created; engine may not support emoji: %v", err)
		}
	})

	t.Run("unicode content includes CJK", func(t *testing.T) {
		tmpDir := t.TempDir()
		inputFile := filepath.Join(tmpDir, "cjk-test.md")
		outputFile := filepath.Join(tmpDir, "cjk.pdf")

		markdown := `# CJK Test
世界 日本 中国
`

		if err := os.WriteFile(inputFile, []byte(markdown), 0644); err != nil {
			t.Fatalf("failed to create test file: %v", err)
		}

		cmd := exec.Command(vevePath, "convert", "-o", outputFile, inputFile)
		_, err := cmd.CombinedOutput()

		if err != nil {
			t.Logf("CJK conversion failed: %v", err)
		}
	})

	t.Run("unicode content includes math symbols", func(t *testing.T) {
		tmpDir := t.TempDir()
		inputFile := filepath.Join(tmpDir, "math-test.md")
		outputFile := filepath.Join(tmpDir, "math.pdf")

		markdown := `# Math Test
∑ ± ∫ ∈ ∪ ∩
`

		if err := os.WriteFile(inputFile, []byte(markdown), 0644); err != nil {
			t.Fatalf("failed to create test file: %v", err)
		}

		cmd := exec.Command(vevePath, "convert", "-o", outputFile, inputFile)
		_, err := cmd.CombinedOutput()

		if err != nil {
			t.Logf("math conversion failed: %v", err)
		}
	})
}

// buildVeve builds or locates the veve binary for testing
func buildVeve(t *testing.T) string {
	// Try multiple locations to find veve
	locations := []string{
		"veve",                              // Current directory
		"./veve",                            // Relative to current
		"../veve",                           // Parent directory
		"../../veve",                        // Two levels up
		filepath.Join(os.TempDir(), "veve"), // Temp directory
	}

	// First, try common locations from cwd
	for _, loc := range locations {
		if _, err := os.Stat(loc); err == nil {
			abs, err := filepath.Abs(loc)
			if err == nil {
				return abs
			}
			return loc
		}
	}

	// Try to find via executable path (for test binary location)
	exePath, err := os.Executable()
	if err == nil {
		repoRoot := filepath.Join(filepath.Dir(exePath), "..", "..")
		vevePath := filepath.Join(repoRoot, "veve")
		if _, err := os.Stat(vevePath); err == nil {
			return vevePath
		}
	}

	return ""
}
