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

type BlockStmt struct {
	Stmts []Stmt
}

type VarDeclarationStmt struct {
	Name        lexing.Token
	Initializer Expr
}

type IfStmt struct {
	ConditionExpr Expr
	IfStatement   Stmt
	ElseStatement Stmt
}

type ForStmt struct {
	InitializerStmt Stmt
	ConditionExpr   Expr
	IncrementExpr   Expr
	Statement       Stmt
}

type BreakStmt struct {
}

type ContinueStmt struct {
}

type FunDeclarationStmt struct {
	Name      lexing.Token
	Params    []lexing.Token
	Statement BlockStmt
}

type ReturnStmt struct {
	ReturnToken lexing.Token
	Expr        Expr
}
