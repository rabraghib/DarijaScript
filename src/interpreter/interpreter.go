package interpreter

import (
	"bufio"
	"fmt"
	"io"
	"os"

	"github.com/rabraghib/darijascript/src/parser"
)

type Evaluator struct {
	env    *Environment
	input  *bufio.Reader
	output io.Writer
	depth  int
}

func NewEvaluator() *Evaluator {
	return NewEvaluatorWithIO(os.Stdin, os.Stdout)
}

func NewEvaluatorWithIO(input io.Reader, output io.Writer) *Evaluator {
	return &Evaluator{env: NewEnvironment(), input: bufio.NewReader(input), output: output}
}

func (eval *Evaluator) EvaluateProgram(program *parser.Program) (interface{}, error) {
	value, err := eval.EvaluateStatements(program.Statements)
	return unwrapReturn(value), err
}

func (eval *Evaluator) EvaluateStatements(statements []parser.Statement) (interface{}, error) {
	// Hoist named functions without mutating the syntax tree. Each captures this scope.
	for _, statement := range statements {
		if fn, ok := statement.(*parser.FunctionDeclaration); ok {
			value := &functionValue{Name: fn.Name.Value, Parameters: fn.Parameters, Body: fn.Body, Env: eval.env}
			if err := eval.env.Set(fn.Name.Value, value); err != nil {
				return nil, atToken(err, fn.Name.Token)
			}
		}
	}
	for _, statement := range statements {
		val, err := eval.evaluateStatement(statement)
		if err != nil {
			return nil, err
		}
		if val != nil {
			return val, nil
		}
	}
	return nil, nil
}

func (eval *Evaluator) evaluateStatement(statement parser.Statement) (interface{}, error) {
	if statement == nil {
		return nil, nil
	}
	switch stmt := statement.(type) {
	case *parser.LetStatement:
		return nil, eval.evaluateLetStatement(stmt)
	case *parser.ExpressionStatement:
		_, err := eval.evaluateExpression(stmt.Expression)
		return nil, err
	case *parser.FunctionDeclaration:
		return nil, nil
	case *parser.BreakStatement:
		return breakSignal{}, nil
	case *parser.ContinueStatement:
		return continueSignal{}, nil
	case *parser.TryStatement:
		value, err := eval.evaluateBlockStatement(stmt.Body)
		if err == nil {
			return value, nil
		}
		return eval.evaluateBlockStatementExtended(stmt.Handler, map[string]interface{}{stmt.Name.Value: err.Error()})
	case *parser.ForStatement:
		collection, err := eval.evaluateExpression(stmt.Collection)
		if err != nil {
			return nil, err
		}
		items, err := iterationValues(collection)
		if err != nil {
			return nil, atToken(err, stmt.Token)
		}
		for _, item := range items {
			value, err := eval.evaluateBlockStatementExtended(stmt.Body, map[string]interface{}{stmt.Name.Value: item})
			if err != nil {
				return nil, err
			}
			switch value.(type) {
			case breakSignal:
				return nil, nil
			case continueSignal:
				continue
			}
			if value != nil {
				return value, nil
			}
		}
	case *parser.IfStatement:
		condition, err := eval.evaluateBoolExpression(stmt.Condition)
		if err != nil {
			return nil, err
		}
		if *condition != stmt.IsConditionReversed {
			return eval.evaluateBlockStatement(stmt.Consequence)
		} else if stmt.Alternative != nil {
			return eval.evaluateBlockStatement(stmt.Alternative)
		}
	case *parser.WhileStatement:
		for {
			condition, err := eval.evaluateBoolExpression(stmt.Condition)
			if err != nil {
				return nil, err
			}
			if !*condition {
				break
			}
			val, err := eval.evaluateBlockStatement(stmt.Consequence)
			switch val.(type) {
			case breakSignal:
				return nil, err
			case continueSignal:
				if err != nil {
					return nil, err
				}
				continue
			}
			if err != nil || val != nil {
				return val, err
			}
		}
	case *parser.BlockStatement:
		return eval.evaluateBlockStatement(stmt)
	case *parser.ReturnStatement:
		value, err := eval.evaluateExpression(stmt.ReturnValue)
		if err != nil {
			return nil, err
		}
		return &returnSignal{value: value}, nil
	case *parser.ThrowStatement:
		val, err := eval.evaluateExpression(stmt.ReturnValue)
		if err != nil {
			return nil, err
		}
		return nil, atToken(fmt.Errorf("WAAAAAAA: %s", displayValue(val)), stmt.Token)
	default:
		return nil, fmt.Errorf("unsupported statement type: %T", statement)
	}
	return nil, nil
}

func (eval *Evaluator) evaluateBoolExpression(expression parser.Expression) (*bool, error) {
	rawExprValue, err := eval.evaluateExpression(expression)
	expressionValue, ok := rawExprValue.(bool)
	if err != nil {
		return nil, err
	}
	if !ok {
		return nil, atToken(fmt.Errorf("expression must evaluate to a boolean"), expressionToken(expression))
	}
	return &expressionValue, nil
}

func (eval *Evaluator) evaluateLetStatement(stmt *parser.LetStatement) error {
	val, err := eval.evaluateExpression(stmt.Value)
	if err != nil {
		return err
	}
	return atToken(eval.env.Set(stmt.Name.Value, val), stmt.Name.Token)
}

func (eval *Evaluator) evaluateBlockStatement(block *parser.BlockStatement) (interface{}, error) {
	return eval.evaluateBlockStatementExtended(block, nil)
}

func (eval *Evaluator) evaluateBlockStatementExtended(block *parser.BlockStatement, extraParams map[string]interface{}) (interface{}, error) {
	outer := eval.env
	eval.env = NewEnclosedEnvironment(outer)
	defer func() { eval.env = outer }()
	for k, v := range extraParams {
		err := eval.env.Set(k, v)
		if err != nil {
			return nil, err
		}
	}
	return eval.EvaluateStatements(block.Statements)
}

type returnSignal struct{ value interface{} }
type breakSignal struct{}
type continueSignal struct{}

func unwrapReturn(value interface{}) interface{} {
	if signal, ok := value.(*returnSignal); ok {
		return signal.value
	}
	return value
}
