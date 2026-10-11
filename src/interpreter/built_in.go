package interpreter

import (
	"fmt"
	"io"
	"math"
	"reflect"
	"runtime"
	"strconv"
	"strings"
	"time"

	"github.com/rabraghib/darijascript/src/parser"
)

func (eval *Evaluator) evaluateBuiltinFunctionCall(callExpression *parser.CallExpression) (interface{}, error) {
	switch callExpression.Function.Value {
	case "n3ess":
		if len(callExpression.Arguments) != 1 {
			return nil, fmt.Errorf("n3ess() takes exactly 1 argument, %d given", len(callExpression.Arguments))
		}
		arg, err := eval.evaluateExpression(callExpression.Arguments[0])
		if err != nil {
			return nil, err
		}
		num, err := convertToInt64(arg)
		if err != nil {
			return nil, err
		}
		time.Sleep(time.Duration(num) * time.Millisecond)
		return nil, nil
	case "abs":
		if len(callExpression.Arguments) != 1 {
			return nil, fmt.Errorf("abs() takes exactly 1 argument, %d given", len(callExpression.Arguments))
		}
		arg, err := eval.evaluateExpression(callExpression.Arguments[0])
		if err != nil {
			return nil, err
		}
		num, err := convertToInt64(arg)
		if err != nil {
			return nil, err
		}
		if num < 0 {
			return -num, nil
		}
		return num, nil
	case "golih":
		if len(callExpression.Arguments) != 1 {
			return nil, fmt.Errorf("golih() takes exactly 1 argument, %d given", len(callExpression.Arguments))
		}
		arg, err := eval.evaluateExpression(callExpression.Arguments[0])
		if err != nil {
			return nil, err
		}
		_, err = fmt.Fprintln(eval.output, displayValue(arg))
		return nil, err
	case "dakhel":
		if runtime.GOOS == "js" {
			return nil, fmt.Errorf("dakhel() is not supported in the browser")
		}
		if len(callExpression.Arguments) > 1 {
			return nil, fmt.Errorf("dakhel() takes at most 1 argument, %d given", len(callExpression.Arguments))
		}
		if len(callExpression.Arguments) == 1 {
			message, err := eval.evaluateExpression(callExpression.Arguments[0])
			if err != nil {
				return nil, err
			}
			if _, err := fmt.Fprint(eval.output, displayValue(message)); err != nil {
				return nil, err
			}
		}
		input, err := eval.input.ReadString('\n')
		if err != nil && (err != io.EOF || len(input) == 0) {
			return nil, err
		}
		return strings.TrimSuffix(strings.TrimSuffix(input, "\n"), "\r"), nil
	case "rdo3adad", "rdoBooleen", "rdoString":
		if len(callExpression.Arguments) != 1 {
			return nil, fmt.Errorf(
				"%s() takes exactly 1 argument, %d given",
				callExpression.Function.Value,
				len(callExpression.Arguments),
			)
		}
		arg, err := eval.evaluateExpression(callExpression.Arguments[0])
		if err != nil {
			return nil, err
		}
		if callExpression.Function.Value == "rdo3adad" {
			return convertToInt64(arg)
		}
		if callExpression.Function.Value == "rdoBooleen" {
			return convertToBool(arg)
		}
		return convertToString(arg)
	case "c", "r", "l":
		list := make([]interface{}, len(callExpression.Arguments))
		for i, argument := range callExpression.Arguments {
			value, err := eval.evaluateExpression(argument)
			if err != nil {
				return nil, err
			}
			list[i] = value
		}
		return list, nil
	case "infini":
		if len(callExpression.Arguments) != 0 {
			return nil, fmt.Errorf("infini() takes exactly 0 arguments, %d given", len(callExpression.Arguments))
		}
		return math.Inf(1), nil
	case "len":
		if len(callExpression.Arguments) != 1 {
			return nil, fmt.Errorf("len() takes exactly 1 argument, %d given", len(callExpression.Arguments))
		}
		arg, err := eval.evaluateExpression(callExpression.Arguments[0])
		if err != nil {
			return nil, err
		}
		switch v := arg.(type) {
		case []interface{}:
			return float64(len(v)), nil
		case string:
			return float64(len(v)), nil
		case map[string]interface{}:
			return float64(len(v)), nil
		default:
			return nil, fmt.Errorf("unsupported type for len(): %T", arg)
		}
	case "dirRow", "dirCol":
		if len(callExpression.Arguments) != 2 {
			return nil, fmt.Errorf("%s() takes exactly 2 arguments, %d given", callExpression.Function.Value, len(callExpression.Arguments))
		}
		n, err := eval.evaluateExpression(callExpression.Arguments[0])
		if err != nil {
			return nil, err
		}
		value, err := eval.evaluateExpression(callExpression.Arguments[1])
		if err != nil {
			return nil, err
		}
		num, err := integer(n, 1_000_000)
		if err != nil {
			return nil, err
		}
		row := make([]interface{}, num)
		for i := range row {
			row[i] = value
		}
		return row, nil
	case "toLetter":
		if len(callExpression.Arguments) != 1 {
			return nil, fmt.Errorf("toLetter() takes exactly 1 argument, %d given", len(callExpression.Arguments))
		}
		n, err := eval.evaluateExpression(callExpression.Arguments[0])
		if err != nil {
			return nil, err
		}
		num, err := integer(n, 25)
		if err != nil {
			return nil, err
		}
		return string(rune('A' + num)), nil
	case "ara":
		if len(callExpression.Arguments) < 2 {
			return nil, fmt.Errorf("ara() takes at least 2 arguments, %d given", len(callExpression.Arguments))
		}
		arg, err := eval.evaluateExpression(callExpression.Arguments[0])
		if err != nil {
			return nil, err
		}
		current := arg
		for i := 1; i < len(callExpression.Arguments); i++ {
			indexArg, err := eval.evaluateExpression(callExpression.Arguments[i])
			if err != nil {
				return nil, err
			}
			switch v := current.(type) {
			case []interface{}:
				index, err := integer(indexArg, len(v)-1)
				if err != nil {
					return nil, err
				}
				current = v[index]

			default:
				return nil, fmt.Errorf("unexpected type at level %d; expected []interface{}", i)
			}
		}

		// `current` now contains the final accessed element
		return current, nil
	case "atih":
		// atih(array, index, value)
		// atih(matrix, row, col, value)
		if len(callExpression.Arguments) < 3 {
			return nil, fmt.Errorf("atih() takes at least 3 arguments, %d given", len(callExpression.Arguments))
		}
		arg, err := eval.evaluateExpression(callExpression.Arguments[0])
		if err != nil {
			return nil, err
		}
		current := arg
		for i := 1; i < len(callExpression.Arguments)-1; i++ {
			values, ok := current.([]interface{})
			if !ok {
				return nil, fmt.Errorf("expected list at index level %d", i)
			}
			indexArg, err := eval.evaluateExpression(callExpression.Arguments[i])
			if err != nil {
				return nil, err
			}
			index, err := integer(indexArg, len(values)-1)
			if err != nil {
				return nil, err
			}
			if i == len(callExpression.Arguments)-2 {
				value, err := eval.evaluateExpression(callExpression.Arguments[i+1])
				if err != nil {
					return nil, err
				}
				values[index] = value
				return nil, nil
			}
			current = values[index]
		}
		return nil, nil
	default:
		return nil, fmt.Errorf("function not found: %s", callExpression.Function.Value)
	}
}

func convertToInt64(value interface{}) (float64, error) {
	switch v := value.(type) {
	case float64:
		return v, nil
	case int:
		return float64(v), nil
	case bool:
		if v {
			return 1, nil
		}
		return 0, nil
	case string:
		return strconv.ParseFloat(v, 64)
	default:
		return 0, fmt.Errorf("unsupported type for integer conversion: %T", value)
	}
}

func convertToBool(value interface{}) (bool, error) {
	switch v := value.(type) {
	case bool:
		return v, nil
	case float64:
		return v != 0, nil
	case string:
		if v == "s7i7" || v == "S7I7" {
			return true, nil
		}
		if v == "ghalt" || v == "GHALT" {
			return false, nil
		}
		return strconv.ParseBool(v)
	default:
		return false, fmt.Errorf("unsupported type for boolean conversion: %T", value)
	}
}

func convertToString(value interface{}) (string, error) {
	switch v := value.(type) {
	case bool:
		if v {
			return "s7i7", nil
		}
		return "ghalt", nil
	case float64:
		return strconv.FormatFloat(v, 'f', -1, 64), nil
	case string:
		return v, nil
	default:
		return displayValue(value), nil
		// return "", fmt.Errorf("unsupported type for string conversion: %T", value)
	}
}

func integer(value interface{}, max int) (int, error) {
	n, err := convertToInt64(value)
	if err != nil {
		return 0, err
	}
	if math.IsNaN(n) || math.IsInf(n, 0) || math.Trunc(n) != n || n < 0 || n > float64(max) {
		return 0, fmt.Errorf("expected an integer in range 0..%d, got %v", max, n)
	}
	return int(n), nil
}

// atih() can create self-referencing lists; format them without infinite recursion.
func displayValue(value interface{}) string { return formatValue(value, map[uintptr]bool{}) }

func formatValue(value interface{}, visiting map[uintptr]bool) string {
	values, ok := value.([]interface{})
	if !ok {
		return fmt.Sprint(value)
	}
	key := reflect.ValueOf(values).Pointer()
	if visiting[key] {
		return "<cycle>"
	}
	visiting[key] = true
	defer delete(visiting, key)
	var text strings.Builder
	text.WriteByte('[')
	for i, item := range values {
		if i > 0 {
			text.WriteByte(' ')
		}
		text.WriteString(formatValue(item, visiting))
	}
	text.WriteByte(']')
	return text.String()
}
