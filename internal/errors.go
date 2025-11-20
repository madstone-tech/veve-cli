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
