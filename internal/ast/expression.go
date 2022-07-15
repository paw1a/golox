package ast

import "github.com/paw1a/golox/internal/lexing"

type Expr interface {
	Printer
}

type BinaryExpr struct {
	LeftExpr  Expr
	Operator  lexing.Token
	RightExpr Expr
}

type UnaryExpr struct {
	Operator  lexing.Token
	RightExpr Expr
}

type LiteralExpr struct {
	LiteralValue interface{}
}

type GroupingExpr struct {
	Expr Expr
}

type VariableExpr struct {
	Name lexing.Token
}

type AssignExpr struct {
	Name        lexing.Token
	Initializer Expr
}
