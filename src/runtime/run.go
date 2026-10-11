// Package runtime provides execution shared by the CLI, REPL, and WASM hosts.
package runtime

import (
	"errors"
	"fmt"
	"regexp"
	"strconv"
	"strings"
	"unicode/utf8"

	"github.com/rabraghib/darijascript/src/interpreter"
	"github.com/rabraghib/darijascript/src/lexer"
	"github.com/rabraghib/darijascript/src/parser"
)

func tokenize(source string) ([]*lexer.Token, error) {
	if !utf8.ValidString(source) {
		return nil, fmt.Errorf("source is not valid UTF-8")
	}
	l := lexer.NewLexer(source)
	tokens := []*lexer.Token{}
	for {
		token, err := l.NextToken()
		if err != nil {
			return nil, err
		}
		tokens = append(tokens, token)
		if token.Type == lexer.TT_EOF {
			break
		}
	}
	return tokens, nil
}

func Parse(source string) (*parser.Program, error) {
	tokens, err := tokenize(source)
	if err != nil {
		return nil, err
	}
	return parser.NewParser(tokens).ParseProgram()
}

func Incomplete(source string) bool {
	tokens, err := tokenize(source)
	if err != nil {
		return strings.HasPrefix(err.Error(), "unterminated string") || strings.HasPrefix(err.Error(), "unterminated comment")
	}
	p := parser.NewParser(tokens)
	_, err = p.ParseProgram()
	return err != nil && p.Incomplete
}

func Run(source string, evaluator *interpreter.Evaluator) (interface{}, error) {
	return RunNamed(source, "<input>", evaluator)
}

func RunNamed(source, name string, evaluator *interpreter.Evaluator) (interface{}, error) {
	program, err := Parse(source)
	if err == nil {
		var value interface{}
		value, err = evaluator.EvaluateProgram(program)
		if err == nil {
			return value, nil
		}
	}
	return nil, &SourceError{Name: name, Source: source, Cause: err}
}

type SourceError struct {
	Name, Source string
	Cause        error
}

func (err *SourceError) Unwrap() error { return err.Cause }
func (err *SourceError) Error() string {
	pos := errorPosition(err.Cause)
	lines := strings.Split(err.Source, "\n")
	message := fmt.Sprintf("%s:%d:%d: %s", err.Name, pos.Line, pos.Column, err.Cause)
	if pos.Line > 0 && pos.Line <= len(lines) {
		line := strings.TrimSuffix(lines[pos.Line-1], "\r")
		// Retain tabs in the prefix so the caret follows the rendered source.
		prefix := []rune(line)
		column := pos.Column - 1
		if column > len(prefix) {
			column = len(prefix)
		}
		var caret strings.Builder
		for _, char := range prefix[:column] {
			if char == '\t' {
				caret.WriteRune('\t')
			} else {
				caret.WriteRune(' ')
			}
		}
		message += "\n  " + line + "\n  " + caret.String() + "^"
	}
	return message
}

var positionPattern = regexp.MustCompile("(?:at |at line )([0-9]+)(?::| column )([0-9]+)")

func errorPosition(err error) lexer.Position {
	var positioned *interpreter.RuntimeError
	if errors.As(err, &positioned) {
		return positioned.Position
	}
	match := positionPattern.FindStringSubmatch(err.Error())
	if len(match) == 3 {
		line, _ := strconv.Atoi(match[1])
		column, _ := strconv.Atoi(match[2])
		return lexer.Position{Line: line, Column: column}
	}
	return lexer.Position{Line: 1, Column: 1}
}

type Diagnostic struct {
	Message string `json:"message"`
	Line    int    `json:"line"`
	Column  int    `json:"column"`
}

func Check(source string) []Diagnostic {
	diagnostics := []Diagnostic{}
	_, err := Parse(source)
	if err != nil {
		pos := errorPosition(err)
		diagnostics = append(diagnostics, Diagnostic{Message: err.Error(), Line: pos.Line, Column: pos.Column})
	}
	return diagnostics
}
