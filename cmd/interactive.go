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
		source := ""
		for {
			if source == "" {
				fmt.Fprint(output, ">>> ")
			} else {
				fmt.Fprint(output, "... ")
			}
			line, readErr := reader.ReadString('\n')
			if readErr != nil && readErr != io.EOF {
				return readErr
			}
			line = strings.TrimSuffix(strings.TrimSuffix(line, "\n"), "\r")
			if line == "exit" && source == "" {
				break
			}
			if strings.TrimSpace(line) != "" || source != "" {
				source += line + "\n"
				if readErr != io.EOF && runtime.Incomplete(source) {
					continue
				}
				if _, err := runtime.Run(source, eval); err != nil {
					fmt.Fprintln(cmd.ErrOrStderr(), "Error:", err)
				}
				source = ""
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
