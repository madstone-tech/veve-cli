package main

import (
	"github.com/spf13/cobra"
)

var convertCmd = &cobra.Command{
	Use:   "convert [input]",
	Short: "Convert markdown to PDF",
	Long:  `Convert a markdown file to PDF with optional theming and styling.`,
	Args:  cobra.ExactArgs(1),
	RunE: func(cmd *cobra.Command, args []string) error {
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

		pdfEngine, err := cmd.Flags().GetString("engine")
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

		remoteImagesTempDir, err := cmd.Flags().GetString("remote-images-temp-dir")
		if err != nil {
			return err
		}

		// Delegate to shared conversion function
		return performConversion(inputFile, outputFile, theme, pdfEngine, quiet, verbose,
			enableRemoteImages, remoteImagesTimeout, remoteImagesMaxRetries,
			remoteImagesTempDir)
	},
}

func init() {
	convertCmd.Flags().StringP("output", "o", "", "output PDF file path (default: input filename with .pdf extension)")
	convertCmd.Flags().StringP("theme", "t", "default", "theme to use for PDF styling")
	convertCmd.Flags().StringP("engine", "e", "xelatex", "PDF rendering engine to use (xelatex, lualatex, weasyprint, prince)")
	convertCmd.Flags().BoolP("enable-remote-images", "r", true, "automatically download and embed remote images in PDF")
	convertCmd.Flags().Int("remote-images-timeout", 10, "timeout in seconds for downloading each remote image")
	convertCmd.Flags().Int("remote-images-max-retries", 3, "maximum number of retries for failed image downloads")
	convertCmd.Flags().String("remote-images-temp-dir", "", "custom temporary directory for downloaded images (default: system temp dir)")
}
