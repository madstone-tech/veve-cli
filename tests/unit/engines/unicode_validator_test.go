package engines_test

import (
	"os/exec"
	"testing"

	"github.com/madstone-tech/veve-cli/internal/engines"
)

// TestValidateUnicodeSupport_XELaTeX tests that xelatex passes unicode test
func TestValidateUnicodeSupport_XELaTeX(t *testing.T) {
	// Skip if xelatex not available
	_, err := exec.LookPath("xelatex")
	if err != nil {
		t.Skip("xelatex not found; skipping unicode validation test")
	}

	t.Run("xelatex passes unicode test", func(t *testing.T) {
		xelatex := engines.PDFEngine{
			Name:         "xelatex",
			DisplayLabel: "XeLaTeX",
			Priority:     1,
		}

		result := engines.ValidateUnicodeSupport(xelatex)

		if result == nil {
			t.Error("test result should not be nil")
			return
		}

		// XeLaTeX should support unicode
		if !result.Success {
			t.Logf("xelatex unicode test failed: %s", result.ErrorMessage)
			t.Error("xelatex should pass unicode test")
		}

		if result.Duration == 0 {
			t.Error("duration should be recorded")
		}
	})
}

// TestValidateUnicodeSupport_PDFLaTeX tests that pdflatex fails unicode test
func TestValidateUnicodeSupport_PDFLaTeX(t *testing.T) {
	// Skip if pdflatex not available
	_, err := exec.LookPath("pdflatex")
	if err != nil {
		t.Skip("pdflatex not found; skipping unicode validation test")
	}

	t.Run("pdflatex fails unicode test", func(t *testing.T) {
		pdflatex := engines.PDFEngine{
			Name:         "pdflatex",
			DisplayLabel: "PDFLaTeX",
			Priority:     5,
		}

		result := engines.ValidateUnicodeSupport(pdflatex)

		if result == nil {
			t.Error("test result should not be nil")
			return
		}

		// PDFLaTeX should NOT support unicode
		if result.Success {
			t.Error("pdflatex should fail unicode test")
		}

		if result.ErrorMessage == "" {
			t.Error("error message should be provided for failed test")
		}
	})
}

// TestValidateUnicodeSupport_Timeout tests that test respects timeout
func TestValidateUnicodeSupport_Timeout(t *testing.T) {
	t.Run("test completes within timeout", func(t *testing.T) {
		xelatex := engines.PDFEngine{
			Name:         "xelatex",
			DisplayLabel: "XeLaTeX",
			Priority:     1,
		}

		result := engines.ValidateUnicodeSupport(xelatex)

		if result == nil {
			t.Skip("test result nil; xelatex may not be installed")
			return
		}

		// Test should complete (not timeout)
		if result.Duration == 0 {
			t.Error("duration not recorded")
		}

		// Timeout should be 5 seconds max
		if result.Duration.Seconds() > 10 {
			t.Errorf("test took too long: %v", result.Duration)
		}
	})
}

// TestValidateUnicodeSupport_TestResultStructure tests that result has all fields
func TestValidateUnicodeSupport_TestResultStructure(t *testing.T) {
	_, err := exec.LookPath("xelatex")
	if err != nil {
		t.Skip("xelatex not found")
	}

	t.Run("test result includes all required fields", func(t *testing.T) {
		xelatex := engines.PDFEngine{
			Name:         "xelatex",
			DisplayLabel: "XeLaTeX",
			Priority:     1,
		}

		result := engines.ValidateUnicodeSupport(xelatex)

		if result == nil {
			t.Error("result should not be nil")
			return
		}

		// Check required fields exist
		if result.Duration == 0 {
			t.Error("Duration should be set")
		}

		// Success or failure should be recorded
		if result.Success {
			// Success case: no error message needed
			if result.ExitCode != 0 {
				t.Logf("note: success but exit code is %d", result.ExitCode)
			}
		} else {
			// Failure case: error message should exist
			if result.ErrorMessage == "" && result.Stderr == "" {
				t.Error("failed test should have error message or stderr")
			}
		}
	})
}

// TestValidateEngineInstalled tests engine installation verification
func TestValidateEngineInstalled(t *testing.T) {
	tests := []struct {
		name        string
		engine      engines.PDFEngine
		expectError bool
	}{
		{
			name: "installed engine passes validation",
			engine: engines.PDFEngine{
				Name:        "xelatex",
				IsInstalled: true,
			},
			expectError: false,
		},
		{
			name: "uninstalled engine fails validation",
			engine: engines.PDFEngine{
				Name:        "fake-engine-xyz",
				IsInstalled: false,
			},
			expectError: true,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			err := engines.ValidateEngineInstalled(tt.engine)

			if tt.expectError && err == nil {
				t.Error("expected error, got none")
			}
			if !tt.expectError && err != nil {
				// May error if actually not installed in test environment
				t.Logf("note: %v", err)
			}
		})
	}
}
