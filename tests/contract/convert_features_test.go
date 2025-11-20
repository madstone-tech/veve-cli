package contract

import (
	"os"
	"os/exec"
	"path/filepath"
	"testing"
)

// TestMarkdownFeaturesRendering tests that various markdown features are properly rendered in PDF.
func TestMarkdownFeaturesRendering(t *testing.T) {
	tmpDir := t.TempDir()

	// Create markdown with various features
	testMDContent := `# Main Heading

## Subheading

### Third Level

**Bold text** and *italic text* and ***bold italic***

Regular paragraph with some content.

- List item 1
- List item 2
  - Nested item
- List item 3

1. Numbered item 1
2. Numbered item 2
3. Numbered item 3

[Link to Google](https://www.google.com)

` + "```go\nfunc main() {\n\tfmt.Println(\"Hello, World!\")\n}\n```" + `

| Header 1 | Header 2 |
|----------|----------|
| Cell 1   | Cell 2   |
| Cell 3   | Cell 4   |

> This is a blockquote.
> It spans multiple lines.

---

Final paragraph.
`
	testMDPath := filepath.Join(tmpDir, "features.md")
	if err := os.WriteFile(testMDPath, []byte(testMDContent), 0o644); err != nil {
		t.Fatalf("failed to create test markdown file: %v", err)
	}

	outputPath := filepath.Join(tmpDir, "features.pdf")

	// Run veve convert
	cmd := exec.Command("veve", testMDPath, "-o", outputPath)
	if output, err := cmd.CombinedOutput(); err != nil {
		t.Logf("veve output:\n%s", string(output))
		t.Fatalf("veve convert failed: %v", err)
	}

	// Verify output PDF was created
	if _, err := os.Stat(outputPath); err != nil {
		t.Fatalf("output PDF not created: %v", err)
	}

	// Verify it's a valid PDF
	pdf, err := os.ReadFile(outputPath)
	if err != nil {
		t.Fatalf("failed to read output PDF: %v", err)
	}

	if len(pdf) < 4 || string(pdf[:4]) != "%PDF" {
		t.Fatal("output file is not a valid PDF")
	}

	// Check that content is in the PDF (not comprehensive, just a sanity check)
	pdfStr := string(pdf)
	if len(pdfStr) == 0 {
		t.Fatal("PDF appears to be empty")
	}

	// PDF should be substantially sized to contain all the markdown features
	if len(pdf) < 5000 {
		t.Logf("warning: PDF size is %d bytes, may be too small", len(pdf))
	}
}

// TestMarkdownWithImages tests conversion with images (if supported by markdown).
// Note: This test assumes images are embedded or referenced correctly.
func TestMarkdownWithImages(t *testing.T) {
	tmpDir := t.TempDir()

	// Create a simple 1x1 PNG image for testing
	// This is the smallest valid PNG file possible
	pngData := []byte{
		0x89, 0x50, 0x4E, 0x47, 0x0D, 0x0A, 0x1A, 0x0A, // PNG signature
		0x00, 0x00, 0x00, 0x0D, 0x49, 0x48, 0x44, 0x52, // IHDR chunk
		0x00, 0x00, 0x00, 0x01, 0x00, 0x00, 0x00, 0x01,
		0x08, 0x02, 0x00, 0x00, 0x00, 0x90, 0x77, 0x53,
		0xDE, 0x00, 0x00, 0x00, 0x0C, 0x49, 0x44, 0x41, // IDAT chunk
		0x54, 0x08, 0x99, 0x63, 0xF8, 0xFF, 0xFF, 0x3F,
		0x00, 0x05, 0xFE, 0x02, 0xB7, 0x00, 0x00, 0x00, // IEND chunk
		0x00, 0x49, 0x45, 0x4E, 0x44, 0xAE, 0x42, 0x60,
		0x82,
	}

	imagePath := filepath.Join(tmpDir, "test.png")
	if err := os.WriteFile(imagePath, pngData, 0o644); err != nil {
		t.Fatalf("failed to create test image: %v", err)
	}

	// Create markdown referencing the image
	testMDContent := `# Document with Image

Here's an image:

![Test Image](./test.png)

And some text after the image.
`
	testMDPath := filepath.Join(tmpDir, "with_image.md")
	if err := os.WriteFile(testMDPath, []byte(testMDContent), 0o644); err != nil {
		t.Fatalf("failed to create test markdown file: %v", err)
	}

	outputPath := filepath.Join(tmpDir, "with_image.pdf")

	// Run veve convert
	cmd := exec.Command("veve", testMDPath, "-o", outputPath)
	if output, err := cmd.CombinedOutput(); err != nil {
		t.Logf("veve output:\n%s", string(output))
		t.Fatalf("veve convert failed: %v", err)
	}

	// Verify output PDF was created
	if _, err := os.Stat(outputPath); err != nil {
		t.Fatalf("output PDF not created: %v", err)
	}
}
