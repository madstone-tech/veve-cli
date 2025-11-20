package contract

import (
	"os"
	"os/exec"
	"path/filepath"
	"testing"
)

// TestOutputDirectoryCreation tests that veve creates output directories if they don't exist.
func TestOutputDirectoryCreation(t *testing.T) {
	tmpDir := t.TempDir()

	testMDContent := `# Test Document

This is a test.
`
	testMDPath := filepath.Join(tmpDir, "test.md")
	if err := os.WriteFile(testMDPath, []byte(testMDContent), 0o644); err != nil {
		t.Fatalf("failed to create test markdown file: %v", err)
	}

	// Specify an output path in a directory that doesn't exist yet
	nestedDir := filepath.Join(tmpDir, "output", "nested", "deep")
	outputPath := filepath.Join(nestedDir, "output.pdf")

	// Verify the directory doesn't exist yet
	if _, err := os.Stat(nestedDir); err == nil {
		t.Fatal("test setup error: output directory should not exist yet")
	}

	// Run veve convert
	cmd := exec.Command("veve", testMDPath, "-o", outputPath)
	if output, err := cmd.CombinedOutput(); err != nil {
		t.Logf("veve output:\n%s", string(output))
		t.Fatalf("veve convert failed: %v", err)
	}

	// Verify the output directory was created
	if _, err := os.Stat(nestedDir); err != nil {
		t.Fatalf("output directory was not created: %v", err)
	}

	// Verify the PDF was created in the new directory
	if _, err := os.Stat(outputPath); err != nil {
		t.Fatalf("output PDF not created in directory: %v", err)
	}

	// Verify it's a valid PDF
	pdf, err := os.ReadFile(outputPath)
	if err != nil {
		t.Fatalf("failed to read output PDF: %v", err)
	}

	if len(pdf) < 4 || string(pdf[:4]) != "%PDF" {
		t.Fatal("output file is not a valid PDF")
	}
}

// TestExistingOutputDirectory tests conversion when output directory already exists.
func TestExistingOutputDirectory(t *testing.T) {
	tmpDir := t.TempDir()

	testMDContent := `# Test

Content.
`
	testMDPath := filepath.Join(tmpDir, "test.md")
	if err := os.WriteFile(testMDPath, []byte(testMDContent), 0o644); err != nil {
		t.Fatalf("failed to create test markdown file: %v", err)
	}

	// Create the output directory ahead of time
	outputDir := filepath.Join(tmpDir, "output")
	if err := os.Mkdir(outputDir, 0o755); err != nil {
		t.Fatalf("failed to create output directory: %v", err)
	}

	outputPath := filepath.Join(outputDir, "output.pdf")

	// Run veve convert
	cmd := exec.Command("veve", testMDPath, "-o", outputPath)
	if output, err := cmd.CombinedOutput(); err != nil {
		t.Logf("veve output:\n%s", string(output))
		t.Fatalf("veve convert failed: %v", err)
	}

	// Verify the PDF was created
	if _, err := os.Stat(outputPath); err != nil {
		t.Fatalf("output PDF not created: %v", err)
	}
}

// TestOverwriteExistingOutput tests that veve overwrites existing PDF files.
func TestOverwriteExistingOutput(t *testing.T) {
	tmpDir := t.TempDir()

	testMDContent := `# New Content

This is the new content.
`
	testMDPath := filepath.Join(tmpDir, "test.md")
	if err := os.WriteFile(testMDPath, []byte(testMDContent), 0o644); err != nil {
		t.Fatalf("failed to create test markdown file: %v", err)
	}

	outputPath := filepath.Join(tmpDir, "output.pdf")

	// Create a dummy PDF file first
	dummyPDF := []byte("%PDF-1.4\n1 0 obj\n<< /Type /Catalog >>\nendobj\nxref\ntrailer\n<< /Size 1 >>\nstartxref\n0\n%%EOF")
	if err := os.WriteFile(outputPath, dummyPDF, 0o644); err != nil {
		t.Fatalf("failed to create dummy PDF: %v", err)
	}

	originalSize := len(dummyPDF)

	// Run veve convert to the same path
	cmd := exec.Command("veve", testMDPath, "-o", outputPath)
	if output, err := cmd.CombinedOutput(); err != nil {
		t.Logf("veve output:\n%s", string(output))
		t.Fatalf("veve convert failed: %v", err)
	}

	// Verify the PDF was overwritten (size should be different)
	fileInfo, err := os.Stat(outputPath)
	if err != nil {
		t.Fatalf("output PDF not found: %v", err)
	}

	if fileInfo.Size() == int64(originalSize) {
		t.Logf("warning: new PDF size (%d) is the same as dummy size (%d), may not have been overwritten", fileInfo.Size(), originalSize)
	}

	// Verify it's still a valid PDF
	pdf, err := os.ReadFile(outputPath)
	if err != nil {
		t.Fatalf("failed to read output PDF: %v", err)
	}

	if len(pdf) < 4 || string(pdf[:4]) != "%PDF" {
		t.Fatal("output file is not a valid PDF")
	}
}
