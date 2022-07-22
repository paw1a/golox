package runtime

type Caller interface {
	Call(interpreter *Interpreter, arguments []interface{}) interface{}
	ParametersCount() int
}
