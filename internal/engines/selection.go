// Package engines provides PDF engine detection, validation, and selection logic.
package engines

import (
	"fmt"
	"sync"
)

// EngineSelector handles automatic engine selection and fallback logic
type EngineSelector struct {
	availableEngines []AvailableEngine
	defaultEngine    *AvailableEngine
	mu               sync.RWMutex
}

// NewEngineSelector creates and initializes an engine selector
// Detects installed engines and validates unicode support
func NewEngineSelector() (*EngineSelector, error) {
	selector := &EngineSelector{}

	// Detect installed engines
	installed, err := DetectInstalledEngines()
	if err != nil {
		return nil, err
	}

	// Validate each engine's unicode support
	for _, engine := range installed {
		testResult := ValidateUnicodeSupport(engine)

		available := AvailableEngine{
			Engine:             engine,
			IsCapableOfUnicode: testResult.Success,
			UnicodeTestResult:  testResult,
			FallbackRank:       engine.Priority,
		}

		selector.availableEngines = append(selector.availableEngines, available)

		// First unicode-capable engine becomes default
		if testResult.Success && selector.defaultEngine == nil {
			selector.defaultEngine = &available
		}
	}

	// If no unicode-capable engine found, return error
	if selector.defaultEngine == nil {
		return nil, fmt.Errorf(
			"no unicode-capable PDF engine found; " +
				"please install one of: xelatex, lualatex, weasyprint, or prince",
		)
	}

	return selector, nil
}

// SelectDefaultEngine returns the default unicode-capable engine
// Respects priority order: xelatex → lualatex → weasyprint → prince
func (es *EngineSelector) SelectDefaultEngine() (*PDFEngine, error) {
	es.mu.RLock()
	defer es.mu.RUnlock()

	if es.defaultEngine == nil {
		return nil, fmt.Errorf("no default engine selected")
	}

	return &es.defaultEngine.Engine, nil
}

// SelectEngine selects an engine by name
// Returns error if engine not available or not unicode-capable
func (es *EngineSelector) SelectEngine(engineName string) (*PDFEngine, error) {
	es.mu.RLock()
	defer es.mu.RUnlock()

	for _, available := range es.availableEngines {
		if available.Engine.Name == engineName {
			if !available.IsCapableOfUnicode {
				return nil, fmt.Errorf(
					"engine '%s' does not support unicode: %s",
					engineName, available.GetErrorMessage(),
				)
			}
			return &available.Engine, nil
		}
	}

	return nil, fmt.Errorf("engine '%s' not found or not installed", engineName)
}

// SelectEngineFallback attempts to use specified engine, falls back to default if fails
// Returns the selected engine and whether fallback was needed
func (es *EngineSelector) SelectEngineFallback(engineName string) (*PDFEngine, bool, error) {
	es.mu.RLock()
	defer es.mu.RUnlock()

	// Try to find requested engine
	for _, available := range es.availableEngines {
		if available.Engine.Name == engineName {
			if available.IsCapableOfUnicode {
				return &available.Engine, false, nil
			}
			// Engine found but not unicode-capable, fallback
			break
		}
	}

	// Use default engine as fallback
	if es.defaultEngine != nil {
		return &es.defaultEngine.Engine, true, nil
	}

	return nil, false, fmt.Errorf("no fallback engine available")
}

// GetAvailableEngines returns list of all available unicode-capable engines
func (es *EngineSelector) GetAvailableEngines() []string {
	es.mu.RLock()
	defer es.mu.RUnlock()

	var engines []string
	for _, available := range es.availableEngines {
		if available.IsCapableOfUnicode {
			engines = append(engines, available.Engine.Name)
		}
	}
	return engines
}

// GetAllEngines returns all detected engines (installed or not)
func (es *EngineSelector) GetAllEngines() []AvailableEngine {
	es.mu.RLock()
	defer es.mu.RUnlock()

	return append([]AvailableEngine{}, es.availableEngines...)
}

// IsEngineAvailable checks if engine is available and unicode-capable
func (es *EngineSelector) IsEngineAvailable(engineName string) bool {
	es.mu.RLock()
	defer es.mu.RUnlock()

	for _, available := range es.availableEngines {
		if available.Engine.Name == engineName && available.IsCapableOfUnicode {
			return true
		}
	}
	return false
}

// GetEngineInfo returns detailed information about an engine
func (es *EngineSelector) GetEngineInfo(engineName string) (*AvailableEngine, error) {
	es.mu.RLock()
	defer es.mu.RUnlock()

	for i, available := range es.availableEngines {
		if available.Engine.Name == engineName {
			return &es.availableEngines[i], nil
		}
	}

	return nil, fmt.Errorf("engine '%s' not found", engineName)
}

// ValidateUserSelection validates if user-specified engine can be used
// Returns the engine or detailed error explaining what went wrong
func (es *EngineSelector) ValidateUserSelection(engineName string) (*PDFEngine, error) {
	if engineName == "" {
		// Empty means use default
		return es.SelectDefaultEngine()
	}

	// Check if requested engine exists and is unicode-capable
	return es.SelectEngine(engineName)
}

// RefreshAvailability re-detects and re-validates all engines
// Useful if system state changes (engines installed/uninstalled)
func (es *EngineSelector) RefreshAvailability() error {
	es.mu.Lock()
	defer es.mu.Unlock()

	es.availableEngines = nil
	es.defaultEngine = nil

	// Re-detect
	installed, err := DetectInstalledEngines()
	if err != nil {
		return err
	}

	// Re-validate
	for _, engine := range installed {
		testResult := ValidateUnicodeSupport(engine)

		available := AvailableEngine{
			Engine:             engine,
			IsCapableOfUnicode: testResult.Success,
			UnicodeTestResult:  testResult,
			FallbackRank:       engine.Priority,
		}

		es.availableEngines = append(es.availableEngines, available)

		if testResult.Success && es.defaultEngine == nil {
			es.defaultEngine = &available
		}
	}

	if es.defaultEngine == nil {
		return fmt.Errorf("no unicode-capable engines found after refresh")
	}

	return nil
}
