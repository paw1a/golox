package runtime

import (
	"fmt"
	"github.com/paw1a/golox/internal/ast"
	"github.com/paw1a/golox/internal/lexing"
)

func (i Interpreter) Execute(stmt ast.Stmt) {
	switch stmt.(type) {
	case ast.ExpressionStmt:
		i.executeExprStmt(stmt.(ast.ExpressionStmt))
	case ast.PrintStmt:
		i.executePrintStmt(stmt.(ast.PrintStmt))
	case ast.VarDeclarationStmt:
		i.executeVarDeclarationStmt(stmt.(ast.VarDeclarationStmt))
	case ast.BlockStmt:
		i.executeBlockStmt(stmt.(ast.BlockStmt))
	case ast.IfStmt:
		i.executeIfStmt(stmt.(ast.IfStmt))
	default:
		runtimeError(lexing.Token{}, "invalid ast type")
	}
}

func (i Interpreter) executeExprStmt(stmt ast.ExpressionStmt) {
	i.Evaluate(stmt.Expr)
}

func (i Interpreter) executePrintStmt(stmt ast.PrintStmt) {
	value := i.Evaluate(stmt.Expr)
	fmt.Printf("%v\n", value)
}

func (i Interpreter) executeVarDeclarationStmt(stmt ast.VarDeclarationStmt) {
	var value interface{}
	if stmt.Initializer != nil {
		value = i.Evaluate(stmt.Initializer)
	}

	i.env.define(stmt.Name.Lexeme, value)
}

func (i Interpreter) executeBlockStmt(blockStmt ast.BlockStmt) {
	enclosingEnv := i.env

	i.env = NewEnvironment(enclosingEnv)
	defer func() {
		i.env = enclosingEnv
	}()

	for _, stmt := range blockStmt.Stmts {
		i.Execute(stmt)
	}
}

func (i Interpreter) executeIfStmt(stmt ast.IfStmt) {
	conditionValue := i.Evaluate(stmt.ConditionExpr)

	if isTruthy(conditionValue) {
		i.Execute(stmt.IfStatement)
		return
	}

	if stmt.ElseStatement != nil {
		i.Execute(stmt.ElseStatement)
	}
}
