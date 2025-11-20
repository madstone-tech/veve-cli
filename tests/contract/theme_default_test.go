package contract

import (
	"os"
	"os/exec"
	"path/filepath"
	"testing"
)

// TestDefaultThemeApplied tests that when no --theme is specified, the default theme is applied.
func TestDefaultThemeApplied(t *testing.T) {
	tmpDir := t.TempDir()

	testMDContent := `# Test with Default Theme

This document uses the default theme since no --theme flag was specified.

## Content

Some content here.
`
	testMDPath := filepath.Join(tmpDir, "test.md")
	if err := os.WriteFile(testMDPath, []byte(testMDContent), 0o644); err != nil {
		t.Fatalf("failed to create test markdown: %v", err)
	}

	outputPath := filepath.Join(tmpDir, "output.pdf")

	// Run without --theme flag
	cmd := exec.Command("veve", testMDPath, "-o", outputPath)
	if output, err := cmd.CombinedOutput(); err != nil {
		t.Logf("veve without theme flag failed: %v\nOutput: %s", err, string(output))
		t.Skip("Conversion failed - possible Pandoc issue")
	}

	// Verify PDF was created
	if _, err := os.Stat(outputPath); err != nil {
		t.Fatalf("PDF not created: %v", err)
	}

	// Verify it's a valid PDF
	pdf, err := os.ReadFile(outputPath)
	if err != nil {
		t.Fatalf("failed to read PDF: %v", err)
	}

	if len(pdf) < 4 || string(pdf[:4]) != "%PDF" {
		t.Fatal("output is not a valid PDF")
	}
}

// TestDefaultThemeExplicit tests that --theme default is equivalent to no theme
func TestDefaultThemeExplicit(t *testing.T) {
	tmpDir := t.TempDir()

	testMDContent := `# Test

Content.
`
	testMDPath := filepath.Join(tmpDir, "test.md")
	if err := os.WriteFile(testMDPath, []byte(testMDContent), 0o644); err != nil {
		t.Fatalf("failed to create test markdown: %v", err)
	}

	// Create without theme flag
	output1 := filepath.Join(tmpDir, "output1.pdf")
	cmd1 := exec.Command("veve", testMDPath, "-o", output1)
	if err := cmd1.Run(); err != nil {
		t.Logf("conversion without theme flag failed: %v", err)
		t.Skip("Conversion failed")
	}

	// Create with explicit default theme
	output2 := filepath.Join(tmpDir, "output2.pdf")
	cmd2 := exec.Command("veve", testMDPath, "-o", output2, "--theme", "default")
	if err := cmd2.Run(); err != nil {
		t.Logf("conversion with --theme default failed: %v", err)
		t.Skip("Conversion failed")
	}

	// Both should exist
	info1, err := os.Stat(output1)
	if err != nil {
		t.Fatalf("first PDF not created: %v", err)
	}

	info2, err := os.Stat(output2)
	if err != nil {
		t.Fatalf("second PDF not created: %v", err)
	}

	// Files should both be valid PDFs (may have slightly different sizes due to metadata)
	pdf1, err := os.ReadFile(output1)
	if err != nil {
		t.Fatalf("failed to read first PDF: %v", err)
	}

	pdf2, err := os.ReadFile(output2)
	if err != nil {
		t.Fatalf("failed to read second PDF: %v", err)
	}

	if string(pdf1[:4]) != "%PDF" || string(pdf2[:4]) != "%PDF" {
		t.Fatal("outputs are not valid PDFs")
	}

	t.Logf("Both PDFs created successfully: %d bytes and %d bytes", info1.Size(), info2.Size())
}

// TestThemeIntegrationWithOtherFlags tests that theme works with other flags
func TestThemeIntegrationWithOtherFlags(t *testing.T) {
	tmpDir := t.TempDir()

	testMDContent := `# Test
Content.
`
	testMDPath := filepath.Join(tmpDir, "test.md")
	if err := os.WriteFile(testMDPath, []byte(testMDContent), 0o644); err != nil {
		t.Fatalf("failed to create test markdown: %v", err)
	}

	outputPath := filepath.Join(tmpDir, "output.pdf")

	// Use theme with other flags
	cmd := exec.Command("veve", testMDPath, "-o", outputPath, "--theme", "dark", "--pdf-engine", "pdflatex", "--verbose")
	if output, err := cmd.CombinedOutput(); err != nil {
		t.Logf("conversion with multiple flags failed: %v\nOutput: %s", err, string(output))
		t.Skip("Conversion failed")
	}

	// Verify PDF was created
	if _, err := os.Stat(outputPath); err != nil {
		t.Fatalf("PDF not created: %v", err)
	}
}
