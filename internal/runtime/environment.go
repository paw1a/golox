package runtime

import (
	"fmt"
	"github.com/paw1a/golox/internal/lexing"
)

type Environment struct {
	Objects map[string]interface{}
}

func (e Environment) define(name string, value interface{}) {
	e.Objects[name] = value
}

func (e Environment) get(name lexing.Token) interface{} {
	value, ok := e.Objects[name.Lexeme]
	if ok {
		return value
	}

	runtimeError(name, fmt.Sprintf("undefined variable '%s'", name.Lexeme))
	return nil
}

func NewEnvironment() Environment {
	return Environment{Objects: make(map[string]interface{})}
}
