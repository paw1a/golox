package runtime

import (
	"fmt"
	"github.com/paw1a/golox/internal/lexing"
)

type Interpreter struct {
	env           *Environment
	global        *Environment
	loopContext   loopContext
	returnContext returnContext
}

type loopContext struct {
	breakFlag    bool
	continueFlag bool
}

type returnContext struct {
	returnFlag  bool
	returnValue interface{}
}

func runtimeError(token lexing.Token, message string) {
	var errorMessage string
	if token.TokenType == lexing.Eof {
		errorMessage = fmt.Sprintf("line %d | at end of input: %s", token.Line, message)
	} else {
		errorMessage = fmt.Sprintf("line %d | at '%s': %s", token.Line, token.Lexeme, message)
	}
	panic(errorMessage)
}

func isTruthy(value interface{}) bool {
	if value == nil {
		return false
	}

	switch value.(type) {
	case float64:
		return value.(float64) != 0
	case bool:
		return value.(bool)
	case string:
		return value.(string) != ""
	}

	return true
}

func isNumber(value interface{}) bool {
	switch value.(type) {
	case float64:
		return true
	}
	return false
}

func isString(value interface{}) bool {
	switch value.(type) {
	case string:
		return true
	}
	return false
}

func NewInterpreter() *Interpreter {
	global := NewEnvironment(nil)
	global.define("clock", ClockFunc{})
	return &Interpreter{
		env:    global,
		global: global,
	}
}
