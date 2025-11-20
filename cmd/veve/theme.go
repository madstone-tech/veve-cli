package main

import (
	"fmt"
	"os"
	"path/filepath"
	"text/tabwriter"

	"github.com/andhi/veve-cli/internal/config"
	"github.com/andhi/veve-cli/internal/theme"
	"github.com/spf13/cobra"
)

var themeCmd = &cobra.Command{
	Use:   "theme",
	Short: "Manage themes",
	Long:  `Manage built-in and custom themes for PDF styling.`,
}

var themeListCmd = &cobra.Command{
	Use:   "list",
	Short: "List all available themes",
	Long:  `Display all available themes (built-in and user-installed).`,
	RunE: func(cmd *cobra.Command, args []string) error {
		// Get XDG paths
		paths, err := config.GetPaths()
		if err != nil {
			return fmt.Errorf("failed to get config paths: %w", err)
		}

		// Create and initialize theme loader
		loader := theme.NewLoader(paths.ThemesDir)
		if err := loader.DiscoverThemes(); err != nil {
			return fmt.Errorf("failed to discover themes: %w", err)
		}

		// Get all themes
		themes := loader.ListThemes()

		// Format as table
		w := tabwriter.NewWriter(os.Stdout, 0, 0, 2, ' ', 0)
		fmt.Fprintln(w, "NAME\tAUTHOR\tDESCRIPTION\tTYPE")
		fmt.Fprintln(w, "----\t------\t-----------\t----")

		for _, t := range themes {
			themeType := "user"
			if t.IsBuiltIn {
				themeType = "built-in"
			}
			fmt.Fprintf(w, "%s\t%s\t%s\t%s\n", t.Name, t.Author, t.Description, themeType)
		}

		w.Flush()
		return nil
	},
}

var themeAddCmd = &cobra.Command{
	Use:   "add [name] [path]",
	Short: "Add a custom theme",
	Long:  `Install a custom theme from a CSS file or zip archive.`,
	Args:  cobra.ExactArgs(2),
	RunE: func(cmd *cobra.Command, args []string) error {
		themeName := args[0]
		source := args[1]

		if themeName == "" {
			return fmt.Errorf("theme name cannot be empty")
		}

		// Get XDG paths
		paths, err := config.GetPaths()
		if err != nil {
			return fmt.Errorf("failed to get config paths: %w", err)
		}

		// Ensure themes directory exists
		if err := paths.EnsureDirectories(); err != nil {
			return fmt.Errorf("failed to create themes directory: %w", err)
		}

		// Download the theme
		downloader := theme.NewDownloader()
		css, err := downloader.Download(source)
		if err != nil {
			return fmt.Errorf("failed to download theme '%s': %w", themeName, err)
		}

		// Save theme to file
		themeFilePath := filepath.Join(paths.ThemesDir, themeName+".css")

		// Parse metadata from the CSS if present
		metadata, _, err := theme.ParseMetadata(css)
		if err != nil {
			metadata = &theme.ThemeMetadata{}
		}
		if metadata == nil {
			metadata = &theme.ThemeMetadata{}
		}

		// Apply defaults
		theme.ApplyMetadataDefaults(metadata, themeName)

		// Reconstruct CSS with metadata if we have it
		cssToSave := css
		if metadata.Name != "" {
			// Rebuild with metadata
			metadataBlock := fmt.Sprintf(`---
name: %s
author: %s
description: %s
version: %s
---
`, metadata.Name, metadata.Author, metadata.Description, metadata.Version)
			cssToSave = metadataBlock + "\n" + css
		}

		// Write theme file
		if err := os.WriteFile(themeFilePath, []byte(cssToSave), 0o644); err != nil {
			return fmt.Errorf("failed to save theme: %w", err)
		}

		fmt.Printf("Theme '%s' installed successfully at %s\n", themeName, themeFilePath)
		return nil
	},
}

var themeRemoveCmd = &cobra.Command{
	Use:   "remove [name]",
	Short: "Remove a custom theme",
	Long:  `Uninstall a custom theme.`,
	Args:  cobra.ExactArgs(1),
	RunE: func(cmd *cobra.Command, args []string) error {
		themeName := args[0]

		// Get XDG paths
		paths, err := config.GetPaths()
		if err != nil {
			return fmt.Errorf("failed to get config paths: %w", err)
		}

		// Get theme loader
		loader := theme.NewLoader(paths.ThemesDir)
		if err := loader.DiscoverThemes(); err != nil {
			return fmt.Errorf("failed to discover themes: %w", err)
		}

		// Check if theme exists
		t, exists := loader.GetRegistry().GetTheme(themeName)
		if !exists {
			return fmt.Errorf("theme not found: %s", themeName)
		}

		// Prevent removal of built-in themes
		if t.IsBuiltIn {
			return fmt.Errorf("cannot remove built-in theme '%s'", themeName)
		}

		// Check for --force flag
		force, err := cmd.Flags().GetBool("force")
		if err != nil {
			return err
		}

		// If not forcing, ask for confirmation
		if !force {
			fmt.Printf("Remove theme '%s'? (y/n) ", themeName)
			var response string
			_, err := fmt.Scanln(&response)
			if err != nil || (response != "y" && response != "Y") {
				fmt.Println("Cancelled.")
				return nil
			}
		}

		// Delete the theme file
		if err := os.Remove(t.FilePath); err != nil {
			return fmt.Errorf("failed to remove theme file: %w", err)
		}

		fmt.Printf("Theme '%s' removed successfully.\n", themeName)
		return nil
	},
}

func init() {
	themeRemoveCmd.Flags().BoolP("force", "f", false, "skip confirmation prompt")
	themeCmd.AddCommand(themeListCmd)
	themeCmd.AddCommand(themeAddCmd)
	themeCmd.AddCommand(themeRemoveCmd)
}
