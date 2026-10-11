package interpreter

import (
	"fmt"

	"github.com/rabraghib/darijascript/src/parser"
)

type Environment struct {
	store     map[string]interface{}
	outer     *Environment
	functions map[string]*parser.FunctionDeclaration
}

func NewEnvironment() *Environment {
	return &Environment{
		outer:     nil,
		store:     make(map[string]interface{}),
		functions: make(map[string]*parser.FunctionDeclaration),
	}
}

func NewEnclosedEnvironment(outer *Environment) *Environment {
	env := NewEnvironment()
	env.outer = outer
	return env
}

func (env *Environment) Get(name string) (interface{}, bool) {
	val, ok := env.store[name]
	if !ok {
		fn, ok := env.functions[name]
		if ok {
			return fn, ok
		}
		if env.outer != nil {
			return env.outer.Get(name)
		}
	}
	return val, ok
}

func (env *Environment) Set(name string, value interface{}) error {
	_, ok := env.store[name]
	if ok {
		return fmt.Errorf("id %s already declared", name)
	}
	env.store[name] = value
	return nil
}

func (env *Environment) Update(name string, value interface{}) error {
	if _, ok := env.store[name]; ok {
		env.store[name] = value
		return nil
	} else if env.outer != nil {
		return env.outer.Update(name, value)
	}
	return fmt.Errorf("identifier not found: %s", name)
}

func (env *Environment) GetFunction(name string) (*parser.FunctionDeclaration, bool, error) {
	fnRaw, ok := env.Get(name)
	if ok {
		fn, ok := fnRaw.(*parser.FunctionDeclaration)
		if !ok {
			return nil, false, fmt.Errorf("identifier %s is not a function", name)
		}
		return fn, ok, nil
	}
	return nil, ok, nil
}

func (env *Environment) SetFunction(name string, fn *parser.FunctionDeclaration) {
	env.functions[name] = fn
}
