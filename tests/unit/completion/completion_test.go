package completion_test

import (
	"testing"
	"time"

	"github.com/madstone-tech/veve-cli/internal/completion"
	"github.com/madstone-tech/veve-cli/internal/engines"
)

// TestGetAvailableEngines tests that completion returns available engines
func TestGetAvailableEngines(t *testing.T) {
	t.Run("returns list of available engines", func(t *testing.T) {
		enginesList := engines.GetAvailableEnginesForCompletion()

		if len(enginesList) == 0 {
			t.Error("should return at least one engine")
		}

		// Verify engines are strings
		for _, e := range enginesList {
			if e == "" {
				t.Error("engine name should not be empty")
			}
		}

		t.Logf("Available engines: %v", enginesList)
	})

	t.Run("includes known engines in results", func(t *testing.T) {
		enginesList := engines.GetAvailableEnginesForCompletion()

		knownEngines := map[string]bool{
			"xelatex":    false,
			"lualatex":   false,
			"weasyprint": false,
			"prince":     false,
		}

		// At least some known engines should be available
		for _, e := range enginesList {
			knownEngines[e] = true
		}

		hasAtLeastOne := false
		for _, found := range knownEngines {
			if found {
				hasAtLeastOne = true
				break
			}
		}

		if !hasAtLeastOne {
			t.Logf("warning: no known engines found in %v", enginesList)
		}
	})
}

// TestCompletionFiltering tests that completion filters by prefix
func TestCompletionFiltering(t *testing.T) {
	t.Run("filters engines by prefix", func(t *testing.T) {
		flag := completion.EngineFlag()
		completions := flag.GetCompletions("w")

		// Should include 'weasyprint' but not others starting with different letters
		foundWeasyprint := false
		for _, c := range completions {
			if c == "weasyprint" {
				foundWeasyprint = true
			}
			// All completions should start with 'w'
			if len(c) > 0 && c[0] != 'w' {
				t.Errorf("completion '%s' does not start with 'w'", c)
			}
		}

		if !foundWeasyprint && len(completions) > 0 {
			t.Logf("warning: weasyprint not in completions (may not be installed)")
		}
	})

	t.Run("returns all engines for empty prefix", func(t *testing.T) {
		flag := completion.EngineFlag()
		completions := flag.GetCompletions("")

		if len(completions) == 0 {
			t.Error("empty prefix should return all accepted values")
		}

		if len(completions) != len(flag.AcceptedValues) {
			t.Errorf("expected %d completions, got %d", len(flag.AcceptedValues), len(completions))
		}
	})

	t.Run("filters xelatex by 'x' prefix", func(t *testing.T) {
		flag := completion.EngineFlag()
		completions := flag.GetCompletions("x")

		foundXelatex := false
		for _, c := range completions {
			if c == "xelatex" {
				foundXelatex = true
			}
			if len(c) > 0 && c[0] != 'x' {
				t.Errorf("completion '%s' does not start with 'x'", c)
			}
		}

		if !foundXelatex {
			t.Error("xelatex should be in completions for 'x' prefix")
		}
	})
}

// TestCompletionPerformance tests that completion responds quickly
func TestCompletionPerformance(t *testing.T) {
	t.Run("completion returns within 100ms", func(t *testing.T) {
		flag := completion.EngineFlag()

		start := time.Now()
		completions := flag.GetCompletions("w")
		elapsed := time.Since(start)

		if elapsed > 100*time.Millisecond {
			t.Errorf("completion took %v (should be <100ms)", elapsed)
		}

		if len(completions) == 0 {
			t.Logf("warning: no completions returned")
		}

		t.Logf("Completion time: %v", elapsed)
	})
}

// TestEngineFlagValidation tests flag validation
func TestEngineFlagValidation(t *testing.T) {
	t.Run("validates engine names", func(t *testing.T) {
		flag := completion.EngineFlag()

		tests := []struct {
			value string
			valid bool
			name  string
		}{
			{"xelatex", true, "valid xelatex"},
			{"lualatex", true, "valid lualatex"},
			{"weasyprint", true, "valid weasyprint"},
			{"prince", true, "valid prince"},
			{"pdflatex", false, "invalid pdflatex"},
			{"fake-engine", false, "invalid engine"},
			{"", false, "empty value"},
		}

		for _, tt := range tests {
			t.Run(tt.name, func(t *testing.T) {
				valid := flag.IsValidValue(tt.value)
				if valid != tt.valid {
					t.Errorf("expected valid=%v for %s, got %v", tt.valid, tt.value, valid)
				}
			})
		}
	})
}

// TestNewCompletionOptions tests completion option creation
func TestNewCompletionOptions(t *testing.T) {
	t.Run("creates engine completion option", func(t *testing.T) {
		opt := completion.NewEngineCompletion("xelatex", true)

		if opt == nil {
			t.Error("should create completion option")
			return
		}

		if opt.Label != "xelatex" {
			t.Errorf("expected label 'xelatex', got '%s'", opt.Label)
		}

		if opt.Type != completion.CompletionTypeEngine {
			t.Errorf("expected type ENGINE, got %v", opt.Type)
		}

		if opt.AdditionalInfo == "" {
			t.Error("additional info should be set for installed engine")
		}
	})

	t.Run("creates flag completion option", func(t *testing.T) {
		opt := completion.NewFlagCompletion("engine", "e", "PDF engine to use")

		if opt == nil {
			t.Error("should create flag completion option")
			return
		}

		if opt.Type != completion.CompletionTypeFlag {
			t.Errorf("expected type FLAG, got %v", opt.Type)
		}
	})
}
