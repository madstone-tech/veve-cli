package testutil

import (
	"fmt"
	"net/http"
	"net/http/httptest"
	"strings"
	"time"
)

// MockHTTPServer provides a test HTTP server for mocking remote image downloads.
type MockHTTPServer struct {
	Server    *httptest.Server
	Responses map[string]MockResponse
}

// MockResponse defines the response for a given URL.
type MockResponse struct {
	Status      int
	ContentType string
	Body        []byte
	Delay       time.Duration
	// If set, this will be called instead of returning a standard response
	Handler func(w http.ResponseWriter, r *http.Request)
}

// NewMockHTTPServer creates a new mock HTTP server for testing.
// It serves responses based on the registered Responses map.
func NewMockHTTPServer() *MockHTTPServer {
	mock := &MockHTTPServer{
		Responses: make(map[string]MockResponse),
	}

	// Create the test server
	mock.Server = httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		// Get the path from the request
		path := r.URL.Path
		if r.URL.RawQuery != "" {
			path += "?" + r.URL.RawQuery
		}

		// Look up the response
		response, exists := mock.Responses[path]

		if !exists {
			// Default to 404 if not registered
			http.Error(w, "Not Found", http.StatusNotFound)
			return
		}

		// Apply delay if specified
		if response.Delay > 0 {
			time.Sleep(response.Delay)
		}

		// Use custom handler if provided
		if response.Handler != nil {
			response.Handler(w, r)
			return
		}

		// Set response headers
		if response.ContentType != "" {
			w.Header().Set("Content-Type", response.ContentType)
		}
		if len(response.Body) > 0 {
			w.Header().Set("Content-Length", fmt.Sprintf("%d", len(response.Body)))
		}

		// Write status code
		w.WriteHeader(response.Status)

		// Write body
		if len(response.Body) > 0 {
			w.Write(response.Body)
		}
	}))

	return mock
}

// RegisterResponse registers a response for a given path.
func (m *MockHTTPServer) RegisterResponse(path string, status int, contentType string, body []byte) {
	m.Responses[path] = MockResponse{
		Status:      status,
		ContentType: contentType,
		Body:        body,
	}
}

// RegisterImage registers a mock image response (common case for testing).
// Creates a minimal valid PNG image (1x1 pixel).
func (m *MockHTTPServer) RegisterImage(path string, imageType string) {
	// Minimal PNG: 1x1 pixel
	pngData := []byte{
		0x89, 0x50, 0x4E, 0x47, 0x0D, 0x0A, 0x1A, 0x0A, // PNG signature
		0x00, 0x00, 0x00, 0x0D, 0x49, 0x48, 0x44, 0x52, // IHDR chunk
		0x00, 0x00, 0x00, 0x01, 0x00, 0x00, 0x00, 0x01, // 1x1 size
		0x08, 0x02, 0x00, 0x00, 0x00, 0x90, 0x77, 0x53, // 8-bit RGB
		0xDE, 0x00, 0x00, 0x00, 0x0C, 0x49, 0x44, 0x41, // IDAT chunk
		0x54, 0x08, 0x99, 0x63, 0xF8, 0xCF, 0xC0, 0x00, // Image data
		0x00, 0x03, 0x01, 0x01, 0x00, 0x18, 0xDD, 0x8D, // More data
		0xB4, 0x00, 0x00, 0x00, 0x00, 0x49, 0x45, 0x4E, // IEND chunk
		0x44, 0xAE, 0x42, 0x60, // End of PNG
	}

	// Minimal JPEG: valid JPEG header
	jpegData := []byte{
		0xFF, 0xD8, 0xFF, 0xE0, 0x00, 0x10, 0x4A, 0x46, // JPEG SOI + APP0
		0x49, 0x46, 0x00, 0x01, 0x01, 0x00, 0x00, 0x01, // JFIF marker
		0x00, 0x01, 0x00, 0x00, 0xFF, 0xD9, // End of Image
	}

	// Minimal GIF: 1x1 pixel
	gifData := []byte{
		0x47, 0x49, 0x46, 0x38, 0x39, 0x61, 0x01, 0x00, // GIF89a header + width
		0x01, 0x00, 0xF0, 0x00, 0x00, 0xFF, 0xFF, 0xFF, // height + color table
		0x00, 0x00, 0x00, 0x2C, 0x00, 0x00, 0x00, 0x00, // Image descriptor
		0x01, 0x00, 0x01, 0x00, 0x00, 0x02, 0x02, 0x44, // Image data
		0x01, 0x00, 0x3B, // Trailer
	}

	// Minimal SVG
	svgData := []byte(`<?xml version="1.0" encoding="UTF-8"?>
<svg width="1" height="1" xmlns="http://www.w3.org/2000/svg">
  <rect width="1" height="1" fill="black"/>
</svg>`)

	var body []byte
	var contentType string

	switch strings.ToLower(imageType) {
	case "png":
		body = pngData
		contentType = "image/png"
	case "jpeg", "jpg":
		body = jpegData
		contentType = "image/jpeg"
	case "gif":
		body = gifData
		contentType = "image/gif"
	case "svg":
		body = svgData
		contentType = "image/svg+xml"
	default:
		body = pngData
		contentType = "image/png"
	}

	m.RegisterResponse(path, http.StatusOK, contentType, body)
}

// RegisterError registers an error response (404, 500, etc.).
func (m *MockHTTPServer) RegisterError(path string, status int, message string) {
	m.Responses[path] = MockResponse{
		Status:      status,
		ContentType: "text/plain",
		Body:        []byte(message),
	}
}

// RegisterWithDelay registers a response that delays before responding.
// Useful for testing timeouts.
func (m *MockHTTPServer) RegisterWithDelay(path string, delay time.Duration, status int, contentType string, body []byte) {
	m.Responses[path] = MockResponse{
		Status:      status,
		ContentType: contentType,
		Body:        body,
		Delay:       delay,
	}
}

// RegisterWithHandler registers a custom handler for a path.
// Useful for complex scenarios like redirects.
func (m *MockHTTPServer) RegisterWithHandler(path string, handler func(w http.ResponseWriter, r *http.Request)) {
	m.Responses[path] = MockResponse{
		Handler: handler,
	}
}

// URL returns the base URL of the mock server.
func (m *MockHTTPServer) URL() string {
	return m.Server.URL
}

// ImageURL returns a full image URL for the given path.
func (m *MockHTTPServer) ImageURL(path string) string {
	return m.URL() + path
}

// Close closes the mock server.
func (m *MockHTTPServer) Close() {
	if m.Server != nil {
		m.Server.Close()
	}
}

// CreateTestImageData creates minimal valid image data for testing.
// Returns (data, contentType) for the specified image type.
func CreateTestImageData(imageType string) ([]byte, string) {
	switch strings.ToLower(imageType) {
	case "png":
		// Minimal PNG: 1x1 pixel
		return []byte{
			0x89, 0x50, 0x4E, 0x47, 0x0D, 0x0A, 0x1A, 0x0A,
			0x00, 0x00, 0x00, 0x0D, 0x49, 0x48, 0x44, 0x52,
			0x00, 0x00, 0x00, 0x01, 0x00, 0x00, 0x00, 0x01,
			0x08, 0x02, 0x00, 0x00, 0x00, 0x90, 0x77, 0x53,
			0xDE, 0x00, 0x00, 0x00, 0x0C, 0x49, 0x44, 0x41,
			0x54, 0x08, 0x99, 0x63, 0xF8, 0xCF, 0xC0, 0x00,
			0x00, 0x03, 0x01, 0x01, 0x00, 0x18, 0xDD, 0x8D,
			0xB4, 0x00, 0x00, 0x00, 0x00, 0x49, 0x45, 0x4E,
			0x44, 0xAE, 0x42, 0x60,
		}, "image/png"

	case "jpeg", "jpg":
		// Minimal JPEG
		return []byte{
			0xFF, 0xD8, 0xFF, 0xE0, 0x00, 0x10, 0x4A, 0x46,
			0x49, 0x46, 0x00, 0x01, 0x01, 0x00, 0x00, 0x01,
			0x00, 0x01, 0x00, 0x00, 0xFF, 0xD9,
		}, "image/jpeg"

	case "gif":
		// Minimal GIF: 1x1 pixel
		return []byte{
			0x47, 0x49, 0x46, 0x38, 0x39, 0x61, 0x01, 0x00,
			0x01, 0x00, 0xF0, 0x00, 0x00, 0xFF, 0xFF, 0xFF,
			0x00, 0x00, 0x00, 0x2C, 0x00, 0x00, 0x00, 0x00,
			0x01, 0x00, 0x01, 0x00, 0x00, 0x02, 0x02, 0x44,
			0x01, 0x00, 0x3B,
		}, "image/gif"

	case "svg":
		// Minimal SVG
		return []byte(`<?xml version="1.0" encoding="UTF-8"?>
<svg width="1" height="1" xmlns="http://www.w3.org/2000/svg">
  <rect width="1" height="1" fill="black"/>
</svg>`), "image/svg+xml"

	default:
		// Default to PNG
		data, _ := CreateTestImageData("png")
		return data, "image/png"
	}
}

// CreateLargeTestImageData creates image data of specified size (in bytes).
// Useful for testing size limits.
func CreateLargeTestImageData(sizeBytes int) []byte {
	// Start with a minimal PNG header
	header := []byte{
		0x89, 0x50, 0x4E, 0x47, 0x0D, 0x0A, 0x1A, 0x0A,
	}

	// Create data of requested size by padding
	data := make([]byte, sizeBytes)
	copy(data, header)

	// Fill the rest with zeros (valid image data will be ignored for size testing)
	for i := len(header); i < sizeBytes; i++ {
		data[i] = 0x00
	}

	return data
}
