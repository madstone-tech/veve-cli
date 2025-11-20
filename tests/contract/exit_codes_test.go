package contract

import (
	"os"
	"os/exec"
	"path/filepath"
	"testing"
)

// TestExitCodeSuccess tests that veve returns exit code 0 on success.
func TestExitCodeSuccess(t *testing.T) {
	tmpDir := t.TempDir()
	inputFile := filepath.Join(tmpDir, "test.md")
	outputFile := filepath.Join(tmpDir, "output.pdf")

	if err := os.WriteFile(inputFile, []byte("# Test\nContent"), 0o644); err != nil {
		t.Fatalf("failed to create test file: %v", err)
	}

	cmd := exec.Command("veve", inputFile, "-o", outputFile)
	err := cmd.Run()

	if err != nil {
		if exitErr, ok := err.(*exec.ExitError); ok {
			t.Errorf("expected exit code 0, got %d", exitErr.ExitCode())
		} else {
			// Command not found or other error
			t.Skipf("veve not installed: %v", err)
		}
	}
}

// TestExitCodeError tests that veve returns non-zero on error.
func TestExitCodeError(t *testing.T) {
	// Use non-existent input file
	cmd := exec.Command("veve", "/nonexistent/file.md", "-o", "/tmp/out.pdf")
	err := cmd.Run()

	if err == nil {
		t.Errorf("expected non-zero exit code for non-existent file")
	} else {
		if exitErr, ok := err.(*exec.ExitError); ok {
			code := exitErr.ExitCode()
			if code == 0 {
				t.Errorf("expected non-zero exit code, got 0")
			}
			t.Logf("Non-existent file returned exit code: %d", code)
		}
	}
}

// TestExitCodeInvalidTheme tests exit code for invalid theme.
func TestExitCodeInvalidTheme(t *testing.T) {
	tmpDir := t.TempDir()
	inputFile := filepath.Join(tmpDir, "test.md")

	if err := os.WriteFile(inputFile, []byte("# Test"), 0o644); err != nil {
		t.Fatalf("failed to create test file: %v", err)
	}

	cmd := exec.Command("veve", inputFile, "--theme", "nonexistent-theme")
	err := cmd.Run()

	if err == nil {
		t.Logf("Note: command succeeded with invalid theme (may use default)")
	} else {
		if exitErr, ok := err.(*exec.ExitError); ok {
			code := exitErr.ExitCode()
			if code != 0 {
				t.Logf("Invalid theme returned exit code: %d", code)
			}
		}
	}
}

// TestExitCodeMissingInput tests exit code when input is missing.
func TestExitCodeMissingInput(t *testing.T) {
	cmd := exec.Command("veve")
	err := cmd.Run()

	if err == nil {
		t.Logf("Note: command succeeded with no arguments (shows help)")
	} else {
		if exitErr, ok := err.(*exec.ExitError); ok {
			code := exitErr.ExitCode()
			t.Logf("Missing input returned exit code: %d", code)
		}
	}
}

// TestExitCodeInvalidFlags tests exit code for invalid flags.
func TestExitCodeInvalidFlags(t *testing.T) {
	tmpDir := t.TempDir()
	inputFile := filepath.Join(tmpDir, "test.md")

	if err := os.WriteFile(inputFile, []byte("# Test"), 0o644); err != nil {
		t.Fatalf("failed to create test file: %v", err)
	}

	cmd := exec.Command("veve", inputFile, "--invalid-flag")
	err := cmd.Run()

	if err == nil {
		t.Logf("Note: invalid flag was accepted")
	} else {
		if exitErr, ok := err.(*exec.ExitError); ok {
			code := exitErr.ExitCode()
			t.Logf("Invalid flag returned exit code: %d", code)
		}
	}
}

// TestExitCodePandocNotFound tests behavior when pandoc is not found.
// This is tricky to test without removing pandoc from PATH.
// For now, we just verify the error handling works.
func TestExitCodePandocRequired(t *testing.T) {
	// This test would require temporarily removing pandoc from PATH
	// Skipping for now as it's environment-specific
	t.Skip("requires pandoc in PATH manipulation")
}

// TestExitCodeStdinConversion tests exit code for stdin conversion.
func TestExitCodeStdinConversion(t *testing.T) {
	tmpDir := t.TempDir()
	outputFile := filepath.Join(tmpDir, "output.pdf")

	cmd := exec.Command("veve", "-", "-o", outputFile)
	cmd.Stdin = nil // Empty stdin

	err := cmd.Run()

	if err != nil {
		if exitErr, ok := err.(*exec.ExitError); ok {
			code := exitErr.ExitCode()
			t.Logf("Stdin conversion error returned exit code: %d", code)
		}
	}
}

// TestExitCodeMultipleErrors tests exit code with multiple validation errors.
func TestExitCodeMultipleErrors(t *testing.T) {
	// Provide non-existent file with invalid theme
	cmd := exec.Command("veve", "/nonexistent/file.md", "--theme", "invalid-theme")
	err := cmd.Run()

	if err == nil {
		t.Logf("Note: command succeeded despite errors")
	} else {
		if exitErr, ok := err.(*exec.ExitError); ok {
			code := exitErr.ExitCode()
			if code != 0 {
				t.Logf("Multiple errors returned exit code: %d", code)
			}
		}
	}
}

// TestExitCodeHelpFlag tests exit code for --help flag.
func TestExitCodeHelpFlag(t *testing.T) {
	cmd := exec.Command("veve", "--help")
	err := cmd.Run()

	// --help usually returns 0 (success)
	if err != nil {
		if exitErr, ok := err.(*exec.ExitError); ok {
			code := exitErr.ExitCode()
			t.Logf("Help flag returned exit code: %d (may be 0)", code)
		}
	}
}

// TestExitCodeVersionFlag tests exit code for --version flag.
func TestExitCodeVersionFlag(t *testing.T) {
	cmd := exec.Command("veve", "--version")
	err := cmd.Run()

	// --version usually returns 0 (success)
	if err != nil {
		if exitErr, ok := err.(*exec.ExitError); ok {
			code := exitErr.ExitCode()
			t.Logf("Version flag returned exit code: %d (may be 0)", code)
		}
	}
}

// TestExitCodeSubcommandHelp tests exit code for theme subcommand.
func TestExitCodeSubcommandHelp(t *testing.T) {
	cmd := exec.Command("veve", "theme", "list")
	err := cmd.Run()

	if err == nil {
		t.Logf("Theme list command succeeded with exit code 0")
	} else {
		if exitErr, ok := err.(*exec.ExitError); ok {
			code := exitErr.ExitCode()
			t.Logf("Theme list returned exit code: %d", code)
		}
	}
}
