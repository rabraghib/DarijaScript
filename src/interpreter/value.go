package interpreter

import (
	"fmt"
	"reflect"
	"sort"
	"strconv"
	"strings"

	"github.com/rabraghib/darijascript/src/parser"
)

// A list owns its slice header so append/remove remain visible through aliases.
type listValue struct{ Items []interface{} }
type functionValue struct {
	Name       string
	Parameters []*parser.Identifier
	Body       *parser.BlockStatement
	Env        *Environment
}
type methodValue struct {
	Receiver interface{}
	Name     string
}

func valueType(value interface{}) string {
	switch value.(type) {
	case nil:
		return "walo"
	case float64:
		return "number"
	case bool:
		return "boolean"
	case string:
		return "string"
	case *listValue:
		return "list"
	case map[string]interface{}:
		return "dictionary"
	case *functionValue, builtinValue, *methodValue:
		return "function"
	default:
		return fmt.Sprintf("%T", value)
	}
}

func indexValue(collection, index interface{}) (interface{}, error) {
	switch value := collection.(type) {
	case *listValue:
		i, err := integer(index, len(value.Items)-1)
		if err != nil {
			return nil, err
		}
		return value.Items[i], nil
	case string:
		chars := []rune(value)
		i, err := integer(index, len(chars)-1)
		if err != nil {
			return nil, err
		}
		return string(chars[i]), nil
	case map[string]interface{}:
		key, ok := index.(string)
		if !ok {
			return nil, fmt.Errorf("dictionary keys must be strings")
		}
		return value[key], nil
	default:
		return nil, fmt.Errorf("cannot index %s", valueType(collection))
	}
}

func assignIndex(collection, index, item interface{}) error {
	switch value := collection.(type) {
	case *listValue:
		i, err := integer(index, len(value.Items)-1)
		if err != nil {
			return err
		}
		value.Items[i] = item
		return nil
	case map[string]interface{}:
		key, ok := index.(string)
		if !ok {
			return fmt.Errorf("dictionary keys must be strings")
		}
		value[key] = item
		return nil
	default:
		return fmt.Errorf("indexed assignment requires a list or dictionary")
	}
}

func sortedKeys(dict map[string]interface{}) []string {
	keys := make([]string, 0, len(dict))
	for key := range dict {
		keys = append(keys, key)
	}
	sort.Strings(keys)
	return keys
}

// Loops visit a snapshot: mutation cannot extend the active traversal.
func iterationValues(collection interface{}) ([]interface{}, error) {
	switch value := collection.(type) {
	case *listValue:
		return append([]interface{}{}, value.Items...), nil
	case map[string]interface{}:
		items := []interface{}{}
		for _, key := range sortedKeys(value) {
			items = append(items, key)
		}
		return items, nil
	case string:
		items := []interface{}{}
		for _, char := range value {
			items = append(items, string(char))
		}
		return items, nil
	default:
		return nil, fmt.Errorf("kol requires a list, dictionary, or string")
	}
}

func displayValue(value interface{}) string { return formatValue(value, map[uintptr]bool{}) }
func formatValue(value interface{}, visiting map[uintptr]bool) string {
	switch v := value.(type) {
	case nil:
		return "walo"
	case bool:
		if v {
			return "s7i7"
		}
		return "ghalt"
	case float64:
		return strconv.FormatFloat(v, 'f', -1, 64)
	case string:
		return v
	case *functionValue:
		return "<fonksyon " + v.Name + ">"
	case builtinValue:
		return "<fonksyon " + string(v) + ">"
	case *methodValue:
		return "<method " + v.Name + ">"
	case *listValue:
		key := reflect.ValueOf(v).Pointer()
		if visiting[key] {
			return "<cycle>"
		}
		visiting[key] = true
		defer delete(visiting, key)
		items := make([]string, len(v.Items))
		for i, item := range v.Items {
			items[i] = formatNested(item, visiting)
		}
		return "[" + strings.Join(items, ", ") + "]"
	case map[string]interface{}:
		key := uintptr(reflect.ValueOf(v).UnsafePointer())
		if visiting[key] {
			return "<cycle>"
		}
		visiting[key] = true
		defer delete(visiting, key)
		items := []string{}
		for _, key := range sortedKeys(v) {
			items = append(items, quoteString(key)+": "+formatNested(v[key], visiting))
		}
		return "{" + strings.Join(items, ", ") + "}"
	default:
		return fmt.Sprint(value)
	}
}

// Keep escaping within the language's supported string escapes.
func quoteString(text string) string {
	replacer := strings.NewReplacer("\\", "\\\\", "\"", "\\\"", "\n", "\\n", "\r", "\\r", "\t", "\\t")
	return "\"" + replacer.Replace(text) + "\""
}

func formatNested(value interface{}, visiting map[uintptr]bool) string {
	if text, ok := value.(string); ok {
		return quoteString(text)
	}
	return formatValue(value, visiting)
}

type equalityPair struct {
	Left, Right uintptr
	Kind        string
}

func equalValues(left, right interface{}) bool {
	return compareValues(left, right, map[equalityPair]bool{})
}
func compareValues(left, right interface{}, seen map[equalityPair]bool) bool {
	if valueType(left) != valueType(right) {
		return false
	}
	switch a := left.(type) {
	case nil:
		return right == nil
	case float64:
		b, ok := right.(float64)
		return ok && a == b
	case bool:
		b, ok := right.(bool)
		return ok && a == b
	case string:
		b, ok := right.(string)
		return ok && a == b
	case builtinValue:
		b, ok := right.(builtinValue)
		return ok && a == b
	case *functionValue:
		b, ok := right.(*functionValue)
		return ok && a == b
	case *methodValue:
		b, ok := right.(*methodValue)
		return ok && a == b
	case *listValue:
		b, ok := right.(*listValue)
		if !ok || len(a.Items) != len(b.Items) {
			return false
		}
		pair := equalityPair{reflect.ValueOf(a).Pointer(), reflect.ValueOf(b).Pointer(), "list"}
		if seen[pair] {
			return true
		}
		seen[pair] = true
		for i := range a.Items {
			if !compareValues(a.Items[i], b.Items[i], seen) {
				return false
			}
		}
		return true
	case map[string]interface{}:
		b, ok := right.(map[string]interface{})
		if !ok || len(a) != len(b) {
			return false
		}
		pair := equalityPair{uintptr(reflect.ValueOf(a).UnsafePointer()), uintptr(reflect.ValueOf(b).UnsafePointer()), "dictionary"}
		if seen[pair] {
			return true
		}
		seen[pair] = true
		for key, value := range a {
			other, ok := b[key]
			if !ok || !compareValues(value, other, seen) {
				return false
			}
		}
		return true
	default:
		return false
	}
}
