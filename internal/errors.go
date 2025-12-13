package internal

import (
	"errors"
	"fmt"
)

// Error codes used throughout veve-cli
const (
	ExitSuccess = 0
	ExitError   = 1
	ExitUsage   = 2
)

// VeveError represents a veve-specific error with formatted output.
type VeveError struct {
	Command    string // The command that failed (e.g., "convert", "theme")
	Action     string // The action that failed (e.g., "read input file", "apply theme")
	Reason     string // The underlying reason for failure
	Suggestion string // A helpful suggestion for the user
	Err        error  // The underlying error (for logging)
}

func (e *VeveError) Error() string {
	msg := fmt.Sprintf("[ERROR] %s: %s failed: %s", e.Command, e.Action, e.Reason)
	if e.Suggestion != "" {
		msg += fmt.Sprintf(" (try: %s)", e.Suggestion)
	}
	return msg
}

func (e *VeveError) Unwrap() error {
	return e.Err
}

// NewVeveError creates a new VeveError with the given parameters.
func NewVeveError(command, action, reason, suggestion string, err error) *VeveError {
	return &VeveError{
		Command:    command,
		Action:     action,
		Reason:     reason,
		Suggestion: suggestion,
		Err:        err,
	}
}

// IsVeveError checks if an error is a VeveError.
func IsVeveError(err error) bool {
	var ve *VeveError
	return errors.As(err, &ve)
}

// Common error constructors for consistency

// InputFileNotFound creates an error for missing input files.
func InputFileNotFound(command string, filePath string) *VeveError {
	return NewVeveError(
		command,
		"read input file",
		"file not found: "+filePath,
		"check file path and permissions",
		nil,
	)
}

// ThemeNotFound creates an error for missing themes.
func ThemeNotFound(command string, themeName string, availableThemes string) *VeveError {
	return NewVeveError(
		command,
		"apply theme",
		fmt.Sprintf("theme not found: %s", themeName),
		fmt.Sprintf("use one of: %s", availableThemes),
		nil,
	)
}

// PandocNotFound creates an error for missing Pandoc installation.
func PandocNotFound() *VeveError {
	return NewVeveError(
		"main",
		"initialize converter",
		"pandoc not found in PATH",
		"install pandoc (https://pandoc.org/installing.html)",
		nil,
	)
}

// ConversionFailed creates an error for conversion failures.
func ConversionFailed(command, inputFile string, err error) *VeveError {
	return NewVeveError(
		command,
		"convert markdown",
		fmt.Sprintf("pandoc conversion failed for %s", inputFile),
		"check input file syntax or try with --verbose for details",
		err,
	)
}

// ConfigLoadFailed creates an error for configuration loading failures.
func ConfigLoadFailed(filePath string, err error) *VeveError {
	return NewVeveError(
		"main",
		"load configuration",
		fmt.Sprintf("failed to load config file: %s", filePath),
		fmt.Sprintf("check config file syntax or delete to use defaults"),
		err,
	)
}

// PDFEngineNotFound creates an error for missing PDF engine.
func PDFEngineNotFound(engineName string) *VeveError {
	return NewVeveError(
		"convert",
		"select PDF engine",
		fmt.Sprintf("engine '%s' not found in PATH", engineName),
		"install a unicode-capable engine: xelatex, weasyprint, or prince",
		nil,
	)
}

// UnicodeNotSupported creates an error for unicode rendering failures.
// Uses platform-specific installation instructions.
func UnicodeNotSupported(engineName, platform string) *VeveError {
	instructions := getPlatformInstallInstructions(engineName, platform)

	return NewVeveError(
		"convert",
		"render unicode/emoji",
		fmt.Sprintf("engine '%s' does not support unicode characters", engineName),
		fmt.Sprintf("install xelatex or weasyprint; %s", instructions),
		nil,
	)
}

// NoUnicodeEngineAvailable creates an error when no unicode-capable engine is found.
func NoUnicodeEngineAvailable() *VeveError {
	return NewVeveError(
		"convert",
		"select PDF engine",
		"no unicode-capable PDF engine found in PATH",
		"install one of: xelatex, lualatex, weasyprint, or prince; see docs for instructions",
		nil,
	)
}

// getPlatformInstallInstructions returns platform-specific install instructions
func getPlatformInstallInstructions(engineName, platform string) string {
	instructions := map[string]map[string]string{
		"xelatex": {
			"darwin":  "On macOS: brew install mactex",
			"linux":   "On Ubuntu/Debian: sudo apt-get install texlive-xetex",
			"windows": "On Windows: Download MiKTeX from https://miktex.org/",
		},
		"weasyprint": {
			"darwin":  "On macOS: brew install weasyprint",
			"linux":   "On Ubuntu/Debian: sudo apt-get install weasyprint",
			"windows": "On Windows: pip install weasyprint",
		},
	}

	if insts, ok := instructions[engineName]; ok {
		if inst, ok := insts[platform]; ok {
			return inst
		}
	}

	return "install xelatex or weasyprint"
}
