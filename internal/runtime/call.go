package runtime

import "github.com/paw1a/golox/internal/ast"

type Caller interface {
	Call(interpreter *Interpreter, arguments []interface{}) interface{}
	ParametersCount() int
}

type Function struct {
	Declaration ast.FunDeclarationStmt
}

func (f Function) Call(interpreter *Interpreter, arguments []interface{}) interface{} {
	innerScope := NewEnvironment(interpreter.global)
	for i, argument := range arguments {
		innerScope.define(f.Declaration.Params[i].Lexeme, argument)
	}

	enclosingEnv := interpreter.env
	interpreter.env = innerScope
	defer func() {
		interpreter.env = enclosingEnv
	}()

	interpreter.executeBlockStmt(f.Declaration.Statement)
	return nil
}

func (f Function) ParametersCount() int {
	return len(f.Declaration.Params)
}
