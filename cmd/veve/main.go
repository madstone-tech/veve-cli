package main

import (
	"fmt"
	"os"
	"os/exec"
	"path/filepath"
	"strings"

	"github.com/madstone-tech/veve-cli/internal"
	"github.com/madstone-tech/veve-cli/internal/config"
	"github.com/madstone-tech/veve-cli/internal/converter"
	"github.com/madstone-tech/veve-cli/internal/logging"
	"github.com/madstone-tech/veve-cli/internal/theme"
	"github.com/spf13/cobra"
)

var (
	version = "0.1.1"
	logger  *logging.Logger
)

var rootCmd = &cobra.Command{
	Use:   "veve [input]",
	Short: "veve - markdown to PDF converter with theme support",
	Long: `veve is a fast, cross-platform CLI tool for converting markdown files to beautiful PDFs.
It supports built-in themes, custom styling, and Pandoc-powered conversion.

Usage:
  veve input.md [-o output.pdf] [--theme theme-name] [flags]
  veve convert input.md [flags]
  veve theme list|add|remove [...]`,
	Version: version,
	Args:    cobra.MaximumNArgs(1),
	PersistentPreRunE: func(cmd *cobra.Command, args []string) error {
		// Check if pandoc is installed
		if _, err := exec.LookPath("pandoc"); err != nil {
			return internal.PandocNotFound()
		}
		return nil
	},
	RunE: func(cmd *cobra.Command, args []string) error {
		// Allow "-" for stdin without requiring it as an explicit argument
		// If no args and no stdin, show help
		if len(args) == 0 {
			// Check if stdin is available
			stat, _ := os.Stdin.Stat()
			if (stat.Mode() & os.ModeCharDevice) != 0 {
				// stdin is a terminal (no piped input)
				return cmd.Help()
			}
			// stdin has piped input, treat as "-"
			args = []string{"-"}
		}

		// If a markdown file is provided, treat it as convert command
		inputFile := args[0]

		// Get flags
		outputFile, err := cmd.Flags().GetString("output")
		if err != nil {
			return err
		}

		theme, err := cmd.Flags().GetString("theme")
		if err != nil {
			return err
		}

		pdfEngine, err := cmd.Flags().GetString("pdf-engine")
		if err != nil {
			return err
		}

		enableRemoteImages, err := cmd.Flags().GetBool("enable-remote-images")
		if err != nil {
			return err
		}

		remoteImagesTimeout, err := cmd.Flags().GetInt("remote-images-timeout")
		if err != nil {
			return err
		}

		remoteImagesMaxRetries, err := cmd.Flags().GetInt("remote-images-max-retries")
		if err != nil {
			return err
		}

		// Delegate to convert logic
		return performConversion(inputFile, outputFile, theme, pdfEngine, quiet, verbose,
			enableRemoteImages, remoteImagesTimeout, remoteImagesMaxRetries)
	},
}

var (
	verbose bool
	quiet   bool
)

func init() {
	rootCmd.PersistentFlags().BoolVar(&verbose, "verbose", false, "enable verbose output")
	rootCmd.PersistentFlags().BoolVar(&quiet, "quiet", false, "suppress non-error output")
	rootCmd.Flags().StringP("output", "o", "", "output PDF file path (default: input filename with .pdf extension)")
	rootCmd.Flags().StringP("theme", "t", "default", "theme to use for PDF styling")
	rootCmd.Flags().StringP("pdf-engine", "e", "pdflatex", "Pandoc PDF engine to use")
	rootCmd.Flags().BoolP("enable-remote-images", "r", true, "automatically download and embed remote images in PDF")
	rootCmd.Flags().Int("remote-images-timeout", 10, "timeout in seconds for downloading each remote image")
	rootCmd.Flags().Int("remote-images-max-retries", 3, "maximum number of retries for failed image downloads")
}

// performConversion is a shared function used by both root command and convert subcommand.
func performConversion(inputFile, outputFile, themeName, pdfEngine string, quiet, verbose bool,
	enableRemoteImages bool, remoteImagesTimeout, remoteImagesMaxRetries int) error {
	// Log if verbose
	logger.Debug("Converting %s to PDF (theme: %s, engine: %s)", inputFile, themeName, pdfEngine)

	// Create converter
	pc, err := converter.NewPandocConverter()
	if err != nil {
		return err
	}

	// Get XDG paths for theme discovery
	paths, err := config.GetPaths()
	if err != nil {
		return fmt.Errorf("failed to get config paths: %w", err)
	}

	// Ensure all necessary directories exist (including themes directory)
	if err := paths.EnsureDirectories(); err != nil {
		logger.Debug("Warning: Failed to create directories: %v", err)
		// Continue anyway - directories may already exist or not be writable
	}

	// Create theme loader
	loader := theme.NewLoader(paths.ThemesDir)

	// Discover available themes
	if err := loader.DiscoverThemes(); err != nil {
		logger.Debug("Error discovering themes: %v (continuing with defaults)", err)
	}

	// Check if theme is a file path (contains / or \ or .css)
	isFilePath := strings.ContainsAny(themeName, "/\\") || strings.HasSuffix(themeName, ".css")

	// Load theme CSS
	var themeFile string
	if isFilePath {
		// Handle file path theme
		css, err := loader.LoadThemeFromPath(themeName)
		if err != nil {
			return fmt.Errorf("failed to load theme from path '%s': %w", themeName, err)
		}

		if css != "" {
			// Write theme CSS to temporary file for Pandoc
			// Extract just the filename without path for temp file naming
			baseName := filepath.Base(themeName)
			if !strings.HasSuffix(baseName, ".css") {
				baseName = baseName + ".css"
			}
			tempThemeFile := filepath.Join(os.TempDir(), fmt.Sprintf("veve-theme-%s", baseName))
			if err := os.WriteFile(tempThemeFile, []byte(css), 0o644); err != nil {
				logger.Warn("Failed to write theme CSS: %v", err)
			} else {
				themeFile = tempThemeFile
				defer os.Remove(tempThemeFile) // Clean up temp file after conversion
			}
		}
	} else {
		// Handle named theme
		selectedTheme, err := loader.LoadTheme(themeName)
		if err != nil {
			// Build helpful error message with available themes
			availableThemes := loader.ListThemes()
			themeNames := make([]string, len(availableThemes))
			for i, t := range availableThemes {
				themeNames[i] = t.Name
			}
			return fmt.Errorf("invalid theme '%s': available themes are: %v", themeName, themeNames)
		}

		// Load theme CSS
		if selectedTheme.Name != "default" || selectedTheme.IsBuiltIn {
			css, err := loader.LoadThemeCSS(themeName)
			if err != nil {
				// If theme not found in loader's CSS, skip it
				logger.Debug("Theme CSS not found for %s: %v", themeName, err)
			} else if css != "" {
				// Write theme CSS to temporary file for Pandoc
				tempThemeFile := filepath.Join(os.TempDir(), fmt.Sprintf("veve-theme-%s.css", themeName))
				if err := os.WriteFile(tempThemeFile, []byte(css), 0o644); err != nil {
					logger.Warn("Failed to write theme CSS: %v", err)
				} else {
					themeFile = tempThemeFile
					defer os.Remove(tempThemeFile) // Clean up temp file after conversion
				}
			}
		}
	}

	// Process remote images if enabled
	var processedInputFile string
	var imageProcessor *converter.ImageProcessor
	if enableRemoteImages {
		// Create temp directory for downloaded images
		tempDir := filepath.Join(os.TempDir(), fmt.Sprintf("veve-images-%d", os.Getpid()))
		imageProcessor = converter.NewImageProcessor(tempDir).
			WithTimeoutSeconds(remoteImagesTimeout).
			WithMaxRetries(remoteImagesMaxRetries)
		defer imageProcessor.Cleanup()

		// Read markdown content
		content, err := os.ReadFile(inputFile)
		if err != nil {
			return fmt.Errorf("failed to read input file: %w", err)
		}

		// Process markdown to download remote images
		processedContent, err := imageProcessor.ProcessMarkdown(string(content))
		if err != nil {
			logger.Debug("Warning: Image processing failed: %v (continuing with original content)", err)
			processedInputFile = inputFile
		} else {
			// Write processed content to temporary file
			tempProcessedFile := filepath.Join(os.TempDir(), fmt.Sprintf("veve-processed-%d.md", os.Getpid()))
			if err := os.WriteFile(tempProcessedFile, []byte(processedContent), 0o644); err != nil {
				logger.Debug("Warning: Failed to write processed markdown: %v (using original)", err)
				processedInputFile = inputFile
			} else {
				processedInputFile = tempProcessedFile
				defer os.Remove(tempProcessedFile) // Clean up temp file after conversion
			}

			// Log image download summary with detailed error reporting
			successful, failed, total := imageProcessor.GetDownloadStats()
			if !quiet {
				if total > 0 {
					if failed == 0 {
						// All succeeded
						logger.Info("Successfully downloaded %d image(s)", successful)
					} else if successful == 0 {
						// All failed
						logger.Warn("Failed to download %d image(s)", failed)
					} else {
						// Partial success
						logger.Info("Downloaded %d of %d image(s)", successful, total)
					}
				}

				// Log detailed error information
				if failed > 0 {
					errorSummary := imageProcessor.GetErrorSummary()
					logger.Warn(errorSummary)
				}
			}
		}
	} else {
		processedInputFile = inputFile
	}

	// Perform conversion
	opts := converter.ConversionOptions{
		InputFile:  processedInputFile,
		OutputFile: outputFile,
		PDFEngine:  pdfEngine,
		Theme:      themeFile,
		Standalone: true,
		Quiet:      quiet,
		Verbose:    verbose,
	}

	if err := pc.Convert(opts); err != nil {
		return err
	}

	// Log success
	resolvedOutput := converter.ResolveOutputPath(inputFile, outputFile)
	if !quiet {
		logger.Info("Successfully converted %s to %s", inputFile, resolvedOutput)
	}

	return nil
}

func main() {
	// Initialize logger
	logger = logging.NewLogger(quiet, verbose)
	logging.SetGlobalLogger(logger)

	// Execute the root command
	if err := rootCmd.Execute(); err != nil {
		// Check if it's a VeveError for proper formatting
		if veveErr, ok := err.(*internal.VeveError); ok {
			fmt.Fprintf(os.Stderr, "%s\n", veveErr.Error())
			os.Exit(internal.ExitError)
		}

		// For Cobra errors and others
		fmt.Fprintf(os.Stderr, "[ERROR] %v\n", err)

		// Determine exit code based on error type
		if _, ok := err.(interface{ ExitCode() int }); ok {
			os.Exit(internal.ExitError)
		}

		os.Exit(internal.ExitError)
	}
}
