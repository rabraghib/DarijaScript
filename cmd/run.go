package cmd

import (
	"fmt"
	"os"

	"github.com/rabraghib/darijascript/src/interpreter"
	"github.com/rabraghib/darijascript/src/runtime"
	"github.com/spf13/cobra"
)

var runCmd = &cobra.Command{
	Use:   "run [file.ds]",
	Short: "Run a DarijaScript program",
	Args:  cobra.ExactArgs(1),
	RunE: func(cmd *cobra.Command, args []string) error {
		source, err := os.ReadFile(args[0])
		if err != nil {
			return fmt.Errorf("reading source file: %w", err)
		}
		eval := interpreter.NewEvaluatorWithIO(cmd.InOrStdin(), cmd.OutOrStdout())
		_, err = runtime.Run(string(source), eval)
		return err
	},
}

func init() { rootCmd.AddCommand(runCmd) }
