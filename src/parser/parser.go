package parser

import (
	"fmt"

	"github.com/rabraghib/darijascript/src/lexer"
)

type Parser struct {
	errors []Error
	tokens []*lexer.Token
	curPos int
	depth  int
}

func NewParser(tokens []*lexer.Token) *Parser {
	filtered := make([]*lexer.Token, 0, len(tokens)+1)
	pos := lexer.Position{Line: 1, Column: 1}
	for _, token := range tokens {
		if token == nil {
			continue
		}
		pos = token.Pos
		if token.Type != lexer.TT_COMMENT {
			filtered = append(filtered, token)
		}
	}
	if len(filtered) == 0 || filtered[len(filtered)-1].Type != lexer.TT_EOF {
		filtered = append(filtered, &lexer.Token{Type: lexer.TT_EOF, Pos: pos})
	}
	return &Parser{tokens: filtered}
}

func (p *Parser) errorsOrNil() error {
	if len(p.errors) == 0 {
		return nil
	}
	errorMsg := "Parser errors: "
	for _, err := range p.errors {
		errorMsg += fmt.Sprintf("\n\t%s at %d:%d", err.Message, err.Pos.Line, err.Pos.Column)
	}
	return fmt.Errorf("%s", errorMsg)
}

func (p *Parser) ParseProgram() (*Program, error) {
	program := &Program{}
	program.Statements = []Statement{}

	for !p.isAtEnd() && len(p.errors) == 0 {
		stmt := p.parseStatement()
		if stmt != nil {
			program.Statements = append(program.Statements, stmt)
		}
	}

	return program, p.errorsOrNil()
}

func (p *Parser) isAtEnd() bool {
	return p.curPos >= len(p.tokens) || p.currentToken().Type == lexer.TT_EOF
}

func (p *Parser) currentToken() *lexer.Token {
	if p.curPos >= len(p.tokens) {
		return p.tokens[len(p.tokens)-1]
	}
	return p.tokens[p.curPos]
}

func (p *Parser) parseStatement() Statement {
	p.depth++
	defer func() { p.depth-- }()
	if p.depth > 256 {
		p.addError("maximum parser nesting exceeded", p.currentToken())
		return nil
	}
	currentToken := p.currentToken()
	switch currentToken.Type {
	case lexer.TT_COMMENT:
		p.curPos++ // Skip comment token
		return nil
	case lexer.TT_IF:
		return p.parseIfStatement(false)
	case lexer.TT_IF_NOT:
		return p.parseIfStatement(true)
	case lexer.TT_WHILE:
		return p.parseWhileStatement()
	case lexer.TT_FUNCTION:
		return p.parseFunctionDeclaration()
	case lexer.TT_RETURN:
		return p.parseReturnStatement()
	case lexer.TT_THROW:
		return p.parseThrowStatement()
	case lexer.TT_LET:
		return p.parseLetStatement()
	default:
		return p.parseExpressionStatement()
	}
}

func (p *Parser) parseIfStatement(isNegative bool) *IfStatement {
	p.consumeToken() // consume TT_IF or TT_IF_NOT token
	ifStatement := &IfStatement{}
	ifStatement.IsConditionReversed = isNegative
	ifStatement.Condition = p.parseExpression()
	ifStatement.Consequence = p.parseBlockStatement()
	if p.matchToken(lexer.TT_ELSE) {
		if p.checkToken(lexer.TT_IF) {
			ifStatement.Alternative = &BlockStatement{
				Statements: []Statement{p.parseStatement()},
			}
		} else if p.checkToken(lexer.TT_IF_NOT) {
			ifStatement.Alternative = &BlockStatement{
				Statements: []Statement{p.parseStatement()},
			}
		} else {
			ifStatement.Alternative = p.parseBlockStatement()
		}
	}
	return ifStatement
}

func (p *Parser) parseWhileStatement() *WhileStatement {
	p.consumeToken() // consume TT_WHILE token
	whileStatement := &WhileStatement{}
	whileStatement.Condition = p.parseExpression()
	whileStatement.Consequence = p.parseBlockStatement()
	return whileStatement
}

func (p *Parser) parseBlockStatement() *BlockStatement {
	p.consumeTokenOfType(lexer.TT_LBRACE)
	blockStatement := &BlockStatement{}
	blockStatement.Statements = []Statement{}
	for !p.checkToken(lexer.TT_RBRACE) && !p.isAtEnd() && len(p.errors) == 0 {
		stmt := p.parseStatement()
		if stmt != nil {
			blockStatement.Statements = append(blockStatement.Statements, stmt)
		}
	}
	p.consumeTokenOfType(lexer.TT_RBRACE)
	return blockStatement
}

func (p *Parser) parseExpressionStatement() *ExpressionStatement {
	expressionStatement := &ExpressionStatement{
		Expression: p.parseExpression(),
	}
	p.consumeTokenOfType(lexer.TT_SEMICOLON)
	return expressionStatement
}

func (p *Parser) parseFunctionDeclaration() *FunctionDeclaration {
	functionDeclaration := &FunctionDeclaration{}
	p.consumeToken() // consume TT_FUNCTION token
	functionDeclaration.Name = p.parseIdentifier()
	p.consumeTokenOfType(lexer.TT_LPAREN)
	functionDeclaration.Parameters = p.parseFunctionParameters()
	functionDeclaration.Body = p.parseBlockStatement()
	return functionDeclaration
}

func (p *Parser) parseReturnStatement() *ReturnStatement {
	p.consumeToken() // consume TT_RETURN token
	returnStatement := &ReturnStatement{}
	if !p.checkToken(lexer.TT_SEMICOLON) {
		returnStatement.ReturnValue = p.parseExpression()
	}
	p.consumeTokenOfType(lexer.TT_SEMICOLON)
	return returnStatement
}

func (p *Parser) parseThrowStatement() Statement {
	p.consumeToken() // consume TT_THROW token
	throwStatement := &ThrowStatement{}
	throwStatement.ReturnValue = p.parseExpression()
	p.consumeTokenOfType(lexer.TT_SEMICOLON)
	return throwStatement
}

func (p *Parser) parseLetStatement() Statement {
	p.consumeToken() // consume TT_LET token
	name := p.parseIdentifier()
	p.consumeTokenOfType(lexer.TT_ASSIGN)
	value := p.parseExpression()
	p.consumeTokenOfType(lexer.TT_SEMICOLON)
	return &LetStatement{
		Name:  name,
		Value: value,
	}
}

func (p *Parser) parseIdentifier() *Identifier {
	token := p.consumeTokenOfType(lexer.TT_IDENTIFIER)
	return &Identifier{
		Value: token.Literal,
		Token: token,
	}
}

func (p *Parser) parseFunctionParameters() []*Identifier {
	parameters := []*Identifier{}
	seen := map[string]bool{}
	if p.checkToken(lexer.TT_IDENTIFIER) {
		parameters = append(parameters, p.parseIdentifier())
		for p.matchToken(lexer.TT_COMMA) {
			parameters = append(parameters, p.parseIdentifier())
		}
	}
	for _, parameter := range parameters {
		if seen[parameter.Value] {
			p.addError("duplicate parameter: "+parameter.Value, parameter.Token)
		}
		seen[parameter.Value] = true
	}
	p.consumeTokenOfType(lexer.TT_RPAREN)
	return parameters
}

// Utility methods
func (p *Parser) consumeToken() *lexer.Token {
	token := p.currentToken()
	if !p.isAtEnd() {
		p.curPos++
	}
	return token
}

func (p *Parser) consumeTokenOfType(tokenType lexer.TokenType) *lexer.Token {
	token := p.consumeToken()
	if token.Type != tokenType {
		p.errors = append(p.errors, Error{
			Message: fmt.Sprintf("Expected token of type %s, got %s", tokenType, token.Type),
			Pos:     token.Pos,
		})
		return token
	}
	return token
}

func (p *Parser) matchToken(tokenType lexer.TokenType) bool {
	if p.checkToken(tokenType) {
		p.consumeToken()
		return true
	}
	return false
}

func (p *Parser) checkToken(tokenType lexer.TokenType) bool {
	return p.currentToken().Type == tokenType
}

func (p *Parser) addError(message string, token *lexer.Token) {
	p.errors = append(p.errors, Error{Message: message, Pos: token.Pos})
}
