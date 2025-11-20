package theme

import (
	"strings"
	"testing"
)

// TestParseMetadataWithYAML tests parsing YAML front matter from CSS.
func TestParseMetadataWithYAML(t *testing.T) {
	content := `---
name: myTheme
author: John Doe
description: A custom theme
version: 1.0.0
---
body { color: blue; }
h1 { font-size: 24pt; }
`

	meta, css, _ := ParseMetadata(content)
	_ = css // css is checked below

	if meta == nil {
		t.Fatal("expected metadata, got nil")
	}

	if meta.Name != "myTheme" {
		t.Errorf("expected name 'myTheme', got '%s'", meta.Name)
	}

	if meta.Author != "John Doe" {
		t.Errorf("expected author 'John Doe', got '%s'", meta.Author)
	}

	if meta.Description != "A custom theme" {
		t.Errorf("expected description 'A custom theme', got '%s'", meta.Description)
	}

	if meta.Version != "1.0.0" {
		t.Errorf("expected version '1.0.0', got '%s'", meta.Version)
	}

	// Check CSS content
	if !contains(css, "color: blue") {
		t.Errorf("CSS missing 'color: blue'")
	}

	if !contains(css, "font-size: 24pt") {
		t.Errorf("CSS missing 'font-size: 24pt'")
	}
}

// TestParseMetadataWithoutYAML tests parsing CSS without YAML front matter.
func TestParseMetadataWithoutYAML(t *testing.T) {
	content := `body { color: red; }
h1 { font-weight: bold; }
`

	meta, css, err := ParseMetadata(content)

	if err != nil {
		t.Fatalf("ParseMetadata failed: %v", err)
	}

	if meta != nil {
		t.Fatal("expected nil metadata for CSS without front matter")
	}

	if css != content {
		t.Errorf("expected CSS to be unchanged")
	}
}

// TestParseMetadataPartial tests parsing with some metadata fields missing.
func TestParseMetadataPartial(t *testing.T) {
	content := `---
name: partial
author: Jane Smith
---
body { }
`

	meta, _, err := ParseMetadata(content)

	if err != nil {
		t.Fatalf("ParseMetadata failed: %v", err)
	}

	if meta == nil {
		t.Fatal("expected metadata")
	}

	if meta.Name != "partial" {
		t.Errorf("expected name 'partial', got '%s'", meta.Name)
	}

	if meta.Author != "Jane Smith" {
		t.Errorf("expected author 'Jane Smith', got '%s'", meta.Author)
	}

	// Empty fields should be empty string
	if meta.Description != "" {
		t.Errorf("expected empty description, got '%s'", meta.Description)
	}

	if meta.Version != "" {
		t.Errorf("expected empty version, got '%s'", meta.Version)
	}
}

// TestParseMetadataWithQuotedValues tests parsing quoted values in YAML.
func TestParseMetadataQuotedValues(t *testing.T) {
	content := `---
name: "quoted-theme"
author: 'Single Quoted Author'
description: "A description with spaces"
---
body { }
`

	meta, _, err := ParseMetadata(content)

	if err != nil {
		t.Fatalf("ParseMetadata failed: %v", err)
	}

	if meta.Name != "quoted-theme" {
		t.Errorf("expected unquoted name 'quoted-theme', got '%s'", meta.Name)
	}

	if meta.Author != "Single Quoted Author" {
		t.Errorf("expected unquoted author, got '%s'", meta.Author)
	}

	if meta.Description != "A description with spaces" {
		t.Errorf("expected unquoted description, got '%s'", meta.Description)
	}
}

// TestValidateCSS tests CSS validation.
func TestValidateCSS(t *testing.T) {
	tests := []struct {
		name      string
		css       string
		shouldErr bool
	}{
		{"valid CSS", "body { color: blue; }", false},
		{"multiple rules", "body { } h1 { color: red; }", false},
		{"empty content", "", true},
		{"unbalanced braces open", "body { color: blue; ", true},
		{"unbalanced braces close", "body color: blue; }", true},
		{"complex CSS", "body { font-family: Arial, sans-serif; color: #333; }", false},
	}

	for _, tt := range tests {
		tt := tt // Shadow for parallel safe
		t.Run(tt.name, func(t *testing.T) {
			err := ValidateCSS(tt.css)
			if (err != nil) != tt.shouldErr {
				t.Errorf("ValidateCSS error = %v, shouldErr %v", err, tt.shouldErr)
			}
		})
	}
}

// TestApplyMetadataDefaults tests applying defaults to metadata.
func TestApplyMetadataDefaults(t *testing.T) {
	meta := &ThemeMetadata{}
	ApplyMetadataDefaults(meta, "mytheme")

	if meta.Name != "mytheme" {
		t.Errorf("expected name 'mytheme', got '%s'", meta.Name)
	}

	if meta.Author != "Unknown" {
		t.Errorf("expected author 'Unknown', got '%s'", meta.Author)
	}

	if meta.Description != "Custom theme" {
		t.Errorf("expected description 'Custom theme', got '%s'", meta.Description)
	}

	if meta.Version != "1.0.0" {
		t.Errorf("expected version '1.0.0', got '%s'", meta.Version)
	}
}

// TestApplyMetadataDefaultsPreservesExisting tests that existing values are preserved.
func TestApplyMetadataDefaultsPreservesExisting(t *testing.T) {
	meta := &ThemeMetadata{
		Name:        "existing",
		Author:      "Original Author",
		Description: "Original Description",
		Version:     "2.0.0",
	}

	ApplyMetadataDefaults(meta, "override")

	if meta.Name != "existing" {
		t.Errorf("expected existing name 'existing', got '%s'", meta.Name)
	}

	if meta.Author != "Original Author" {
		t.Errorf("expected author preserved")
	}

	if meta.Description != "Original Description" {
		t.Errorf("expected description preserved")
	}

	if meta.Version != "2.0.0" {
		t.Errorf("expected version preserved")
	}
}

// TestParseMetadataEdgeCase_MissingClosingMarker tests malformed YAML.
func TestParseMetadataEdgeCase_MissingClosingMarker(t *testing.T) {
	content := `---
name: unclosed
author: Test
body { color: blue; }
`

	meta, _, _ := ParseMetadata(content)

	// Should treat as no metadata since closing --- not found
	if meta != nil {
		t.Fatal("expected nil metadata for malformed YAML")
	}
}

// TestParseMetadataEdgeCase_EmptyYAML tests empty YAML block.
func TestParseMetadataEdgeCase_EmptyYAML(t *testing.T) {
	content := `---
---
body { }
`

	meta, _, err := ParseMetadata(content)

	if err != nil {
		t.Fatalf("ParseMetadata failed: %v", err)
	}

	if meta == nil {
		t.Fatal("expected metadata object")
	}

	// All fields should be empty
	if meta.Name != "" || meta.Author != "" || meta.Description != "" || meta.Version != "" {
		t.Error("expected all metadata fields to be empty")
	}
}

// TestParseMetadataWithComments tests YAML with comments.
func TestParseMetadataWithComments(t *testing.T) {
	content := `---
# This is a comment
name: commented
# Another comment
author: Author Name
description: Theme description
---
body { }
`

	meta, _, err := ParseMetadata(content)

	if err != nil {
		t.Fatalf("ParseMetadata failed: %v", err)
	}

	if meta.Name != "commented" {
		t.Errorf("expected name 'commented', got '%s'", meta.Name)
	}

	if meta.Author != "Author Name" {
		t.Errorf("expected author 'Author Name', got '%s'", meta.Author)
	}
}

// Helper function to check if string contains substring
func contains(s, substr string) bool {
	return strings.Contains(s, substr)
}
