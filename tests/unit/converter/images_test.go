package converter_test

import (
	"net/http"
	"path/filepath"
	"os"
	"testing"
	"time"

	"github.com/madstone-tech/veve-cli/internal/converter"
	"github.com/madstone-tech/veve-cli/tests/testutil"
)

// ============================================================================
// T011: Image Detection Unit Tests
// ============================================================================

func TestDetectRemoteImages(t *testing.T) {
	tests := []struct {
		name            string
		content         string
		expectedURLs    []string
		expectedCount   int
		testDescription string
	}{
		{
			name: "single_remote_image",
			content: `# Test
![alt text](https://example.com/image.png)`,
			expectedURLs:    []string{"https://example.com/image.png"},
			expectedCount:   1,
			testDescription: "Should detect a single remote image",
		},
		{
			name: "multiple_remote_images",
			content: `# Test
![first](https://example.com/first.png)
Some text
![second](https://example.com/second.jpg)`,
			expectedURLs:    []string{"https://example.com/first.png", "https://example.com/second.jpg"},
			expectedCount:   2,
			testDescription: "Should detect multiple remote images",
		},
		{
			name: "mixed_local_and_remote",
			content: `# Test
![local](/images/local.png)
![remote](https://example.com/remote.png)
![relative](./images/relative.png)`,
			expectedURLs:    []string{"https://example.com/remote.png"},
			expectedCount:   1,
			testDescription: "Should detect only remote images, ignoring local paths",
		},
		{
			name: "http_and_https",
			content: `# Test
![http](http://example.com/image.png)
![https](https://example.com/image.jpg)`,
			expectedURLs:    []string{"http://example.com/image.png", "https://example.com/image.jpg"},
			expectedCount:   2,
			testDescription: "Should detect both HTTP and HTTPS images",
		},
		{
			name: "duplicate_urls",
			content: `# Test
![first](https://example.com/image.png)
![second](https://example.com/image.png)
![third](https://example.com/image.png)`,
			expectedURLs:    []string{"https://example.com/image.png"},
			expectedCount:   1,
			testDescription: "Should deduplicate repeated URLs",
		},
		{
			name: "no_remote_images",
			content: `# Test
No images here!
![local](/images/test.png)
![relative](./test.jpg)`,
			expectedURLs:    []string{},
			expectedCount:   0,
			testDescription: "Should return empty list when no remote images",
		},
		{
			name: "complex_markdown",
			content: `# Title
Some paragraph with ![inline](https://example.com/inline.png) image.

Another paragraph with [link](https://example.com).

![standalone](https://example.com/standalone.jpg)

And ![another](https://cdn.example.com/image.gif) one.`,
			expectedURLs:    []string{"https://example.com/inline.png", "https://example.com/standalone.jpg", "https://cdn.example.com/image.gif"},
			expectedCount:   3,
			testDescription: "Should handle complex markdown with mixed content",
		},
		{
			name:            "empty_content",
			content:         "",
			expectedURLs:    []string{},
			expectedCount:   0,
			testDescription: "Should handle empty content",
		},
		{
			name:            "image_with_url_params",
			content:         `![test](https://example.com/image.png?width=800&height=600)`,
			expectedURLs:    []string{"https://example.com/image.png?width=800&height=600"},
			expectedCount:   1,
			testDescription: "Should detect images with URL parameters",
		},
		{
			name:            "image_with_spaces_in_alt_text",
			content:         `![This is a long alt text with spaces](https://example.com/image.png)`,
			expectedURLs:    []string{"https://example.com/image.png"},
			expectedCount:   1,
			testDescription: "Should detect images with spaces in alt text",
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			tempDir := t.TempDir()
			processor := converter.NewImageProcessor(tempDir)

			urls := processor.DetectRemoteImages(tt.content)

			if len(urls) != tt.expectedCount {
				t.Errorf("%s: got %d URLs, want %d", tt.testDescription, len(urls), tt.expectedCount)
			}

			for i, expectedURL := range tt.expectedURLs {
				if i >= len(urls) {
					t.Errorf("%s: missing URL at index %d: %s", tt.testDescription, i, expectedURL)
					continue
				}
				if urls[i] != expectedURL {
					t.Errorf("%s: URL mismatch at index %d. Got %s, want %s", tt.testDescription, i, urls[i], expectedURL)
				}
			}
		})
	}
}

// ============================================================================
// T012: URL Validation Unit Tests
// ============================================================================

func TestIsRemoteURL(t *testing.T) {
	tests := []struct {
		url      string
		expected bool
		testName string
	}{
		{url: "https://example.com/image.png", expected: true, testName: "HTTPS URL"},
		{url: "http://example.com/image.png", expected: true, testName: "HTTP URL"},
		{url: "/local/path/image.png", expected: false, testName: "absolute local path"},
		{url: "./relative/image.png", expected: false, testName: "relative path"},
		{url: "image.png", expected: false, testName: "filename only"},
		{url: "ftp://example.com/image.png", expected: false, testName: "FTP protocol"},
		{url: "", expected: false, testName: "empty string"},
		{url: "HTTPS://EXAMPLE.COM/IMAGE.PNG", expected: true, testName: "uppercase HTTPS"},
		{url: "HTTP://EXAMPLE.COM/IMAGE.PNG", expected: true, testName: "uppercase HTTP"},
		{url: "https://", expected: true, testName: "HTTPS with no domain"},
		{url: "http://", expected: true, testName: "HTTP with no domain"},
	}

	for _, tt := range tests {
		t.Run(tt.testName, func(t *testing.T) {
			tempDir := t.TempDir()
			processor := converter.NewImageProcessor(tempDir)

			result := processor.IsRemoteURL(tt.url)
			if result != tt.expected {
				t.Errorf("IsRemoteURL(%q) = %v, want %v", tt.url, result, tt.expected)
			}
		})
	}
}

func TestIsImageContentType(t *testing.T) {
	tests := []struct {
		contentType string
		expected    bool
		testName    string
	}{
		{contentType: "image/png", expected: true, testName: "PNG image"},
		{contentType: "image/jpeg", expected: true, testName: "JPEG image"},
		{contentType: "image/gif", expected: true, testName: "GIF image"},
		{contentType: "image/webp", expected: true, testName: "WEBP image"},
		{contentType: "image/svg+xml", expected: true, testName: "SVG image"},
		{contentType: "image/bmp", expected: true, testName: "BMP image"},
		{contentType: "image/tiff", expected: true, testName: "TIFF image"},
		{contentType: "image/x-icon", expected: true, testName: "ICO image"},
		{contentType: "image/png; charset=utf-8", expected: true, testName: "PNG with charset"},
		{contentType: "text/html", expected: false, testName: "HTML content"},
		{contentType: "text/plain", expected: false, testName: "plain text"},
		{contentType: "application/json", expected: false, testName: "JSON"},
		{contentType: "application/pdf", expected: false, testName: "PDF"},
		{contentType: "", expected: false, testName: "empty content type"},
		{contentType: "IMAGE/PNG", expected: true, testName: "uppercase image type"},
	}

	for _, tt := range tests {
		t.Run(tt.testName, func(t *testing.T) {
			tempDir := t.TempDir()
			processor := converter.NewImageProcessor(tempDir)

			result := processor.IsImageContentType(tt.contentType)
			if result != tt.expected {
				t.Errorf("IsImageContentType(%q) = %v, want %v", tt.contentType, result, tt.expected)
			}
		})
	}
}

func TestGetExtensionFromContentType(t *testing.T) {
	tests := []struct {
		contentType string
		expected    string
		testName    string
	}{
		{contentType: "image/jpeg", expected: ".jpg", testName: "JPEG"},
		{contentType: "image/jpg", expected: ".jpg", testName: "JPG"},
		{contentType: "image/png", expected: ".png", testName: "PNG"},
		{contentType: "image/gif", expected: ".gif", testName: "GIF"},
		{contentType: "image/webp", expected: ".webp", testName: "WEBP"},
		{contentType: "image/svg+xml", expected: ".svg", testName: "SVG"},
		{contentType: "image/bmp", expected: ".bmp", testName: "BMP"},
		{contentType: "image/tiff", expected: ".tiff", testName: "TIFF"},
		{contentType: "image/png; charset=utf-8", expected: ".png", testName: "PNG with charset"},
		{contentType: "image/unknown", expected: ".img", testName: "unknown type"},
		{contentType: "", expected: ".img", testName: "empty type"},
	}

	for _, tt := range tests {
		t.Run(tt.testName, func(t *testing.T) {
			tempDir := t.TempDir()
			processor := converter.NewImageProcessor(tempDir)

			result := processor.GetExtensionFromContentType(tt.contentType)
			if result != tt.expected {
				t.Errorf("GetExtensionFromContentType(%q) = %q, want %q", tt.contentType, result, tt.expected)
			}
		})
	}
}

func TestValidateImageSize(t *testing.T) {
	tests := []struct {
		name       string
		setupFunc  func(*converter.ImageProcessor) // Optional setup
		imageSize  int64
		shouldFail bool
		testDesc   string
	}{
		{
			name:       "small_image",
			imageSize:  1024 * 1024, // 1MB
			shouldFail: false,
			testDesc:   "Should allow images under 100MB",
		},
		{
			name:       "max_image_size",
			imageSize:  100 * 1024 * 1024, // 100MB (exactly at limit)
			shouldFail: false,
			testDesc:   "Should allow images at exactly 100MB limit",
		},
		{
			name:       "exceeds_image_limit",
			imageSize:  101 * 1024 * 1024, // 101MB
			shouldFail: true,
			testDesc:   "Should reject images exceeding 100MB",
		},
		{
			name: "session_limit_multiple_images",
			setupFunc: func(ip *converter.ImageProcessor) {
				// Simulate prior downloads
				ip.RecordDownload(400 * 1024 * 1024)
			},
			imageSize:  150 * 1024 * 1024, // Would exceed 500MB session limit
			shouldFail: true,
			testDesc:   "Should reject when session limit exceeded",
		},
		{
			name: "session_limit_near_boundary",
			setupFunc: func(ip *converter.ImageProcessor) {
				ip.RecordDownload(400 * 1024 * 1024)
			},
			imageSize:  100 * 1024 * 1024,
			shouldFail: false,
			testDesc:   "Should allow if session limit not exceeded",
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			tempDir := t.TempDir()
			processor := converter.NewImageProcessor(tempDir)

			if tt.setupFunc != nil {
				tt.setupFunc(processor)
			}

			err := processor.ValidateImageSize(tt.imageSize)

			if (err != nil) != tt.shouldFail {
				t.Errorf("%s: got error %v, shouldFail %v", tt.testDesc, err, tt.shouldFail)
			}
		})
	}
}

// ============================================================================
// T013: Single Image Download Unit Tests
// ============================================================================

func TestDownloadImageOnce(t *testing.T) {
	tests := []struct {
		name        string
		setupServer func(*testutil.MockHTTPServer) string
		shouldFail  bool
		testDesc    string
	}{
		{
			name: "successful_png_download",
			setupServer: func(mock *testutil.MockHTTPServer) string {
				mock.RegisterImage("/test.png", "png")
				return mock.ImageURL("/test.png")
			},
			shouldFail: false,
			testDesc:   "Should successfully download PNG image",
		},
		{
			name: "successful_jpeg_download",
			setupServer: func(mock *testutil.MockHTTPServer) string {
				mock.RegisterImage("/test.jpg", "jpeg")
				return mock.ImageURL("/test.jpg")
			},
			shouldFail: false,
			testDesc:   "Should successfully download JPEG image",
		},
		{
			name: "successful_gif_download",
			setupServer: func(mock *testutil.MockHTTPServer) string {
				mock.RegisterImage("/test.gif", "gif")
				return mock.ImageURL("/test.gif")
			},
			shouldFail: false,
			testDesc:   "Should successfully download GIF image",
		},
		{
			name: "404_not_found",
			setupServer: func(mock *testutil.MockHTTPServer) string {
				return mock.URL() + "/nonexistent.png"
			},
			shouldFail: true,
			testDesc:   "Should fail on 404 Not Found",
		},
		{
			name: "non_image_content_type",
			setupServer: func(mock *testutil.MockHTTPServer) string {
				path := "/notimage.txt"
				mock.RegisterResponse(path, http.StatusOK, "text/plain", []byte("not an image"))
				return mock.ImageURL(path)
			},
			shouldFail: true,
			testDesc:   "Should fail on non-image content type",
		},
		{
			name: "timeout_exceeded",
			setupServer: func(mock *testutil.MockHTTPServer) string {
				path := "/slow.png"
				pngData, _ := testutil.CreateTestImageData("png")
				mock.RegisterWithDelay(path, 5*time.Second, http.StatusOK, "image/png", pngData)
				return mock.ImageURL(path)
			},
			shouldFail: true,
			testDesc:   "Should fail on timeout",
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			tempDir := t.TempDir()
			processor := converter.NewImageProcessor(tempDir).WithTimeoutSeconds(2)

			mock := testutil.NewMockHTTPServer()
			defer mock.Close()

			imageURL := tt.setupServer(mock)

			localPath, err := processor.DownloadImageOnce(imageURL)

			if (err != nil) != tt.shouldFail {
				t.Errorf("%s: got error %v, shouldFail %v", tt.testDesc, err, tt.shouldFail)
				return
			}

			if !tt.shouldFail {
				// Verify file exists
				if _, err := os.Stat(localPath); err != nil {
					t.Errorf("%s: downloaded file does not exist at %s: %v", tt.testDesc, localPath, err)
				}

				// Verify it's in the imageMap
				imageMap := processor.GetImageMap()
				if imageMap[imageURL] != localPath {
					t.Errorf("%s: image not in map or path mismatch", tt.testDesc)
				}
			}
		})
	}
}

func TestDownloadImageOnceCaching(t *testing.T) {
	tempDir := t.TempDir()
	processor := converter.NewImageProcessor(tempDir)

	mock := testutil.NewMockHTTPServer()
	defer mock.Close()

	mock.RegisterImage("/test.png", "png")
	imageURL := mock.ImageURL("/test.png")

	// First download
	path1, err := processor.DownloadImageOnce(imageURL)
	if err != nil {
		t.Fatalf("First download failed: %v", err)
	}

	// Second download should return cached path
	path2, err := processor.DownloadImageOnce(imageURL)
	if err != nil {
		t.Fatalf("Second download failed: %v", err)
	}

	if path1 != path2 {
		t.Errorf("Caching failed: got different paths for same URL: %s vs %s", path1, path2)
	}
}

// ============================================================================
// T014: Markdown Rewriting Unit Tests
// ============================================================================

func TestRewriteMarkdownImageURLs(t *testing.T) {
	tests := []struct {
		name     string
		content  string
		imageMap map[string]string
		expected string
		testDesc string
	}{
		{
			name: "single_image_rewrite",
			content: `# Test
![alt](https://example.com/test.png)`,
			imageMap: map[string]string{
				"https://example.com/test.png": "/tmp/veve-image-abc123.png",
			},
			expected: `# Test
![alt](/tmp/veve-image-abc123.png)`,
			testDesc: "Should rewrite single remote image URL",
		},
		{
			name: "multiple_images_rewrite",
			content: `![first](https://example.com/first.png)
Some text
![second](https://example.com/second.jpg)`,
			imageMap: map[string]string{
				"https://example.com/first.png":  "/tmp/veve-image-aaa.png",
				"https://example.com/second.jpg": "/tmp/veve-image-bbb.jpg",
			},
			expected: `![first](/tmp/veve-image-aaa.png)
Some text
![second](/tmp/veve-image-bbb.jpg)`,
			testDesc: "Should rewrite multiple image URLs",
		},
		{
			name: "mixed_rewrites",
			content: `![remote](https://example.com/remote.png)
![local](/local/path.png)`,
			imageMap: map[string]string{
				"https://example.com/remote.png": "/tmp/veve-image-111.png",
			},
			expected: `![remote](/tmp/veve-image-111.png)
![local](/local/path.png)`,
			testDesc: "Should rewrite remote but keep local paths",
		},
		{
			name:     "no_images",
			content:  `# Just text, no images`,
			imageMap: map[string]string{},
			expected: `# Just text, no images`,
			testDesc: "Should handle content with no images",
		},
		{
			name:     "image_not_in_map",
			content:  `![notdownloaded](https://example.com/notdownloaded.png)`,
			imageMap: map[string]string{},
			expected: `![notdownloaded](https://example.com/notdownloaded.png)`,
			testDesc: "Should leave URL unchanged if not in imageMap",
		},
		{
			name:    "image_with_url_parameters",
			content: `![test](https://example.com/image.png?width=800)`,
			imageMap: map[string]string{
				"https://example.com/image.png?width=800": "/tmp/veve-image-xyz.png",
			},
			expected: `![test](/tmp/veve-image-xyz.png)`,
			testDesc: "Should rewrite URLs with query parameters",
		},
		{
			name: "duplicate_images",
			content: `![first](https://example.com/same.png)
![second](https://example.com/same.png)`,
			imageMap: map[string]string{
				"https://example.com/same.png": "/tmp/veve-image-dup.png",
			},
			expected: `![first](/tmp/veve-image-dup.png)
![second](/tmp/veve-image-dup.png)`,
			testDesc: "Should rewrite duplicate images with same local path",
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			tempDir := t.TempDir()
			processor := converter.NewImageProcessor(tempDir)

			// Populate the imageMap
			for url, localPath := range tt.imageMap {
				processor.SetImageMap(url, localPath)
			}

			result := processor.RewriteMarkdownImageURLs(tt.content)

			if result != tt.expected {
				t.Errorf("%s:\nGot:\n%s\nWant:\n%s", tt.testDesc, result, tt.expected)
			}
		})
	}
}

// ============================================================================
// T024: Transient Error Classification Unit Tests
// ============================================================================

func TestIsTransientError(t *testing.T) {
	tests := []struct {
		name        string
		err         error
		statusCode  int
		isTransient bool
		testDesc    string
	}{
		// Transient errors - should retry
		{
			name:        "timeout_error",
			err:         timeoutError{},
			statusCode:  0,
			isTransient: true,
			testDesc:    "Network timeouts are transient",
		},
		{
			name:        "http_408_request_timeout",
			err:         nil,
			statusCode:  408,
			isTransient: true,
			testDesc:    "HTTP 408 is transient",
		},
		{
			name:        "http_429_rate_limit",
			err:         nil,
			statusCode:  429,
			isTransient: true,
			testDesc:    "HTTP 429 (rate limit) is transient",
		},
		{
			name:        "http_503_service_unavailable",
			err:         nil,
			statusCode:  503,
			isTransient: true,
			testDesc:    "HTTP 503 (service unavailable) is transient",
		},
		{
			name:        "http_504_gateway_timeout",
			err:         nil,
			statusCode:  504,
			isTransient: true,
			testDesc:    "HTTP 504 (gateway timeout) is transient",
		},

		// Permanent errors - should not retry
		{
			name:        "http_404_not_found",
			err:         nil,
			statusCode:  404,
			isTransient: false,
			testDesc:    "HTTP 404 is permanent",
		},
		{
			name:        "http_403_forbidden",
			err:         nil,
			statusCode:  403,
			isTransient: false,
			testDesc:    "HTTP 403 is permanent",
		},
		{
			name:        "http_401_unauthorized",
			err:         nil,
			statusCode:  401,
			isTransient: false,
			testDesc:    "HTTP 401 is permanent",
		},
		{
			name:        "http_400_bad_request",
			err:         nil,
			statusCode:  400,
			isTransient: false,
			testDesc:    "HTTP 400 is permanent",
		},
		{
			name:        "http_500_server_error",
			err:         nil,
			statusCode:  500,
			isTransient: false,
			testDesc:    "HTTP 500 is permanent (not in transient list)",
		},
		{
			name:        "http_301_redirect",
			err:         nil,
			statusCode:  301,
			isTransient: false,
			testDesc:    "HTTP 301 is permanent",
		},
		{
			name:        "no_error_success",
			err:         nil,
			statusCode:  200,
			isTransient: false,
			testDesc:    "HTTP 200 is not an error",
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			tempDir := t.TempDir()
			processor := converter.NewImageProcessor(tempDir)

			result := processor.IsTransientError(tt.err, tt.statusCode)
			if result != tt.isTransient {
				t.Errorf("%s: got %v, want %v", tt.testDesc, result, tt.isTransient)
			}
		})
	}
}

// timeoutError is a mock timeout error for testing
type timeoutError struct{}

func (e timeoutError) Error() string   { return "timeout" }
func (e timeoutError) Timeout() bool   { return true }
func (e timeoutError) Temporary() bool { return true }

// ============================================================================
// T025: Retry Backoff Calculation Unit Tests
// ============================================================================

func TestCalculateBackoff(t *testing.T) {
	tests := []struct {
		name         string
		attempt      int
		minExpected  float64
		maxExpected  float64
		testDesc     string
	}{
		{
			name:         "attempt_0",
			attempt:      0,
			minExpected:  0.0,
			maxExpected:  1.0,
			testDesc:     "Attempt 0 should backoff 0-1 seconds (2^0)",
		},
		{
			name:         "attempt_1",
			attempt:      1,
			minExpected:  0.0,
			maxExpected:  2.0,
			testDesc:     "Attempt 1 should backoff 0-2 seconds (2^1)",
		},
		{
			name:         "attempt_2",
			attempt:      2,
			minExpected:  0.0,
			maxExpected:  4.0,
			testDesc:     "Attempt 2 should backoff 0-4 seconds (2^2)",
		},
		{
			name:         "attempt_3",
			attempt:      3,
			minExpected:  0.0,
			maxExpected:  8.0,
			testDesc:     "Attempt 3 should backoff 0-8 seconds (2^3)",
		},
		{
			name:         "attempt_4_capped",
			attempt:      4,
			minExpected:  0.0,
			maxExpected:  10.0,
			testDesc:     "Attempt 4+ should be capped at 10 seconds (2^4 = 16, capped)",
		},
		{
			name:         "attempt_5_capped",
			attempt:      5,
			minExpected:  0.0,
			maxExpected:  10.0,
			testDesc:     "Attempt 5+ should be capped at 10 seconds",
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			tempDir := t.TempDir()
			processor := converter.NewImageProcessor(tempDir)

			// Test multiple times to check randomness
			results := make([]float64, 10)
			for i := 0; i < 10; i++ {
				results[i] = processor.CalculateBackoff(tt.attempt)
				if results[i] < tt.minExpected || results[i] > tt.maxExpected {
					t.Errorf("%s: attempt %d iteration %d got %f, expected %f-%f",
						tt.testDesc, tt.attempt, i, results[i], tt.minExpected, tt.maxExpected)
				}
			}

			// Check for randomness (if max > 0, should have different values)
			if tt.maxExpected > 0 {
				hasVariation := false
				for i := 1; i < len(results); i++ {
					if results[i] != results[i-1] {
						hasVariation = true
						break
					}
				}
				if !hasVariation {
					t.Logf("Warning: %s produced same value every time (may be unlucky randomness)", tt.testDesc)
				}
			}
		})
	}
}

// ============================================================================
// T026 & T027: Download Retry Tests
// ============================================================================

func TestDownloadWithRetryTransientFailure(t *testing.T) {
	tempDir := t.TempDir()
	processor := converter.NewImageProcessor(tempDir).WithMaxRetries(3)

	mock := testutil.NewMockHTTPServer()
	defer mock.Close()

	// Register a custom handler that returns 503 on first request, then 200
	attemptCount := 0
	mock.RegisterWithHandler("/test.png", func(w http.ResponseWriter, r *http.Request) {
		attemptCount++
		if attemptCount == 1 {
			// First attempt: service unavailable (transient)
			http.Error(w, "Service Unavailable", http.StatusServiceUnavailable)
		} else {
			// Second attempt: success
			w.Header().Set("Content-Type", "image/png")
			w.WriteHeader(http.StatusOK)
			pngData, _ := testutil.CreateTestImageData("png")
			w.Write(pngData)
		}
	})

	imageURL := mock.ImageURL("/test.png")

	// Should succeed on retry
	localPath, err := processor.DownloadWithRetry(imageURL)
	if err != nil {
		t.Fatalf("DownloadWithRetry failed: %v", err)
	}

	if localPath == "" {
		t.Error("Expected non-empty local path")
	}

	if attemptCount != 2 {
		t.Errorf("Expected 2 attempts (1 failure + 1 success), got %d", attemptCount)
	}

	// Verify file exists
	if _, err := os.Stat(localPath); err != nil {
		t.Errorf("Downloaded file should exist: %v", err)
	}
}

func TestDownloadWithRetryPermanentFailure(t *testing.T) {
	tempDir := t.TempDir()
	processor := converter.NewImageProcessor(tempDir).WithMaxRetries(3)

	mock := testutil.NewMockHTTPServer()
	defer mock.Close()

	// Register handler that always returns 404
	attemptCount := 0
	mock.RegisterWithHandler("/notfound.png", func(w http.ResponseWriter, r *http.Request) {
		attemptCount++
		http.Error(w, "Not Found", http.StatusNotFound)
	})

	imageURL := mock.ImageURL("/notfound.png")

	// Should fail immediately without retries
	localPath, err := processor.DownloadWithRetry(imageURL)
	if err == nil {
		t.Error("Expected error for 404 (permanent failure)")
	}

	if localPath != "" {
		t.Errorf("Expected empty path for failed download, got: %s", localPath)
	}

	if attemptCount != 1 {
		t.Errorf("Expected 1 attempt (no retries for 404), got %d", attemptCount)
	}
}

// ============================================================================
// T028: Error Message Formatting Unit Tests
// ============================================================================

func TestErrorMessageFormatting(t *testing.T) {
	tempDir := t.TempDir()
	processor := converter.NewImageProcessor(tempDir)

	mock := testutil.NewMockHTTPServer()
	defer mock.Close()

	// Test various error types
	testCases := []struct {
		name          string
		path          string
		setupHandler  func(string, *testutil.MockHTTPServer)
		expectedError bool
	}{
		{
			name: "404_not_found",
			path: "/notfound.png",
			setupHandler: func(path string, mock *testutil.MockHTTPServer) {
				mock.RegisterError(path, http.StatusNotFound, "Not Found")
			},
			expectedError: true,
		},
		{
			name: "500_server_error",
			path: "/error.png",
			setupHandler: func(path string, mock *testutil.MockHTTPServer) {
				mock.RegisterError(path, http.StatusInternalServerError, "Internal Server Error")
			},
			expectedError: true,
		},
		{
			name: "403_forbidden",
			path: "/forbidden.png",
			setupHandler: func(path string, mock *testutil.MockHTTPServer) {
				mock.RegisterError(path, http.StatusForbidden, "Forbidden")
			},
			expectedError: true,
		},
	}

	for _, tc := range testCases {
		t.Run(tc.name, func(t *testing.T) {
			tc.setupHandler(tc.path, mock)

			imageURL := mock.ImageURL(tc.path)
			_, err := processor.DownloadImageOnce(imageURL)

			if !tc.expectedError && err != nil {
				t.Errorf("Expected no error, got: %v", err)
			}
			if tc.expectedError && err == nil {
				t.Error("Expected error but got none")
			}

			// Check downloadErrors map is populated
			downloadErrors := processor.GetDownloadErrors()
			if tc.expectedError {
				if _, exists := downloadErrors[imageURL]; !exists {
					t.Error("Expected error in downloadErrors map")
				}
			}
		})
	}
}

// ============================================================================
// T040: Cleanup File Removal Unit Tests
// ============================================================================

func TestCleanupFileRemoval(t *testing.T) {
	tempDir := t.TempDir()
	processor := converter.NewImageProcessor(tempDir)

	// Create some dummy temp files and add them to imageMap
	file1, _ := os.CreateTemp(tempDir, "test1-*.png")
	file2, _ := os.CreateTemp(tempDir, "test2-*.jpg")
	file1Path := file1.Name()
	file2Path := file2.Name()
	file1.Close()
	file2.Close()

	// Add files to imageMap
	processor.SetImageMap("https://example.com/image1.png", file1Path)
	processor.SetImageMap("https://example.com/image2.jpg", file2Path)

	// Verify files exist
	if _, err := os.Stat(file1Path); err != nil {
		t.Fatalf("File 1 should exist: %v", err)
	}
	if _, err := os.Stat(file2Path); err != nil {
		t.Fatalf("File 2 should exist: %v", err)
	}

	// Call cleanup
	err := processor.Cleanup()
	if err != nil {
		t.Errorf("Cleanup returned error: %v", err)
	}

	// Verify files are deleted
	if _, err := os.Stat(file1Path); err == nil {
		t.Error("File 1 should be deleted after cleanup")
	}
	if _, err := os.Stat(file2Path); err == nil {
		t.Error("File 2 should be deleted after cleanup")
	}
}

// ============================================================================
// T041: Cleanup with Missing Files Unit Tests
// ============================================================================

func TestCleanupWithMissingFiles(t *testing.T) {
	tempDir := t.TempDir()
	processor := converter.NewImageProcessor(tempDir)

	// Add paths to non-existent files
	processor.SetImageMap("https://example.com/missing1.png", filepath.Join(tempDir, "missing1.png"))
	processor.SetImageMap("https://example.com/missing2.jpg", filepath.Join(tempDir, "missing2.jpg"))

	// Cleanup should not error even though files don't exist
	err := processor.Cleanup()
	if err != nil {
		t.Errorf("Cleanup should not error with missing files, got: %v", err)
	}

	// Verify imageMap is still accessible (cleanup doesn't clear it)
	imageMap := processor.GetImageMap()
	if len(imageMap) != 2 {
		t.Errorf("ImageMap should still have entries after cleanup")
	}
}

// ============================================================================
// T042: Disk Space Tracking Unit Tests
// ============================================================================

func TestDiskSpaceTracking(t *testing.T) {

	tests := []struct {
		name            string
		downloads       []int64 // sizes in bytes
		shouldAllSucceed bool
		testDesc        string
	}{
		{
			name:            "within_limit",
			downloads:       []int64{50 * 1024 * 1024, 150 * 1024 * 1024, 200 * 1024 * 1024},
			shouldAllSucceed: true,
			testDesc:        "Downloads within 500MB limit should all succeed",
		},
		{
			name:            "exceed_limit",
			downloads:       []int64{50 * 1024 * 1024, 150 * 1024 * 1024, 200 * 1024 * 1024, 150 * 1024 * 1024},
			shouldAllSucceed: false,
			testDesc:        "Fourth download should fail (exceeds 500MB limit)",
		},
		{
			name:            "at_limit",
			downloads:       []int64{250 * 1024 * 1024, 250 * 1024 * 1024},
			shouldAllSucceed: true,
			testDesc:        "Downloads totaling exactly 500MB should succeed",
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			tempDir := t.TempDir()
			processor := converter.NewImageProcessor(tempDir)

			successCount := 0
			for _, size := range tt.downloads {
				err := processor.ValidateImageSize(size)
				if err == nil {
					successCount++
					// Simulate successful download
					processor.RecordDownload(size)
				}
			}

			if tt.shouldAllSucceed {
				if successCount != len(tt.downloads) {
					t.Errorf("%s: expected all %d downloads to succeed, %d succeeded",
						tt.testDesc, len(tt.downloads), successCount)
				}
			} else {
				if successCount == len(tt.downloads) {
					t.Errorf("%s: expected some downloads to fail, all succeeded",
						tt.testDesc)
				}
			}
		})
	}
}

// ============================================================================
// T043: Per-Image Size Validation Unit Tests
// ============================================================================

func TestPerImageSizeValidation(t *testing.T) {
	tempDir := t.TempDir()
	processor := converter.NewImageProcessor(tempDir)

	mock := testutil.NewMockHTTPServer()
	defer mock.Close()

	tests := []struct {
		name         string
		setupHandler func(string, *testutil.MockHTTPServer)
		shouldFail   bool
		testDesc     string
	}{
		{
			name: "normal_size_image",
			setupHandler: func(path string, mock *testutil.MockHTTPServer) {
				mock.RegisterImage(path, "png")
			},
			shouldFail: false,
			testDesc:   "Normal size images should download",
		},
		{
			name: "oversized_image",
			setupHandler: func(path string, mock *testutil.MockHTTPServer) {
				// Create a large body (over 100MB)
				largeBody := make([]byte, 101*1024*1024)
				mock.RegisterResponse(path, http.StatusOK, "image/png", largeBody)
			},
			shouldFail: true,
			testDesc:   "Images over 100MB should be rejected",
		},
		{
			name: "boundary_size_100mb",
			setupHandler: func(path string, mock *testutil.MockHTTPServer) {
				// Exactly 100MB should be OK
				body := make([]byte, 100*1024*1024)
				mock.RegisterResponse(path, http.StatusOK, "image/png", body)
			},
			shouldFail: false,
			testDesc:   "Images exactly 100MB should be allowed",
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			path := "/" + tt.name + ".png"
			tt.setupHandler(path, mock)

			imageURL := mock.ImageURL(path)
			_, err := processor.DownloadImageOnce(imageURL)

			if (err != nil) != tt.shouldFail {
				t.Errorf("%s: got error %v, shouldFail %v", tt.testDesc, err, tt.shouldFail)
			}
		})
	}
}
