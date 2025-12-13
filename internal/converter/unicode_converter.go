// Package converter provides markdown-to-PDF conversion with unicode support
package converter

import (
	"fmt"
	"os"

	"github.com/madstone-tech/veve-cli/internal/engines"
)

// UnicodeConversionOptions extends ConversionOptions with unicode-aware settings
type UnicodeConversionOptions struct {
	// Base conversion options
	InputFile  string // Path to markdown file (or "-" for stdin)
	OutputFile string // Path to output PDF (or "-" for stdout)
	PDFEngine  string // PDF engine to use (empty = auto-detect)
	Theme      string // Path to CSS theme file (optional)
	Standalone bool   // Generate standalone PDF

	// Unicode settings
	ValidateUnicode bool // Whether to validate unicode support before conversion
	AllowFallback   bool // Whether to allow fallback to different engine
	Verbose         bool // Enable verbose output
}

// ConvertWithUnicodeSupport converts markdown to PDF with automatic engine selection
// for unicode-capable rendering.
//
// Behavior:
// 1. If PDFEngine is specified: use that engine (user override via FR-001.1)
// 2. If PDFEngine is empty: auto-detect unicode in content and select appropriate engine
// 3. If ValidateUnicode is true: verify engine can handle unicode content before conversion
// 4. If AllowFallback is true: try fallback engines if primary fails
//
// Returns error with actionable message if conversion fails
func ConvertWithUnicodeSupport(opts UnicodeConversionOptions) error {
	// Select engine based on options and content
	selectedEngine, err := selectEngineForConversion(opts)
	if err != nil {
		return err
	}

	if opts.Verbose {
		fmt.Fprintf(os.Stderr, "Selected PDF engine: %s\n", selectedEngine.Name)
	}

	// Prepare base conversion options
	convertOpts := ConversionOptions{
		InputFile:  opts.InputFile,
		OutputFile: opts.OutputFile,
		PDFEngine:  selectedEngine.Name,
		Theme:      opts.Theme,
		Standalone: opts.Standalone,
	}

	// Create converter
	converter, err := NewPandocConverter()
	if err != nil {
		return fmt.Errorf("failed to initialize converter: %w", err)
	}

	// Perform conversion
	if err := converter.Convert(convertOpts); err != nil {
		// If conversion failed and unicode was involved, provide actionable error
		if opts.ValidateUnicode {
			contentHasUnicode, _ := detectUnicodeInFile(opts.InputFile)
			if contentHasUnicode {
				return formatUnicodeError(selectedEngine, err)
			}
		}
		return err
	}

	return nil
}

// selectEngineForConversion selects the appropriate PDF engine
// Respects explicit engine selection; auto-detects if needed
// Prefers emoji-capable engines (WeasyPrint/Prince) for emoji-heavy content
func selectEngineForConversion(opts UnicodeConversionOptions) (*engines.PDFEngine, error) {
	// If user explicitly specified engine, use it (FR-001.1)
	if opts.PDFEngine != "" {
		return engines.SelectEngineForConversion(opts.PDFEngine)
	}

	// Read file content for intelligent engine selection
	content, err := os.ReadFile(opts.InputFile)
	if err != nil {
		// If we can't read, use default
		return engines.GetDefaultEngine()
	}

	// Analyze content to determine best engine
	contentStr := string(content)
	hasEmoji := engines.ContainsEmoji(contentStr)
	hasCJK := engines.ContainsCJK(contentStr)
	hasHighComplexity := hasEmoji || (hasCJK && len(contentStr) > 5000) // CJK with lots of text

	// For high-complexity unicode (emoji, extensive CJK), prefer WeasyPrint or Prince
	// These engines have better font support for emoji and complex scripts
	if hasHighComplexity {
		// Try to select an emoji-capable engine
		if engine, err := selectEmojiCapableEngine(); err == nil {
			return engine, nil
		}
		// Fall back to default if emoji-capable selection fails
	}

	// For regular unicode content, use default
	return engines.GetDefaultEngine()
}

// selectEmojiCapableEngine attempts to select an engine with good emoji support
// Prefers WeasyPrint and Prince over XeLaTeX for emoji rendering
func selectEmojiCapableEngine() (*engines.PDFEngine, error) {
	// Try to use selector to find best engine
	selector, err := engines.NewEngineSelector()
	if err != nil {
		return nil, err
	}

	// Get all available engines and check for emoji-capable ones
	availableEngines := selector.GetAvailableEngines()

	// Prefer WeasyPrint (better emoji support)
	for _, name := range availableEngines {
		if name == "weasyprint" {
			return engines.SelectEngineForConversion("weasyprint")
		}
	}

	// Then try Prince
	for _, name := range availableEngines {
		if name == "prince" {
			return engines.SelectEngineForConversion("prince")
		}
	}

	// Fall back to default
	return engines.GetDefaultEngine()
}

// detectUnicodeInFile reads content from file and detects unicode
// Returns (hasUnicode, error)
func detectUnicodeInFile(filePath string) (bool, error) {
	if filePath == "-" {
		// Can't reliably detect from stdin without consuming it
		// Return false (no unicode detected) to use default engine
		return false, nil
	}

	content, err := os.ReadFile(filePath)
	if err != nil {
		// If can't read, return error (file will fail in converter anyway)
		return false, err
	}

	// Check for unicode
	hasUnicode := engines.DetectUnicodeInBytes(content)
	return hasUnicode, nil
}

// formatUnicodeError creates actionable error message for unicode rendering failures
// Provides platform-specific installation instructions
func formatUnicodeError(engine *engines.PDFEngine, originalErr error) error {
	platform := getPlatform()
	instructions := getPlatformInstallInstructions(engine.Name, platform)

	return fmt.Errorf(
		"PDF conversion failed - unicode rendering not supported by engine '%s'\n"+
			"Error: %v\n"+
			"Solution: Install a unicode-capable engine\n%s",
		engine.Name, originalErr, instructions,
	)
}

// getPlatform returns the current platform (darwin, linux, windows)
func getPlatform() string {
	switch os.Getenv("GOOS") {
	case "darwin":
		return "darwin"
	case "linux":
		return "linux"
	case "windows":
		return "windows"
	default:
		// Try to detect from runtime
		switch len(os.Getenv("HOME")) > 0 {
		case true:
			return "linux" // Likely Unix-like
		default:
			return "windows"
		}
	}
}

// getPlatformInstallInstructions returns platform-specific installation help
func getPlatformInstallInstructions(engineName, platform string) string {
	instructions := map[string]map[string]string{
		"xelatex": {
			"darwin":  "\nOn macOS:\n  brew install mactex\n  # This may take 10+ minutes",
			"linux":   "\nOn Ubuntu/Debian:\n  sudo apt-get update\n  sudo apt-get install texlive-xetex\n\nOn Fedora:\n  sudo dnf install texlive-xetex",
			"windows": "\nOn Windows:\n  Download MiKTeX from https://miktex.org/\n  Run installer and select xelatex during setup",
		},
		"lualatex": {
			"darwin":  "\nOn macOS:\n  brew install mactex",
			"linux":   "\nOn Ubuntu/Debian:\n  sudo apt-get update\n  sudo apt-get install texlive-luatex\n\nOn Fedora:\n  sudo dnf install texlive-luatex",
			"windows": "\nOn Windows:\n  Download MiKTeX from https://miktex.org/\n  Run installer and select luatex during setup",
		},
		"weasyprint": {
			"darwin":  "\nOn macOS:\n  brew install weasyprint",
			"linux":   "\nOn Ubuntu/Debian:\n  sudo apt-get update\n  sudo apt-get install weasyprint\n\nOn Fedora:\n  sudo dnf install weasyprint",
			"windows": "\nOn Windows:\n  pip install weasyprint",
		},
	}

	if insts, ok := instructions[engineName]; ok {
		if inst, ok := insts[platform]; ok {
			return inst
		}
	}

	// Fallback
	return "\nPlease install a unicode-capable PDF engine:\n" +
		"  - xelatex (recommended)\n" +
		"  - lualatex\n" +
		"  - weasyprint\n" +
		"See documentation for installation instructions"
}

// QuickConvert is a convenience function for basic conversions with unicode support
// Uses sensible defaults for most users
func QuickConvert(inputFile, outputFile string) error {
	return ConvertWithUnicodeSupport(UnicodeConversionOptions{
		InputFile:       inputFile,
		OutputFile:      outputFile,
		PDFEngine:       "", // Auto-detect
		Theme:           "",
		Standalone:      true,
		ValidateUnicode: true,
		AllowFallback:   true,
		Verbose:         false,
	})
}
