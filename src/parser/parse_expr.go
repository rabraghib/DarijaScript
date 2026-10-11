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
	case lexer.TT_LPAREN, lexer.TT_LBRACKET, lexer.TT_DOT:
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
			left = p.parseCall(left, token)
			continue
		}
		if token.Type == lexer.TT_LBRACKET {
			index := p.parseExpression()
			p.consumeTokenOfType(lexer.TT_RBRACKET)
			left = &IndexExpression{Left: left, Index: index, Token: token}
			continue
		}
		if token.Type == lexer.TT_DOT {
			name := p.parseIdentifier()
			left = &MemberExpression{Left: left, Name: name.Value, Token: token}
			continue
		}
		next := precedence(token.Type) + 1
		if token.Type == lexer.TT_ASSIGN {
			next--
		}
		right := p.parsePrecedence(next)
		if token.Type == lexer.TT_ASSIGN {
			switch left.(type) {
			case *Identifier, *IndexExpression:
			default:
				p.addError("assignment requires a variable or collection index", token)
				return nil
			}
			left = &AssignmentExpression{Target: left, Value: right, Token: token}
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
	case lexer.TT_NULL:
		return &NullLiteral{Token: token}
	case lexer.TT_LBRACKET:
		return &ArrayLiteral{Elements: p.parseArguments(lexer.TT_RBRACKET), Token: token}
	case lexer.TT_LBRACE:
		hash := &HashLiteral{Token: token}
		for !p.checkToken(lexer.TT_RBRACE) && !p.isAtEnd() && len(p.errors) == 0 {
			key := p.parseExpression()
			p.consumeTokenOfType(lexer.TT_COLON)
			value := p.parseExpression()
			hash.Pairs = append(hash.Pairs, HashPair{Key: key, Value: value})
			if !p.matchToken(lexer.TT_COMMA) {
				break
			}
		}
		p.consumeTokenOfType(lexer.TT_RBRACE)
		return hash
	case lexer.TT_FUNCTION:
		p.consumeTokenOfType(lexer.TT_LPAREN)
		parameters := p.parseFunctionParameters()
		return &FunctionLiteral{Parameters: parameters, Body: p.parseFunctionBody(), Token: token}
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

func (p *Parser) parseCall(function Expression, token *lexer.Token) Expression {
	return &CallExpression{Function: function, Token: token, Arguments: p.parseArguments(lexer.TT_RPAREN)}
}

func (p *Parser) parseArguments(closing lexer.TokenType) []Expression {
	arguments := []Expression{}
	for !p.isAtEnd() && len(p.errors) == 0 {
		if p.checkToken(closing) {
			break
		}
		arguments = append(arguments, p.parseExpression())
		if !p.matchToken(lexer.TT_COMMA) {
			break
		}
	}
	p.consumeTokenOfType(closing)
	return arguments
}
