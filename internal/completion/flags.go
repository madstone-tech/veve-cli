// Package completion provides shell autocompletion support for CLI flags and values.
package completion

// CLIFlag represents a command-line flag with its configuration
type CLIFlag struct {
	// Name is the flag name (e.g., "engine", "theme")
	Name string

	// ShortForm is single-character short form (e.g., "e" for -e)
	ShortForm string

	// Description is help text (used by --help and completion)
	Description string

	// AcceptedValues is list of valid values (e.g., ["xelatex", "weasyprint"])
	AcceptedValues []string

	// ValueType is type of value: ENUM, STRING, PATH, FILE, INTEGER
	ValueType ValueType

	// IsRequired indicates whether flag must be provided
	IsRequired bool

	// DefaultValue is default if flag not provided
	DefaultValue string

	// MutuallyExclusive lists other flags that conflict with this one
	MutuallyExclusive []string

	// DeprecationNotice is warning if flag is deprecated
	DeprecationNotice string
}

// ValueType represents the type of value a flag accepts
type ValueType string

const (
	// ValueTypeEnum accepts one of predefined values
	ValueTypeEnum ValueType = "ENUM"

	// ValueTypeString accepts any string value
	ValueTypeString ValueType = "STRING"

	// ValueTypePath accepts file/directory path
	ValueTypePath ValueType = "PATH"

	// ValueTypeFile accepts file path
	ValueTypeFile ValueType = "FILE"

	// ValueTypeInteger accepts integer values
	ValueTypeInteger ValueType = "INTEGER"
)

// IsValidValue checks if value is in AcceptedValues
func (cf *CLIFlag) IsValidValue(value string) bool {
	for _, v := range cf.AcceptedValues {
		if v == value {
			return true
		}
	}
	return false
}

// GetCompletions returns filtered completions for the given partial input
func (cf *CLIFlag) GetCompletions(toComplete string) []string {
	var completions []string
	for _, v := range cf.AcceptedValues {
		if len(toComplete) == 0 || hasPrefix(v, toComplete) {
			completions = append(completions, v)
		}
	}
	return completions
}

// hasPrefix is a helper for case-sensitive prefix matching
func hasPrefix(s, prefix string) bool {
	return len(s) >= len(prefix) && s[:len(prefix)] == prefix
}

// EngineFlag returns the CLI flag definition for --engine
func EngineFlag() *CLIFlag {
	return &CLIFlag{
		Name:              "engine",
		ShortForm:         "e",
		Description:       "PDF rendering engine to use (xelatex, weasyprint, prince)",
		AcceptedValues:    []string{"xelatex", "lualatex", "weasyprint", "prince"},
		ValueType:         ValueTypeEnum,
		IsRequired:        false,
		DefaultValue:      "xelatex",
		MutuallyExclusive: []string{},
	}
}
