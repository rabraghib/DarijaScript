package runtime

import (
	"strings"

	"github.com/rabraghib/darijascript/src/lexer"
)

// Format validates first and retains comments and string values.
func Format(source string) (string, error) {
	if _, err := Parse(source); err != nil {
		return "", err
	}
	tokens, err := tokenize(source)
	if err != nil {
		return "", err
	}
	var output strings.Builder
	line := ""
	indent := 0
	flush := func() {
		if line != "" {
			output.WriteString(strings.Repeat("    ", indent) + strings.TrimSpace(line) + "\n")
			line = ""
		}
	}
	for i, token := range tokens {
		literal := token.Literal
		switch token.Type {
		case lexer.TT_EOF:
			flush()
		case lexer.TT_COMMENT:
			if line != "" {
				line += " "
			}
			line += literal
			flush()
		case lexer.TT_STRING:
			literal = "\"" + strings.NewReplacer("\\", "\\\\", "\"", "\\\"", "\n", "\\n", "\r", "\\r", "\t", "\\t").Replace(literal) + "\""
			if line != "" && tokens[i-1].Type != lexer.TT_LPAREN && tokens[i-1].Type != lexer.TT_LBRACKET {
				line += " "
			}
			line += literal
		case lexer.TT_LBRACE:
			if line != "" {
				line += " "
			}
			line += "{"
			flush()
			indent++
		case lexer.TT_RBRACE:
			flush()
			indent--
			line = "}"
			next := tokens[i+1].Type
			switch next {
			case lexer.TT_ELSE, lexer.TT_CATCH, lexer.TT_SEMICOLON, lexer.TT_COMMA, lexer.TT_RPAREN, lexer.TT_RBRACKET,
				lexer.TT_LPAREN, lexer.TT_LBRACKET, lexer.TT_DOT, lexer.TT_COLON, lexer.TT_ASSIGN,
				lexer.TT_PLUS, lexer.TT_MINUS, lexer.TT_MULTIPLY, lexer.TT_DIVIDE, lexer.TT_MODULO,
				lexer.TT_EQUAL, lexer.TT_NOT_EQUAL, lexer.TT_AND, lexer.TT_OR,
				lexer.TT_LESS_THAN, lexer.TT_LESS_EQUAL, lexer.TT_GREATER_THAN, lexer.TT_GREATER_EQUAL:
			default:
				flush()
			}
		case lexer.TT_SEMICOLON:
			line += ";"
			flush()
		case lexer.TT_COMMA:
			line += ","
		case lexer.TT_COLON:
			line += ":"
		case lexer.TT_DOT:
			line += "."
		case lexer.TT_RPAREN, lexer.TT_RBRACKET:
			line += literal
		default:
			if line != "" {
				previous := tokens[i-1].Type
				space := previous != lexer.TT_LPAREN && previous != lexer.TT_LBRACKET && previous != lexer.TT_DOT
				if token.Type == lexer.TT_LPAREN {
					space = previous == lexer.TT_IF || previous == lexer.TT_IF_NOT || previous == lexer.TT_WHILE || previous == lexer.TT_FOR || previous == lexer.TT_CATCH
				}
				if token.Type == lexer.TT_LBRACKET {
					space = previous != lexer.TT_IDENTIFIER && previous != lexer.TT_RPAREN && previous != lexer.TT_RBRACKET && previous != lexer.TT_RBRACE && previous != lexer.TT_LPAREN && previous != lexer.TT_LBRACKET
				}
				if space {
					line += " "
				}
			}
			line += literal
		}
	}
	return output.String(), nil
}
