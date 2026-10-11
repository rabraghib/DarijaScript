package lexer

import (
	"fmt"
	"regexp"
	"unicode"
	"unicode/utf8"
)

type Lexer struct {
	input        string
	position     int
	readPosition int
	ch           rune
	line         int
	column       int
}

func NewLexer(input string) *Lexer {
	l := &Lexer{input: input, line: 1, column: 0}
	l.readChar()
	return l
}

func isLetter(ch rune) bool {
	return unicode.IsLetter(ch) || unicode.IsMark(ch) || ch == '_'
}

func isDigit(ch rune) bool {
	return ch >= '0' && ch <= '9'
}

func newToken(tokenType TokenType, literal string, pos Position) *Token {
	return &Token{Type: tokenType, Literal: literal, Pos: pos}
}

func (l *Lexer) IsEOL() bool {
	l.skipWhitespace()
	return l.position >= len(l.input)
}

func (l *Lexer) NextToken() (*Token, error) {
	l.skipWhitespace()
	var position = Position{l.line, l.column}

	if l.ch == '#' || (l.ch == '/' && l.peekChar() == '/') {
		if l.ch == '/' {
			l.readChar()
			l.readChar()
		} else {
			l.readChar()
		}
		return newToken(TT_COMMENT, l.readComment(), position), nil
	}

	if l.ch == '/' && l.peekChar() == '*' {
		l.readChar()
		l.readChar()
		literal, err := l.readMultiLineComment()
		if err != nil {
			return nil, fmt.Errorf("%w at %d:%d", err, position.Line, position.Column)
		}
		return newToken(TT_COMMENT, literal, position), nil
	}

	if l.ch == '"' {
		literal, err := l.readString()
		if err != nil {
			return nil, fmt.Errorf("%w at %d:%d", err, position.Line, position.Column)
		}
		return newToken(TT_STRING, literal, position), nil
	}

	if isDigit(l.ch) {
		if l.ch == '9' && l.peekChar() == 'a' {
			var identifier = l.readIdentifier()
			if identifier == "9ayed" {
				return newToken(TT_LET, identifier, position), nil
			}
			return nil, fmt.Errorf("invalid identifier %q at %d:%d", identifier, position.Line, position.Column)
		} else {
			return newToken(TT_NUMBER, l.readNumber(), position), nil
		}
	}

	if isLetter(l.ch) {
		var identifier = l.readIdentifier()

		switch identifier {
		case "ilakan":
			return newToken(TT_IF, identifier, position), nil
		case "ilamakanch":
			return newToken(TT_IF_NOT, identifier, position), nil
		case "sinn":
			return newToken(TT_ELSE, identifier, position), nil
		case "ma7dBa9i":
			return newToken(TT_WHILE, identifier, position), nil
		case "fonksyon":
			return newToken(TT_FUNCTION, identifier, position), nil
		case "true", "s7i7":
			return newToken(TT_TRUE, identifier, position), nil
		case "false", "ghalt":
			return newToken(TT_FALSE, identifier, position), nil
		case "rjje3":
			return newToken(TT_RETURN, identifier, position), nil
		default:
			if throwKeyword.MatchString(identifier) {
				return newToken(TT_THROW, identifier, position), nil
			}
			return newToken(TT_IDENTIFIER, identifier, position), nil
		}
	}

	switch l.ch {
	case '+':
		l.readChar()
		return newToken(TT_PLUS, "+", position), nil
	case '-':
		l.readChar()
		return newToken(TT_MINUS, "-", position), nil
	case '*':
		l.readChar()
		return newToken(TT_MULTIPLY, "*", position), nil
	case '/':
		l.readChar()
		return newToken(TT_DIVIDE, "/", position), nil
	case '%':
		l.readChar()
		return newToken(TT_MODULO, "%", position), nil
	case '=':
		l.readChar()
		if l.ch == '=' {
			l.readChar()
			return newToken(TT_EQUAL, "==", position), nil
		}
		return newToken(TT_ASSIGN, "=", position), nil
	case '&':
		l.readChar()
		if l.ch == '&' {
			l.readChar()
			return newToken(TT_AND, "&&", position), nil
		}
		return nil, fmt.Errorf("unknown token: & at %d:%d", position.Line, position.Column)
	case '|':
		l.readChar()
		if l.ch == '|' {
			l.readChar()
			return newToken(TT_OR, "||", position), nil
		}
		return nil, fmt.Errorf("unknown token: | at %d:%d", position.Line, position.Column)
	case '!':
		l.readChar()
		if l.ch == '=' {
			l.readChar()
			return newToken(TT_NOT_EQUAL, "!=", position), nil
		}
		return newToken(TT_NOT, "!", position), nil
	case '<':
		l.readChar()
		if l.ch == '=' {
			l.readChar()
			return newToken(TT_LESS_EQUAL, "<=", position), nil
		}
		return newToken(TT_LESS_THAN, "<", position), nil
	case '>':
		l.readChar()
		if l.ch == '=' {
			l.readChar()
			return newToken(TT_GREATER_EQUAL, ">=", position), nil
		}
		return newToken(TT_GREATER_THAN, ">", position), nil
	case '(':
		l.readChar()
		return newToken(TT_LPAREN, "(", position), nil
	case ')':
		l.readChar()
		return newToken(TT_RPAREN, ")", position), nil
	case '{':
		l.readChar()
		return newToken(TT_LBRACE, "{", position), nil
	case '}':
		l.readChar()
		return newToken(TT_RBRACE, "}", position), nil
	case '[':
		l.readChar()
		return newToken(TT_LBRACKET, "[", position), nil
	case ']':
		l.readChar()
		return newToken(TT_RBRACKET, "]", position), nil
	case ',':
		l.readChar()
		return newToken(TT_COMMA, ",", position), nil
	case ';':
		l.readChar()
		return newToken(TT_SEMICOLON, ";", position), nil
	case ':':
		l.readChar()
		return newToken(TT_COLON, ":", position), nil

	}

	if l.position >= len(l.input) {
		return newToken(TT_EOF, "", position), nil
	}

	return nil, fmt.Errorf("unknown token: %s at line %d column %d", string(l.ch), l.line, l.column)
}

func (l *Lexer) peekChar() rune {
	if l.readPosition >= len(l.input) {
		return 0
	}
	ch, _ := utf8.DecodeRuneInString(l.input[l.readPosition:])
	return ch
}

var throwKeyword = regexp.MustCompile("^WAA[A]*$")
