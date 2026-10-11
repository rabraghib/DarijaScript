package src

import (
	"fmt"
	"os"
	"testing"

	"github.com/rabraghib/darijascript/src/lexer"
	"github.com/rabraghib/darijascript/src/parser"
)

func TestParser(t *testing.T) {
	// Test case 1
	// Test case input
	input := `
		9ayed x = 5;
		9ayed y = 10;
		9ayed foobar = s7i7;
		golih("Somme: " + (x + y));
	`

	// Expected output
	expected := &parser.Program{
		Statements: []parser.Statement{
			&parser.LetStatement{
				Name: &parser.Identifier{
					Token: &lexer.Token{Type: lexer.TT_IDENTIFIER, Literal: "x"},
					Value: "x",
				},
				Value: &parser.NumberLiteral{
					Token: &lexer.Token{Type: lexer.TT_NUMBER, Literal: "5"},
					Value: 5,
				},
			},
			&parser.LetStatement{
				Name: &parser.Identifier{
					Token: &lexer.Token{Type: lexer.TT_IDENTIFIER, Literal: "y"},
					Value: "y",
				},
				Value: &parser.NumberLiteral{
					Token: &lexer.Token{Type: lexer.TT_NUMBER, Literal: "10"},
					Value: 10,
				},
			},
			&parser.LetStatement{
				Name: &parser.Identifier{
					Token: &lexer.Token{Type: lexer.TT_IDENTIFIER, Literal: "foobar"},
					Value: "foobar",
				},
				Value: &parser.BooleanLiteral{
					Token: &lexer.Token{Type: lexer.TT_TRUE, Literal: "s7i7"},
					Value: true,
				},
			},
			&parser.ExpressionStatement{
				Expression: &parser.CallExpression{
					Token: &lexer.Token{Type: lexer.TT_LPAREN, Literal: "("},
					Function: &parser.Identifier{
						Token: &lexer.Token{Type: lexer.TT_IDENTIFIER, Literal: "golih"},
						Value: "golih",
					},
					Arguments: []parser.Expression{
						&parser.InfixExpression{
							Token: &lexer.Token{Type: lexer.TT_PLUS, Literal: "+"},
							Left: &parser.StringLiteral{
								Token: &lexer.Token{Type: lexer.TT_STRING, Literal: "Somme: "},
								Value: "Somme: ",
							},
							Operator: "+",
							Right: &parser.InfixExpression{
								Token: &lexer.Token{Type: lexer.TT_PLUS, Literal: "+"},
								Left: &parser.Identifier{
									Token: &lexer.Token{Type: lexer.TT_IDENTIFIER, Literal: "x"},
									Value: "x",
								},
								Operator: "+",
								Right: &parser.Identifier{
									Value: "y",
									Token: &lexer.Token{Type: lexer.TT_IDENTIFIER, Literal: "y"},
								},
							},
						},
					},
				},
			},
		},
	}

	runParserTest(t, input, expected)
}

func runParserTest(t *testing.T, input string, expected *parser.Program) {
	l := lexer.NewLexer(input)
	tokens := []*lexer.Token{}
	for !l.IsEOL() {
		token, err := l.NextToken()
		if err != nil {
			fmt.Println("Error lexing source code:", err)
			os.Exit(1)
		}
		tokens = append(tokens, token)
	}
	p := parser.NewParser(tokens)
	program, err := p.ParseProgram()
	if err != nil {
		t.Fatalf("Error parsing program: %v", err)
	}
	if len(program.Statements) != len(expected.Statements) {
		t.Fatalf("Expected %d statements, got %d", len(expected.Statements), len(program.Statements))
	}
	for i, statement := range program.Statements {
		checkStatement(t, statement, expected.Statements[i])
	}
}

func checkStatement(t *testing.T, statement parser.Statement, expected parser.Statement) {
	switch s := statement.(type) {
	case *parser.LetStatement:
		checkLetStatement(t, s, expected.(*parser.LetStatement))
	case *parser.ExpressionStatement:
		checkExpressionStatement(t, s, expected.(*parser.ExpressionStatement))
	default:
		t.Fatalf("Unexpected statement type: %T", s)
	}
}

func checkLetStatement(t *testing.T, statement *parser.LetStatement, expected *parser.LetStatement) {
	checkIdentifier(t, statement.Name, expected.Name)
	checkExpression(t, statement.Value, expected.Value)
}

func checkExpressionStatement(t *testing.T, statement *parser.ExpressionStatement, expected *parser.ExpressionStatement) {
	checkExpression(t, statement.Expression, expected.Expression)
}

func checkExpression(t *testing.T, expression parser.Expression, expected parser.Expression) {
	switch e := expression.(type) {
	case *parser.NumberLiteral:
		checkNumberLiteral(t, e, expected.(*parser.NumberLiteral))
	case *parser.BooleanLiteral:
		checkBooleanLiteral(t, e, expected.(*parser.BooleanLiteral))
	case *parser.StringLiteral:
		checkStringLiteral(t, e, expected.(*parser.StringLiteral))
	case *parser.Identifier:
		checkIdentifier(t, e, expected.(*parser.Identifier))
	case *parser.CallExpression:
		checkCallExpression(t, e, expected.(*parser.CallExpression))
	case *parser.InfixExpression:
		checkInfixExpression(t, e, expected.(*parser.InfixExpression))
	default:
		t.Fatalf("Unexpected expression type: %T", e)
	}
}

func checkNumberLiteral(t *testing.T, expression *parser.NumberLiteral, expected *parser.NumberLiteral) {
	if expression.Value != expected.Value {
		t.Fatalf("Expected number value %f, got %f", expected.Value, expression.Value)
	}
}

func checkBooleanLiteral(t *testing.T, expression *parser.BooleanLiteral, expected *parser.BooleanLiteral) {
	if expression.Value != expected.Value {
		t.Fatalf("Expected boolean value %t, got %t", expected.Value, expression.Value)
	}
}

func checkStringLiteral(t *testing.T, expression *parser.StringLiteral, expected *parser.StringLiteral) {
	if expression.Value != expected.Value {
		t.Fatalf("Expected string value %q, got %q", expected.Value, expression.Value)
	}
}

func checkIdentifier(t *testing.T, expression *parser.Identifier, expected *parser.Identifier) {
	if expression.Value != expected.Value {
		t.Fatalf("Expected identifier value %q, got %q", expected.Value, expression.Value)
	}
}

func checkCallExpression(t *testing.T, expression *parser.CallExpression, expected *parser.CallExpression) {
	checkExpression(t, expression.Function, expected.Function)
	if len(expression.Arguments) != len(expected.Arguments) {
		t.Fatalf("Expected %d arguments, got %d", len(expected.Arguments), len(expression.Arguments))
	}
	for i, argument := range expression.Arguments {
		checkExpression(t, argument, expected.Arguments[i])
	}
}

func checkInfixExpression(t *testing.T, expression *parser.InfixExpression, expected *parser.InfixExpression) {
	if expression.Operator != expected.Operator {
		t.Fatalf("Expected operator %q, got %q", expected.Operator, expression.Operator)
	}
	checkExpression(t, expression.Left, expected.Left)
	checkExpression(t, expression.Right, expected.Right)
}

// func TestParserErrors(t *testing.T) {
// 	// Test case 1
// 	// Test case input
// 	input := `
// 		9ayed x 5;
// 	`

// 	// Expected error
// 	expectedError := "Expected next token to be =, got NUMBER instead"

// 	// Run the test
// 	runParserErrorTest(t, input, expectedError)
// }
