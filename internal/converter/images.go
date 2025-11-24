package converter

import (
	"context"
	"fmt"
	"io"
	"math"
	"math/rand"
	"net"
	"net/http"
	"os"
	"regexp"
	"strings"
	"sync"
	"time"
)

// ImageProcessor handles downloading remote images and processing markdown content.
type ImageProcessor struct {
	// Core fields
	tempDir    string
	imageMap   map[string]string // URL -> local path mapping
	httpClient *http.Client

	// Configuration fields
	maxConcurrentDownloads int
	maxBytesPerSession     int64
	timeoutSeconds         int
	maxRetries             int

	// Runtime state
	downloadErrors       map[string]string // URL -> error message
	totalBytesDownloaded int64
	mu                   sync.Mutex // Protects shared state: imageMap, downloadErrors, totalBytesDownloaded
}

// NewImageProcessor creates a new ImageProcessor instance.
// tempDir is the directory where downloaded images will be stored.
func NewImageProcessor(tempDir string) *ImageProcessor {
	return &ImageProcessor{
		tempDir:                tempDir,
		imageMap:               make(map[string]string),
		downloadErrors:         make(map[string]string),
		httpClient:             &http.Client{}, // Per-request timeout will be set in context
		maxConcurrentDownloads: 5,
		maxBytesPerSession:     500 * 1024 * 1024, // 500MB per spec
		timeoutSeconds:         10,                // Per request timeout
		maxRetries:             3,                 // Per spec
	}
}

// WithTimeoutSeconds sets custom timeout for image downloads.
func (ip *ImageProcessor) WithTimeoutSeconds(seconds int) *ImageProcessor {
	if seconds > 0 {
		ip.timeoutSeconds = seconds
	}
	return ip
}

// WithMaxRetries sets custom max retries.
func (ip *ImageProcessor) WithMaxRetries(retries int) *ImageProcessor {
	if retries >= 0 {
		ip.maxRetries = retries
	}
	return ip
}

// ============================================================================
// PHASE 2 FOUNDATIONAL FUNCTIONS
// ============================================================================

// isRemoteURL checks if a URL is a remote HTTP(S) URL.
func isRemoteURL(imageURL string) bool {
	lowerURL := strings.ToLower(imageURL)
	return strings.HasPrefix(lowerURL, "http://") || strings.HasPrefix(lowerURL, "https://")
}

// IsRemoteURL checks if a URL is a remote HTTP(S) URL. Public version for testing.
func (ip *ImageProcessor) IsRemoteURL(imageURL string) bool {
	return isRemoteURL(imageURL)
}

// hashURL creates a simple hash from the URL string.
// This is not cryptographically secure but sufficient for filename uniqueness.
func hashURL(imageURL string) string {
	h := 0
	for i, c := range imageURL {
		h = h*31 + int(c)
		// Keep it manageable by modulo
		if i%10 == 0 {
			h = h % 1000000
		}
	}
	return fmt.Sprintf("%x", h)
}

// isTransientError checks if an error is transient (should be retried)
// or permanent (should not be retried).
func isTransientError(err error, statusCode int) bool {
	// Network timeouts are transient
	if netErr, ok := err.(net.Error); ok && netErr.Timeout() {
		return true
	}

	// Check for transient HTTP status codes
	return statusCode == 408 || statusCode == 429 || statusCode == 503 || statusCode == 504
}

// validateHTTPRequest validates an HTTP request and response.
// Checks status code and content type.
func validateHTTPRequest(resp *http.Response) error {
	// Check HTTP status
	if resp.StatusCode != http.StatusOK {
		return fmt.Errorf("HTTP %d %s", resp.StatusCode, resp.Status)
	}

	// Validate content type
	contentType := resp.Header.Get("Content-Type")
	if !isImageContentType(contentType) {
		return fmt.Errorf("invalid content type: %s (expected image/*)", contentType)
	}

	return nil
}

// isImageContentType checks if the content type is an image type.
func isImageContentType(contentType string) bool {
	// Extract base type (before semicolon if present)
	baseType := strings.Split(contentType, ";")[0]
	baseType = strings.TrimSpace(baseType)
	lowerType := strings.ToLower(baseType)
	return strings.HasPrefix(lowerType, "image/")
}

// IsImageContentType checks if the content type is an image type. Public version for testing.
func (ip *ImageProcessor) IsImageContentType(contentType string) bool {
	return isImageContentType(contentType)
}

// ValidateImageSize checks if the image size is within limits.
// Returns error if size exceeds per-image limit (100MB).
func (ip *ImageProcessor) ValidateImageSize(contentLength int64) error {
	const maxImageSize = 100 * 1024 * 1024 // 100MB per image
	if contentLength > maxImageSize {
		return fmt.Errorf("image too large: %d bytes (max %d)", contentLength, maxImageSize)
	}

	// Check session limit
	ip.mu.Lock()
	defer ip.mu.Unlock()
	if ip.totalBytesDownloaded+contentLength > ip.maxBytesPerSession {
		return fmt.Errorf("session size limit exceeded: %d + %d > %d",
			ip.totalBytesDownloaded, contentLength, ip.maxBytesPerSession)
	}

	return nil
}

// generateFileName generates a unique filename for the downloaded image.
// Uses the URL hash to create a unique name and appends the appropriate extension.
func generateFileName(imageURL string, contentType string) string {
	ext := getExtensionFromContentType(contentType)
	hash := hashURL(imageURL)
	return fmt.Sprintf("veve-image-%s%s", hash, ext)
}

// getExtensionFromContentType returns the file extension based on content type.
func getExtensionFromContentType(contentType string) string {
	contentType = strings.Split(contentType, ";")[0] // Remove charset info
	contentType = strings.TrimSpace(contentType)

	switch contentType {
	case "image/jpeg", "image/jpg":
		return ".jpg"
	case "image/png":
		return ".png"
	case "image/gif":
		return ".gif"
	case "image/webp":
		return ".webp"
	case "image/svg+xml":
		return ".svg"
	case "image/bmp":
		return ".bmp"
	case "image/tiff":
		return ".tiff"
	default:
		return ".img" // fallback extension
	}
}

// GetExtensionFromContentType returns file extension based on content type. Public for testing.
func (ip *ImageProcessor) GetExtensionFromContentType(contentType string) string {
	return getExtensionFromContentType(contentType)
}

// ============================================================================
// MARKDOWN PROCESSING (T008)
// ============================================================================

// DetectRemoteImages extracts all remote image URLs from markdown content.
// Returns a list of unique remote URLs, ignoring duplicates and local paths.
func (ip *ImageProcessor) DetectRemoteImages(content string) []string {
	imageRegex := regexp.MustCompile(`!\[([^\]]*)\]\(([^)]+)\)`)
	matches := imageRegex.FindAllStringSubmatch(content, -1)

	seen := make(map[string]bool)
	var urls []string

	for _, match := range matches {
		if len(match) < 3 {
			continue
		}
		imageURL := match[2]

		// Only include remote URLs, avoid duplicates
		if isRemoteURL(imageURL) && !seen[imageURL] {
			urls = append(urls, imageURL)
			seen[imageURL] = true
		}
	}

	return urls
}

// ============================================================================
// CONCURRENCY & CLEANUP INFRASTRUCTURE (T009, T010)
// ============================================================================

// GetImageMap returns the mapping of downloaded image URLs to local paths.
// Useful for testing and logging. Returns a copy to prevent external modification.
func (ip *ImageProcessor) GetImageMap() map[string]string {
	ip.mu.Lock()
	defer ip.mu.Unlock()

	result := make(map[string]string)
	for k, v := range ip.imageMap {
		result[k] = v
	}
	return result
}

// SetImageMap sets a URL-to-local-path mapping. Used for testing.
func (ip *ImageProcessor) SetImageMap(url, localPath string) {
	ip.mu.Lock()
	defer ip.mu.Unlock()
	ip.imageMap[url] = localPath
}

// RecordDownload records bytes downloaded. Used for testing size limits.
func (ip *ImageProcessor) RecordDownload(bytes int64) {
	ip.mu.Lock()
	defer ip.mu.Unlock()
	ip.totalBytesDownloaded += bytes
}

// GetDownloadErrors returns the mapping of failed image URLs to error messages.
func (ip *ImageProcessor) GetDownloadErrors() map[string]string {
	ip.mu.Lock()
	defer ip.mu.Unlock()

	result := make(map[string]string)
	for k, v := range ip.downloadErrors {
		result[k] = v
	}
	return result
}

// Cleanup removes all temporary image files created by this processor.
// Uses best-effort approach: logs warnings but doesn't block on errors.
func (ip *ImageProcessor) Cleanup() error {
	ip.mu.Lock()
	imagesToClean := make([]string, 0, len(ip.imageMap))
	for _, localPath := range ip.imageMap {
		imagesToClean = append(imagesToClean, localPath)
	}
	ip.mu.Unlock()

	// Remove all downloaded image files
	for _, localPath := range imagesToClean {
		if err := os.Remove(localPath); err != nil && !os.IsNotExist(err) {
			fmt.Fprintf(os.Stderr, "[WARN] Failed to remove temp image file %s: %v\n", localPath, err)
		}
	}

	// Try to remove the temp directory itself
	if err := os.RemoveAll(ip.tempDir); err != nil && !os.IsNotExist(err) {
		fmt.Fprintf(os.Stderr, "[WARN] Failed to remove temp image directory: %v\n", err)
	}

	// Always return nil: cleanup failures don't block conversion
	return nil
}

// ============================================================================
// HELPER FUNCTIONS FOR CONCURRENCY & RETRY
// ============================================================================

// calculateBackoff calculates exponential backoff with jitter.
// Returns seconds to wait before retrying.
// Formula: random(0, min(10, 2^attempt))
func (ip *ImageProcessor) calculateBackoff(attempt int) float64 {
	baseBackoff := math.Min(10, math.Pow(2, float64(attempt)))
	jitteredWait := rand.Float64() * baseBackoff
	return jitteredWait
}

// ============================================================================
// PLACEHOLDER FUNCTIONS (implemented in later phases)
// ============================================================================

// ProcessMarkdown processes markdown content to download remote images and replace URLs.
// Returns the processed markdown content with local image paths.
// PLACEHOLDER: Full implementation in Phase 3 with concurrent downloads
func (ip *ImageProcessor) ProcessMarkdown(content string) (string, error) {
	// Create temp directory if it doesn't exist
	if err := os.MkdirAll(ip.tempDir, 0o755); err != nil {
		return "", fmt.Errorf("failed to create temp directory for images: %w", err)
	}

	// Phase 3 will add full concurrent download logic here
	return content, nil
}

// DownloadImageOnce downloads a single image without retries.
// Returns the local file path where the image was saved.
// If the image is already cached in imageMap, returns the cached path immediately.
func (ip *ImageProcessor) DownloadImageOnce(imageURL string) (string, error) {
	// Check cache first
	ip.mu.Lock()
	if cachedPath, exists := ip.imageMap[imageURL]; exists {
		ip.mu.Unlock()
		return cachedPath, nil
	}
	ip.mu.Unlock()

	// Create context with timeout
	ctx, cancel := context.WithTimeout(context.Background(), time.Duration(ip.timeoutSeconds)*time.Second)
	defer cancel()

	// Create HTTP request
	req, err := http.NewRequestWithContext(ctx, "GET", imageURL, nil)
	if err != nil {
		return "", fmt.Errorf("failed to create request for %s: %w", imageURL, err)
	}

	// Execute request
	resp, err := ip.httpClient.Do(req)
	if err != nil {
		return "", fmt.Errorf("failed to download %s: %w", imageURL, err)
	}
	defer resp.Body.Close()

	// Validate response
	if err := validateHTTPRequest(resp); err != nil {
		return "", fmt.Errorf("invalid HTTP response from %s: %w", imageURL, err)
	}

	// Validate size
	contentLength := resp.ContentLength
	if contentLength == -1 {
		// If Content-Length is not set, we need to read the body to determine size
		// For now, allow it and check during write
		contentLength = 0
	}
	if contentLength > 0 {
		if err := ip.ValidateImageSize(contentLength); err != nil {
			return "", fmt.Errorf("image size validation failed for %s: %w", imageURL, err)
		}
	}

	// Generate filename and create temp file
	fileName := generateFileName(imageURL, resp.Header.Get("Content-Type"))
	tempFile, err := os.CreateTemp(ip.tempDir, fileName)
	if err != nil {
		return "", fmt.Errorf("failed to create temp file for %s: %w", imageURL, err)
	}
	defer tempFile.Close()

	// Copy response body to file with size tracking
	writtenBytes, err := io.Copy(tempFile, resp.Body)
	if err != nil {
		// Clean up failed download
		os.Remove(tempFile.Name())
		return "", fmt.Errorf("failed to write image from %s: %w", imageURL, err)
	}

	// Validate size after download if not provided in header
	if contentLength == 0 {
		if err := ip.ValidateImageSize(writtenBytes); err != nil {
			os.Remove(tempFile.Name())
			return "", fmt.Errorf("image too large from %s: %w", imageURL, err)
		}
	}

	localPath := tempFile.Name()

	// Update state
	ip.mu.Lock()
	ip.imageMap[imageURL] = localPath
	ip.totalBytesDownloaded += writtenBytes
	ip.mu.Unlock()

	return localPath, nil
}

// downloadWithRetry downloads an image with retry logic.
// PLACEHOLDER: Full implementation in Phase 4
func (ip *ImageProcessor) downloadWithRetry(imageURL string) error {
	return fmt.Errorf("not implemented in Phase 2 (foundation)")
}

// RewriteMarkdownImageURLs rewrites markdown image references to use local paths.
// For each markdown image ![alt](url), if url is in the imageMap, replaces it with the local path.
// Otherwise, leaves the original URL unchanged.
func (ip *ImageProcessor) RewriteMarkdownImageURLs(content string) string {
	// Regex to match markdown image syntax: ![alt text](url)
	imageRegex := regexp.MustCompile(`!\[([^\]]*)\]\(([^)]+)\)`)

	// Get a snapshot of the image map
	ip.mu.Lock()
	imageMapSnapshot := make(map[string]string)
	for k, v := range ip.imageMap {
		imageMapSnapshot[k] = v
	}
	ip.mu.Unlock()

	// Replace matched URLs with local paths if available
	result := imageRegex.ReplaceAllStringFunc(content, func(match string) string {
		// Extract the URL from the match
		submatches := imageRegex.FindStringSubmatch(match)
		if len(submatches) < 3 {
			return match
		}

		altText := submatches[1]
		imageURL := submatches[2]

		// Check if we have a local path for this URL
		if localPath, exists := imageMapSnapshot[imageURL]; exists {
			return fmt.Sprintf("![%s](%s)", altText, localPath)
		}

		// Return original if not in map
		return match
	})

	return result
}
