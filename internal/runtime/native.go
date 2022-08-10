package runtime

import (
	"fmt"
	"github.com/paw1a/golox/internal/lexing"
	"os"
	"strings"
	"time"
)

type ClockFunc struct {
}

func (f ClockFunc) Call(interpreter *Interpreter, arguments []interface{}) interface{} {
	return float64(time.Now().UnixMilli())
}

func (f ClockFunc) ParametersCount() int {
	return 0
}

type ExitFunc struct {
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

type AppendFunc struct {
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

type LenFunc struct {
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

type PrintFunc struct {
}

func (f PrintFunc) Call(interpreter *Interpreter, arguments []interface{}) interface{} {
	switch arguments[0].(type) {
	case string:
		arguments[0] = strings.ReplaceAll(arguments[0].(string), `\n`, "\n")
		fmt.Printf(arguments[0].(string), arguments[1:]...)
		return nil
	}

	runtimeError(lexing.Token{}, "printf expect format string at first argument")
	return nil
}

func (f PrintFunc) ParametersCount() int {
	return -1
}

type SleepFunc struct {
}

func (f SleepFunc) Call(interpreter *Interpreter, arguments []interface{}) interface{} {
	arg0 := arguments[0]
	if isNumber(arg0) {
		time.Sleep(time.Duration(arg0.(float64)) * time.Millisecond)
		return nil
	}

	runtimeError(lexing.Token{}, "sleep expect number argument")
	return nil
}

func (f SleepFunc) ParametersCount() int {
	return 1
}

type ClearFunc struct {
}

func (f ClearFunc) Call(interpreter *Interpreter, arguments []interface{}) interface{} {
	fmt.Print("\033[2J")
	return nil
}

func (f ClearFunc) ParametersCount() int {
	return 0
}
