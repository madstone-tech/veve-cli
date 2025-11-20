package theme

import (
	"encoding/json"
	"os"
	"time"
)

// Theme represents metadata about a theme.
type Theme struct {
	Name        string    `json:"name"`        // Theme identifier (e.g., "dark")
	DisplayName string    `json:"displayName"` // Human-readable name
	Description string    `json:"description"` // Short description
	Author      string    `json:"author"`      // Theme author
	Version     string    `json:"version"`     // Theme version
	FilePath    string    `json:"filePath"`    // Path to the CSS file
	IsBuiltIn   bool      `json:"isBuiltIn"`   // Whether this is a built-in theme
	CreatedAt   time.Time `json:"createdAt"`   // When the theme was added
}

// Registry manages all available themes (built-in + user-installed).
type Registry struct {
	Themes map[string]Theme `json:"themes"`
}

// NewRegistry creates a new empty theme registry.
func NewRegistry() *Registry {
	return &Registry{
		Themes: make(map[string]Theme),
	}
}

// AddTheme adds or updates a theme in the registry.
func (r *Registry) AddTheme(theme Theme) {
	r.Themes[theme.Name] = theme
}

// RemoveTheme removes a theme from the registry.
func (r *Registry) RemoveTheme(name string) bool {
	if _, exists := r.Themes[name]; exists {
		delete(r.Themes, name)
		return true
	}
	return false
}

// GetTheme retrieves a theme by name.
func (r *Registry) GetTheme(name string) (Theme, bool) {
	theme, exists := r.Themes[name]
	return theme, exists
}

// ListThemes returns all themes in the registry.
func (r *Registry) ListThemes() []Theme {
	themes := make([]Theme, 0, len(r.Themes))
	for _, theme := range r.Themes {
		themes = append(themes, theme)
	}
	return themes
}

// LoadFromFile loads the theme registry from a JSON file.
func (r *Registry) LoadFromFile(filePath string) error {
	data, err := os.ReadFile(filePath)
	if err != nil {
		if os.IsNotExist(err) {
			// File doesn't exist yet; return empty registry
			return nil
		}
		return err
	}

	return json.Unmarshal(data, r)
}

// SaveToFile saves the theme registry to a JSON file.
func (r *Registry) SaveToFile(filePath string) error {
	data, err := json.MarshalIndent(r, "", "  ")
	if err != nil {
		return err
	}

	return os.WriteFile(filePath, data, 0o644)
}
