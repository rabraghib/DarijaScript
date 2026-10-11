package lexer

import "fmt"

type TokenType int

const (
	TT_NUMBER TokenType = iota
	TT_STRING
	TT_IDENTIFIER
	TT_IF
	TT_IF_NOT
	TT_ELSE
	TT_COMMENT
	TT_WHILE
	TT_FUNCTION
	TT_TRUE
	TT_FALSE
	TT_ASSIGN
	TT_LET
	TT_RETURN
	TT_THROW

	TT_PLUS
	TT_MINUS
	TT_MULTIPLY
	TT_DIVIDE
	TT_MODULO
	TT_EQUAL
	TT_NOT_EQUAL
	TT_LESS_THAN
	TT_LESS_EQUAL
	TT_GREATER_THAN
	TT_GREATER_EQUAL
	TT_AND
	TT_OR
	TT_NOT

	TT_LPAREN   // (
	TT_RPAREN   // )
	TT_LBRACE   // {
	TT_RBRACE   // }
	TT_LBRACKET // [
	TT_RBRACKET // ]
	TT_COMMA
	TT_COLON
	TT_SEMICOLON
	TT_EOF
	TT_NULL
	TT_FOR
	TT_BREAK
	TT_CONTINUE
	TT_TRY
	TT_CATCH
	TT_DOT
)

type Position struct {
	Line   int
	Column int
}

type Token struct {
	Type    TokenType
	Literal string
	Pos     Position
}

func (t *Token) String() string {
	return fmt.Sprintf(
		"<type:%v, l:c:%v:%v, str:%v>",
		t.Type.String(), t.Pos.Line, t.Pos.Column, t.Literal,
	)
}

func (t TokenType) String() string {
	switch t {
	case TT_NUMBER:
		return "TT_NUMBER"
	case TT_STRING:
		return "TT_STRING"
	case TT_IDENTIFIER:
		return "TT_IDENTIFIER"
	case TT_IF:
		return "TT_IF"
	case TT_IF_NOT:
		return "TT_IF_NOT"
	case TT_ELSE:
		return "TT_ELSE"
	case TT_COMMENT:
		return "TT_COMMENT"
	case TT_WHILE:
		return "TT_WHILE"
	case TT_FUNCTION:
		return "TT_FUNCTION"
	case TT_TRUE:
		return "TT_TRUE"
	case TT_FALSE:
		return "TT_FALSE"
	case TT_ASSIGN:
		return "TT_ASSIGN"
	case TT_LET:
		return "TT_LET"
	case TT_RETURN:
		return "TT_RETURN"
	case TT_THROW:
		return "TT_THROW"
	case TT_PLUS:
		return "TT_PLUS"
	case TT_MINUS:
		return "TT_MINUS"
	case TT_MULTIPLY:
		return "TT_MULTIPLY"
	case TT_DIVIDE:
		return "TT_DIVIDE"
	case TT_MODULO:
		return "TT_MODULO"
	case TT_EQUAL:
		return "TT_EQUAL"
	case TT_NOT_EQUAL:
		return "TT_NOT_EQUAL"
	case TT_LESS_THAN:
		return "TT_LESS_THAN"
	case TT_LESS_EQUAL:
		return "TT_LESS_EQUAL"
	case TT_GREATER_THAN:
		return "TT_GREATER_THAN"
	case TT_GREATER_EQUAL:
		return "TT_GREATER_EQUAL"
	case TT_AND:
		return "TT_AND"
	case TT_OR:
		return "TT_OR"
	case TT_NOT:
		return "TT_NOT"
	case TT_LPAREN:
		return "TT_LPAREN"
	case TT_RPAREN:
		return "TT_RPAREN"
	case TT_LBRACE:
		return "TT_LBRACE"
	case TT_RBRACE:
		return "TT_RBRACE"
	case TT_LBRACKET:
		return "TT_LBRACKET"
	case TT_RBRACKET:
		return "TT_RBRACKET"
	case TT_COMMA:
		return "TT_COMMA"
	case TT_COLON:
		return "TT_COLON"
	case TT_SEMICOLON:
		return "TT_SEMICOLON"
	case TT_EOF:
		return "TT_EOF"
	case TT_NULL:
		return "walo"
	case TT_FOR:
		return "kol"
	case TT_BREAK:
		return "w9ef"
	case TT_CONTINUE:
		return "kmmel"
	case TT_TRY:
		return "jerreb"
	case TT_CATCH:
		return "ched"
	case TT_DOT:
		return "."
	default:
		return "Unknown"
	}
}
