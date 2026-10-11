package parser

import "github.com/rabraghib/darijascript/src/lexer"

type Program struct {
	Statements []Statement
}

type Statement interface {
	statementNode()
}

type Expression interface {
	expressionNode()
}

type statementMarker struct{}

func (s *statementMarker) statementNode() {}

type expressionMarker struct{}

func (e *expressionMarker) expressionNode() {}

type IfStatement struct {
	statementMarker
	Condition           Expression
	IsConditionReversed bool
	Consequence         *BlockStatement
	Alternative         *BlockStatement
}

type WhileStatement struct {
	statementMarker
	Condition   Expression
	Consequence *BlockStatement
}

type BlockStatement struct {
	statementMarker
	Statements []Statement
}

type ExpressionStatement struct {
	statementMarker
	Expression Expression
}

type FunctionDeclaration struct {
	statementMarker
	Name       *Identifier
	Parameters []*Identifier
	Body       *BlockStatement
}

type LetStatement struct {
	statementMarker
	Name  *Identifier
	Value Expression
}

type ReturnStatement struct {
	statementMarker
	ReturnValue Expression
}

type ThrowStatement struct {
	statementMarker
	ReturnValue Expression
	Token       *lexer.Token
}

type Identifier struct {
	expressionMarker
	Value string
	Token *lexer.Token
}

type NumberLiteral struct {
	expressionMarker
	Value float64
	Token *lexer.Token
}

type StringLiteral struct {
	expressionMarker
	Value string
	Token *lexer.Token
}

type BooleanLiteral struct {
	expressionMarker
	Value bool
	Token *lexer.Token
}

type PrefixExpression struct {
	expressionMarker
	Operator string
	Right    Expression
	Token    *lexer.Token
}

type AssignmentExpression struct {
	expressionMarker
	Target Expression
	Value  Expression
	Token  *lexer.Token
}

type InfixExpression struct {
	expressionMarker
	Left     Expression
	Operator string
	Right    Expression
	Token    *lexer.Token
}

type CallExpression struct {
	expressionMarker
	Function  Expression
	Arguments []Expression
	Token     *lexer.Token
}

type ArrayLiteral struct {
	expressionMarker
	Elements []Expression
	Token    *lexer.Token
}

type IndexExpression struct {
	expressionMarker
	Left  Expression
	Index Expression
	Token *lexer.Token
}

type HashLiteral struct {
	expressionMarker
	Pairs []HashPair
	Token *lexer.Token
}

type HashPair struct{ Key, Value Expression }

type NullLiteral struct {
	expressionMarker
	Token *lexer.Token
}
type MemberExpression struct {
	expressionMarker
	Left  Expression
	Name  string
	Token *lexer.Token
}
type FunctionLiteral struct {
	expressionMarker
	Parameters []*Identifier
	Body       *BlockStatement
	Token      *lexer.Token
}
type ForStatement struct {
	statementMarker
	Name       *Identifier
	Collection Expression
	Body       *BlockStatement
	Token      *lexer.Token
}
type BreakStatement struct {
	statementMarker
	Token *lexer.Token
}
type ContinueStatement struct {
	statementMarker
	Token *lexer.Token
}
type TryStatement struct {
	statementMarker
	Body    *BlockStatement
	Name    *Identifier
	Handler *BlockStatement
}

type Error struct {
	expressionMarker
	Message string
	Pos     lexer.Position
}
