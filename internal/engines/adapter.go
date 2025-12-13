// Package engines provides PDF engine detection, validation, and selection logic.
package engines

import (
	"fmt"
	"sync"
)

// GlobalSelector is a singleton instance for engine selection across the CLI
var (
	globalSelector *EngineSelector
	selectorOnce   sync.Once
	selectorErr    error
)

// GetDefaultEngine returns the default unicode-capable engine
// Uses singleton pattern for efficiency (detection runs once per CLI invocation)
func GetDefaultEngine() (*PDFEngine, error) {
	selectorOnce.Do(func() {
		globalSelector, selectorErr = NewEngineSelector()
	})

	if selectorErr != nil {
		return nil, selectorErr
	}

	return globalSelector.SelectDefaultEngine()
}

// SelectEngineForConversion selects an engine for conversion
// If engineName is empty, uses default; otherwise uses specified engine
// Respects FR-001.1: explicit flag overrides automatic selection
func SelectEngineForConversion(engineName string) (*PDFEngine, error) {
	selectorOnce.Do(func() {
		globalSelector, selectorErr = NewEngineSelector()
	})

	if selectorErr != nil {
		return nil, selectorErr
	}

	if engineName == "" {
		// Use default selection
		return globalSelector.SelectDefaultEngine()
	}

	// Use explicit selection
	engine, err := globalSelector.SelectEngine(engineName)
	if err != nil {
		return nil, err
	}

	return engine, nil
}

// ValidateEngineForContent validates if engine can handle unicode content
// Returns error with actionable message if engine cannot handle unicode
func ValidateEngineForContent(engine *PDFEngine, hasUnicode bool) error {
	if !hasUnicode {
		// No unicode content, any engine is fine
		return nil
	}

	selectorOnce.Do(func() {
		globalSelector, selectorErr = NewEngineSelector()
	})

	if selectorErr != nil {
		return fmt.Errorf("cannot validate engine: %w", selectorErr)
	}

	// Check if engine is unicode-capable
	if !globalSelector.IsEngineAvailable(engine.Name) {
		return fmt.Errorf(
			"engine '%s' does not support unicode; use one of: %v",
			engine.Name, globalSelector.GetAvailableEngines(),
		)
	}

	return nil
}

// GetAvailableEnginesForCompletion returns list of engine names for shell completion
func GetAvailableEnginesForCompletion() []string {
	selectorOnce.Do(func() {
		globalSelector, selectorErr = NewEngineSelector()
	})

	if selectorErr != nil {
		// If engine detection fails, return hardcoded list for completion
		return []string{"xelatex", "lualatex", "weasyprint", "prince"}
	}

	available := globalSelector.GetAvailableEngines()
	if len(available) == 0 {
		// Fallback to all known engines if none detected as unicode-capable
		return []string{"xelatex", "lualatex", "weasyprint", "prince"}
	}

	return available
}

// ResetGlobalSelector clears the singleton (useful for testing)
func ResetGlobalSelector() {
	globalSelector = nil
	selectorErr = nil
	selectorOnce = sync.Once{}
}
