package parsing

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
