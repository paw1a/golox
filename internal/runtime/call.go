package runtime

import (
	"github.com/paw1a/golox/internal/ast"
)

type Caller interface {
	Call(interpreter *Interpreter, arguments []interface{}) interface{}
	ParametersCount() int
}

type Function struct {
	Declaration ast.FunDeclarationStmt
	Closure     *Environment
}

func (f Function) Call(interpreter *Interpreter, arguments []interface{}) interface{} {
	innerScope := NewEnvironment(f.Closure)
	for i, argument := range arguments {
		innerScope.define(f.Declaration.Params[i].Lexeme, argument)
	}

	enclosingEnv := interpreter.env
	interpreter.env = innerScope
	defer func() {
		interpreter.env = enclosingEnv
	}()

	interpreter.executeBlockStmt(f.Declaration.Statement)
	if interpreter.returnContext.returnFlag {
		interpreter.returnContext.returnFlag = false
		return interpreter.returnContext.returnValue
	}

	return nil
}

func (f Function) ParametersCount() int {
	return len(f.Declaration.Params)
}

type LambdaFunction struct {
	LambdaExpr ast.LambdaExpr
	Closure    *Environment
}

func (f LambdaFunction) Call(interpreter *Interpreter, arguments []interface{}) interface{} {
	innerScope := NewEnvironment(f.Closure)
	for i, argument := range arguments {
		innerScope.define(f.LambdaExpr.Params[i].Lexeme, argument)
	}

	enclosingEnv := interpreter.env
	interpreter.env = innerScope
	defer func() {
		interpreter.env = enclosingEnv
	}()

	interpreter.executeBlockStmt(f.LambdaExpr.Statement)
	if interpreter.returnContext.returnFlag {
		interpreter.returnContext.returnFlag = false
		return interpreter.returnContext.returnValue
	}

	return nil
}

func (f LambdaFunction) ParametersCount() int {
	return len(f.LambdaExpr.Params)
}
