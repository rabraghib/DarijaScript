package cmd

import (
	"os"

	"github.com/spf13/cobra"
)

var rootCmd = &cobra.Command{
	Use:   "darijascript",
	Short: "DarijaScript Interpreter CLI",
	Long: `DarijaScript is a programming language that is based on the Moroccan Arabic language.
This application is a CLI that can be used to run DarijaScript programs.`,
	SilenceUsage: true,
}

func Execute() {
	if err := rootCmd.Execute(); err != nil {
		os.Exit(1)
	}
}
