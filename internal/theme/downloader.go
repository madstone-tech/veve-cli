package theme

import (
	"archive/zip"
	"fmt"
	"io"
	"net/http"
	"net/url"
	"os"
	"path/filepath"
	"strings"
	"time"
)

// Downloader handles downloading and extracting theme files from URLs or local paths.
type Downloader struct {
	timeout time.Duration
}

// NewDownloader creates a new downloader with default timeout.
func NewDownloader() *Downloader {
	return &Downloader{
		timeout: 30 * time.Second,
	}
}

// Download downloads a theme from a URL or local file path.
// Returns the CSS content if successful.
func (d *Downloader) Download(source string) (string, error) {
	// Check if it's a local file path
	if !isURL(source) {
		return d.downloadFromFile(source)
	}

	// Validate URL
	if err := validateURL(source); err != nil {
		return "", err
	}

	// Determine file type from URL
	if strings.HasSuffix(strings.ToLower(source), ".zip") {
		return d.downloadAndExtractZip(source)
	}

	// Default to CSS file
	return d.downloadCSSFile(source)
}

// isURL checks if a string looks like a URL.
func isURL(source string) bool {
	return strings.HasPrefix(source, "http://") || strings.HasPrefix(source, "https://")
}

// validateURL validates that a URL is properly formatted and uses HTTPS.
func validateURL(urlStr string) error {
	if urlStr == "" {
		return fmt.Errorf("URL cannot be empty")
	}

	// Parse URL
	u, err := url.Parse(urlStr)
	if err != nil {
		return fmt.Errorf("invalid URL: %w", err)
	}

	// Require HTTPS for security
	if u.Scheme != "https" {
		return fmt.Errorf("only HTTPS URLs are supported for security (got %s)", u.Scheme)
	}

	// Validate that it's a proper URL with a host
	if u.Host == "" {
		return fmt.Errorf("URL must include a host")
	}

	return nil
}

// downloadFromFile downloads a theme from a local file path.
func (d *Downloader) downloadFromFile(filePath string) (string, error) {
	// Expand ~ to home directory
	if strings.HasPrefix(filePath, "~") {
		home, err := os.UserHomeDir()
		if err != nil {
			return "", fmt.Errorf("failed to expand home directory: %w", err)
		}
		filePath = filepath.Join(home, filePath[1:])
	}

	// Resolve to absolute path
	absPath, err := filepath.Abs(filePath)
	if err != nil {
		return "", fmt.Errorf("failed to resolve file path: %w", err)
	}

	// Check if file exists
	if _, err := os.Stat(absPath); err != nil {
		return "", fmt.Errorf("file not found: %s", filePath)
	}

	// Read the file
	content, err := os.ReadFile(absPath)
	if err != nil {
		return "", fmt.Errorf("failed to read file: %w", err)
	}

	// Parse metadata and return CSS
	_, css, err := ParseMetadata(string(content))
	if err != nil {
		// If metadata parsing fails, treat entire content as CSS
		css = string(content)
	}

	return css, nil
}

// downloadCSSFile downloads a single CSS file from a URL.
func (d *Downloader) downloadCSSFile(urlStr string) (string, error) {
	client := &http.Client{
		Timeout: d.timeout,
	}

	resp, err := client.Get(urlStr)
	if err != nil {
		return "", fmt.Errorf("failed to download file: %w", err)
	}
	defer resp.Body.Close()

	// Check for HTTP errors
	if resp.StatusCode != http.StatusOK {
		return "", fmt.Errorf("download failed with status %d: %s", resp.StatusCode, http.StatusText(resp.StatusCode))
	}

	// Read content
	content, err := io.ReadAll(resp.Body)
	if err != nil {
		return "", fmt.Errorf("failed to read downloaded content: %w", err)
	}

	// Validate it looks like CSS
	contentStr := string(content)
	if err := ValidateCSS(contentStr); err != nil {
		return "", fmt.Errorf("downloaded file doesn't appear to be valid CSS: %w", err)
	}

	// Parse metadata and return CSS
	_, css, err := ParseMetadata(contentStr)
	if err != nil {
		// If metadata parsing fails, treat entire content as CSS
		css = contentStr
	}

	return css, nil
}

// downloadAndExtractZip downloads a zip file and extracts the first CSS/LaTeX file found.
func (d *Downloader) downloadAndExtractZip(urlStr string) (string, error) {
	client := &http.Client{
		Timeout: d.timeout,
	}

	resp, err := client.Get(urlStr)
	if err != nil {
		return "", fmt.Errorf("failed to download zip file: %w", err)
	}
	defer resp.Body.Close()

	if resp.StatusCode != http.StatusOK {
		return "", fmt.Errorf("download failed with status %d", resp.StatusCode)
	}

	// Read entire response into memory
	zipContent, err := io.ReadAll(resp.Body)
	if err != nil {
		return "", fmt.Errorf("failed to read zip file: %w", err)
	}

	// Extract CSS or LaTeX files from zip
	return d.extractFromZip(zipContent)
}

// extractFromZip extracts theme content from a zip archive.
func (d *Downloader) extractFromZip(zipContent []byte) (string, error) {
	// Create a reader for the zip content
	reader := strings.NewReader(string(zipContent))
	zr, err := zip.NewReader(reader, int64(len(zipContent)))
	if err != nil {
		return "", fmt.Errorf("invalid zip file: %w", err)
	}

	// Look for CSS, LaTeX, or Markdown files
	validExtensions := map[string]bool{
		".css": true,
		".tex": true,
		".md":  true,
	}

	for _, file := range zr.File {
		// Skip directories
		if file.FileInfo().IsDir() {
			continue
		}

		// Check extension
		ext := filepath.Ext(strings.ToLower(file.Name))
		if !validExtensions[ext] {
			continue
		}

		// Read file content
		f, err := file.Open()
		if err != nil {
			continue
		}

		content, err := io.ReadAll(f)
		f.Close()

		if err != nil {
			continue
		}

		contentStr := string(content)

		// For CSS, validate it
		if ext == ".css" {
			if err := ValidateCSS(contentStr); err != nil {
				continue
			}

			// Parse metadata and return CSS
			_, css, err := ParseMetadata(contentStr)
			if err != nil {
				css = contentStr
			}
			return css, nil
		}

		// For LaTeX and Markdown, return as-is after validation
		if ext == ".tex" {
			if err := ValidateLaTeX(contentStr); err != nil {
				continue
			}
			return contentStr, nil
		}

		// For Markdown, just return it
		if ext == ".md" {
			return contentStr, nil
		}
	}

	return "", fmt.Errorf("no valid CSS, LaTeX, or Markdown files found in zip archive")
}

// ValidateFileContent checks if content looks like valid CSS, LaTeX, or Markdown.
func ValidateFileContent(content string) error {
	if strings.TrimSpace(content) == "" {
		return fmt.Errorf("file content is empty")
	}

	// Try CSS validation
	if err := ValidateCSS(content); err == nil {
		return nil
	}

	// Try LaTeX validation
	if err := ValidateLaTeX(content); err == nil {
		return nil
	}

	// For Markdown, minimal validation - just check it's not binary
	if !isBinary(content) {
		return nil
	}

	return fmt.Errorf("content doesn't appear to be CSS, LaTeX, or Markdown")
}

// isBinary checks if content appears to be binary data.
func isBinary(content string) bool {
	// Look for null bytes which indicate binary data
	return strings.ContainsRune(content, '\x00')
}
