package interpreter

import (
	"fmt"
	"math"

	"github.com/rabraghib/darijascript/src/parser"
)

func (eval *Evaluator) evaluateExpression(expression parser.Expression) (value interface{}, err error) {
	eval.depth++
	defer func() { eval.depth--; err = atToken(err, expressionToken(expression)) }()
	if eval.depth > 256 {
		return nil, fmt.Errorf("maximum evaluation nesting exceeded")
	}
	if expression == nil {
		return nil, nil
	}
	switch expr := expression.(type) {
	case *parser.Identifier:
		val, ok := eval.env.Get(expr.Value)
		if !ok {
			if _, exists := builtinNames[expr.Value]; exists {
				return builtinValue(expr.Value), nil
			}
			return nil, fmt.Errorf("identifier not found: %s", expr.Value)
		}
		return val, nil
	case *parser.NumberLiteral:
		return expr.Value, nil
	case *parser.StringLiteral:
		return expr.Value, nil
	case *parser.BooleanLiteral:
		return expr.Value, nil
	case *parser.NullLiteral:
		return nil, nil
	case *parser.ArrayLiteral:
		if len(expr.Elements) > 1_000_000 {
			return nil, fmt.Errorf("maximum list length exceeded")
		}
		items, err := eval.evaluateArguments(expr.Elements)
		return &listValue{Items: items}, err
	case *parser.HashLiteral:
		dict := map[string]interface{}{}
		for _, pair := range expr.Pairs {
			key, err := eval.evaluateExpression(pair.Key)
			if err != nil {
				return nil, err
			}
			text, ok := key.(string)
			if !ok {
				return nil, fmt.Errorf("dictionary keys must be strings")
			}
			value, err := eval.evaluateExpression(pair.Value)
			if err != nil {
				return nil, err
			}
			dict[text] = value
		}
		return dict, nil
	case *parser.IndexExpression:
		left, err := eval.evaluateExpression(expr.Left)
		if err != nil {
			return nil, err
		}
		index, err := eval.evaluateExpression(expr.Index)
		if err != nil {
			return nil, err
		}
		return indexValue(left, index)
	case *parser.MemberExpression:
		left, err := eval.evaluateExpression(expr.Left)
		if err != nil {
			return nil, err
		}
		if !hasMethod(left, expr.Name) {
			return nil, fmt.Errorf("unknown method %s for %s", expr.Name, valueType(left))
		}
		return &methodValue{Receiver: left, Name: expr.Name}, nil
	case *parser.FunctionLiteral:
		return &functionValue{Name: "<anonymous>", Parameters: expr.Parameters, Body: expr.Body, Env: eval.env}, nil
	case *parser.PrefixExpression:
		rightValue, err := eval.evaluateExpression(expr.Right)
		if err != nil {
			return nil, err
		}
		return eval.evaluatePrefixExpression(expr.Operator, rightValue)
	case *parser.AssignmentExpression:
		// Resolve an indexed target once, before evaluating the right-hand side.
		var collection, index interface{}
		if target, ok := expr.Target.(*parser.IndexExpression); ok {
			collection, err = eval.evaluateExpression(target.Left)
			if err != nil {
				return nil, err
			}
			index, err = eval.evaluateExpression(target.Index)
			if err != nil {
				return nil, err
			}
		}
		value, err := eval.evaluateExpression(expr.Value)
		if err != nil {
			return nil, err
		}
		if target, ok := expr.Target.(*parser.Identifier); ok {
			err = eval.env.Update(target.Value, value)
		} else {
			err = assignIndex(collection, index, value)
		}
		if err != nil {
			return nil, err
		}
		return value, nil
	case *parser.InfixExpression:
		leftValue, err := eval.evaluateExpression(expr.Left)
		if err != nil {
			return nil, err
		}
		if expr.Operator == "&&" || expr.Operator == "||" {
			left, err := toBool(leftValue)
			if err != nil {
				return nil, err
			}
			if expr.Operator == "&&" && !left {
				return false, nil
			}
			if expr.Operator == "||" && left {
				return true, nil
			}
		}
		rightValue, err := eval.evaluateExpression(expr.Right)
		if err != nil {
			return nil, err
		}
		return eval.evaluateInfixExpression(expr.Operator, leftValue, rightValue)
	case *parser.CallExpression:
		return eval.evaluateCallExpression(expr)
	default:
		return nil, fmt.Errorf("unsupported expression type: %T", expression)
	}
}

func (eval *Evaluator) evaluatePrefixExpression(operator string, rightValue interface{}) (interface{}, error) {
	switch operator {
	case "!":
		value, err := toBool(rightValue)
		if err != nil {
			return nil, err
		}
		return !value, nil
	case "-":
		value, err := toNumber(rightValue)
		if err != nil {
			return nil, err
		}
		return -value, nil
	case "+":
		return toNumber(rightValue)
	default:
		return nil, fmt.Errorf("unsupported prefix operator: %s", operator)
	}
}

func (eval *Evaluator) evaluateInfixExpression(operator string, leftValue interface{}, rightValue interface{}) (interface{}, error) {
	switch operator {
	case "+":
		_, leftString := leftValue.(string)
		_, rightString := rightValue.(string)
		if leftString || rightString {
			return displayValue(leftValue) + displayValue(rightValue), nil
		}
		left, err := toNumber(leftValue)
		if err != nil {
			return nil, err
		}
		right, err := toNumber(rightValue)
		return left + right, err
	case "-":
		left, err := toNumber(leftValue)
		if err != nil {
			return nil, err
		}
		right, err := toNumber(rightValue)
		if err != nil {
			return nil, err
		}
		return left - right, nil
	case "*":
		left, err := toNumber(leftValue)
		if err != nil {
			return nil, err
		}
		right, err := toNumber(rightValue)
		if err != nil {
			return nil, err
		}
		return left * right, nil
	case "/":
		left, err := toNumber(leftValue)
		if err != nil {
			return nil, err
		}
		right, err := toNumber(rightValue)
		if err != nil {
			return nil, err
		}
		if right == 0 {
			return nil, fmt.Errorf("division by zero")
		}
		return left / right, nil
	case "%":
		left, err := toNumber(leftValue)
		if err != nil {
			return nil, err
		}
		right, err := toNumber(rightValue)
		if err != nil {
			return nil, err
		}
		if math.IsNaN(left) || math.IsInf(left, 0) || math.IsNaN(right) || math.IsInf(right, 0) || math.Trunc(left) != left || math.Trunc(right) != right {
			return nil, fmt.Errorf("modulo requires finite integer operands")
		}
		if right == 0 {
			return nil, fmt.Errorf("modulo by zero")
		}
		return math.Mod(left, right), nil
	case "<":
		left, err := toNumber(leftValue)
		if err != nil {
			return nil, err
		}
		right, err := toNumber(rightValue)
		if err != nil {
			return nil, err
		}
		return left < right, nil
	case "<=":
		left, err := toNumber(leftValue)
		if err != nil {
			return nil, err
		}
		right, err := toNumber(rightValue)
		if err != nil {
			return nil, err
		}
		return left <= right, nil
	case ">":
		left, err := toNumber(leftValue)
		if err != nil {
			return nil, err
		}
		right, err := toNumber(rightValue)
		if err != nil {
			return nil, err
		}
		return left > right, nil
	case ">=":
		left, err := toNumber(leftValue)
		if err != nil {
			return nil, err
		}
		right, err := toNumber(rightValue)
		if err != nil {
			return nil, err
		}
		return left >= right, nil
	case "==":
		return equalValues(leftValue, rightValue), nil
	case "!=":
		return !equalValues(leftValue, rightValue), nil
	case "&&":
		left, err := toBool(leftValue)
		if err != nil {
			return nil, err
		}
		right, err := toBool(rightValue)
		if err != nil {
			return nil, err
		}
		return left && right, nil
	case "||":
		left, err := toBool(leftValue)
		if err != nil {
			return nil, err
		}
		right, err := toBool(rightValue)
		if err != nil {
			return nil, err
		}
		return left || right, nil
	default:
		return nil, fmt.Errorf("unsupported infix operator: %s", operator)
	}
}

func (eval *Evaluator) evaluateCallExpression(call *parser.CallExpression) (value interface{}, err error) {
	callee, err := eval.evaluateExpression(call.Function)
	if err != nil {
		return nil, err
	}
	arguments, err := eval.evaluateArguments(call.Arguments)
	if err != nil {
		return nil, err
	}
	switch fn := callee.(type) {
	case builtinValue:
		return eval.callBuiltin(string(fn), arguments)
	case *methodValue:
		return callMethod(fn.Receiver, fn.Name, arguments)
	case *functionValue:
		defer func() { err = callError(err, fn.Name, call.Token) }()
		if len(fn.Parameters) != len(arguments) {
			return nil, fmt.Errorf("wrong number of arguments: expected %d, got %d", len(fn.Parameters), len(arguments))
		}
		outer := eval.env
		eval.env = NewEnclosedEnvironment(fn.Env)
		defer func() { eval.env = outer }()
		for i, parameter := range fn.Parameters {
			if err := eval.env.Set(parameter.Value, arguments[i]); err != nil {
				return nil, err
			}
		}
		value, err := eval.EvaluateStatements(fn.Body.Statements)
		return unwrapReturn(value), err
	default:
		return nil, fmt.Errorf("cannot call %s", valueType(callee))
	}
}

func (eval *Evaluator) evaluateArguments(arguments []parser.Expression) ([]interface{}, error) {
	values := make([]interface{}, len(arguments))
	for i, argument := range arguments {
		value, err := eval.evaluateExpression(argument)
		if err != nil {
			return nil, err
		}
		values[i] = value
	}
	return values, nil
}

func toBool(value interface{}) (bool, error) {
	switch v := value.(type) {
	case bool:
		return v, nil
	default:
		return false, fmt.Errorf("unsupported type for boolean conversion: %T", value)
	}
}

func toNumber(value interface{}) (float64, error) {
	switch v := value.(type) {
	case int64:
		return float64(v), nil
	case int32:
		return float64(v), nil
	case int16:
		return float64(v), nil
	case int8:
		return float64(v), nil
	case int:
		return float64(v), nil
	case float64:
		return v, nil
	default:
		return 0, fmt.Errorf("unsupported type for number conversion: %T", value)
	}
}
