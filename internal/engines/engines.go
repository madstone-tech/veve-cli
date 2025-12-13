// Package engines provides PDF engine detection, validation, and selection logic
// for unicode-capable markdown to PDF conversion.
package engines

import "time"

// PDFEngine represents a PDF rendering engine available on the system.
type PDFEngine struct {
	// Name is the identifier name (e.g., "xelatex", "weasyprint", "prince")
	Name string

	// DisplayLabel is a user-friendly name for display (e.g., "XeLaTeX")
	DisplayLabel string

	// Priority is the selection priority (1=highest, used in fallback chain)
	Priority int

	// UnicodeSupport indicates whether engine handles full unicode character set
	UnicodeSupport bool

	// EmojiSupport indicates whether engine renders standard emoji correctly
	EmojiSupport bool

	// IsInstalled indicates whether engine is present in system PATH
	IsInstalled bool

	// Version is the detected engine version (for debugging)
	Version string

	// InstallationInstructions is platform-specific installation help text
	InstallationInstructions string
}

// AvailableEngine represents a runtime representation of detected engines with capabilities
type AvailableEngine struct {
	// Engine is a reference to base engine definition
	Engine PDFEngine

	// IsCapableOfUnicode indicates detected unicode capability (tested at runtime)
	IsCapableOfUnicode bool

	// UnicodeTestResult details of unicode detection test
	UnicodeTestResult *TestResult

	// FallbackRank is position in fallback chain (1=first tried)
	FallbackRank int
}

// TestResult represents the outcome of a unicode capability test
type TestResult struct {
	// Success indicates if test passed
	Success bool

	// ExitCode from engine process
	ExitCode int

	// Stderr output from engine
	Stderr string

	// Duration of test execution
	Duration time.Duration

	// ErrorMessage if test failed
	ErrorMessage string
}

// CanHandle checks if engine can handle the given text content
// Returns true if the engine is capable of unicode/emoji based on test results
func (ae *AvailableEngine) CanHandle(text string) bool {
	return ae.IsCapableOfUnicode
}

// GetErrorMessage formats an actionable error message if test failed
func (ae *AvailableEngine) GetErrorMessage() string {
	if ae.UnicodeTestResult == nil || ae.UnicodeTestResult.Success {
		return ""
	}

	if ae.UnicodeTestResult.ErrorMessage != "" {
		return ae.UnicodeTestResult.ErrorMessage
	}

	return "Engine failed unicode capability test"
}

// PriorityOrder defines the engine selection priority (highest to lowest)
var PriorityOrder = []string{
	"xelatex",    // Priority 1: Native UTF-8 support, widely available
	"lualatex",   // Priority 2: Similar capabilities, slightly slower
	"weasyprint", // Priority 3: For users without LaTeX, requires Python
	"prince",     // Priority 4: Commercial option, excellent support
}

// DefaultEngineDefinitions provides the set of supported engines
func DefaultEngineDefinitions() map[string]PDFEngine {
	return map[string]PDFEngine{
		"xelatex": {
			Name:           "xelatex",
			DisplayLabel:   "XeLaTeX",
			Priority:       1,
			UnicodeSupport: true,
			EmojiSupport:   true,
			IsInstalled:    false, // Will be detected
			Version:        "",
			InstallationInstructions: "" +
				"macOS: brew install mactex\n" +
				"Ubuntu/Debian: sudo apt-get install texlive-xetex\n" +
				"Fedora: sudo dnf install texlive-xetex\n" +
				"Windows: Download from https://miktex.org/",
		},
		"lualatex": {
			Name:           "lualatex",
			DisplayLabel:   "LuaLaTeX",
			Priority:       2,
			UnicodeSupport: true,
			EmojiSupport:   true,
			IsInstalled:    false,
			Version:        "",
			InstallationInstructions: "" +
				"macOS: brew install mactex\n" +
				"Ubuntu/Debian: sudo apt-get install texlive-luatex\n" +
				"Fedora: sudo dnf install texlive-luatex\n" +
				"Windows: Download from https://miktex.org/",
		},
		"weasyprint": {
			Name:           "weasyprint",
			DisplayLabel:   "WeasyPrint",
			Priority:       3,
			UnicodeSupport: true,
			EmojiSupport:   true,
			IsInstalled:    false,
			Version:        "",
			InstallationInstructions: "" +
				"macOS: brew install weasyprint\n" +
				"Ubuntu/Debian: sudo apt-get install weasyprint\n" +
				"Fedora: sudo dnf install weasyprint\n" +
				"Windows: pip install weasyprint",
		},
		"prince": {
			Name:           "prince",
			DisplayLabel:   "Prince XML",
			Priority:       4,
			UnicodeSupport: true,
			EmojiSupport:   true,
			IsInstalled:    false,
			Version:        "",
			InstallationInstructions: "" +
				"Download from https://www.princexml.com/download/\n" +
				"Prince is a commercial tool with a free trial.",
		},
	}
}
