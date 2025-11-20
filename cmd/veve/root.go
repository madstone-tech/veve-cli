package main

// rootCmd is the main command and should be initialized in main.go
// This file provides additional root command setup if needed.

func init() {
	rootCmd.AddCommand(convertCmd)
	rootCmd.AddCommand(themeCmd)
}
