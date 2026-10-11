package parser

import (
	"fmt"
	"strconv"

	"github.com/rabraghib/darijascript/src/lexer"
)

// Higher numbers bind more tightly. Only assignment associates to the right.
func precedence(kind lexer.TokenType) int {
	switch kind {
	case lexer.TT_ASSIGN:
		return 1
	case lexer.TT_OR:
		return 2
	case lexer.TT_AND:
		return 3
	case lexer.TT_EQUAL, lexer.TT_NOT_EQUAL:
		return 4
	case lexer.TT_LESS_THAN, lexer.TT_LESS_EQUAL, lexer.TT_GREATER_THAN, lexer.TT_GREATER_EQUAL:
		return 5
	case lexer.TT_PLUS, lexer.TT_MINUS:
		return 6
	case lexer.TT_MULTIPLY, lexer.TT_DIVIDE, lexer.TT_MODULO:
		return 7
	case lexer.TT_LPAREN:
		return 9
	default:
		return 0
	}
}

func (p *Parser) parseExpression() Expression { return p.parsePrecedence(1) }

func (p *Parser) parsePrecedence(min int) Expression {
	p.depth++
	defer func() { p.depth-- }()
	if p.depth > 256 {
		p.addError("maximum parser nesting exceeded", p.currentToken())
		return nil
	}
	left := p.parsePrefix()
	for len(p.errors) == 0 && precedence(p.currentToken().Type) >= min {
		token := p.consumeToken()
		if token.Type == lexer.TT_LPAREN {
			function, ok := left.(*Identifier)
			if !ok {
				p.addError("only named functions can be called", token)
				return nil
			}
			left = p.parseCall(function, token)
			continue
		}
		next := precedence(token.Type) + 1
		if token.Type == lexer.TT_ASSIGN {
			next--
		}
		right := p.parsePrecedence(next)
		if token.Type == lexer.TT_ASSIGN {
			variable, ok := left.(*Identifier)
			if !ok {
				p.addError("assignment requires a variable", token)
				return nil
			}
			left = &AssignmentExpression{Variable: *variable, Value: right, Token: token}
		} else {
			left = &InfixExpression{Left: left, Operator: token.Literal, Right: right, Token: token}
		}
	}
	return left
}

func (p *Parser) parsePrefix() Expression {
	token := p.consumeToken()
	switch token.Type {
	case lexer.TT_IDENTIFIER:
		return &Identifier{Value: token.Literal, Token: token}
	case lexer.TT_NUMBER:
		value, err := strconv.ParseFloat(token.Literal, 64)
		if err != nil {
			p.addError("invalid number", token)
			return nil
		}
		return &NumberLiteral{Value: value, Token: token}
	case lexer.TT_STRING:
		return &StringLiteral{Value: token.Literal, Token: token}
	case lexer.TT_TRUE, lexer.TT_FALSE:
		return &BooleanLiteral{Value: token.Type == lexer.TT_TRUE, Token: token}
	case lexer.TT_PLUS, lexer.TT_MINUS, lexer.TT_NOT:
		return &PrefixExpression{Operator: token.Literal, Right: p.parsePrecedence(8), Token: token}
	case lexer.TT_LPAREN:
		expression := p.parseExpression()
		p.consumeTokenOfType(lexer.TT_RPAREN)
		return expression
	default:
		p.addError(fmt.Sprintf("expected expression, got %s", token.Type), token)
		return nil
	}
}

func (p *Parser) parseCall(function *Identifier, token *lexer.Token) Expression {
	call := &CallExpression{Function: function, Token: token, Arguments: []Expression{}}
	if p.matchToken(lexer.TT_RPAREN) {
		return call
	}
	for !p.isAtEnd() && len(p.errors) == 0 {
		call.Arguments = append(call.Arguments, p.parseExpression())
		if !p.matchToken(lexer.TT_COMMA) {
			break
		}
	}
	p.consumeTokenOfType(lexer.TT_RPAREN)
	return call
}
