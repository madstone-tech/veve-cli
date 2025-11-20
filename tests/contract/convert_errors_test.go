package contract

import (
	"os"
	"os/exec"
	"path/filepath"
	"testing"
)

// TestMissingInputFile tests that veve exits with error code 1 when input file doesn't exist.
func TestMissingInputFile(t *testing.T) {
	// Use a path that doesn't exist
	nonexistentFile := "/tmp/veve_test_nonexistent_file_" + randomString(10) + ".md"

	cmd := exec.Command("veve", nonexistentFile)
	err := cmd.Run()

	// Should have exited with error
	if err == nil {
		t.Fatal("expected veve to fail with nonexistent input file, but it succeeded")
	}

	// Check exit code is 1 (error)
	if exitErr, ok := err.(*exec.ExitError); ok {
		if exitErr.ExitCode() != 1 {
			t.Fatalf("expected exit code 1, got %d", exitErr.ExitCode())
		}
	} else {
		t.Fatalf("expected ExitError, got %T", err)
	}
}

// TestInvalidMarkdown tests handling of invalid markdown syntax.
// Note: Pandoc is quite forgiving, so this test may not fail; it's more of a sanity check.
func TestInvalidMarkdown(t *testing.T) {
	tmpDir := t.TempDir()

	// Create markdown with potential issues
	testMDContent := `# Test

This is still valid markdown, even with odd formatting

[Broken link](
`
	testMDPath := filepath.Join(tmpDir, "invalid.md")
	if err := os.WriteFile(testMDPath, []byte(testMDContent), 0o644); err != nil {
		t.Fatalf("failed to create test markdown file: %v", err)
	}

	outputPath := filepath.Join(tmpDir, "invalid.pdf")

	// Run veve convert
	cmd := exec.Command("veve", testMDPath, "-o", outputPath)
	output, err := cmd.CombinedOutput()

	// Pandoc is forgiving, so this might succeed or fail depending on content
	// The important thing is that veve handles it gracefully
	if err != nil {
		t.Logf("veve failed (this may be expected): %v\nOutput: %s", err, string(output))
	}

	// If it succeeded, verify the PDF was created
	if err == nil {
		if _, err := os.Stat(outputPath); err != nil {
			t.Fatalf("output PDF not created: %v", err)
		}
	}
}

// TestPermissionDenied tests handling when output directory is not writable.
func TestPermissionDenied(t *testing.T) {
	// Skip this test on systems where we can't control permissions
	if os.Geteuid() == 0 {
		t.Skip("skipping permission test when running as root")
	}

	tmpDir := t.TempDir()

	testMDContent := `# Test

Content.
`
	testMDPath := filepath.Join(tmpDir, "test.md")
	if err := os.WriteFile(testMDPath, []byte(testMDContent), 0o644); err != nil {
		t.Fatalf("failed to create test markdown file: %v", err)
	}

	// Create a directory with no write permissions
	restrictedDir := filepath.Join(tmpDir, "restricted")
	if err := os.Mkdir(restrictedDir, 0o555); err != nil {
		t.Fatalf("failed to create restricted directory: %v", err)
	}
	defer os.Chmod(restrictedDir, 0o755) // Restore permissions for cleanup

	outputPath := filepath.Join(restrictedDir, "output.pdf")

	// Run veve convert to restricted directory
	cmd := exec.Command("veve", testMDPath, "-o", outputPath)
	err := cmd.Run()

	// Should have exited with error
	if err == nil {
		t.Fatal("expected veve to fail when writing to restricted directory")
	}

	// Check exit code is 1 (error)
	if exitErr, ok := err.(*exec.ExitError); ok {
		if exitErr.ExitCode() != 1 {
			t.Fatalf("expected exit code 1, got %d", exitErr.ExitCode())
		}
	}
}

// Helper function to generate random string for unique file names
func randomString(length int) string {
	const charset = "abcdefghijklmnopqrstuvwxyzABCDEFGHIJKLMNOPQRSTUVWXYZ0123456789"
	b := make([]byte, length)
	for i := range b {
		b[i] = charset[i%len(charset)]
	}
	return string(b)
}
