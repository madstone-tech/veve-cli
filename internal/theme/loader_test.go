package theme

import (
	"os"
	"path/filepath"
	"testing"
	"time"
)

// TestDiscoverBuiltInThemes tests that all built-in themes are discovered.
func TestDiscoverBuiltInThemes(t *testing.T) {
	loader := NewLoader("")

	err := loader.DiscoverThemes()
	if err != nil {
		t.Fatalf("DiscoverThemes failed: %v", err)
	}

	expectedThemes := []string{"default", "dark", "academic"}
	for _, name := range expectedThemes {
		theme, exists := loader.GetRegistry().GetTheme(name)
		if !exists {
			t.Errorf("expected theme %s not found", name)
		}
		if !theme.IsBuiltIn {
			t.Errorf("theme %s should be marked as built-in", name)
		}
	}
}

// TestLoadBuiltInTheme tests loading CSS for a built-in theme.
func TestLoadBuiltInTheme(t *testing.T) {
	loader := NewLoader("")

	tests := []struct {
		name          string
		shouldHaveCSS bool
	}{
		{"default", true},
		{"dark", true},
		{"academic", true},
		{"nonexistent", false},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			css, err := loader.LoadThemeCSS(tt.name)
			if tt.shouldHaveCSS {
				if err != nil && !tt.shouldHaveCSS {
					t.Fatalf("LoadThemeCSS failed: %v", err)
				}
				if css == "" && tt.shouldHaveCSS {
					t.Fatal("expected CSS content but got empty string")
				}
				// Verify it's valid CSS-like content
				if len(css) < 10 {
					t.Errorf("CSS too short: %d bytes", len(css))
				}
			} else {
				if err == nil && !tt.shouldHaveCSS {
					t.Fatal("expected error but got nil")
				}
			}
		})
	}
}

// TestLoadTheme tests loading a theme by name.
func TestLoadTheme(t *testing.T) {
	loader := NewLoader("")
	loader.DiscoverThemes()

	tests := []struct {
		themeName   string
		shouldExist bool
	}{
		{"default", true},
		{"dark", true},
		{"academic", true},
		{"nonexistent", false},
		{"MyTheme", false},
	}

	for _, tt := range tests {
		t.Run(tt.themeName, func(t *testing.T) {
			theme, err := loader.LoadTheme(tt.themeName)
			if tt.shouldExist {
				if err != nil {
					t.Fatalf("LoadTheme failed: %v", err)
				}
				if theme.Name != tt.themeName {
					t.Errorf("expected theme name %s, got %s", tt.themeName, theme.Name)
				}
			} else {
				if err == nil {
					t.Fatal("expected error but got nil")
				}
			}
		})
	}
}

// TestListThemes tests listing all available themes.
func TestListThemes(t *testing.T) {
	loader := NewLoader("")
	loader.DiscoverThemes()

	themes := loader.ListThemes()

	// Should have at least the 3 built-in themes
	if len(themes) < 3 {
		t.Fatalf("expected at least 3 themes, got %d", len(themes))
	}

	// Themes should be sorted by name
	for i := 1; i < len(themes); i++ {
		if themes[i].Name < themes[i-1].Name {
			t.Errorf("themes not sorted: %s before %s", themes[i].Name, themes[i-1].Name)
		}
	}

	// Check that built-in themes are present
	builtInFound := map[string]bool{
		"default":  false,
		"dark":     false,
		"academic": false,
	}

	for _, theme := range themes {
		if _, ok := builtInFound[theme.Name]; ok {
			builtInFound[theme.Name] = true
		}
	}

	for name, found := range builtInFound {
		if !found {
			t.Errorf("built-in theme %s not found in list", name)
		}
	}
}

// TestDiscoverUserThemes tests discovering user-installed themes.
func TestDiscoverUserThemes(t *testing.T) {
	// Create temporary directory for user themes
	tmpDir := t.TempDir()

	// Create a dummy CSS file
	themePath := filepath.Join(tmpDir, "custom.css")
	css := `body { color: red; }`
	if err := os.WriteFile(themePath, []byte(css), 0o644); err != nil {
		t.Fatalf("failed to create test theme: %v", err)
	}

	// Create loader and discover themes
	loader := NewLoader(tmpDir)
	if err := loader.DiscoverThemes(); err != nil {
		t.Fatalf("DiscoverThemes failed: %v", err)
	}

	// Check that custom theme was found
	customTheme, exists := loader.GetRegistry().GetTheme("custom")
	if !exists {
		t.Fatal("custom theme not discovered")
	}

	if customTheme.IsBuiltIn {
		t.Error("custom theme should not be marked as built-in")
	}

	if customTheme.FilePath != themePath {
		t.Errorf("expected file path %s, got %s", themePath, customTheme.FilePath)
	}
}

// TestLoadUserThemeCSS tests loading CSS from a user theme file.
func TestLoadUserThemeCSS(t *testing.T) {
	tmpDir := t.TempDir()

	// Create a CSS file
	css := `body { color: blue; font-size: 14pt; }`
	themePath := filepath.Join(tmpDir, "myTheme.css")
	if err := os.WriteFile(themePath, []byte(css), 0o644); err != nil {
		t.Fatalf("failed to create test theme: %v", err)
	}

	// Create loader and discover
	loader := NewLoader(tmpDir)
	loader.DiscoverThemes()

	// Load the CSS
	loadedCSS, err := loader.LoadThemeCSS("myTheme")
	if err != nil {
		t.Fatalf("LoadThemeCSS failed: %v", err)
	}

	if loadedCSS != css {
		t.Errorf("expected CSS %q, got %q", css, loadedCSS)
	}
}

// TestThemeMetadata tests that theme metadata is properly populated.
func TestThemeMetadata(t *testing.T) {
	loader := NewLoader("")
	loader.DiscoverThemes()

	// Check default theme
	defaultTheme, exists := loader.GetRegistry().GetTheme("default")
	if !exists {
		t.Fatal("default theme not found")
	}

	if defaultTheme.DisplayName != "Default" {
		t.Errorf("expected DisplayName 'Default', got %q", defaultTheme.DisplayName)
	}

	if defaultTheme.Author != "veve-cli" {
		t.Errorf("expected Author 'veve-cli', got %q", defaultTheme.Author)
	}

	if len(defaultTheme.Description) == 0 {
		t.Error("expected Description to be populated")
	}

	if defaultTheme.Version == "" {
		t.Error("expected Version to be populated")
	}
}

// TestLoadThemeNotFound tests behavior when theme doesn't exist.
func TestLoadThemeNotFound(t *testing.T) {
	loader := NewLoader("")
	loader.DiscoverThemes()

	_, err := loader.LoadTheme("nonexistent-xyz")
	if err == nil {
		t.Fatal("expected error for non-existent theme")
	}
}

// TestThemeCreatedAt tests that theme metadata includes timestamps.
func TestThemeMetadataStructure(t *testing.T) {
	loader := NewLoader("")
	loader.DiscoverThemes()

	themes := loader.ListThemes()
	if len(themes) == 0 {
		t.Fatal("no themes discovered")
	}

	// Check first theme has proper structure
	theme := themes[0]
	if theme.Name == "" {
		t.Error("theme Name is empty")
	}
	if theme.DisplayName == "" {
		t.Error("theme DisplayName is empty")
	}
	if theme.Author == "" {
		t.Error("theme Author is empty")
	}

	// For built-in themes, CreatedAt should be set
	if theme.IsBuiltIn && theme.CreatedAt.IsZero() {
		t.Logf("warning: CreatedAt is zero for built-in theme %s", theme.Name)
	}
}

// TestConcurrentThemeDiscovery tests that theme discovery is thread-safe.
func TestConcurrentThemeDiscovery(t *testing.T) {
	tmpDir := t.TempDir()

	// Create a CSS file
	themePath := filepath.Join(tmpDir, "test.css")
	if err := os.WriteFile(themePath, []byte("body { }"), 0o644); err != nil {
		t.Fatalf("failed to create test theme: %v", err)
	}

	// Run discovery concurrently
	done := make(chan bool, 2)
	for i := 0; i < 2; i++ {
		go func() {
			loader := NewLoader(tmpDir)
			_ = loader.DiscoverThemes()
			themes := loader.ListThemes()
			if len(themes) < 4 { // 3 built-in + 1 custom
				t.Errorf("expected at least 4 themes, got %d", len(themes))
			}
			done <- true
		}()
	}

	// Wait for goroutines
	for i := 0; i < 2; i++ {
		<-done
	}
}

// TestThemeTimestamp tests that theme timestamps are reasonable.
func TestThemeTimestampFormat(t *testing.T) {
	loader := NewLoader("")
	loader.DiscoverThemes()

	themes := loader.ListThemes()
	now := time.Now()

	for _, theme := range themes {
		// CreatedAt should be in the past (or now)
		if theme.CreatedAt.After(now.Add(time.Second)) {
			t.Errorf("theme %s has future timestamp: %v", theme.Name, theme.CreatedAt)
		}
	}
}
