package main

import (
	"fmt"
	"os"
	"os/exec"

	"github.com/andhi/veve-cli/internal"
	"github.com/andhi/veve-cli/internal/converter"
	"github.com/andhi/veve-cli/internal/logging"
	"github.com/spf13/cobra"
)

var (
	version = "0.1.0"
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
		// If no args provided, show help
		if len(args) == 0 {
			return cmd.Help()
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

		// Delegate to convert logic
		return performConversion(inputFile, outputFile, theme, pdfEngine, quiet, verbose)
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
}

// performConversion is a shared function used by both root command and convert subcommand.
func performConversion(inputFile, outputFile, theme, pdfEngine string, quiet, verbose bool) error {
	// Log if verbose
	logger.Debug("Converting %s to PDF (theme: %s, engine: %s)", inputFile, theme, pdfEngine)

	// Create converter
	pc, err := converter.NewPandocConverter()
	if err != nil {
		return err
	}

	// TODO: For MVP, we're skipping theme support until it's fully implemented.
	// Users can still style PDFs using Pandoc's CSS support with future enhancements.
	themeFile := ""
	if theme != "default" {
		// For non-default themes, would load from registry
		// For now, skip theme support
		logger.Warn("Theme support is not yet implemented. Using Pandoc defaults.")
	}

	// Perform conversion
	opts := converter.ConversionOptions{
		InputFile:  inputFile,
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
