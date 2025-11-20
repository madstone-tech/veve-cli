package theme

import (
	"fmt"
	"os"
	"path/filepath"
	"sort"
	"strings"
)

// Loader handles loading themes from built-in and user-installed locations.
type Loader struct {
	builtInThemes map[string]Theme
	userThemesDir string
	registry      *Registry
}

// NewLoader creates a new theme loader.
func NewLoader(userThemesDir string) *Loader {
	return &Loader{
		builtInThemes: make(map[string]Theme),
		userThemesDir: userThemesDir,
		registry:      NewRegistry(),
	}
}

// AddBuiltInTheme registers a built-in theme.
func (l *Loader) AddBuiltInTheme(theme Theme) {
	l.builtInThemes[theme.Name] = theme
}

// DiscoverThemes discovers all available themes (built-in + user-installed).
func (l *Loader) DiscoverThemes() error {
	// Start fresh
	l.registry = NewRegistry()

	// Add all built-in themes
	for _, theme := range l.builtInThemes {
		l.registry.AddTheme(theme)
	}

	// Discover user-installed themes
	if _, err := os.Stat(l.userThemesDir); err == nil {
		entries, err := os.ReadDir(l.userThemesDir)
		if err != nil {
			return err
		}

		for _, entry := range entries {
			if entry.IsDir() {
				continue
			}

			// Only process .css files
			if !strings.HasSuffix(entry.Name(), ".css") {
				continue
			}

			// Extract theme name from filename (without .css extension)
			themeName := strings.TrimSuffix(entry.Name(), ".css")
			filePath := filepath.Join(l.userThemesDir, entry.Name())

			theme := Theme{
				Name:        themeName,
				DisplayName: themeName,
				Description: "Custom user theme",
				Author:      "Unknown",
				Version:     "1.0.0",
				FilePath:    filePath,
				IsBuiltIn:   false,
			}

			// User themes override built-in themes with the same name
			l.registry.AddTheme(theme)
		}
	}

	return nil
}

// LoadTheme loads a theme by name, checking built-in and user-installed themes.
// User-installed themes take precedence over built-in themes with the same name.
func (l *Loader) LoadTheme(name string) (Theme, error) {
	theme, exists := l.registry.GetTheme(name)
	if !exists {
		available := l.ListThemes()
		availableNames := make([]string, len(available))
		for i, t := range available {
			availableNames[i] = t.Name
		}
		return Theme{}, fmt.Errorf("theme not found: %s (available: %s)", name, strings.Join(availableNames, ", "))
	}

	return theme, nil
}

// LoadThemeCSS loads the CSS content for a theme.
// For built-in themes, returns the embedded CSS.
// For user-installed themes, reads from the file system.
func (l *Loader) LoadThemeCSS(theme Theme) (string, error) {
	if theme.IsBuiltIn {
		// Get embedded CSS from themes package
		// This requires themes/embed.go to be in the main package or imported
		return "", fmt.Errorf("built-in theme CSS loading not yet implemented")
	}

	// Read from file system
	css, err := os.ReadFile(theme.FilePath)
	if err != nil {
		return "", fmt.Errorf("failed to read theme file: %w", err)
	}

	return string(css), nil
}

// ListThemes returns all available themes, sorted by name.
func (l *Loader) ListThemes() []Theme {
	themes := l.registry.ListThemes()
	sort.Slice(themes, func(i, j int) bool {
		return themes[i].Name < themes[j].Name
	})
	return themes
}

// GetRegistry returns the underlying theme registry.
func (l *Loader) GetRegistry() *Registry {
	return l.registry
}
