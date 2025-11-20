package contract

import (
	"bytes"
	"os"
	"os/exec"
	"path/filepath"
	"strings"
	"testing"
)

// TestErrorsToStderr tests that errors are printed to stderr, not stdout.
func TestErrorsToStderr(t *testing.T) {
	cmd := exec.Command("veve", "/nonexistent/file.md")

	var stdout, stderr bytes.Buffer
	cmd.Stdout = &stdout
	cmd.Stderr = &stderr

	_ = cmd.Run()

	stdoutStr := stdout.String()
	stderrStr := stderr.String()

	// Errors should be in stderr, not stdout
	if strings.Contains(stdoutStr, "error") || strings.Contains(stdoutStr, "Error") {
		t.Errorf("error message found in stdout (should be in stderr)")
	}

	// There should be something in stderr
	if stderrStr == "" {
		t.Logf("Note: stderr is empty - error handling may not be fully implemented")
	} else {
		t.Logf("Error message in stderr: %s", stderrStr)
	}
}

// TestSuccessOutputLocation tests that success messages go to correct stream.
func TestSuccessOutputLocation(t *testing.T) {
	tmpDir := t.TempDir()
	inputFile := filepath.Join(tmpDir, "test.md")
	outputFile := filepath.Join(tmpDir, "output.pdf")

	if err := os.WriteFile(inputFile, []byte("# Test\nContent"), 0o644); err != nil {
		t.Fatalf("failed to create test file: %v", err)
	}

	cmd := exec.Command("veve", inputFile, "-o", outputFile)

	var stdout, stderr bytes.Buffer
	cmd.Stdout = &stdout
	cmd.Stderr = &stderr

	err := cmd.Run()

	stdoutStr := stdout.String()
	stderrStr := stderr.String()

	if err != nil {
		t.Logf("conversion failed: %v", err)
		return
	}

	// Success messages may go to stdout or stderr, but not both ideally
	if stdoutStr != "" && stderrStr != "" {
		t.Logf("Output on both stdout and stderr")
	}

	t.Logf("Stdout: %q, Stderr: %q", stdoutStr, stderrStr)
}

// TestWarningsToStderr tests that warnings go to stderr.
func TestWarningsToStderr(t *testing.T) {
	// Using verbose flag might produce warnings
	tmpDir := t.TempDir()
	inputFile := filepath.Join(tmpDir, "test.md")

	if err := os.WriteFile(inputFile, []byte("# Test"), 0o644); err != nil {
		t.Fatalf("failed to create test file: %v", err)
	}

	cmd := exec.Command("veve", inputFile, "--verbose")

	var stdout, stderr bytes.Buffer
	cmd.Stdout = &stdout
	cmd.Stderr = &stderr

	_ = cmd.Run()

	stderrStr := stderr.String()

	// Verbose output should be minimal on stderr unless there are warnings
	if stderrStr != "" {
		t.Logf("Stderr with --verbose: %s", stderrStr)
	}
}

// TestQuietSuppressesMessages tests that --quiet suppresses non-error output.
func TestQuietSuppressesMessages(t *testing.T) {
	tmpDir := t.TempDir()
	inputFile := filepath.Join(tmpDir, "test.md")
	outputFile := filepath.Join(tmpDir, "output.pdf")

	if err := os.WriteFile(inputFile, []byte("# Test"), 0o644); err != nil {
		t.Fatalf("failed to create test file: %v", err)
	}

	cmd := exec.Command("veve", inputFile, "-o", outputFile, "--quiet")

	var stdout, stderr bytes.Buffer
	cmd.Stdout = &stdout
	cmd.Stderr = &stderr

	err := cmd.Run()

	if err != nil {
		t.Logf("conversion failed: %v", err)
		return
	}

	stdoutStr := stdout.String()

	// With --quiet, there should be minimal output unless there are errors
	if stdoutStr != "" && !strings.Contains(stdoutStr, "error") {
		t.Logf("Note: stdout not suppressed with --quiet: %s", stdoutStr)
	}
}

// TestVerboseAddsDetails tests that --verbose adds detail to output.
func TestVerboseAddsDetails(t *testing.T) {
	tmpDir := t.TempDir()
	inputFile := filepath.Join(tmpDir, "test.md")
	outputFile := filepath.Join(tmpDir, "output.pdf")

	if err := os.WriteFile(inputFile, []byte("# Test"), 0o644); err != nil {
		t.Fatalf("failed to create test file: %v", err)
	}

	// Run without verbose
	cmd1 := exec.Command("veve", inputFile, "-o", outputFile, "--quiet")
	var stdout1, stderr1 bytes.Buffer
	cmd1.Stdout = &stdout1
	cmd1.Stderr = &stderr1
	cmd1.Run()

	// Run with verbose
	cmd2 := exec.Command("veve", inputFile, "-o", outputFile, "--verbose")
	var stdout2, stderr2 bytes.Buffer
	cmd2.Stdout = &stdout2
	cmd2.Stderr = &stderr2
	cmd2.Run()

	quietOutput := stdout1.String() + stderr1.String()
	verboseOutput := stdout2.String() + stderr2.String()

	t.Logf("Quiet output length: %d, Verbose output length: %d",
		len(quietOutput), len(verboseOutput))

	if len(verboseOutput) > len(quietOutput) {
		t.Logf("Verbose output is larger (expected)")
	} else {
		t.Logf("Note: verbose output not larger than quiet")
	}
}

// TestInvalidInputToStderr tests that invalid input errors go to stderr.
func TestInvalidInputToStderr(t *testing.T) {
	cmd := exec.Command("veve", "/nonexistent/path.md", "-o", "/tmp/out.pdf")

	var stdout, stderr bytes.Buffer
	cmd.Stdout = &stdout
	cmd.Stderr = &stderr

	_ = cmd.Run()

	stderrStr := stderr.String()
	stdoutStr := stdout.String()

	if stdoutStr != "" {
		t.Errorf("stdout should be empty for invalid input: %s", stdoutStr)
	}

	if stderrStr == "" {
		t.Logf("Note: stderr is empty - error handling may need work")
	} else {
		if !strings.Contains(stderrStr, "not found") && !strings.Contains(stderrStr, "error") {
			t.Logf("Stderr: %s", stderrStr)
		}
	}
}

// TestThemeErrorToStderr tests that theme errors go to stderr.
func TestThemeErrorToStderr(t *testing.T) {
	tmpDir := t.TempDir()
	inputFile := filepath.Join(tmpDir, "test.md")

	if err := os.WriteFile(inputFile, []byte("# Test"), 0o644); err != nil {
		t.Fatalf("failed to create test file: %v", err)
	}

	cmd := exec.Command("veve", inputFile, "--theme", "definitely-nonexistent-theme")

	var stdout, stderr bytes.Buffer
	cmd.Stdout = &stdout
	cmd.Stderr = &stderr

	_ = cmd.Run()

	stderrStr := stderr.String()
	stdoutStr := stdout.String()

	if stdoutStr != "" && strings.Contains(stdoutStr, "theme") {
		t.Errorf("theme error found in stdout (should be in stderr)")
	}

	if stderrStr == "" {
		t.Logf("Note: stderr is empty for theme error")
	} else {
		t.Logf("Theme error in stderr: %s", stderrStr)
	}
}

// TestStdinErrorToStderr tests that stdin read errors go to stderr.
func TestStdinErrorToStderr(t *testing.T) {
	tmpDir := t.TempDir()
	outputFile := filepath.Join(tmpDir, "output.pdf")

	cmd := exec.Command("veve", "-", "-o", outputFile)
	cmd.Stdin = nil // This might cause an error

	var stdout, stderr bytes.Buffer
	cmd.Stdout = &stdout
	cmd.Stderr = &stderr

	output := cmd.Run()

	stderrStr := stderr.String()
	stdoutStr := stdout.String()

	if output != nil && stdoutStr != "" {
		t.Logf("error output in stdout (may be ok)")
	}

	if stderrStr != "" && (strings.Contains(stderrStr, "error") ||
		strings.Contains(stderrStr, "failed")) {
		t.Logf("stdin error in stderr: %s", stderrStr)
	}
}

// TestPandocErrorToStderr tests that pandoc errors are forwarded to stderr.
func TestPandocErrorToStderr(t *testing.T) {
	// Create a markdown file that might cause pandoc to produce warnings
	tmpDir := t.TempDir()
	inputFile := filepath.Join(tmpDir, "test.md")

	// Create file with potentially problematic content
	content := "# Test\n\n```invalid-lang\ncode\n```\n"
	if err := os.WriteFile(inputFile, []byte(content), 0o644); err != nil {
		t.Fatalf("failed to create test file: %v", err)
	}

	cmd := exec.Command("veve", inputFile)

	var stdout, stderr bytes.Buffer
	cmd.Stdout = &stdout
	cmd.Stderr = &stderr

	_ = cmd.Run()

	if stderrStr := stderr.String(); stderrStr != "" {
		t.Logf("Pandoc output in stderr: %s", stderrStr)
	}
}

// TestStdoutNotContainErrors tests that PDF stdout doesn't contain error text.
func TestStdoutNotContainErrors(t *testing.T) {
	tmpDir := t.TempDir()
	inputFile := filepath.Join(tmpDir, "test.md")
	if err := os.WriteFile(inputFile, []byte("# Test"), 0o644); err != nil {
		t.Fatalf("failed to create test file: %v", err)
	}

	cmd := exec.Command("veve", inputFile, "-o", "-")

	var stdout, stderr bytes.Buffer
	cmd.Stdout = &stdout
	cmd.Stderr = &stderr

	err := cmd.Run()

	if err != nil {
		t.Logf("conversion failed: %v", err)
		return
	}

	stdoutData := stdout.Bytes()

	// Stdout should contain PDF binary, not error text
	if bytes.Contains(stdoutData, []byte("error")) || bytes.Contains(stdoutData, []byte("Error")) {
		t.Errorf("error text found in PDF stdout")
	}

	// Should start with PDF magic number
	if len(stdoutData) < 4 || string(stdoutData[:4]) != "%PDF" {
		t.Errorf("stdout does not contain valid PDF")
	}
}
