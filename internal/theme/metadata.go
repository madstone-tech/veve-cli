package theme

import (
	"fmt"
	"strings"
)

// ThemeMetadata represents parsed YAML front matter from a theme CSS file.
type ThemeMetadata struct {
	Name        string
	Author      string
	Description string
	Version     string
}

// ParseMetadata extracts YAML front matter from a CSS file content.
// Format:
//
//	---
//	name: theme-name
//	author: Author Name
//	description: Theme description
//	version: 1.0.0
//	---
//	/* CSS content here */
//
// Returns the metadata and remaining CSS content.
// If no YAML front matter is found, returns nil metadata and the full content as CSS.
func ParseMetadata(content string) (*ThemeMetadata, string, error) {
	lines := strings.Split(content, "\n")

	// Check if content starts with ---
	if len(lines) < 2 || strings.TrimSpace(lines[0]) != "---" {
		// No metadata, return content as-is
		return nil, content, nil
	}

	// Find closing ---
	var endIdx int
	found := false
	for i := 1; i < len(lines); i++ {
		if strings.TrimSpace(lines[i]) == "---" {
			endIdx = i
			found = true
			break
		}
	}

	if !found {
		// Malformed YAML, treat entire content as CSS
		return nil, content, nil
	}

	// Parse YAML front matter
	yamlLines := lines[1:endIdx]
	metadata := &ThemeMetadata{}

	for _, line := range yamlLines {
		line = strings.TrimSpace(line)
		if line == "" || strings.HasPrefix(line, "#") {
			continue
		}

		parts := strings.SplitN(line, ":", 2)
		if len(parts) != 2 {
			continue
		}

		key := strings.TrimSpace(parts[0])
		value := strings.TrimSpace(parts[1])
		// Remove quotes if present
		value = strings.Trim(value, "\"'")

		switch key {
		case "name":
			metadata.Name = value
		case "author":
			metadata.Author = value
		case "description":
			metadata.Description = value
		case "version":
			metadata.Version = value
		}
	}

	// Extract remaining CSS content
	cssLines := lines[endIdx+1:]
	css := strings.Join(cssLines, "\n")
	css = strings.TrimSpace(css)

	return metadata, css, nil
}

// ValidateCSS performs basic validation of CSS content.
// Checks for:
// - Non-empty content
// - Balanced braces
// - No obviously malformed syntax
func ValidateCSS(css string) error {
	if strings.TrimSpace(css) == "" {
		return fmt.Errorf("CSS content is empty")
	}

	// Count braces
	openBraces := strings.Count(css, "{")
	closeBraces := strings.Count(css, "}")

	if openBraces != closeBraces {
		return fmt.Errorf("unbalanced braces: %d { and %d }", openBraces, closeBraces)
	}

	return nil
}

// ValidateLaTeX performs basic validation of LaTeX content (if used in CSS).
// This is a simplified check for common LaTeX patterns.
func ValidateLaTeX(content string) error {
	// Check for common LaTeX issues
	if strings.Count(content, "\\begin") != strings.Count(content, "\\end") {
		return fmt.Errorf("unbalanced LaTeX begin/end blocks")
	}

	return nil
}

// ApplyMetadataDefaults updates theme metadata with defaults from filename if fields are empty.
func ApplyMetadataDefaults(meta *ThemeMetadata, themeName string) {
	if meta.Name == "" {
		meta.Name = themeName
	}
	if meta.Author == "" {
		meta.Author = "Unknown"
	}
	if meta.Description == "" {
		meta.Description = "Custom theme"
	}
	if meta.Version == "" {
		meta.Version = "1.0.0"
	}
}
