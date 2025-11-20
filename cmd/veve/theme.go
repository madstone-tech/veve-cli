package main

import (
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
		// TODO: Implement theme list logic
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
