package engines_test

import (
	"os/exec"
	"testing"

	"github.com/madstone-tech/veve-cli/internal/engines"
)

// TestEngineSelector_DefaultSelection tests automatic default selection
func TestEngineSelector_DefaultSelection(t *testing.T) {
	// Skip if no engines available
	_, err := exec.LookPath("xelatex")
	if err != nil {
		t.Skip("xelatex not found; skipping engine selector tests")
	}

	t.Run("selects default engine when available", func(t *testing.T) {
		selector, err := engines.NewEngineSelector()
		if err != nil {
			t.Fatalf("failed to create selector: %v", err)
		}

		engine, err := selector.SelectDefaultEngine()
		if err != nil {
			t.Errorf("failed to select default engine: %v", err)
			return
		}

		if engine == nil {
			t.Error("default engine should not be nil")
			return
		}

		if engine.Name == "" {
			t.Error("engine name should not be empty")
		}

		if !engine.IsInstalled {
			t.Error("default engine should be installed")
		}
	})
}

// TestEngineSelector_SelectSpecificEngine tests selecting a named engine
func TestEngineSelector_SelectSpecificEngine(t *testing.T) {
	selector, err := engines.NewEngineSelector()
	if err != nil {
		t.Skip("no unicode engines available; skipping test")
	}

	t.Run("selects specific engine by name", func(t *testing.T) {
		available := selector.GetAvailableEngines()
		if len(available) == 0 {
			t.Skip("no available engines")
		}

		// Try to select first available engine
		engine, err := selector.SelectEngine(available[0])
		if err != nil {
			t.Errorf("failed to select engine %s: %v", available[0], err)
			return
		}

		if engine == nil {
			t.Error("selected engine should not be nil")
			return
		}

		if engine.Name != available[0] {
			t.Errorf("selected engine name mismatch: want %s, got %s", available[0], engine.Name)
		}
	})

	t.Run("fails to select non-existent engine", func(t *testing.T) {
		engine, err := selector.SelectEngine("fake-engine-xyz")

		if err == nil {
			t.Error("should return error for non-existent engine")
		}
		if engine != nil {
			t.Error("should not return engine for non-existent engine")
		}
	})
}

// TestEngineSelector_FallbackChain tests fallback when primary fails
func TestEngineSelector_FallbackChain(t *testing.T) {
	selector, err := engines.NewEngineSelector()
	if err != nil {
		t.Skip("no engines available")
	}

	t.Run("uses fallback when primary unavailable", func(t *testing.T) {
		// Request an engine that might not be available
		engine, fallback, err := selector.SelectEngineFallback("weasyprint")

		// Should succeed with either the requested engine or fallback
		if err != nil {
			t.Logf("fallback selection failed: %v (acceptable if no engines available)", err)
			return
		}

		if engine == nil {
			t.Error("should return an engine")
			return
		}

		// Either got requested engine (fallback=false) or a fallback (fallback=true)
		if !fallback {
			// Got weasyprint
			if engine.Name != "weasyprint" {
				t.Errorf("expected weasyprint, got %s", engine.Name)
			}
		} else {
			// Got fallback engine
			if engine.Name == "" {
				t.Error("fallback engine name should not be empty")
			}
		}
	})
}

// TestEngineSelector_GetAvailableEngines tests listing available engines
func TestEngineSelector_GetAvailableEngines(t *testing.T) {
	selector, err := engines.NewEngineSelector()
	if err != nil {
		t.Skip("no engines available")
	}

	t.Run("lists all available unicode-capable engines", func(t *testing.T) {
		available := selector.GetAvailableEngines()

		if len(available) == 0 {
			t.Error("should have at least one available engine")
			return
		}

		// Verify each engine name is non-empty
		for _, name := range available {
			if name == "" {
				t.Error("engine name should not be empty")
			}
		}
	})
}

// TestEngineSelector_IsEngineAvailable tests availability check
func TestEngineSelector_IsEngineAvailable(t *testing.T) {
	selector, err := engines.NewEngineSelector()
	if err != nil {
		t.Skip("no engines available")
	}

	available := selector.GetAvailableEngines()
	if len(available) == 0 {
		t.Skip("no available engines")
	}

	t.Run("correctly identifies available engine", func(t *testing.T) {
		if !selector.IsEngineAvailable(available[0]) {
			t.Errorf("engine %s should be available", available[0])
		}
	})

	t.Run("correctly identifies unavailable engine", func(t *testing.T) {
		if selector.IsEngineAvailable("fake-engine-xyz") {
			t.Error("fake engine should not be available")
		}
	})
}

// TestEngineSelector_Priority tests that selection respects priority order
func TestEngineSelector_Priority(t *testing.T) {
	selector, err := engines.NewEngineSelector()
	if err != nil {
		t.Skip("no engines available")
	}

	t.Run("default engine has lowest priority number", func(t *testing.T) {
		defaultEngine, err := selector.SelectDefaultEngine()
		if err != nil {
			t.Fatalf("failed to get default: %v", err)
		}

		allEngines := selector.GetAllEngines()
		if len(allEngines) == 0 {
			t.Skip("no engines")
		}

		// Find default in all engines and verify it has lowest priority
		var defaultPriority int
		for _, e := range allEngines {
			if e.Engine.Name == defaultEngine.Name {
				defaultPriority = e.Engine.Priority
				break
			}
		}

		// Check that default has lower priority number than others
		for _, e := range allEngines {
			if e.IsCapableOfUnicode && e.Engine.Priority < defaultPriority {
				t.Errorf("default engine has priority %d but engine %s has %d",
					defaultPriority, e.Engine.Name, e.Engine.Priority)
			}
		}
	})
}

// TestEngineSelector_Refresh tests refreshing engine availability
func TestEngineSelector_Refresh(t *testing.T) {
	selector, err := engines.NewEngineSelector()
	if err != nil {
		t.Skip("no engines available")
	}

	t.Run("refresh updates engine availability", func(t *testing.T) {
		before := selector.GetAvailableEngines()

		// Refresh availability
		err := selector.RefreshAvailability()
		if err != nil && len(before) > 0 {
			// Only error if we had engines before
			t.Errorf("refresh failed: %v", err)
			return
		}

		after := selector.GetAvailableEngines()

		// Should have same engines after refresh
		if len(before) != len(after) {
			t.Logf("engine count changed: %d -> %d", len(before), len(after))
		}
	})
}
