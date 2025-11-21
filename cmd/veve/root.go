package main

import (
	"github.com/spf13/cobra"
)

// rootCmd is the main command and should be initialized in main.go
// This file provides additional root command setup if needed.

func init() {
	rootCmd.AddCommand(convertCmd)
	rootCmd.AddCommand(themeCmd)
	rootCmd.AddCommand(completionCmd)
}

// completionCmd provides shell completion generation
var completionCmd = &cobra.Command{
	Use:   "completion [bash|zsh|fish]",
	Short: "Generate shell completion scripts",
	Long: `Generate shell completion scripts for veve.

To load completions:

Bash:
  source <(veve completion bash)

  # To load completions for each session, execute once:
  veve completion bash > /usr/local/etc/bash_completion.d/veve

Zsh:
  source <(veve completion zsh)

  # To load completions for each session, execute once:
  veve completion zsh > "${fpath[1]}/_veve"

Fish:
  veve completion fish | source

  # To load completions for each session, execute once:
  veve completion fish > ~/.config/fish/completions/veve.fish
`,
	ValidArgs: []string{"bash", "zsh", "fish"},
	Args:      cobra.MatchAll(cobra.ExactArgs(1), cobra.OnlyValidArgs),
	RunE: func(cmd *cobra.Command, args []string) error {
		switch args[0] {
		case "bash":
			return rootCmd.GenBashCompletion(cmd.OutOrStdout())
		case "zsh":
			return rootCmd.GenZshCompletion(cmd.OutOrStdout())
		case "fish":
			return rootCmd.GenFishCompletion(cmd.OutOrStdout(), true)
		}
		return nil
	},
}
