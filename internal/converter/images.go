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
	"strconv"
	"strings"
	"sync"
	"time"
)

// ImageProcessor handles downloading remote images and processing markdown content.
// It detects HTTP/HTTPS image URLs in markdown, downloads them concurrently with retry logic,
// and rewrites the markdown to use local file paths. Thread-safe with concurrent download support.
//
// Configuration:
//   - maxConcurrentDownloads: Number of concurrent downloads (default 5)
//   - timeoutSeconds: Timeout per image download in seconds (default 10)
//   - maxRetries: Maximum retry attempts for transient errors (default 3)
//   - maxBytesPerSession: Total download limit per session (default 500MB)
//
// Features:
//   - Automatic detection of remote image URLs in markdown
//   - Concurrent downloads with semaphore pattern
//   - Retry logic with exponential backoff for transient errors
//   - Graceful degradation: failed images don't block conversion
//   - Resource limits: per-image (100MB) and per-session (500MB)
//   - Best-effort cleanup of temporary files
//
// Thread Safety:
//   - All shared state protected by mu (sync.Mutex)
//   - Safe for concurrent use
//
// Example:
//
//	processor := NewImageProcessor(tempDir).
//		WithTimeoutSeconds(15).
//		WithMaxRetries(3)
//	defer processor.Cleanup()
//
//	processedMD, err := processor.ProcessMarkdown(markdown)
//	if err != nil {
//		// Handle error, but continue with best-effort processing
//	}
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

// NewImageProcessor creates a new ImageProcessor instance with default configuration.
//
// Parameters:
//   - tempDir: Directory where downloaded images will be stored. Must be writable.
//
// Default Configuration:
//   - maxConcurrentDownloads: 5
//   - timeoutSeconds: 10
//   - maxRetries: 3
//   - maxBytesPerSession: 500MB
//
// The returned processor can be further configured using:
//   - WithTimeoutSeconds() to set per-request timeout
//   - WithMaxRetries() to set retry attempts
//
// Example:
//
//	processor := NewImageProcessor("/tmp/veve-images").
//		WithTimeoutSeconds(20).
//		WithMaxRetries(5)
//	defer processor.Cleanup()
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

// IsTransientError is the public version for testing.
func (ip *ImageProcessor) IsTransientError(err error, statusCode int) bool {
	return isTransientError(err, statusCode)
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

// GetDownloadStats returns statistics about downloads.
// Returns (successful count, failed count, total attempted).
func (ip *ImageProcessor) GetDownloadStats() (successful, failed, total int) {
	ip.mu.Lock()
	defer ip.mu.Unlock()

	successful = len(ip.imageMap)
	failed = len(ip.downloadErrors)
	total = successful + failed
	return
}

// GetErrorSummary returns a formatted error summary for user output.
// Format: "[WARN] Failed to download N images:\n  - URL1: reason1\n  - URL2: reason2"
func (ip *ImageProcessor) GetErrorSummary() string {
	errors := ip.GetDownloadErrors()
	if len(errors) == 0 {
		return ""
	}

	var sb strings.Builder
	sb.WriteString(fmt.Sprintf("[WARN] Failed to download %d image(s):\n", len(errors)))
	for url, errMsg := range errors {
		// Extract just the relevant part of the error message
		reason := errMsg
		if strings.Contains(errMsg, ":") {
			parts := strings.Split(errMsg, ":")
			reason = strings.TrimSpace(parts[len(parts)-1])
		}
		sb.WriteString(fmt.Sprintf("  - %s\n    Reason: %s\n", url, reason))
	}
	return sb.String()
}

// Cleanup removes all temporary image files created by this processor.
// Uses best-effort approach: logs warnings but doesn't block on errors.
//
// Behavior:
//   - Removes all downloaded image files from imageMap
//   - Removes the temporary directory
//   - Logs warnings to stderr if files can't be deleted
//   - Always returns nil (cleanup failures don't block conversion)
//
// Recommended Usage:
//
//	defer processor.Cleanup() // Immediately after creating processor
//
// Thread Safety:
//   - Safe to call concurrently
//   - Snapshots imageMap before cleanup (doesn't hold locks during file removal)
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

// CalculateBackoff is the public version for testing.
func (ip *ImageProcessor) CalculateBackoff(attempt int) float64 {
	return ip.calculateBackoff(attempt)
}

// ============================================================================
// PLACEHOLDER FUNCTIONS (implemented in later phases)
// ============================================================================

// ProcessMarkdown processes markdown content to download remote images and replace URLs.
// Returns the processed markdown content with local image paths.
// Downloads all remote images concurrently (up to maxConcurrentDownloads) with retry logic,
// then rewrites the markdown to use local paths.
// Note: Images that fail to download are left with original URLs.
// Errors are collected but don't prevent conversion (graceful degradation).
func (ip *ImageProcessor) ProcessMarkdown(content string) (string, error) {
	// Create temp directory if it doesn't exist
	if err := os.MkdirAll(ip.tempDir, 0o755); err != nil {
		return "", fmt.Errorf("failed to create temp directory for images: %w", err)
	}

	// Detect all remote image URLs
	imageURLs := ip.DetectRemoteImages(content)

	// If no remote images, return content as-is
	if len(imageURLs) == 0 {
		return content, nil
	}

	// Download images concurrently with semaphore pattern and retry logic
	downloadErrors := ip.downloadImagesWithSemaphore(imageURLs)

	// Store download errors for access and reporting
	ip.mu.Lock()
	for url, err := range downloadErrors {
		ip.downloadErrors[url] = err.Error()
	}
	ip.mu.Unlock()

	// Rewrite markdown with downloaded image paths
	// Images that failed to download will keep original URLs
	processedContent := ip.RewriteMarkdownImageURLs(content)

	// Return processed content even if some downloads failed
	// Errors are collected in downloadErrors for reporting
	return processedContent, nil
}

// downloadImagesWithSemaphore downloads multiple images concurrently using a semaphore pattern.
// Uses retry logic for transient errors.
// Returns a map of URLs that failed to download with their error messages.
func (ip *ImageProcessor) downloadImagesWithSemaphore(urls []string) map[string]error {
	// Create a semaphore to limit concurrent downloads
	semaphore := make(chan struct{}, ip.maxConcurrentDownloads)

	// WaitGroup for synchronization
	var wg sync.WaitGroup
	downloadErrors := make(map[string]error)
	var errorsMu sync.Mutex

	for _, url := range urls {
		wg.Add(1)

		go func(imageURL string) {
			defer wg.Done()

			// Acquire semaphore slot
			semaphore <- struct{}{}
			defer func() { <-semaphore }()

			// Attempt download with retry logic
			_, err := ip.downloadWithRetry(imageURL)
			if err != nil {
				errorsMu.Lock()
				downloadErrors[imageURL] = err
				errorsMu.Unlock()
			}
		}(url)
	}

	wg.Wait()
	return downloadErrors
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
		errMsg := fmt.Sprintf("failed to create request: %v", err)
		ip.mu.Lock()
		ip.downloadErrors[imageURL] = errMsg
		ip.mu.Unlock()
		return "", fmt.Errorf("failed to create request for %s: %w", imageURL, err)
	}

	// Execute request
	resp, err := ip.httpClient.Do(req)
	if err != nil {
		errMsg := fmt.Sprintf("failed to download: %v", err)
		ip.mu.Lock()
		ip.downloadErrors[imageURL] = errMsg
		ip.mu.Unlock()
		return "", fmt.Errorf("failed to download %s: %w", imageURL, err)
	}
	defer resp.Body.Close()

	// Validate response
	if err := validateHTTPRequest(resp); err != nil {
		errMsg := fmt.Sprintf("invalid HTTP response from %s: %v", imageURL, err)
		ip.mu.Lock()
		ip.downloadErrors[imageURL] = errMsg
		ip.mu.Unlock()
		return "", fmt.Errorf("%s", errMsg)
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
			errMsg := fmt.Sprintf("image size validation failed: %v", err)
			ip.mu.Lock()
			ip.downloadErrors[imageURL] = errMsg
			ip.mu.Unlock()
			return "", fmt.Errorf("image size validation failed for %s: %w", imageURL, err)
		}
	}

	// Generate filename and create temp file
	fileName := generateFileName(imageURL, resp.Header.Get("Content-Type"))
	tempFile, err := os.CreateTemp(ip.tempDir, fileName)
	if err != nil {
		errMsg := fmt.Sprintf("failed to create temp file: %v", err)
		ip.mu.Lock()
		ip.downloadErrors[imageURL] = errMsg
		ip.mu.Unlock()
		return "", fmt.Errorf("failed to create temp file for %s: %w", imageURL, err)
	}
	defer tempFile.Close()

	// Copy response body to file with size tracking
	writtenBytes, err := io.Copy(tempFile, resp.Body)
	if err != nil {
		// Clean up failed download
		os.Remove(tempFile.Name())
		errMsg := fmt.Sprintf("failed to write image: %v", err)
		ip.mu.Lock()
		ip.downloadErrors[imageURL] = errMsg
		ip.mu.Unlock()
		return "", fmt.Errorf("failed to write image from %s: %w", imageURL, err)
	}

	// Validate size after download if not provided in header
	if contentLength == 0 {
		if err := ip.ValidateImageSize(writtenBytes); err != nil {
			os.Remove(tempFile.Name())
			errMsg := fmt.Sprintf("image too large: %v", err)
			ip.mu.Lock()
			ip.downloadErrors[imageURL] = errMsg
			ip.mu.Unlock()
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
// Retries on transient errors (timeouts, 5xx, rate limits).
// Fails immediately on permanent errors (4xx except 408).
func (ip *ImageProcessor) downloadWithRetry(imageURL string) (string, error) {
	var lastErr error

	for attempt := 0; attempt <= ip.maxRetries; attempt++ {
		// Try to download
		localPath, err := ip.DownloadImageOnce(imageURL)
		if err == nil {
			return localPath, nil
		}

		// Check if error is transient
		// Extract status code from error message if possible
		statusCode := 0
		if errMsg := err.Error(); strings.Contains(errMsg, "HTTP") {
			// Try to extract status code from error message
			parts := strings.Fields(errMsg)
			for i, part := range parts {
				if part == "HTTP" && i+1 < len(parts) {
					// Next field should be status code
					if code, parseErr := strconv.Atoi(parts[i+1]); parseErr == nil {
						statusCode = code
					}
				}
			}
		}

		lastErr = err
		isTransient := isTransientError(err, statusCode)

		// If permanent error or last attempt, return error
		if !isTransient || attempt >= ip.maxRetries {
			return "", err
		}

		// Calculate backoff and wait
		backoffSeconds := ip.calculateBackoff(attempt)
		time.Sleep(time.Duration(backoffSeconds*1000) * time.Millisecond)
	}

	return "", lastErr
}

// DownloadWithRetry is the public version for testing.
func (ip *ImageProcessor) DownloadWithRetry(imageURL string) (string, error) {
	return ip.downloadWithRetry(imageURL)
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
