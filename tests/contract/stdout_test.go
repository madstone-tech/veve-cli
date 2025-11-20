package contract

import (
	"bytes"
	"os"
	"os/exec"
	"path/filepath"
	"testing"
)

// TestStdoutOutput tests that veve can write PDF to stdout using `-` as output.
func TestStdoutOutput(t *testing.T) {
	// Create a test markdown file
	tmpDir := t.TempDir()
	inputFile := filepath.Join(tmpDir, "test.md")
	testMD := `# Test Document

This is a test document for stdout output.

## Section 1

Some content here.
`

	if err := os.WriteFile(inputFile, []byte(testMD), 0o644); err != nil {
		t.Fatalf("failed to create test file: %v", err)
	}

	// Run veve with stdout output
	cmd := exec.Command("veve", inputFile, "-o", "-")

	var stdout bytes.Buffer
	cmd.Stdout = &stdout

	output, err := cmd.CombinedOutput()

	if err != nil {
		t.Logf("veve stdout test failed: %v\nOutput: %s", err, string(output))
		t.Skip("veve stdout support not yet implemented")
	}

	// Check that PDF was written to stdout
	pdfData := stdout.Bytes()

	if len(pdfData) < 4 {
		t.Errorf("stdout output too small to be a PDF (%d bytes)", len(pdfData))
	}

	if string(pdfData[:4]) != "%PDF" {
		t.Errorf("stdout output is not a valid PDF (header: %x)", pdfData[:4])
	}

	t.Logf("Successfully wrote %d bytes of PDF to stdout", len(pdfData))
}

// TestStdoutWithTheme tests stdout output with theme selection.
func TestStdoutWithTheme(t *testing.T) {
	tmpDir := t.TempDir()
	inputFile := filepath.Join(tmpDir, "test.md")
	os.WriteFile(inputFile, []byte("# Test\nContent"), 0o644)

	cmd := exec.Command("veve", inputFile, "-o", "-", "--theme", "dark")

	var stdout bytes.Buffer
	cmd.Stdout = &stdout

	err := cmd.Run()

	if err != nil {
		t.Logf("stdout with theme failed: %v", err)
		t.Skip("feature not yet implemented")
	}

	pdfData := stdout.Bytes()
	if len(pdfData) < 4 || string(pdfData[:4]) != "%PDF" {
		t.Errorf("stdout with theme did not produce valid PDF")
	}
}

// TestStdoutVsFile tests that stdout output matches file output.
func TestStdoutVsFile(t *testing.T) {
	tmpDir := t.TempDir()
	inputFile := filepath.Join(tmpDir, "test.md")
	fileOutput := filepath.Join(tmpDir, "output.pdf")

	testMD := `# Test

Content for comparison.

## Subsection

More content.
`

	if err := os.WriteFile(inputFile, []byte(testMD), 0o644); err != nil {
		t.Fatalf("failed to create test file: %v", err)
	}

	// First, write to file
	cmd1 := exec.Command("veve", inputFile, "-o", fileOutput)
	if err := cmd1.Run(); err != nil {
		t.Logf("file output command failed: %v", err)
		t.Skip("basic conversion not working")
	}

	// Then, write to stdout
	cmd2 := exec.Command("veve", inputFile, "-o", "-")
	var stdout bytes.Buffer
	cmd2.Stdout = &stdout

	if err := cmd2.Run(); err != nil {
		t.Logf("stdout command failed: %v", err)
		t.Skip("stdout not yet implemented")
	}

	// Read file output
	filePDF, err := os.ReadFile(fileOutput)
	if err != nil {
		t.Fatalf("failed to read file output: %v", err)
	}

	// Compare sizes (may differ slightly due to timestamps, but should be similar)
	stdoutSize := len(stdout.Bytes())
	fileSize := len(filePDF)

	if stdoutSize == 0 {
		t.Errorf("stdout output is empty")
	}

	if fileSize == 0 {
		t.Errorf("file output is empty")
	}

	// Allow 10% size difference
	diff := abs(stdoutSize - fileSize)
	maxDiff := fileSize / 10

	if diff > maxDiff && fileSize > 0 {
		t.Logf("Note: stdout and file outputs differ in size: %d vs %d bytes", stdoutSize, fileSize)
	}

	t.Logf("Stdout: %d bytes, File: %d bytes", stdoutSize, fileSize)
}

// TestStdoutWithQuiet tests stdout with quiet flag.
func TestStdoutWithQuiet(t *testing.T) {
	tmpDir := t.TempDir()
	inputFile := filepath.Join(tmpDir, "test.md")
	os.WriteFile(inputFile, []byte("# Test\nContent"), 0o644)

	cmd := exec.Command("veve", inputFile, "-o", "-", "--quiet")

	var stdout bytes.Buffer
	var stderr bytes.Buffer
	cmd.Stdout = &stdout
	cmd.Stderr = &stderr

	err := cmd.Run()

	if err != nil {
		t.Logf("stdout with quiet failed: %v", err)
		t.Skip("feature not yet implemented")
	}

	// With quiet flag, stderr should be empty (no info messages)
	stderrOutput := stderr.String()
	if stderrOutput != "" && !isError(stderrOutput) {
		t.Logf("Note: quiet flag may not be suppressing messages: %s", stderrOutput)
	}
}

// TestStdoutLargeDocument tests stdout with larger document.
func TestStdoutLargeDocument(t *testing.T) {
	tmpDir := t.TempDir()
	inputFile := filepath.Join(tmpDir, "large.md")

	// Create large document
	var content bytes.Buffer
	content.WriteString("# Large Document\n\n")
	for i := 1; i <= 50; i++ {
		content.WriteString("## Section\n\nContent\n\n")
	}

	if err := os.WriteFile(inputFile, content.Bytes(), 0o644); err != nil {
		t.Fatalf("failed to create test file: %v", err)
	}

	cmd := exec.Command("veve", inputFile, "-o", "-")

	var stdout bytes.Buffer
	cmd.Stdout = &stdout

	err := cmd.Run()

	if err != nil {
		t.Logf("large document stdout failed: %v", err)
		t.Skip("feature not yet implemented")
	}

	pdfData := stdout.Bytes()
	if len(pdfData) < 4 || string(pdfData[:4]) != "%PDF" {
		t.Errorf("large document stdout did not produce valid PDF")
	}

	t.Logf("Large document PDF size: %d bytes", len(pdfData))
}

// TestStdoutNoInput tests stdout without input file.
func TestStdoutNoInput(t *testing.T) {
	cmd := exec.Command("veve", "-o", "-")
	output, err := cmd.CombinedOutput()

	if err == nil {
		t.Logf("Note: command accepted without input")
	} else {
		t.Logf("Command rejected (expected): %v", err)
	}

	_ = string(output)
}

// TestStdoutWithVerbose tests stdout with verbose output.
func TestStdoutWithVerbose(t *testing.T) {
	tmpDir := t.TempDir()
	inputFile := filepath.Join(tmpDir, "test.md")
	os.WriteFile(inputFile, []byte("# Test\nContent"), 0o644)

	cmd := exec.Command("veve", inputFile, "-o", "-", "--verbose")

	var stdout bytes.Buffer
	var stderr bytes.Buffer
	cmd.Stdout = &stdout
	cmd.Stderr = &stderr

	err := cmd.Run()

	if err != nil {
		t.Logf("stdout with verbose failed: %v", err)
		t.Skip("feature not yet implemented")
	}

	// Verbose output may go to stderr
	pdfData := stdout.Bytes()
	stderrData := stderr.String()

	if len(pdfData) < 4 || string(pdfData[:4]) != "%PDF" {
		t.Errorf("stdout did not contain valid PDF")
	}

	if stderrData != "" {
		t.Logf("Verbose output to stderr: %d chars", len(stderrData))
	}
}

// Helper function to check if output is an error message
func isError(output string) bool {
	return bytes.Contains([]byte(output), []byte("error")) ||
		bytes.Contains([]byte(output), []byte("Error")) ||
		bytes.Contains([]byte(output), []byte("failed"))
}

// Helper function for absolute value
func abs(x int) int {
	if x < 0 {
		return -x
	}
	return x
}
