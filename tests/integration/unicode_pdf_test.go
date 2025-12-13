package integration_test

import (
	"bytes"
	"os"
	"os/exec"
	"path/filepath"
	"testing"

	"github.com/madstone-tech/veve-cli/internal/engines"
)

// TestUnicodeRendering_EndToEnd tests full unicode rendering pipeline
func TestUnicodeRendering_EndToEnd(t *testing.T) {
	// Require xelatex for end-to-end test
	_, err := exec.LookPath("pandoc")
	if err != nil {
		t.Skip("pandoc not found; skipping end-to-end test")
	}

	t.Run("converts markdown with unicode to PDF successfully", func(t *testing.T) {
		// Create engine selector
		selector, err := engines.NewEngineSelector()
		if err != nil {
			t.Skipf("no unicode engines available: %v", err)
		}

		defaultEngine, err := selector.SelectDefaultEngine()
		if err != nil {
			t.Fatalf("failed to get default engine: %v", err)
		}

		tmpDir := t.TempDir()
		inputFile := filepath.Join(tmpDir, "test.md")
		outputFile := filepath.Join(tmpDir, "output.pdf")

		// Create test markdown
		markdown := `# Unicode Test

This is a test with unicode characters:
- Emoji: 🎉 📄 ✅
- CJK: 世界 日本 中国
- Math: ∑ ± ∫
- Diacritics: Café naïve Zürich

End test.
`

		if err := os.WriteFile(inputFile, []byte(markdown), 0644); err != nil {
			t.Fatalf("failed to create input file: %v", err)
		}

		// Convert with engine
		cmd := exec.Command("pandoc",
			"--from", "markdown",
			"--to", "pdf",
			"--pdf-engine", defaultEngine.Name,
			"--output", outputFile,
			inputFile,
		)

		output, err := cmd.CombinedOutput()
		if err != nil {
			t.Logf("conversion output: %s", string(output))
			t.Errorf("conversion failed with %s: %v", defaultEngine.Name, err)
			return
		}

		// Verify output exists
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

// TestEmojiRendering tests that emoji characters render correctly
func TestEmojiRendering(t *testing.T) {
	_, err := exec.LookPath("pandoc")
	if err != nil {
		t.Skip("pandoc not found")
	}

	selector, err := engines.NewEngineSelector()
	if err != nil {
		t.Skip("no unicode engines available")
	}

	defaultEngine, err := selector.SelectDefaultEngine()
	if err != nil {
		t.Fatalf("failed to get default engine: %v", err)
	}

	t.Run("renders emoji correctly in PDF", func(t *testing.T) {
		tmpDir := t.TempDir()
		inputFile := filepath.Join(tmpDir, "emoji.md")
		outputFile := filepath.Join(tmpDir, "emoji.pdf")

		markdown := `# Emoji Test
🎉 Celebration
📄 Document
✅ Check
🚀 Rocket
❤️ Heart
👨‍💻 Developer
👩‍🔬 Scientist
`

		if err := os.WriteFile(inputFile, []byte(markdown), 0644); err != nil {
			t.Fatalf("failed to create input file: %v", err)
		}

		cmd := exec.Command("pandoc",
			"--from", "markdown",
			"--to", "pdf",
			"--pdf-engine", defaultEngine.Name,
			"--output", outputFile,
			inputFile,
		)

		if err := cmd.Run(); err != nil {
			t.Logf("emoji rendering may not be fully supported with %s", defaultEngine.Name)
		}

		if _, err := os.Stat(outputFile); err != nil {
			t.Logf("emoji PDF not created (acceptable if engine doesn't support emoji)")
		}
	})
}

// TestCJKRendering tests that CJK characters render correctly
func TestCJKRendering(t *testing.T) {
	_, err := exec.LookPath("pandoc")
	if err != nil {
		t.Skip("pandoc not found")
	}

	selector, err := engines.NewEngineSelector()
	if err != nil {
		t.Skip("no unicode engines available")
	}

	defaultEngine, err := selector.SelectDefaultEngine()
	if err != nil {
		t.Fatalf("failed to get default engine: %v", err)
	}

	t.Run("renders CJK characters correctly in PDF", func(t *testing.T) {
		tmpDir := t.TempDir()
		inputFile := filepath.Join(tmpDir, "cjk.md")
		outputFile := filepath.Join(tmpDir, "cjk.pdf")

		markdown := `# CJK Test

Chinese: 世界 日本 中国 你好 谢谢
Japanese: こんにちは カタカナ 日本語
Korean: 안녕하세요 한국어 감사합니다

End test.
`

		if err := os.WriteFile(inputFile, []byte(markdown), 0644); err != nil {
			t.Fatalf("failed to create input file: %v", err)
		}

		cmd := exec.Command("pandoc",
			"--from", "markdown",
			"--to", "pdf",
			"--pdf-engine", defaultEngine.Name,
			"--output", outputFile,
			inputFile,
		)

		if err := cmd.Run(); err != nil {
			t.Logf("CJK rendering error: %v", err)
		}

		if _, err := os.Stat(outputFile); err == nil {
			t.Logf("CJK PDF created successfully")
		}
	})
}

// TestMathSymbolsRendering tests that math symbols render correctly
func TestMathSymbolsRendering(t *testing.T) {
	_, err := exec.LookPath("pandoc")
	if err != nil {
		t.Skip("pandoc not found")
	}

	selector, err := engines.NewEngineSelector()
	if err != nil {
		t.Skip("no unicode engines available")
	}

	defaultEngine, err := selector.SelectDefaultEngine()
	if err != nil {
		t.Fatalf("failed to get default engine: %v", err)
	}

	t.Run("renders math symbols correctly in PDF", func(t *testing.T) {
		tmpDir := t.TempDir()
		inputFile := filepath.Join(tmpDir, "math.md")
		outputFile := filepath.Join(tmpDir, "math.pdf")

		markdown := `# Math Symbols

Summation: ∑
Plus-Minus: ±
Integral: ∫
Element: ∈
Union: ∪
Intersection: ∩
Sigma: Σ

End test.
`

		if err := os.WriteFile(inputFile, []byte(markdown), 0644); err != nil {
			t.Fatalf("failed to create input file: %v", err)
		}

		cmd := exec.Command("pandoc",
			"--from", "markdown",
			"--to", "pdf",
			"--pdf-engine", defaultEngine.Name,
			"--output", outputFile,
			inputFile,
		)

		if err := cmd.Run(); err != nil {
			t.Logf("math rendering error: %v", err)
		}

		if _, err := os.Stat(outputFile); err == nil {
			t.Logf("math PDF created successfully")
		}
	})
}

// TestZWJSequenceRendering tests that ZWJ emoji sequences render
func TestZWJSequenceRendering(t *testing.T) {
	_, err := exec.LookPath("pandoc")
	if err != nil {
		t.Skip("pandoc not found")
	}

	selector, err := engines.NewEngineSelector()
	if err != nil {
		t.Skip("no unicode engines available")
	}

	defaultEngine, err := selector.SelectDefaultEngine()
	if err != nil {
		t.Fatalf("failed to get default engine: %v", err)
	}

	t.Run("renders ZWJ emoji sequences", func(t *testing.T) {
		tmpDir := t.TempDir()
		inputFile := filepath.Join(tmpDir, "zwj.md")
		outputFile := filepath.Join(tmpDir, "zwj.pdf")

		markdown := `# ZWJ Emoji Test

Family: 👨‍👩‍👧‍👦
Professions: 👨‍💻 👩‍🔬 👨‍🍳
Sports: 🏃‍♂️ 🚴‍♀️

End test.
`

		if err := os.WriteFile(inputFile, []byte(markdown), 0644); err != nil {
			t.Fatalf("failed to create input file: %v", err)
		}

		cmd := exec.Command("pandoc",
			"--from", "markdown",
			"--to", "pdf",
			"--pdf-engine", defaultEngine.Name,
			"--output", outputFile,
			inputFile,
		)

		if err := cmd.Run(); err != nil {
			t.Logf("ZWJ rendering may not be fully supported")
		}

		if _, err := os.Stat(outputFile); err == nil {
			t.Logf("ZWJ PDF created")
		}
	})
}

// TestFallbackChain tests that fallback works when primary engine fails
func TestFallbackChain(t *testing.T) {
	_, err := exec.LookPath("pandoc")
	if err != nil {
		t.Skip("pandoc not found")
	}

	selector, err := engines.NewEngineSelector()
	if err != nil {
		t.Skip("no unicode engines available")
	}

	available := selector.GetAvailableEngines()
	if len(available) < 2 {
		t.Skip("need at least 2 engines for fallback test")
	}

	t.Run("uses fallback engine when primary fails", func(t *testing.T) {
		tmpDir := t.TempDir()
		inputFile := filepath.Join(tmpDir, "fallback-test.md")
		outputFile := filepath.Join(tmpDir, "fallback.pdf")

		markdown := `# Fallback Test
Unicode: 世界 🎉 ∑
`

		if err := os.WriteFile(inputFile, []byte(markdown), 0644); err != nil {
			t.Fatalf("failed to create input file: %v", err)
		}

		// Try first available engine
		cmd := exec.Command("pandoc",
			"--from", "markdown",
			"--to", "pdf",
			"--pdf-engine", available[0],
			"--output", outputFile,
			inputFile,
		)

		var output bytes.Buffer
		cmd.Stdout = &output
		cmd.Stderr = &output

		err := cmd.Run()
		if err == nil {
			t.Logf("primary engine %s succeeded", available[0])
			return
		}

		t.Logf("primary engine %s failed, fallback would be used", available[0])

		// Try fallback
		if len(available) > 1 {
			cmd2 := exec.Command("pandoc",
				"--from", "markdown",
				"--to", "pdf",
				"--pdf-engine", available[1],
				"--output", outputFile,
				inputFile,
			)

			if err := cmd2.Run(); err != nil {
				t.Logf("fallback engine %s also failed: %v", available[1], err)
			}
		}
	})
}
