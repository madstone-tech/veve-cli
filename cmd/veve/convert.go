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

		pdfEngine, err := cmd.Flags().GetString("pdf-engine")
		if err != nil {
			return err
		}

		// Delegate to shared conversion function
		return performConversion(inputFile, outputFile, theme, pdfEngine, quiet, verbose)
	},
}

func init() {
	convertCmd.Flags().StringP("output", "o", "", "output PDF file path (default: input filename with .pdf extension)")
	convertCmd.Flags().StringP("theme", "t", "default", "theme to use for PDF styling")
	convertCmd.Flags().StringP("pdf-engine", "e", "pdflatex", "Pandoc PDF engine to use")
}
