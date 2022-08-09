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
	Variable    Expr
	Initializer Expr
}

type TernaryExpr struct {
	Condition Expr
	TrueExpr  Expr
	FalseExpr Expr
}

type LogicalExpr struct {
	LeftExpr  Expr
	Operator  lexing.Token
	RightExpr Expr
}

type CallExpr struct {
	Callee    Expr
	Paren     lexing.Token
	Arguments []Expr
}

type ArrayExpr struct {
	Elements []Expr
}

type IndexExpr struct {
	Array     Expr
	Bracket   lexing.Token
	IndexExpr Expr
}

type LambdaExpr struct {
	Params    []lexing.Token
	Statement BlockStmt
}
