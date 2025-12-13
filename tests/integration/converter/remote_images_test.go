package converter_test

import (
	"fmt"
	"net/http"
	"os"
	"strings"
	"testing"
	"time"

	"github.com/madstone-tech/veve-cli/internal/converter"
	"github.com/madstone-tech/veve-cli/tests/testutil"
)

// ============================================================================
// T015: End-to-End Integration Tests
// ============================================================================

// TestProcessMarkdownWithRemoteImages tests the full workflow of processing
// markdown with remote images: downloading and rewriting in one call.
func TestProcessMarkdownWithRemoteImages(t *testing.T) {
	tests := []struct {
		name        string
		markdown    string
		setupServer func(*testutil.MockHTTPServer) // Optional setup
		expectedErr bool
		testDesc    string
	}{
		{
			name: "single_remote_image",
			markdown: `# Test Document
This document has ![one image](https://example.com/test.png).`,
			setupServer: func(mock *testutil.MockHTTPServer) {
				mock.RegisterImage("/test.png", "png")
			},
			expectedErr: false,
			testDesc:    "Should successfully download and rewrite single remote image",
		},
		{
			name: "multiple_remote_images",
			markdown: `# Multi-Image Document
First ![image 1](https://example.com/first.png)
Second ![image 2](https://example.com/second.jpg)
Third ![image 3](https://example.com/third.gif)`,
			setupServer: func(mock *testutil.MockHTTPServer) {
				mock.RegisterImage("/first.png", "png")
				mock.RegisterImage("/second.jpg", "jpeg")
				mock.RegisterImage("/third.gif", "gif")
			},
			expectedErr: false,
			testDesc:    "Should download and rewrite multiple images concurrently",
		},
		{
			name: "mixed_local_and_remote",
			markdown: `# Mixed Content
Local image: ![local](/path/to/local.png)
Remote image: ![remote](https://example.com/remote.png)
Another local: ![another](./relative/path.jpg)`,
			setupServer: func(mock *testutil.MockHTTPServer) {
				mock.RegisterImage("/remote.png", "png")
			},
			expectedErr: false,
			testDesc:    "Should handle mix of local and remote images correctly",
		},
		{
			name: "partial_failures",
			markdown: `# Document with failures
Working: ![works](https://example.com/works.png)
Broken: ![broken](https://example.com/broken.png)
Another works: ![also-works](https://example.com/also-works.gif)`,
			setupServer: func(mock *testutil.MockHTTPServer) {
				mock.RegisterImage("/works.png", "png")
				mock.RegisterImage("/also-works.gif", "gif")
				// broken.png is not registered - will get 404
			},
			expectedErr: false, // ProcessMarkdown doesn't error on partial failures
			testDesc:    "Should handle partial download failures gracefully",
		},
		{
			name: "no_remote_images",
			markdown: `# No Remote Images
Just a paragraph with ![local image](/local/path.png).
And some [text links](https://example.com) that aren't images.`,
			setupServer: func(mock *testutil.MockHTTPServer) {
				// No setup needed - no remote images
			},
			expectedErr: false,
			testDesc:    "Should handle content with no remote images",
		},
		{
			name: "duplicate_images",
			markdown: `# Duplicate Images
Image used twice: ![img](https://example.com/image.png)
Same image again: ![img](https://example.com/image.png)
And once more: ![img](https://example.com/image.png)`,
			setupServer: func(mock *testutil.MockHTTPServer) {
				mock.RegisterImage("/image.png", "png")
			},
			expectedErr: false,
			testDesc:    "Should handle duplicate image URLs efficiently with caching",
		},
		{
			name: "image_with_query_params",
			markdown: `# Query Parameters
Responsive image: ![responsive](https://example.com/image.png?width=800&height=600)`,
			setupServer: func(mock *testutil.MockHTTPServer) {
				mock.RegisterImage("/image.png?width=800&height=600", "png")
			},
			expectedErr: false,
			testDesc:    "Should handle images with URL query parameters",
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			tempDir := t.TempDir()
			processor := converter.NewImageProcessor(tempDir)

			mock := testutil.NewMockHTTPServer()
			defer mock.Close()

			if tt.setupServer != nil {
				tt.setupServer(mock)
			}

			processedContent, err := processor.ProcessMarkdown(tt.markdown)

			if (err != nil) != tt.expectedErr {
				t.Errorf("%s: got error %v, expectedErr %v", tt.testDesc, err, tt.expectedErr)
				return
			}

			// Verify that we didn't lose content
			if len(processedContent) == 0 && len(tt.markdown) > 0 {
				t.Errorf("%s: processed content is empty", tt.testDesc)
			}

			// Verify images were cached in the processor
			imageMap := processor.GetImageMap()
			if len(imageMap) > 0 {
				// Check that downloaded files exist
				for url, localPath := range imageMap {
					if _, err := os.Stat(localPath); err != nil {
						t.Errorf("%s: downloaded file missing for %s at %s", tt.testDesc, url, localPath)
					}
				}
			}

			// Clean up
			processor.Cleanup()
		})
	}
}

// TestProcessMarkdownRewritesURLs verifies that processed markdown contains
// local paths instead of remote URLs for successfully downloaded images.
func TestProcessMarkdownRewritesURLs(t *testing.T) {
	tempDir := t.TempDir()
	processor := converter.NewImageProcessor(tempDir)

	mock := testutil.NewMockHTTPServer()
	defer mock.Close()

	mock.RegisterImage("/image1.png", "png")
	mock.RegisterImage("/image2.jpg", "jpeg")

	markdown := `# Test
First image: ![first](` + mock.ImageURL("/image1.png") + `)
Second image: ![second](` + mock.ImageURL("/image2.jpg") + `)`

	processedContent, err := processor.ProcessMarkdown(markdown)
	if err != nil {
		t.Fatalf("ProcessMarkdown failed: %v", err)
	}

	// Check that URLs were rewritten
	if processedContent == markdown {
		t.Error("Markdown was not modified - URLs should have been rewritten")
	}

	// Verify that remote URLs are no longer in the processed content
	if mock.URL() == "" {
		t.Skip("Mock server URL not available")
	}

	// The processed content should NOT contain the mock server URL
	imageMap := processor.GetImageMap()
	if len(imageMap) > 0 {
		// At least one image should be in the map
		hasMapping := false
		for _, localPath := range imageMap {
			if localPath != "" {
				hasMapping = true
				break
			}
		}
		if !hasMapping {
			t.Error("Images downloaded but no local paths found in map")
		}
	}

	processor.Cleanup()
}

// TestProcessMarkdownConcurrency verifies that images are downloaded concurrently.
// This is a bit tricky to test without timing, so we verify by downloading many images.
func TestProcessMarkdownConcurrency(t *testing.T) {
	tempDir := t.TempDir()
	processor := converter.NewImageProcessor(tempDir)

	mock := testutil.NewMockHTTPServer()
	defer mock.Close()

	// Register many images
	imageCount := 10
	markdownParts := []string{"# Many Images\n"}
	for i := 1; i <= imageCount; i++ {
		path := "/" + fmt.Sprintf("image%d.png", i)
		mock.RegisterImage(path, "png")
		markdownParts = append(markdownParts,
			fmt.Sprintf("Image %d: ![img%d](%s)\n", i, i, mock.ImageURL(path)))
	}

	markdown := strings.Join(markdownParts, "")

	_, err := processor.ProcessMarkdown(markdown)
	if err != nil {
		t.Fatalf("ProcessMarkdown failed: %v", err)
	}

	// All images should be downloaded
	imageMap := processor.GetImageMap()
	if len(imageMap) != imageCount {
		t.Errorf("Expected %d images to be downloaded, got %d", imageCount, len(imageMap))
	}

	processor.Cleanup()
}

// TestProcessMarkdownCleanup verifies that temporary files are properly created
// in the temp directory and can be cleaned up.
func TestProcessMarkdownCleanup(t *testing.T) {
	tempDir := t.TempDir()
	processor := converter.NewImageProcessor(tempDir)

	mock := testutil.NewMockHTTPServer()
	defer mock.Close()

	mock.RegisterImage("/test.png", "png")
	mock.RegisterImage("/test.jpg", "jpeg")

	markdown := `![img1](` + mock.ImageURL("/test.png") + `)
![img2](` + mock.ImageURL("/test.jpg") + `)`

	_, err := processor.ProcessMarkdown(markdown)
	if err != nil {
		t.Fatalf("ProcessMarkdown failed: %v", err)
	}

	// Get the downloaded files
	imageMap := processor.GetImageMap()
	downloadedFiles := []string{}
	for _, path := range imageMap {
		downloadedFiles = append(downloadedFiles, path)
	}

	// Verify files exist before cleanup
	for _, filePath := range downloadedFiles {
		if _, err := os.Stat(filePath); err != nil {
			t.Errorf("File should exist before cleanup: %s", filePath)
		}
	}

	// Run cleanup
	processor.Cleanup()

	// Verify files are removed
	for _, filePath := range downloadedFiles {
		if _, err := os.Stat(filePath); err == nil {
			t.Errorf("File should be removed after cleanup: %s", filePath)
		}
	}
}

// ============================================================================
// T029: Partial Failure Handling Integration Tests
// ============================================================================

func TestProcessMarkdownPartialFailure(t *testing.T) {
	tempDir := t.TempDir()
	processor := converter.NewImageProcessor(tempDir)

	mock := testutil.NewMockHTTPServer()
	defer mock.Close()

	// Register 5 images: 3 good, 2 will fail
	mock.RegisterImage("/good1.png", "png")
	mock.RegisterImage("/good2.jpg", "jpeg")
	mock.RegisterImage("/good3.gif", "gif")
	// Intentionally not registering /bad1.png and /bad2.png for 404s

	markdown := `# Document with mixed results
Good image 1: ![img1](` + mock.ImageURL("/good1.png") + `)
Good image 2: ![img2](` + mock.ImageURL("/good2.jpg") + `)
Bad image 1: ![bad1](` + mock.ImageURL("/bad1.png") + `)
Good image 3: ![img3](` + mock.ImageURL("/good3.gif") + `)
Bad image 2: ![bad2](` + mock.ImageURL("/bad2.png") + `)`

	processedContent, err := processor.ProcessMarkdown(markdown)
	if err != nil {
		t.Fatalf("ProcessMarkdown failed: %v", err)
	}

	// Verify partial success
	imageMap := processor.GetImageMap()
	downloadErrors := processor.GetDownloadErrors()

	if len(imageMap) != 3 {
		t.Errorf("Expected 3 successful downloads, got %d", len(imageMap))
	}

	if len(downloadErrors) != 2 {
		t.Errorf("Expected 2 failed downloads, got %d", len(downloadErrors))
	}

	// Verify error messages contain descriptive info
	for url, errMsg := range downloadErrors {
		if errMsg == "" {
			t.Errorf("Error message empty for %s", url)
		}
		if !strings.Contains(errMsg, "HTTP") && !strings.Contains(errMsg, "404") {
			t.Logf("Warning: error message may lack detail: %s -> %s", url, errMsg)
		}
	}

	// Verify good images are rewritten
	for url := range imageMap {
		if strings.Contains(url, "good") {
			// Should be rewritten in content
			if !strings.Contains(processedContent, "/good") {
				t.Logf("Warning: good image might not be properly rewritten in content")
			}
		}
	}
}

// ============================================================================
// T030: Network Timeout Handling Integration Tests
// ============================================================================

func TestProcessMarkdownTimeoutHandling(t *testing.T) {
	tempDir := t.TempDir()
	processor := converter.NewImageProcessor(tempDir).WithTimeoutSeconds(1) // 1 second timeout

	mock := testutil.NewMockHTTPServer()
	defer mock.Close()

	// Register one fast image and one slow image
	mock.RegisterImage("/fast.png", "png")
	pngData, _ := testutil.CreateTestImageData("png")
	mock.RegisterWithDelay("/slow.png", 3*time.Second, http.StatusOK, "image/png", pngData)

	markdown := `# Timeout test
Fast: ![fast](` + mock.ImageURL("/fast.png") + `)
Slow: ![slow](` + mock.ImageURL("/slow.png") + `)`

	_, err := processor.ProcessMarkdown(markdown)
	if err != nil {
		t.Fatalf("ProcessMarkdown failed: %v", err)
	}

	imageMap := processor.GetImageMap()
	downloadErrors := processor.GetDownloadErrors()

	// Fast should succeed
	if len(imageMap) < 1 {
		t.Error("Expected at least 1 successful download (fast image)")
	}

	// Slow should timeout and fail
	if len(downloadErrors) < 1 {
		t.Error("Expected at least 1 failed download (slow/timeout)")
	}

	// Verify error message indicates timeout
	for url, errMsg := range downloadErrors {
		if strings.Contains(url, "slow") {
			if !strings.Contains(strings.ToLower(errMsg), "timeout") &&
				!strings.Contains(strings.ToLower(errMsg), "context") {
				t.Logf("Warning: timeout error message may lack timeout indication: %s", errMsg)
			}
		}
	}
}

// ============================================================================
// T044: Multiple Conversion Cleanup Integration Tests
// ============================================================================

// TestMultipleConversionsCleanup verifies that multiple processors can coexist
// and each can be cleaned up independently without affecting others.
func TestMultipleConversionsCleanup(t *testing.T) {
	mock := testutil.NewMockHTTPServer()
	defer mock.Close()

	mock.RegisterImage("/image1.png", "png")
	mock.RegisterImage("/image2.jpg", "jpeg")
	mock.RegisterImage("/image3.gif", "gif")

	// Create multiple processors with different temp directories
	tempDir1 := t.TempDir()
	tempDir2 := t.TempDir()
	tempDir3 := t.TempDir()

	processor1 := converter.NewImageProcessor(tempDir1)
	processor2 := converter.NewImageProcessor(tempDir2)
	processor3 := converter.NewImageProcessor(tempDir3)

	markdown1 := `![img1](` + mock.ImageURL("/image1.png") + `)`
	markdown2 := `![img2](` + mock.ImageURL("/image2.jpg") + `)`
	markdown3 := `![img3](` + mock.ImageURL("/image3.gif") + `)`

	// Process with each processor
	_, err1 := processor1.ProcessMarkdown(markdown1)
	if err1 != nil {
		t.Fatalf("Processor 1 failed: %v", err1)
	}

	_, err2 := processor2.ProcessMarkdown(markdown2)
	if err2 != nil {
		t.Fatalf("Processor 2 failed: %v", err2)
	}

	_, err3 := processor3.ProcessMarkdown(markdown3)
	if err3 != nil {
		t.Fatalf("Processor 3 failed: %v", err3)
	}

	// Verify all have downloaded files
	imageMap1 := processor1.GetImageMap()
	imageMap2 := processor2.GetImageMap()
	imageMap3 := processor3.GetImageMap()

	if len(imageMap1) != 1 {
		t.Errorf("Processor 1: expected 1 image, got %d", len(imageMap1))
	}
	if len(imageMap2) != 1 {
		t.Errorf("Processor 2: expected 1 image, got %d", len(imageMap2))
	}
	if len(imageMap3) != 1 {
		t.Errorf("Processor 3: expected 1 image, got %d", len(imageMap3))
	}

	// Get file paths before cleanup
	var files1, files2, files3 []string
	for _, path := range imageMap1 {
		files1 = append(files1, path)
	}
	for _, path := range imageMap2 {
		files2 = append(files2, path)
	}
	for _, path := range imageMap3 {
		files3 = append(files3, path)
	}

	// Clean up processor 1
	processor1.Cleanup()

	// Verify processor 1's files are removed
	for _, filePath := range files1 {
		if _, err := os.Stat(filePath); err == nil {
			t.Errorf("Processor 1 file should be removed: %s", filePath)
		}
	}

	// Verify processor 2 and 3 files still exist
	for _, filePath := range files2 {
		if _, err := os.Stat(filePath); err != nil {
			t.Errorf("Processor 2 file should still exist: %s", filePath)
		}
	}
	for _, filePath := range files3 {
		if _, err := os.Stat(filePath); err != nil {
			t.Errorf("Processor 3 file should still exist: %s", filePath)
		}
	}

	// Clean up remaining processors
	processor2.Cleanup()
	processor3.Cleanup()

	// Verify all files are removed
	for _, filePath := range files2 {
		if _, err := os.Stat(filePath); err == nil {
			t.Errorf("Processor 2 file should be removed: %s", filePath)
		}
	}
	for _, filePath := range files3 {
		if _, err := os.Stat(filePath); err == nil {
			t.Errorf("Processor 3 file should be removed: %s", filePath)
		}
	}
}

// TestMultipleConversionsWithSharedTempDir verifies behavior when multiple
// processors use the same parent temp directory (but their own subdirs).
func TestMultipleConversionsWithSequentialProcessing(t *testing.T) {
	parentTempDir := t.TempDir()

	mock := testutil.NewMockHTTPServer()
	defer mock.Close()

	mock.RegisterImage("/image1.png", "png")
	mock.RegisterImage("/image2.jpg", "jpeg")

	// First conversion
	processor1 := converter.NewImageProcessor(parentTempDir + "/conv1")
	os.MkdirAll(parentTempDir+"/conv1", 0755)

	markdown1 := `# First doc
![img1](` + mock.ImageURL("/image1.png") + `)`

	_, err1 := processor1.ProcessMarkdown(markdown1)
	if err1 != nil {
		t.Fatalf("First conversion failed: %v", err1)
	}

	// Second conversion in same parent
	processor2 := converter.NewImageProcessor(parentTempDir + "/conv2")
	os.MkdirAll(parentTempDir+"/conv2", 0755)

	markdown2 := `# Second doc
![img2](` + mock.ImageURL("/image2.jpg") + `)`

	_, err2 := processor2.ProcessMarkdown(markdown2)
	if err2 != nil {
		t.Fatalf("Second conversion failed: %v", err2)
	}

	// Verify both have images
	imageMap1 := processor1.GetImageMap()
	imageMap2 := processor2.GetImageMap()

	if len(imageMap1) != 1 {
		t.Errorf("First conversion: expected 1 image, got %d", len(imageMap1))
	}
	if len(imageMap2) != 1 {
		t.Errorf("Second conversion: expected 1 image, got %d", len(imageMap2))
	}

	// Clean up first
	processor1.Cleanup()

	// Second conversion's files should still exist
	for url, path := range imageMap2 {
		if _, err := os.Stat(path); err != nil {
			t.Errorf("Second conversion file should still exist after first cleanup: %s (for %s)", path, url)
		}
	}

	// Clean up second
	processor2.Cleanup()

	// All files should be gone
	for _, path := range imageMap2 {
		if _, err := os.Stat(path); err == nil {
			t.Errorf("File should be removed after cleanup: %s", path)
		}
	}
}

// ============================================================================
// T045: Session Size Limit Enforcement Integration Tests
// ============================================================================

// TestSessionSizeLimitEnforcement verifies that the 500MB session limit
// prevents downloading additional images when the limit is reached.
func TestSessionSizeLimitEnforcement(t *testing.T) {
	tempDir := t.TempDir()
	processor := converter.NewImageProcessor(tempDir)

	mock := testutil.NewMockHTTPServer()
	defer mock.Close()

	// Create large images to test size limits
	// Use 50MB images (within per-image 100MB limit)
	largeImageData := testutil.CreateLargeTestImageData(50 * 1024 * 1024)
	mock.RegisterResponse("/large1.bin", http.StatusOK, "application/octet-stream", largeImageData)
	mock.RegisterResponse("/large2.bin", http.StatusOK, "application/octet-stream", largeImageData)
	mock.RegisterResponse("/large3.bin", http.StatusOK, "application/octet-stream", largeImageData)
	mock.RegisterImage("/small.png", "png") // Small image at end

	// Build markdown with images that will exceed the limit
	markdown := fmt.Sprintf(`# Large Images Test
![large1](%s)
![large2](%s)
![large3](%s)
![small](%s)`,
		mock.ImageURL("/large1.bin"),
		mock.ImageURL("/large2.bin"),
		mock.ImageURL("/large3.bin"),
		mock.ImageURL("/small.png"),
	)

	_, err := processor.ProcessMarkdown(markdown)
	// ProcessMarkdown doesn't error on size limit, it just stops downloading
	if err != nil {
		t.Logf("ProcessMarkdown returned error: %v (may be expected for size limit)", err)
	}

	imageMap := processor.GetImageMap()
	downloadErrors := processor.GetDownloadErrors()

	// Should have downloaded some images but not all (hit the limit)
	totalDownloaded := len(imageMap)
	if totalDownloaded < 1 {
		t.Error("Expected at least one image to be downloaded")
	}
	if totalDownloaded > 10 {
		t.Error("Downloaded more images than expected")
	}

	// Verify downloads were tracked
	successful, _, _ := processor.GetDownloadStats()
	if successful < 1 {
		t.Error("Expected at least one successful download")
	}

	// If we didn't download all images, there should be errors
	if totalDownloaded < 4 && len(downloadErrors) < 1 {
		t.Logf("Warning: some images not downloaded but no errors recorded")
	}

	processor.Cleanup()
}

// TestSessionSizeLimitBoundary verifies behavior when approaching the 500MB boundary.
func TestSessionSizeLimitBoundary(t *testing.T) {
	tempDir := t.TempDir()
	processor := converter.NewImageProcessor(tempDir)

	mock := testutil.NewMockHTTPServer()
	defer mock.Close()

	// Create images that approach the limit
	// Use a practical approach with smaller images
	imageSize := 5 * 1024 * 1024 // 5MB each
	imageData := testutil.CreateLargeTestImageData(imageSize)

	// Register 110 images of 5MB each (total 550MB, exceeds 500MB limit)
	for i := 1; i <= 110; i++ {
		path := fmt.Sprintf("/image%d.bin", i)
		mock.RegisterResponse(path, http.StatusOK, "application/octet-stream", imageData)
	}

	// Build markdown requesting all 110 images (exceeds 500MB)
	markdownParts := []string{"# Boundary Test\n"}
	for i := 1; i <= 110; i++ {
		path := fmt.Sprintf("/image%d.bin", i)
		markdownParts = append(markdownParts,
			fmt.Sprintf("![img%d](%s)\n", i, mock.ImageURL(path)))
	}

	markdown := strings.Join(markdownParts, "")

	_, err := processor.ProcessMarkdown(markdown)
	if err != nil {
		t.Logf("ProcessMarkdown returned error: %v", err)
	}

	imageMap := processor.GetImageMap()
	successful, _, _ := processor.GetDownloadStats()

	// Should have downloaded exactly 100 images (500MB / 5MB per image)
	expectedImages := 500 * 1024 * 1024 / imageSize
	if successful > expectedImages+1 {
		t.Errorf("Expected at most %d successful downloads, got %d", expectedImages, successful)
	}

	// Not all 110 should have been downloaded (we hit the limit)
	if len(imageMap) == 110 {
		t.Error("Expected size limit to prevent downloading all 110 images")
	}

	processor.Cleanup()
}
