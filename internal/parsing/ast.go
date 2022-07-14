package parsing

import (
	"github.com/paw1a/golox/internal/lexing"
)

type Expr interface {
	Evaluate() interface{}
	Printer
}

type BinaryExpr struct {
	LeftExpr  Expr
	Operator  lexing.Token
	RightExpr Expr
}

func (node BinaryExpr) Evaluate() interface{} {
	return nil
}

type UnaryExpr struct {
	Operator  lexing.Token
	RightExpr Expr
}

func (node UnaryExpr) Evaluate() interface{} {
	return nil
}

type LiteralExpr struct {
	LiteralValue interface{}
}

func (node LiteralExpr) Evaluate() interface{} {
	return nil
}

type GroupingExpr struct {
	Expr Expr
}

func (node GroupingExpr) Evaluate() interface{} {
	return nil
}
