package runtime

import (
	"fmt"
	"github.com/paw1a/golox/internal/ast"
)

func Execute(stmt ast.Stmt) {
	switch stmt.(type) {
	case ast.ExpressionStmt:
		executeExprStmt(stmt.(ast.ExpressionStmt))
	case ast.PrintStmt:
		executePrintStmt(stmt.(ast.PrintStmt))
	case ast.VarDeclarationStmt:
		executeVarDeclarationStmt(stmt.(ast.VarDeclarationStmt))
	}
}

func executeExprStmt(stmt ast.ExpressionStmt) {
	Evaluate(stmt.Expr)
}

func executePrintStmt(stmt ast.PrintStmt) {
	value := Evaluate(stmt.Expr)
	fmt.Printf("%v\n", value)
}

func executeVarDeclarationStmt(stmt ast.VarDeclarationStmt) {

}
