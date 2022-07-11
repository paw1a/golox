package parsing

import "github.com/paw1a/golox/internal/lexing"

type AstNode interface{}

type BinaryAstNode struct {
	LeftExpr  AstNode
	Operator  lexing.Token
	RightExpr AstNode
}

type UnaryAstNode struct {
	Operator  lexing.Token
	RightExpr AstNode
}

type LiteralAstNode struct {
	LiteralValue interface{}
}

type GroupingAstNode struct {
	Expr AstNode
}
