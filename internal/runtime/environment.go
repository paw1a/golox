package runtime

import (
	"fmt"
	"github.com/paw1a/golox/internal/lexing"
)

type Environment struct {
	objects map[string]interface{}
}

func (e Environment) define(name string, value interface{}) {
	e.objects[name] = value
}

func (e Environment) get(name lexing.Token) interface{} {
	value, ok := e.objects[name.Lexeme]
	if ok {
		return value
	}

	runtimeError(name, fmt.Sprintf("undefined variable '%s'", name.Lexeme))
	return nil
}

func (e Environment) assign(name lexing.Token, value interface{}) {
	if _, ok := e.objects[name.Lexeme]; ok {
		e.objects[name.Lexeme] = value
		return
	}

	runtimeError(name, fmt.Sprintf("undefined variable '%s'", name.Lexeme))
}

func NewEnvironment() *Environment {
	return &Environment{objects: make(map[string]interface{})}
}
