package converter

import (
	"bytes"
	"fmt"
	"os"
	"os/exec"
	"path/filepath"
	"strings"
)

// PandocConverter wraps Pandoc for markdown-to-PDF conversion.
type PandocConverter struct {
	PandocPath string // Full path to pandoc executable
}

// NewPandocConverter creates a new Pandoc converter.
// Returns an error if pandoc is not found in PATH.
func NewPandocConverter() (*PandocConverter, error) {
	pandocPath, err := exec.LookPath("pandoc")
	if err != nil {
		return nil, fmt.Errorf("pandoc not found in PATH: %w", err)
	}

	return &PandocConverter{
		PandocPath: pandocPath,
	}, nil
}

// ConversionOptions holds options for markdown-to-PDF conversion.
type ConversionOptions struct {
	InputFile  string // Path to markdown file
	OutputFile string // Path to output PDF (optional; defaults to input with .pdf extension)
	PDFEngine  string // PDF engine (pdflatex, xelatex, etc.)
	Theme      string // Path to CSS theme file (optional)
	Standalone bool   // Generate standalone PDF
	Quiet      bool   // Suppress output messages
	Verbose    bool   // Enable verbose output
}

// ValidateInputFile checks if the input markdown file exists and is readable.
func ValidateInputFile(filePath string) error {
	if filePath == "" {
		return fmt.Errorf("input file path is empty")
	}

	info, err := os.Stat(filePath)
	if err != nil {
		if os.IsNotExist(err) {
			return fmt.Errorf("input file not found: %s", filePath)
		}
		return fmt.Errorf("cannot access input file: %w", err)
	}

	if info.IsDir() {
		return fmt.Errorf("input path is a directory, not a file: %s", filePath)
	}

	return nil
}

// ResolveOutputPath resolves the output PDF path.
// If outputPath is empty, derives it from inputPath by replacing extension with .pdf.
func ResolveOutputPath(inputPath, outputPath string) string {
	if outputPath != "" {
		return outputPath
	}

	// Replace markdown extension with .pdf
	ext := filepath.Ext(inputPath)
	if ext != "" {
		return strings.TrimSuffix(inputPath, ext) + ".pdf"
	}

	return inputPath + ".pdf"
}

// EnsureOutputDirectory creates all parent directories for the output file if they don't exist.
func EnsureOutputDirectory(outputPath string) error {
	outputDir := filepath.Dir(outputPath)

	// If the directory is empty or ".", don't try to create it
	if outputDir == "" || outputDir == "." {
		return nil
	}

	if err := os.MkdirAll(outputDir, 0o755); err != nil {
		return fmt.Errorf("failed to create output directory %s: %w", outputDir, err)
	}

	return nil
}

// Convert converts a markdown file to PDF using Pandoc.
func (pc *PandocConverter) Convert(opts ConversionOptions) error {
	// Validate input file exists
	if err := ValidateInputFile(opts.InputFile); err != nil {
		return fmt.Errorf("input validation failed: %w", err)
	}

	// Resolve output path if not provided
	outputPath := ResolveOutputPath(opts.InputFile, opts.OutputFile)

	// Ensure output directory exists
	if err := EnsureOutputDirectory(outputPath); err != nil {
		return err
	}

	// Build pandoc command
	args := []string{
		opts.InputFile,
		"-o", outputPath,
		"--pdf-engine", opts.PDFEngine,
	}

	// Add standalone flag for better PDF output
	if opts.Standalone {
		args = append(args, "--standalone")
	}

	// Add theme/CSS if provided
	if opts.Theme != "" {
		// Check if theme file exists
		if _, err := os.Stat(opts.Theme); err != nil {
			return fmt.Errorf("theme file not found: %s: %w", opts.Theme, err)
		}
		args = append(args, "--css", opts.Theme)
	}

	// Create command
	cmd := exec.Command(pc.PandocPath, args...)

	// Capture stderr for error reporting
	var stderr bytes.Buffer
	cmd.Stderr = &stderr

	// Run conversion
	if err := cmd.Run(); err != nil {
		stderrMsg := stderr.String()
		if stderrMsg != "" {
			return fmt.Errorf("pandoc conversion failed: %w\nPandoc stderr: %s", err, stderrMsg)
		}
		return fmt.Errorf("pandoc conversion failed: %w", err)
	}

	return nil
}
