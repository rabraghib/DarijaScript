// Package runtime provides execution shared by the CLI, REPL, and WASM hosts.
package runtime

import (
	"fmt"
	"unicode/utf8"

	"github.com/rabraghib/darijascript/src/interpreter"
	"github.com/rabraghib/darijascript/src/lexer"
	"github.com/rabraghib/darijascript/src/parser"
)

func Run(source string, evaluator *interpreter.Evaluator) (interface{}, error) {
	if !utf8.ValidString(source) {
		return nil, fmt.Errorf("source is not valid UTF-8")
	}
	l := lexer.NewLexer(source)
	var tokens []*lexer.Token
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
	program, err := parser.NewParser(tokens).ParseProgram()
	if err != nil {
		return nil, err
	}
	return evaluator.EvaluateProgram(program)
}
