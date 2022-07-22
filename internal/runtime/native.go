package runtime

import (
	"time"
)

const MaxParamsCount = 255

type ClockFunc struct {
}

func (f ClockFunc) Call(interpreter *Interpreter, arguments []interface{}) interface{} {
	return float64(time.Now().UnixMilli())
}

func (f ClockFunc) ParametersCount() int {
	return 0
}
