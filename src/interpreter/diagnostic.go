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
	Frames   []string
}

func (err *RuntimeError) Error() string {
	message := fmt.Sprintf("%s at %d:%d", err.Cause, err.Position.Line, err.Position.Column)
	for _, frame := range err.Frames {
		message += "\n  called from " + frame
	}
	return message
}
func (err *RuntimeError) Unwrap() error { return err.Cause }

func callError(err error, name string, token *lexer.Token) error {
	err = atToken(err, token)
	var positioned *RuntimeError
	if errors.As(err, &positioned) {
		positioned.Frames = append(positioned.Frames, fmt.Sprintf("%s at %d:%d", name, token.Pos.Line, token.Pos.Column))
	}
	return err
}

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
	case *parser.NullLiteral:
		return expr.Token
	case *parser.ArrayLiteral:
		return expr.Token
	case *parser.HashLiteral:
		return expr.Token
	case *parser.IndexExpression:
		return expr.Token
	case *parser.MemberExpression:
		return expr.Token
	case *parser.FunctionLiteral:
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
