package cmd

import (
	"bufio"
	"fmt"
	"io"
	"strings"

	"github.com/rabraghib/darijascript/src/interpreter"
	"github.com/rabraghib/darijascript/src/runtime"
	"github.com/spf13/cobra"
)

var interactiveCmd = &cobra.Command{
	Use:   "interactive",
	Short: "Run the DarijaScript interpreter in interactive mode",
	Args:  cobra.NoArgs,
	RunE: func(cmd *cobra.Command, args []string) error {
		output := cmd.OutOrStdout()
		reader := bufio.NewReader(cmd.InOrStdin())
		eval := interpreter.NewEvaluatorWithIO(reader, output)
		fmt.Fprintln(output, "Welcome to the DarijaScript interpreter!\n\nType 'exit' to exit the interpreter.")
		for {
			fmt.Fprint(output, ">>> ")
			line, readErr := reader.ReadString('\n')
			if readErr != nil && readErr != io.EOF {
				return readErr
			}
			line = strings.TrimSuffix(strings.TrimSuffix(line, "\n"), "\r")
			if line == "exit" {
				break
			}
			if strings.TrimSpace(line) != "" {
				if _, err := runtime.Run(line, eval); err != nil {
					fmt.Fprintln(cmd.ErrOrStderr(), "Error:", err)
				}
			}
			if readErr == io.EOF {
				break
			}
		}
		fmt.Fprintln(output, "\nGoodbye!")
		return nil
	},
}

func init() { rootCmd.AddCommand(interactiveCmd) }
