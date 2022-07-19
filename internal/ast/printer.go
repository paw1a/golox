package ast

import (
	"bytes"
	"fmt"
)

type Printer interface {
	Print() string
}

func (expr BinaryExpr) Print() string {
	var buffer bytes.Buffer

	buffer.WriteString(" (")
	buffer.WriteString(expr.Operator.Lexeme)
	buffer.WriteString(expr.LeftExpr.Print())
	buffer.WriteString(expr.RightExpr.Print())
	buffer.WriteString(") ")

	return buffer.String()
}

func (expr UnaryExpr) Print() string {
	var buffer bytes.Buffer

	buffer.WriteString(" (")
	buffer.WriteString(expr.Operator.Lexeme)
	buffer.WriteString(expr.RightExpr.Print())
	buffer.WriteString(") ")

	return buffer.String()
}

func (expr LiteralExpr) Print() string {
	return fmt.Sprintf(" %v ", expr.LiteralValue)
}

func (expr GroupingExpr) Print() string {
	var buffer bytes.Buffer

	buffer.WriteString(" (")
	buffer.WriteString("group ")
	buffer.WriteString(expr.Expr.Print())
	buffer.WriteString(") ")

	return buffer.String()
}

func (expr VariableExpr) Print() string {
	return expr.Name.Lexeme
}

func (expr AssignExpr) Print() string {
	var buffer bytes.Buffer

	buffer.WriteString(" (")
	buffer.WriteString("= ")
	buffer.WriteString(expr.Name.Lexeme)
	buffer.WriteString(expr.Initializer.Print())
	buffer.WriteString(") ")

	return buffer.String()
}

func (stmt ExpressionStmt) Print() string {
	return ""
}

func (stmt PrintStmt) Print() string {
	return ""
}

func (stmt VarDeclarationStmt) Print() string {
	return ""
}

func (stmt BlockStmt) Print() string {
	return ""
}
