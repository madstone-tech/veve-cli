package contract

import (
	"os"
	"os/exec"
	"path/filepath"
	"testing"
)

// TestBasicConversion tests that veve can convert a simple markdown file to PDF.
// This is a contract test: it verifies the CLI output expectations.
func TestBasicConversion(t *testing.T) {
	// Create temporary directory for test files
	tmpDir := t.TempDir()

	// Create a simple markdown test file
	testMDContent := `# Hello World

This is a test markdown file.

## Section 2

Some more content here.
`
	testMDPath := filepath.Join(tmpDir, "test.md")
	if err := os.WriteFile(testMDPath, []byte(testMDContent), 0o644); err != nil {
		t.Fatalf("failed to create test markdown file: %v", err)
	}

	// Set output path
	outputPath := filepath.Join(tmpDir, "test.pdf")

	// Run veve convert
	cmd := exec.Command("veve", testMDPath, "-o", outputPath)
	if output, err := cmd.CombinedOutput(); err != nil {
		t.Logf("veve output:\n%s", string(output))
		t.Fatalf("veve convert failed: %v", err)
	}

	// Verify output PDF was created
	if _, err := os.Stat(outputPath); err != nil {
		t.Fatalf("output PDF not created: %v", err)
	}

	// Check that the PDF file has content (not empty)
	fileInfo, err := os.Stat(outputPath)
	if err != nil {
		t.Fatalf("failed to stat output PDF: %v", err)
	}

	if fileInfo.Size() == 0 {
		t.Fatal("output PDF is empty")
	}

	// Verify it's a valid PDF by checking magic bytes
	pdf, err := os.ReadFile(outputPath)
	if err != nil {
		t.Fatalf("failed to read output PDF: %v", err)
	}

	// PDF files start with %PDF
	if len(pdf) < 4 || string(pdf[:4]) != "%PDF" {
		t.Fatal("output file is not a valid PDF (missing PDF magic bytes)")
	}
}

// TestBasicConversionWithDefaultOutput tests conversion without explicit output path.
// The output should be created in the same directory as the input with .pdf extension.
func TestBasicConversionWithDefaultOutput(t *testing.T) {
	tmpDir := t.TempDir()

	// Create a simple markdown test file
	testMDContent := `# Test

Content here.
`
	testMDPath := filepath.Join(tmpDir, "document.md")
	if err := os.WriteFile(testMDPath, []byte(testMDContent), 0o644); err != nil {
		t.Fatalf("failed to create test markdown file: %v", err)
	}

	// Run veve convert without -o flag
	cmd := exec.Command("veve", testMDPath)
	if output, err := cmd.CombinedOutput(); err != nil {
		t.Logf("veve output:\n%s", string(output))
		t.Fatalf("veve convert failed: %v", err)
	}

	// Verify default output was created (document.pdf in same directory)
	expectedOutput := filepath.Join(tmpDir, "document.pdf")
	if _, err := os.Stat(expectedOutput); err != nil {
		t.Fatalf("expected PDF not created at default location %s: %v", expectedOutput, err)
	}
}
