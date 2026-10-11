package cmd

import (
	"bytes"
	"encoding/json"
	"os"
	"path/filepath"
	"strings"
	"testing"

	"github.com/rabraghib/darijascript/src/runtime"
	"github.com/spf13/cobra"
)

func executeCommand(command *cobra.Command, args []string, source string) (string, error) {
	var output bytes.Buffer
	command.SetArgs(args)
	command.SetIn(strings.NewReader(source))
	command.SetOut(&output)
	command.SetErr(&output)
	command.SilenceErrors = true
	command.SilenceUsage = true
	err := command.Execute()
	return output.String(), err
}

func TestCheckJSON(t *testing.T) {
	output, err := executeCommand(newCheckCommand(), []string{"--json"}, "golih(dakhel());")
	if err != nil || output != "[]\n" {
		t.Fatalf("check executed source or failed: %q, %v", output, err)
	}
	output, err = executeCommand(newCheckCommand(), []string{"--json"}, "9ayed x = ;")
	if err != nil {
		t.Fatal(err)
	}
	var diagnostics []runtime.Diagnostic
	if err := json.Unmarshal([]byte(output), &diagnostics); err != nil {
		t.Fatal(err)
	}
	if len(diagnostics) != 1 || diagnostics[0].Line != 1 || diagnostics[0].Column != 11 {
		t.Fatalf("wrong diagnostics: %s", output)
	}
}

func TestCheckFile(t *testing.T) {
	file := filepath.Join(t.TempDir(), "broken.ds")
	if err := os.WriteFile(file, []byte("9ayed x = ;"), 0600); err != nil {
		t.Fatal(err)
	}
	_, err := executeCommand(newCheckCommand(), []string{file}, "")
	if err == nil || !strings.Contains(err.Error(), "broken.ds:1:11") {
		t.Fatalf("got %v", err)
	}
}

func TestFormatCommand(t *testing.T) {
	output, err := executeCommand(newFormatCommand(), nil, "golih(dakhel());")
	if err != nil || output != "golih(dakhel());\n" {
		t.Fatalf("got %q, %v", output, err)
	}
	file := filepath.Join(t.TempDir(), "source.ds")
	if err := os.WriteFile(file, []byte("9ayed x=[1,2];"), 0600); err != nil {
		t.Fatal(err)
	}
	if _, err := executeCommand(newFormatCommand(), []string{"--write", file}, ""); err != nil {
		t.Fatal(err)
	}
	formatted, err := os.ReadFile(file)
	if err != nil || string(formatted) != "9ayed x = [1, 2];\n" {
		t.Fatalf("got %q, %v", formatted, err)
	}
	original := "9ayed x = ;"
	if err := os.WriteFile(file, []byte(original), 0600); err != nil {
		t.Fatal(err)
	}
	if _, err := executeCommand(newFormatCommand(), []string{"--write", file}, ""); err == nil {
		t.Fatal("formatted invalid syntax")
	}
	retained, _ := os.ReadFile(file)
	if string(retained) != original {
		t.Fatal("changed invalid source")
	}
}

func TestMultilineREPL(t *testing.T) {
	var output, errors bytes.Buffer
	interactiveCmd.SetIn(strings.NewReader("fonksyon add(x) {\nrjje3 x + 1;\n}\ngolih(add(2));\n9ayed a = [\n1,\n2\n];\ngolih(a);\nexit\n"))
	interactiveCmd.SetOut(&output)
	interactiveCmd.SetErr(&errors)
	defer func() { interactiveCmd.SetIn(nil); interactiveCmd.SetOut(nil); interactiveCmd.SetErr(nil) }()
	if err := interactiveCmd.RunE(interactiveCmd, nil); err != nil {
		t.Fatal(err)
	}
	if errors.Len() != 0 || !strings.Contains(output.String(), "... ") || !strings.Contains(output.String(), "3\n") || !strings.Contains(output.String(), "[1, 2]\n") {
		t.Fatalf("output %q; errors %q", output.String(), errors.String())
	}
}
