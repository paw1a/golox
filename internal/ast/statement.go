package ast

import (
	"github.com/paw1a/golox/internal/lexing"
)

type Stmt interface {
	Printer
}

type ExpressionStmt struct {
	Expr Expr
}

type PrintStmt struct {
	Expr Expr
}

type VarDeclarationStmt struct {
	Name        lexing.Token
	Initializer Expr
}
