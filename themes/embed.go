package themes

import (
	_ "embed"
)

//go:embed default.css
var DefaultCSS string

//go:embed dark.css
var DarkCSS string

//go:embed academic.css
var AcademicCSS string

// GetBuiltInTheme returns the CSS content for a built-in theme by name.
func GetBuiltInTheme(name string) (string, bool) {
	switch name {
	case "default":
		return DefaultCSS, true
	case "dark":
		return DarkCSS, true
	case "academic":
		return AcademicCSS, true
	default:
		return "", false
	}
}

// BuiltInThemes returns a list of built-in theme names.
func BuiltInThemes() []string {
	return []string{"default", "dark", "academic"}
}
