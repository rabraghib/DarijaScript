package interpreter

import (
	"fmt"
	"strings"
)

var methodArities = map[string][2]int{
	"append": {1, 1}, "remove": {1, 1}, "slice": {1, 2},
	"substring": {1, 2}, "contains": {1, 1}, "keys": {0, 0},
	"join": {1, 1}, "split": {1, 1}, "trim": {0, 0}, "indexOf": {1, 1},
}

func hasMethod(receiver interface{}, name string) bool {
	switch receiver.(type) {
	case *listValue:
		return name == "append" || name == "remove" || name == "slice" || name == "contains" || name == "join"
	case map[string]interface{}:
		return name == "keys" || name == "contains" || name == "remove"
	case string:
		return name == "split" || name == "trim" || name == "slice" || name == "substring" || name == "contains" || name == "indexOf"
	default:
		return false
	}
}

func sliceBounds(args []interface{}, length int) (int, int, error) {
	start, err := integer(args[0], length)
	if err != nil {
		return 0, 0, err
	}
	end := length
	if len(args) == 2 {
		end, err = integer(args[1], length)
		if err != nil {
			return 0, 0, err
		}
	}
	if end < start {
		return 0, 0, fmt.Errorf("slice end must be at least its start")
	}
	return start, end, nil
}

func stringArgument(value interface{}) (string, error) {
	text, ok := value.(string)
	if !ok {
		return "", fmt.Errorf("expected a string, got %s", valueType(value))
	}
	return text, nil
}

func callMethod(receiver interface{}, name string, args []interface{}) (interface{}, error) {
	bounds := methodArities[name]
	if len(args) < bounds[0] || len(args) > bounds[1] {
		return nil, fmt.Errorf("wrong number of arguments to .%s: got %d", name, len(args))
	}
	switch value := receiver.(type) {
	case *listValue:
		switch name {
		case "append":
			if len(value.Items) >= 1_000_000 {
				return nil, fmt.Errorf("maximum list length exceeded")
			}
			value.Items = append(value.Items, args[0])
			return nil, nil
		case "remove":
			index, err := integer(args[0], len(value.Items)-1)
			if err != nil {
				return nil, err
			}
			removed := value.Items[index]
			copy(value.Items[index:], value.Items[index+1:])
			value.Items[len(value.Items)-1] = nil
			value.Items = value.Items[:len(value.Items)-1]
			return removed, nil
		case "slice":
			start, end, err := sliceBounds(args, len(value.Items))
			if err != nil {
				return nil, err
			}
			return &listValue{Items: append([]interface{}{}, value.Items[start:end]...)}, nil
		case "contains":
			for _, item := range value.Items {
				if equalValues(item, args[0]) {
					return true, nil
				}
			}
			return false, nil
		case "join":
			separator, err := stringArgument(args[0])
			if err != nil {
				return nil, err
			}
			parts := make([]string, len(value.Items))
			for i, item := range value.Items {
				parts[i], err = stringArgument(item)
				if err != nil {
					return nil, err
				}
			}
			return strings.Join(parts, separator), nil
		}
	case map[string]interface{}:
		switch name {
		case "keys":
			keys, _ := iterationValues(value)
			return &listValue{Items: keys}, nil
		case "contains", "remove":
			key, err := stringArgument(args[0])
			if err != nil {
				return nil, err
			}
			item, exists := value[key]
			if name == "contains" {
				return exists, nil
			}
			delete(value, key)
			return item, nil
		}
	case string:
		switch name {
		case "trim":
			return strings.TrimSpace(value), nil
		case "slice", "substring":
			chars := []rune(value)
			start, end, err := sliceBounds(args, len(chars))
			if err != nil {
				return nil, err
			}
			return string(chars[start:end]), nil
		case "split", "contains", "indexOf":
			text, err := stringArgument(args[0])
			if err != nil {
				return nil, err
			}
			switch name {
			case "split":
				parts := strings.Split(value, text)
				items := make([]interface{}, len(parts))
				for i, part := range parts {
					items[i] = part
				}
				return &listValue{Items: items}, nil
			case "contains":
				return strings.Contains(value, text), nil
			case "indexOf":
				index := strings.Index(value, text)
				if index < 0 {
					return float64(-1), nil
				}
				return float64(len([]rune(value[:index]))), nil
			}
		}
	}
	return nil, fmt.Errorf("unknown method %s for %s", name, valueType(receiver))
}
