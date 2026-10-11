package cmd

import (
	"encoding/json"
	"fmt"
	"io"
	"os"

	"github.com/rabraghib/darijascript/src/runtime"
	"github.com/spf13/cobra"
)

func readSource(cmd *cobra.Command, args []string) (string, string, error) {
	name := "<stdin>"
	var source []byte
	var err error
	if len(args) == 1 {
		name = args[0]
		source, err = os.ReadFile(name)
	} else {
		source, err = io.ReadAll(cmd.InOrStdin())
	}
	return string(source), name, err
}

func newCheckCommand() *cobra.Command {
	var asJSON bool
	command := &cobra.Command{Use: "check [file.ds]", Short: "Check syntax without executing (reads stdin when no file is given)", Args: cobra.MaximumNArgs(1)}
	command.Flags().BoolVar(&asJSON, "json", false, "Return editor diagnostics as JSON")
	command.RunE = func(cmd *cobra.Command, args []string) error {
		source, name, err := readSource(cmd, args)
		if err != nil {
			return err
		}
		diagnostics := runtime.Check(source)
		if asJSON {
			return json.NewEncoder(cmd.OutOrStdout()).Encode(diagnostics)
		}
		if len(diagnostics) > 0 {
			_, err = runtime.Parse(source)
			return &runtime.SourceError{Name: name, Source: source, Cause: err}
		}
		_, err = fmt.Fprintln(cmd.OutOrStdout(), "Syntax OK")
		return err
	}
	return command
}

func newFormatCommand() *cobra.Command {
	var write bool
	command := &cobra.Command{Use: "format [file.ds]", Short: "Format source (reads stdin when no file is given)", Args: cobra.MaximumNArgs(1)}
	command.Flags().BoolVarP(&write, "write", "w", false, "Update the source file in place")
	command.RunE = func(cmd *cobra.Command, args []string) error {
		if write && len(args) != 1 {
			return fmt.Errorf("--write requires a source file")
		}
		source, name, err := readSource(cmd, args)
		if err != nil {
			return err
		}
		formatted, err := runtime.Format(source)
		if err != nil {
			return &runtime.SourceError{Name: name, Source: source, Cause: err}
		}
		if write {
			info, err := os.Stat(name)
			if err != nil {
				return err
			}
			return os.WriteFile(name, []byte(formatted), info.Mode())
		}
		_, err = io.WriteString(cmd.OutOrStdout(), formatted)
		return err
	}
	return command
}

func init() { rootCmd.AddCommand(newCheckCommand(), newFormatCommand()) }
