package runtime

import (
	"fmt"
	"github.com/paw1a/golox/internal/ast"
	"github.com/paw1a/golox/internal/lexing"
)

type Interpreter struct {
	env           *Environment
	global        *Environment
	locals        map[ast.Expr]int
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

func (i *Interpreter) Resolve(expr ast.Expr, index int) {
	i.locals[expr] = index
}

func (i *Interpreter) lookUpVariable(token lexing.Token, expr ast.Expr) interface{} {
	distance, ok := i.locals[expr]
	if ok {
		return i.env.getAt(distance, token.Lexeme)
	} else {
		return i.global.get(token)
	}
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
	global.define("exit", ExitFunc{})
	global.define("append", AppendFunc{})
	global.define("len", LenFunc{})
	return &Interpreter{
		env:    global,
		global: global,
		locals: make(map[ast.Expr]int),
	}
}
