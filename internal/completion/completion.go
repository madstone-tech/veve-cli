// Package completion provides shell autocompletion support for CLI flags and values.
package completion

// CompletionOption represents a single suggestion in shell autocompletion
type CompletionOption struct {
	// Label is the suggestion text (displayed in completion)
	Label string

	// Type is the completion type: FLAG, ENGINE, THEME, SUBCOMMAND
	Type CompletionType

	// Description is brief explanation (shown in zsh/fish, hidden in bash)
	// Per spec FR-005/FR-006, descriptions are NOT shown in completion output
	Description string

	// AdditionalInfo is extra context (e.g., installation status)
	AdditionalInfo string
}

// CompletionType represents the type of completion option
type CompletionType string

const (
	// CompletionTypeFlag for CLI flag suggestions
	CompletionTypeFlag CompletionType = "FLAG"

	// CompletionTypeEngine for PDF engine suggestions
	CompletionTypeEngine CompletionType = "ENGINE"

	// CompletionTypeTheme for theme suggestions
	CompletionTypeTheme CompletionType = "THEME"

	// CompletionTypeSubcommand for subcommand suggestions
	CompletionTypeSubcommand CompletionType = "SUBCOMMAND"
)

// NewEngineCompletion creates a completion option for a PDF engine
func NewEngineCompletion(name string, installed bool) *CompletionOption {
	status := "not installed"
	if installed {
		status = "installed"
	}

	return &CompletionOption{
		Label:          name,
		Type:           CompletionTypeEngine,
		Description:    "",
		AdditionalInfo: status,
	}
}

// NewFlagCompletion creates a completion option for a flag
func NewFlagCompletion(name, shortForm, description string) *CompletionOption {
	flagName := "--" + name
	if shortForm != "" {
		flagName += " / -" + shortForm
	}

	return &CompletionOption{
		Label:          flagName,
		Type:           CompletionTypeFlag,
		Description:    description,
		AdditionalInfo: "",
	}
}
