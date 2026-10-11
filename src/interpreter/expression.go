package interpreter

import (
	"fmt"
	"math"
	"reflect"
	"strconv"

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
			return nil, fmt.Errorf("identifier not found: %s", expr.Value)
		}
		return val, nil
	case *parser.NumberLiteral:
		return expr.Value, nil
	case *parser.StringLiteral:
		return expr.Value, nil
	case *parser.BooleanLiteral:
		return expr.Value, nil
	case *parser.PrefixExpression:
		rightValue, err := eval.evaluateExpression(expr.Right)
		if err != nil {
			return nil, err
		}
		return eval.evaluatePrefixExpression(expr.Operator, rightValue)
	case *parser.AssignmentExpression:
		value, err := eval.evaluateExpression(expr.Value)
		if err != nil {
			return nil, err
		}
		if err := eval.env.Update(expr.Variable.Value, value); err != nil {
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
		left, err1 := toNumber(leftValue)
		right, err2 := toNumber(rightValue)
		if err1 == nil && err2 == nil {
			return left + right, nil
		}
		leftStr, err1 := convertToString(leftValue)
		rightStr, err2 := convertToString(rightValue)
		if err1 != nil || err2 != nil {
			return nil, fmt.Errorf("unsupported infix operator: %s %s %s", leftValue, operator, rightValue)
		}
		return leftStr + rightStr, nil
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
		return reflect.DeepEqual(leftValue, rightValue), nil
	case "!=":
		return !reflect.DeepEqual(leftValue, rightValue), nil
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

func (eval *Evaluator) evaluateCallExpression(callExpression *parser.CallExpression) (interface{}, error) {
	function, ok, err := eval.env.GetFunction(callExpression.Function.Value)
	if err != nil {
		return nil, err
	}
	if !ok {
		return eval.evaluateBuiltinFunctionCall(callExpression)
	}
	if len(function.Parameters) != len(callExpression.Arguments) {
		return nil, fmt.Errorf("wrong number of arguments: expected %d, got %d", len(function.Parameters), len(callExpression.Arguments))
	}

	extraParams := map[string]interface{}{}
	for i, argument := range callExpression.Arguments {
		value, err := eval.evaluateExpression(argument)
		if err != nil {
			return nil, err
		}
		extraParams[function.Parameters[i].Value] = value
	}
	value, err := eval.evaluateBlockStatementExtended(function.Body, extraParams)
	return unwrapReturn(value), err
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
	case string:
		return strconv.ParseFloat(v, 64)
	default:
		return 0, fmt.Errorf("unsupported type for number conversion: %T", value)
	}
}
