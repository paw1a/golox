package runtime

import (
	"fmt"
	"github.com/paw1a/golox/internal/lexing"
	"os"
	"strings"
	"time"
)

const MaxParamsCount = 255

type ClockFunc struct {
}

type ExitFunc struct {
}

type AppendFunc struct {
}

type LenFunc struct {
}

type PrintFunc struct {
}

func (f ClockFunc) Call(interpreter *Interpreter, arguments []interface{}) interface{} {
	return float64(time.Now().UnixMilli())
}

func (f ClockFunc) ParametersCount() int {
	return 0
}

func (f ExitFunc) Call(interpreter *Interpreter, arguments []interface{}) interface{} {
	arg0 := arguments[0]

	if isNumber(arg0) {
		exitCode := int(arg0.(float64))
		os.Exit(exitCode)
	} else {
		runtimeError(lexing.Token{}, "exit code must be integer number")
	}

	return nil
}

func (f ExitFunc) ParametersCount() int {
	return 1
}

func (f AppendFunc) Call(interpreter *Interpreter, arguments []interface{}) interface{} {
	arg0 := arguments[0]
	switch arg0.(type) {
	case []interface{}:
		return append(arg0.([]interface{}), arguments[1])
	}

	runtimeError(lexing.Token{}, "append func expect array argument first")
	return nil
}

func (f AppendFunc) ParametersCount() int {
	return 2
}

func (f LenFunc) Call(interpreter *Interpreter, arguments []interface{}) interface{} {
	arg0 := arguments[0]
	switch arg0.(type) {
	case []interface{}:
		return float64(len(arg0.([]interface{})))
	}

	runtimeError(lexing.Token{}, "len func expect array argument")
	return nil
}

func (f LenFunc) ParametersCount() int {
	return 1
}

func (f PrintFunc) Call(interpreter *Interpreter, arguments []interface{}) interface{} {
	switch arguments[0].(type) {
	case string:
		arguments[0] = strings.Replace(arguments[0].(string), `\n`, "\n", -1)
		fmt.Printf(arguments[0].(string), arguments[1:]...)
		return nil
	}

	runtimeError(lexing.Token{}, "printf expect format string at first argument")
	return nil
}

func (f PrintFunc) ParametersCount() int {
	return -1
}
