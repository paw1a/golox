package runtime

import (
	"fmt"
	"github.com/paw1a/golox/internal/ast"
)

func (i Interpreter) Execute(stmt ast.Stmt) {
	switch stmt.(type) {
	case ast.ExpressionStmt:
		i.executeExprStmt(stmt.(ast.ExpressionStmt))
	case ast.PrintStmt:
		i.executePrintStmt(stmt.(ast.PrintStmt))
	case ast.VarDeclarationStmt:
		i.executeVarDeclarationStmt(stmt.(ast.VarDeclarationStmt))
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
