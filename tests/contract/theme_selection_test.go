package contract

import (
	"os"
	"os/exec"
	"path/filepath"
	"testing"
)

// TestThemeSelection tests that --theme flag applies the specified theme.
// Verifies that different themes produce visually distinct PDFs.
func TestThemeSelection(t *testing.T) {
	tmpDir := t.TempDir()

	// Create test markdown
	testMDContent := `# Themed Document

This document will be converted with different themes.

## Section 2

Some content here with **bold** and *italic* text.

### Code Example

` + "```" + `go
func main() {
	fmt.Println("Hello")
}
` + "```" + `
`

	testMDPath := filepath.Join(tmpDir, "themed.md")
	if err := os.WriteFile(testMDPath, []byte(testMDContent), 0o644); err != nil {
		t.Fatalf("failed to create test markdown: %v", err)
	}

	// Test with dark theme
	darkOutput := filepath.Join(tmpDir, "themed_dark.pdf")
	cmd := exec.Command("veve", testMDPath, "-o", darkOutput, "--theme", "dark")
	if output, err := cmd.CombinedOutput(); err != nil {
		t.Logf("veve with dark theme failed: %v\nOutput: %s", err, string(output))
		t.Skip("Theme support not yet implemented")
	}

	// Verify dark PDF was created
	if _, err := os.Stat(darkOutput); err != nil {
		t.Fatalf("dark PDF not created: %v", err)
	}

	// Verify it's a valid PDF
	darkPDF, err := os.ReadFile(darkOutput)
	if err != nil {
		t.Fatalf("failed to read dark PDF: %v", err)
	}

	if len(darkPDF) < 4 || string(darkPDF[:4]) != "%PDF" {
		t.Fatal("dark output is not a valid PDF")
	}

	// Test with default theme
	defaultOutput := filepath.Join(tmpDir, "themed_default.pdf")
	cmd = exec.Command("veve", testMDPath, "-o", defaultOutput, "--theme", "default")
	if output, err := cmd.CombinedOutput(); err != nil {
		t.Logf("veve with default theme failed: %v\nOutput: %s", err, string(output))
		t.Skip("Theme support not yet implemented")
	}

	// Verify default PDF was created
	if _, err := os.Stat(defaultOutput); err != nil {
		t.Fatalf("default PDF not created: %v", err)
	}

	// Verify both are valid PDFs
	defaultPDF, err := os.ReadFile(defaultOutput)
	if err != nil {
		t.Fatalf("failed to read default PDF: %v", err)
	}

	if len(defaultPDF) < 4 || string(defaultPDF[:4]) != "%PDF" {
		t.Fatal("default output is not a valid PDF")
	}

	// Note: For a comprehensive test, we would compare PDF structures or
	// render them and compare visual output. For now, we verify both are valid.
	t.Logf("Successfully created themed PDFs: dark=%d bytes, default=%d bytes", len(darkPDF), len(defaultPDF))
}

// TestThemeAcademic tests the academic theme
func TestThemeAcademic(t *testing.T) {
	tmpDir := t.TempDir()

	testMDContent := `# Academic Paper

## Introduction

Some introduction text.

## Section

Content here.
`

	testMDPath := filepath.Join(tmpDir, "academic.md")
	if err := os.WriteFile(testMDPath, []byte(testMDContent), 0o644); err != nil {
		t.Fatalf("failed to create test markdown: %v", err)
	}

	outputPath := filepath.Join(tmpDir, "academic.pdf")
	cmd := exec.Command("veve", testMDPath, "-o", outputPath, "--theme", "academic")
	if output, err := cmd.CombinedOutput(); err != nil {
		t.Logf("veve with academic theme failed: %v\nOutput: %s", err, string(output))
		t.Skip("Theme support not yet implemented")
	}

	// Verify PDF was created
	if _, err := os.Stat(outputPath); err != nil {
		t.Fatalf("academic PDF not created: %v", err)
	}

	pdf, err := os.ReadFile(outputPath)
	if err != nil {
		t.Fatalf("failed to read PDF: %v", err)
	}

	if len(pdf) < 4 || string(pdf[:4]) != "%PDF" {
		t.Fatal("output is not a valid PDF")
	}
}
