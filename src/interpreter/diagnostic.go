package interpreter

import (
	"errors"
	"fmt"

	"github.com/rabraghib/darijascript/src/lexer"
	"github.com/rabraghib/darijascript/src/parser"
)

type RuntimeError struct {
	Position lexer.Position
	Cause    error
}

func (err *RuntimeError) Error() string {
	return fmt.Sprintf("%s at %d:%d", err.Cause, err.Position.Line, err.Position.Column)
}
func (err *RuntimeError) Unwrap() error { return err.Cause }

func atToken(err error, token *lexer.Token) error {
	var positioned *RuntimeError
	if err == nil || token == nil || errors.As(err, &positioned) {
		return err
	}
	return &RuntimeError{Position: token.Pos, Cause: err}
}

func expressionToken(expression parser.Expression) *lexer.Token {
	switch expr := expression.(type) {
	case *parser.Identifier:
		return expr.Token
	case *parser.NumberLiteral:
		return expr.Token
	case *parser.StringLiteral:
		return expr.Token
	case *parser.BooleanLiteral:
		return expr.Token
	case *parser.PrefixExpression:
		return expr.Token
	case *parser.InfixExpression:
		return expr.Token
	case *parser.AssignmentExpression:
		return expr.Token
	case *parser.CallExpression:
		return expr.Token
	default:
		return nil
	}
}
