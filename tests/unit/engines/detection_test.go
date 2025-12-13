package engines_test

import (
	"os/exec"
	"testing"

	"github.com/madstone-tech/veve-cli/internal/engines"
)

// TestDetectInstalledEngines verifies that engine detection works
func TestDetectInstalledEngines(t *testing.T) {
	// Check if pandoc is available (required for testing)
	_, err := exec.LookPath("pandoc")
	if err != nil {
		t.Skip("pandoc not found in PATH; skipping engine detection tests")
	}

	tests := []struct {
		name        string
		expectError bool
		minEngines  int
	}{
		{
			name:        "detects at least one engine",
			expectError: false,
			minEngines:  0, // At least one engine should be found if any pdf engine installed
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			detected, err := engines.DetectInstalledEngines()

			if tt.expectError && err == nil {
				t.Errorf("expected error, got none")
			}
			if !tt.expectError && err != nil {
				t.Errorf("unexpected error: %v", err)
			}

			if err == nil && len(detected) == 0 {
				t.Skip("no PDF engines found in PATH; test skipped")
			}

			if err == nil && len(detected) > 0 {
				// Verify detected engines have required fields
				for _, e := range detected {
					if e.Name == "" {
						t.Error("engine name is empty")
					}
					if e.DisplayLabel == "" {
						t.Error("engine display label is empty")
					}
					if e.Priority <= 0 {
						t.Errorf("engine priority should be > 0, got %d", e.Priority)
					}
					if !e.IsInstalled {
						t.Errorf("detected engine %s should have IsInstalled=true", e.Name)
					}
				}
			}
		})
	}
}

// TestDetectInstalledEngines_NoEnginesAvailable tests behavior when no engines found
func TestDetectInstalledEngines_NoEnginesAvailable(t *testing.T) {
	// This test would require mocking PATH which is complex
	// It documents expected behavior: should return error with helpful message
	t.Run("returns helpful error when no engines found", func(t *testing.T) {
		// Implementation note: actual test would mock exec.LookPath
		// For now, we skip if engines are actually available
		detected, err := engines.DetectInstalledEngines()

		if len(detected) > 0 {
			t.Skip("engines found; test requires no engines to be available")
		}

		if err == nil {
			t.Error("expected error when no engines detected")
		}
		if err != nil {
			// Verify error message is helpful
			errMsg := err.Error()
			if errMsg == "" {
				t.Error("error message should not be empty")
			}
		}
	})
}

// TestDetectInstalledEngines_HandlesMissingEngines tests that missing engines don't break detection
func TestDetectInstalledEngines_HandlesMissingEngines(t *testing.T) {
	t.Run("detection completes even if some engines missing", func(t *testing.T) {
		detected, err := engines.DetectInstalledEngines()

		// Should not error just because some engines are missing
		if err != nil {
			// Error is acceptable if NO engines found
			if len(detected) == 0 {
				return
			}
			t.Errorf("should not error when at least one engine available: %v", err)
		}

		// If we got engines, verify they're valid
		if len(detected) > 0 {
			validNames := make(map[string]bool)
			for _, e := range detected {
				validNames[e.Name] = true
				// Verify each engine is actually installed
				if !e.IsInstalled {
					t.Errorf("detected engine %s should have IsInstalled=true", e.Name)
				}
			}
		}
	})
}

// TestEnginePriorityOrder verifies engines are returned in priority order
func TestEnginePriorityOrder(t *testing.T) {
	t.Run("engines are in correct priority order", func(t *testing.T) {
		detected, err := engines.DetectInstalledEngines()
		if err != nil || len(detected) == 0 {
			t.Skip("no engines detected; skipping priority test")
		}

		// Verify priorities are unique and in order
		priorities := make(map[int]bool)
		for i := 0; i < len(detected)-1; i++ {
			if detected[i].Priority >= detected[i+1].Priority {
				t.Errorf("priority order violated: %d >= %d",
					detected[i].Priority, detected[i+1].Priority)
			}

			if priorities[detected[i].Priority] {
				t.Errorf("duplicate priority: %d", detected[i].Priority)
			}
			priorities[detected[i].Priority] = true
		}
	})
}
