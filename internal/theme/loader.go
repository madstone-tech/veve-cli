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

	// Ensure user themes directory exists (auto-create if needed)
	if _, err := l.EnsureThemesDir(); err != nil {
		// Log the issue but continue with discovery (built-in themes still available)
		// This is not fatal since built-in themes will still work
	}

	// Add built-in themes with metadata
	builtInThemeMetadata := map[string]Theme{
		"default": {
			Name:        "default",
			DisplayName: "Default",
			Description: "Clean, professional default theme with blue accents",
			Author:      "veve-cli",
			Version:     "1.0.0",
			FilePath:    "", // Embedded
			IsBuiltIn:   true,
		},
		"dark": {
			Name:        "dark",
			DisplayName: "Dark",
			Description: "Dark theme with blue accents, easy on the eyes",
			Author:      "veve-cli",
			Version:     "1.0.0",
			FilePath:    "", // Embedded
			IsBuiltIn:   true,
		},
		"academic": {
			Name:        "academic",
			DisplayName: "Academic",
			Description: "Formal academic paper style with Times New Roman",
			Author:      "veve-cli",
			Version:     "1.0.0",
			FilePath:    "", // Embedded
			IsBuiltIn:   true,
		},
	}

	for _, theme := range builtInThemeMetadata {
		l.registry.AddTheme(theme)
	}

	// Discover user-installed themes (overrides built-in)
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
// If the theme name looks like a file path (contains / or \), loads from that path.
func (l *Loader) LoadThemeCSS(themeName string) (string, error) {
	// Check if the theme name is a file path
	if strings.ContainsAny(themeName, "/\\") {
		return l.LoadThemeFromPath(themeName)
	}

	// First check built-in themes via embed.go
	builtInCSS := l.loadBuiltInThemeCSS(themeName)
	if builtInCSS != "" {
		return builtInCSS, nil
	}

	// Then check user-installed themes
	theme, exists := l.registry.GetTheme(themeName)
	if !exists {
		return "", fmt.Errorf("theme not found: %s", themeName)
	}

	if !theme.IsBuiltIn {
		// Read from file system for user themes
		content, err := os.ReadFile(theme.FilePath)
		if err != nil {
			return "", fmt.Errorf("failed to read theme file: %w", err)
		}

		// Parse metadata if present
		_, css, err := ParseMetadata(string(content))
		if err != nil {
			// Continue even if metadata parsing fails
			css = string(content)
		}

		return css, nil
	}

	return "", fmt.Errorf("theme CSS not found: %s", themeName)
}

// LoadThemeFromPath loads a theme CSS file from a file system path.
// This allows using themes from arbitrary locations via --theme /path/to/theme.css
func (l *Loader) LoadThemeFromPath(filePath string) (string, error) {
	// Expand ~ to home directory
	if strings.HasPrefix(filePath, "~") {
		home, err := os.UserHomeDir()
		if err != nil {
			return "", fmt.Errorf("failed to get home directory: %w", err)
		}
		filePath = filepath.Join(home, filePath[1:])
	}

	// Make path absolute if it's relative
	absPath, err := filepath.Abs(filePath)
	if err != nil {
		return "", fmt.Errorf("failed to resolve theme path: %w", err)
	}

	// Read the file
	content, err := os.ReadFile(absPath)
	if err != nil {
		return "", fmt.Errorf("failed to read theme file %s: %w", filePath, err)
	}

	// Parse metadata if present
	_, css, err := ParseMetadata(string(content))
	if err != nil {
		// Continue even if metadata parsing fails, use full content
		css = string(content)
	}

	// Validate CSS
	if err := ValidateCSS(css); err != nil {
		return "", fmt.Errorf("theme validation failed for %s: %w", filePath, err)
	}

	return css, nil
}

// ValidateTheme validates a theme CSS file for correctness.
func (l *Loader) ValidateTheme(filePath string) error {
	// Read the file
	content, err := os.ReadFile(filePath)
	if err != nil {
		return fmt.Errorf("failed to read theme file: %w", err)
	}

	// Parse metadata
	_, css, err := ParseMetadata(string(content))
	if err != nil {
		return fmt.Errorf("failed to parse metadata: %w", err)
	}

	// If no CSS was extracted, that's an error
	if strings.TrimSpace(css) == "" {
		return fmt.Errorf("theme file contains no CSS content")
	}

	// Validate CSS syntax
	if err := ValidateCSS(css); err != nil {
		return fmt.Errorf("CSS validation failed: %w", err)
	}

	return nil
}

// loadBuiltInThemeCSS loads CSS from the embedded themes (themes/embed.go).
// This is a simplified implementation; actual themes need to be exported from themes package.
func (l *Loader) loadBuiltInThemeCSS(themeName string) string {
	// Map of theme names to embedded CSS content
	// In a real implementation, this would import from themes package
	switch themeName {
	case "default":
		return `
body { font-family: "Segoe UI", Tahoma, sans-serif; }
h1 { color: #2c3e50; border-bottom: 3px solid #3498db; }
h2 { color: #2c3e50; border-bottom: 2px solid #bdc3c7; }
code { background-color: #f4f4f4; padding: 2px 6px; border-radius: 3px; }
`
	case "dark":
		return `
body { font-family: "Segoe UI", Tahoma, sans-serif; color: #e0e0e0; background-color: #1e1e1e; }
h1 { color: #64b5f6; border-bottom: 3px solid #64b5f6; }
h2 { color: #64b5f6; border-bottom: 2px solid #424242; }
code { background-color: #2d2d2d; color: #81c784; }
`
	case "academic":
		return `
body { font-family: "Times New Roman", Times, serif; }
h1 { font-size: 18pt; text-align: center; border-bottom: 1px solid #000; }
h2 { font-size: 14pt; border-bottom: 1px solid #000; }
code { font-family: "Courier New", monospace; }
`
	}
	return ""
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

// EnsureThemesDir ensures the user themes directory exists.
// Creates the directory with standard permissions (0755) if it doesn't exist.
// Returns the absolute path to the themes directory and any error encountered.
func (l *Loader) EnsureThemesDir() (string, error) {
	// Make sure userThemesDir is set
	if l.userThemesDir == "" {
		return "", fmt.Errorf("themes directory path not configured")
	}

	// Create directory if it doesn't exist
	if err := os.MkdirAll(l.userThemesDir, 0o755); err != nil {
		return "", fmt.Errorf("failed to create themes directory %s: %w", l.userThemesDir, err)
	}

	// Return absolute path
	absPath, err := filepath.Abs(l.userThemesDir)
	if err != nil {
		return "", fmt.Errorf("failed to resolve themes directory path: %w", err)
	}

	return absPath, nil
}
