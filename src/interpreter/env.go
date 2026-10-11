package interpreter

import (
	"fmt"
)

type Environment struct {
	store map[string]interface{}
	outer *Environment
}

func NewEnvironment() *Environment {
	return &Environment{
		outer: nil,
		store: make(map[string]interface{}),
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
