package interpreter

import (
	"fmt"
	"io"
	"math"
	"math/rand"
	"strconv"
	"strings"
	"time"
	"unicode/utf8"
)

type builtinValue string

// Bounds are argument counts; -1 means there is no upper bound.
var builtinNames = map[string][2]int{
	"golih": {1, 1}, "dakhel": {0, 1}, "n3ess": {1, 1},
	"rdo3adad": {1, 1}, "rdoBooleen": {1, 1}, "rdoString": {1, 1},
	"len": {1, 1}, "infini": {0, 0}, "rdoChar": {1, 1},
	"abs": {1, 1}, "sqrt": {1, 1}, "pow": {2, 2},
	"floor": {1, 1}, "ceil": {1, 1}, "round": {1, 1},
	"min": {1, -1}, "max": {1, -1}, "random": {0, 0},
}

func (eval *Evaluator) callBuiltin(name string, args []interface{}) (interface{}, error) {
	bounds := builtinNames[name]
	if len(args) < bounds[0] || (bounds[1] >= 0 && len(args) > bounds[1]) {
		return nil, fmt.Errorf("wrong number of arguments to %s: got %d", name, len(args))
	}
	switch name {
	case "golih":
		_, err := fmt.Fprintln(eval.output, displayValue(args[0]))
		return nil, err
	case "dakhel":
		if len(args) == 1 {
			if _, err := fmt.Fprint(eval.output, displayValue(args[0])); err != nil {
				return nil, err
			}
		}
		input, err := eval.input.ReadString('\n')
		if err != nil && (err != io.EOF || len(input) == 0) {
			return nil, err
		}
		return strings.TrimSuffix(strings.TrimSuffix(input, "\n"), "\r"), nil
	case "rdo3adad":
		return convertToNumber(args[0])
	case "rdoBooleen":
		return convertToBool(args[0])
	case "rdoString":
		return displayValue(args[0]), nil
	case "len":
		switch value := args[0].(type) {
		case *listValue:
			return float64(len(value.Items)), nil
		case map[string]interface{}:
			return float64(len(value)), nil
		case string:
			return float64(utf8.RuneCountInString(value)), nil
		default:
			return nil, fmt.Errorf("len requires a list, dictionary, or string")
		}
	case "infini":
		return math.Inf(1), nil
	case "random":
		return rand.Float64(), nil
	case "rdoChar":
		index, err := integer(args[0], 25)
		if err != nil {
			return nil, err
		}
		return string(rune('A' + index)), nil
	case "n3ess":
		milliseconds, err := integer(args[0], 86_400_000)
		if err != nil {
			return nil, err
		}
		time.Sleep(time.Duration(milliseconds) * time.Millisecond)
		return nil, nil
	}
	numbers := make([]float64, len(args))
	for i, arg := range args {
		number, err := toNumber(arg)
		if err != nil {
			return nil, err
		}
		if math.IsNaN(number) || math.IsInf(number, 0) {
			return nil, fmt.Errorf("%s requires finite numbers", name)
		}
		numbers[i] = number
	}
	var result float64
	switch name {
	case "abs":
		result = math.Abs(numbers[0])
	case "floor":
		result = math.Floor(numbers[0])
	case "ceil":
		result = math.Ceil(numbers[0])
	case "round":
		result = math.Round(numbers[0])
	case "sqrt":
		result = math.Sqrt(numbers[0])
	case "pow":
		result = math.Pow(numbers[0], numbers[1])
	case "min", "max":
		result = numbers[0]
		for _, number := range numbers[1:] {
			if name == "min" {
				result = math.Min(result, number)
			} else {
				result = math.Max(result, number)
			}
		}
	default:
		return nil, fmt.Errorf("function not found: %s", name)
	}
	if math.IsNaN(result) || math.IsInf(result, 0) {
		return nil, fmt.Errorf("%s result is outside the finite numeric domain", name)
	}
	return result, nil
}

func convertToNumber(value interface{}) (float64, error) {
	switch v := value.(type) {
	case float64:
		return v, nil
	case bool:
		if v {
			return 1, nil
		}
		return 0, nil
	case string:
		n, err := strconv.ParseFloat(strings.TrimSpace(v), 64)
		if err != nil || math.IsNaN(n) || math.IsInf(n, 0) {
			return 0, fmt.Errorf("cannot convert %q to a finite number", v)
		}
		return n, nil
	default:
		return 0, fmt.Errorf("cannot convert %s to a number", valueType(value))
	}
}

func convertToBool(value interface{}) (bool, error) {
	switch v := value.(type) {
	case bool:
		return v, nil
	case float64:
		return v != 0, nil
	case string:
		switch strings.ToLower(v) {
		case "s7i7", "true":
			return true, nil
		case "ghalt", "false":
			return false, nil
		}
		return false, fmt.Errorf("cannot convert %q to a boolean", v)
	default:
		return false, fmt.Errorf("cannot convert %s to a boolean", valueType(value))
	}
}

func integer(value interface{}, max int) (int, error) {
	n, err := toNumber(value)
	if err != nil {
		return 0, err
	}
	if math.IsNaN(n) || math.IsInf(n, 0) || math.Trunc(n) != n || n < 0 || n > float64(max) {
		return 0, fmt.Errorf("expected an integer in range 0..%d, got %v", max, n)
	}
	return int(n), nil
}
