package converter

import (
	"bytes"
	"fmt"
	"math/rand"
	"os"
	"os/exec"
	"path/filepath"
	"strings"
	"time"
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
	InputFile  string // Path to markdown file (or "-" for stdin)
	OutputFile string // Path to output PDF (optional; defaults to input with .pdf extension, or "-" for stdout)
	PDFEngine  string // PDF engine (pdflatex, xelatex, etc.)
	Theme      string // Path to CSS theme file (optional)
	Standalone bool   // Generate standalone PDF
	Quiet      bool   // Suppress output messages
	Verbose    bool   // Enable verbose output
}

// ValidateInputFile checks if the input markdown file exists and is readable.
// If filePath is "-", it's treated as stdin (no validation needed).
func ValidateInputFile(filePath string) error {
	if filePath == "" {
		return fmt.Errorf("input file path is empty")
	}

	// "-" represents stdin
	if filePath == "-" {
		return nil
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
// Supports "-" for stdin (input) and stdout (output).
func (pc *PandocConverter) Convert(opts ConversionOptions) error {
	// Validate input file exists
	if err := ValidateInputFile(opts.InputFile); err != nil {
		return fmt.Errorf("input validation failed: %w", err)
	}

	// Determine if we're using stdin/stdout
	isStdin := opts.InputFile == "-"
	isStdout := opts.OutputFile == "-"

	// Resolve output path if not provided (only if not using stdout)
	var outputPath string
	if !isStdout {
		outputPath = ResolveOutputPath(opts.InputFile, opts.OutputFile)
		// Ensure output directory exists
		if err := EnsureOutputDirectory(outputPath); err != nil {
			return err
		}
	} else {
		// For stdout, use a temp file that we'll read and output
		outputPath = filepath.Join(os.TempDir(), "veve-stdout-"+tempRandString()+".pdf")
	}

	// Build pandoc command
	var args []string

	// Handle input
	if isStdin {
		// Read from stdin - don't add input file argument
		// Pandoc will read from stdin if no input file is specified
	} else {
		args = append(args, opts.InputFile)
	}

	// Add output argument
	args = append(args, "-o", outputPath)
	args = append(args, "--pdf-engine", opts.PDFEngine)

	// Add standalone flag for better PDF output
	if opts.Standalone {
		args = append(args, "--standalone")
	}

	// Add theme/CSS if provided
	if opts.Theme != "" {
		// Check if it looks like a file path (contains / or \)
		if strings.Contains(opts.Theme, string(filepath.Separator)) || strings.Contains(opts.Theme, "/") {
			// It's a file path - verify it exists
			if _, err := os.Stat(opts.Theme); err != nil {
				return fmt.Errorf("theme file not found: %s: %w", opts.Theme, err)
			}
			args = append(args, "--css", opts.Theme)
		}
	}

	// Create command
	cmd := exec.Command(pc.PandocPath, args...)

	// If reading from stdin, connect standard input
	if isStdin {
		cmd.Stdin = os.Stdin
	}

	// Capture stderr for error reporting
	var stderr bytes.Buffer
	cmd.Stderr = &stderr

	// If outputting to stdout, prepare to capture stdout
	var stdout bytes.Buffer
	if isStdout {
		cmd.Stdout = &stdout
	}

	// Run conversion
	if err := cmd.Run(); err != nil {
		stderrMsg := stderr.String()
		if stderrMsg != "" {
			return fmt.Errorf("pandoc conversion failed: %w\nPandoc stderr: %s", err, stderrMsg)
		}
		return fmt.Errorf("pandoc conversion failed: %w", err)
	}

	// If outputting to stdout, read the temp file and write to os.Stdout
	if isStdout {
		pdfContent, err := os.ReadFile(outputPath)
		if err != nil {
			return fmt.Errorf("failed to read PDF from temp file: %w", err)
		}
		_, err = os.Stdout.Write(pdfContent)
		if err != nil {
			return fmt.Errorf("failed to write PDF to stdout: %w", err)
		}
		// Clean up temp file
		os.Remove(outputPath)
	}

	return nil
}

// tempRandString generates a random string for temp file names.
func tempRandString() string {
	const charset = "abcdefghijklmnopqrstuvwxyz0123456789"
	b := make([]byte, 8)
	for i := range b {
		b[i] = charset[rand.Intn(len(charset))]
	}
	return string(b)
}

func init() {
	rand.Seed(time.Now().UnixNano())
}
