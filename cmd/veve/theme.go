package main

import (
	"fmt"
	"os"
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
		// TODO: Implement theme add logic
		return nil
	},
}

var themeRemoveCmd = &cobra.Command{
	Use:   "remove [name]",
	Short: "Remove a custom theme",
	Long:  `Uninstall a custom theme.`,
	Args:  cobra.ExactArgs(1),
	RunE: func(cmd *cobra.Command, args []string) error {
		// TODO: Implement theme remove logic
		return nil
	},
}

func init() {
	themeCmd.AddCommand(themeListCmd)
	themeCmd.AddCommand(themeAddCmd)
	themeCmd.AddCommand(themeRemoveCmd)
}
