package parsing

import (
	"bytes"
	"fmt"
)

type Printer interface {
	Print() string
}

func (node BinaryExpr) Print() string {
	var buffer bytes.Buffer

	buffer.WriteString(" (")
	buffer.WriteString(node.Operator.Lexeme)
	buffer.WriteString(node.LeftExpr.Print())
	buffer.WriteString(node.RightExpr.Print())
	buffer.WriteString(") ")

	return buffer.String()
}

func (node UnaryExpr) Print() string {
	var buffer bytes.Buffer

	buffer.WriteString(" (")
	buffer.WriteString(node.Operator.Lexeme)
	buffer.WriteString(node.RightExpr.Print())
	buffer.WriteString(") ")

	return buffer.String()
}

func (node LiteralExpr) Print() string {
	return fmt.Sprintf(" %v ", node.LiteralValue)
}

func (node GroupingExpr) Print() string {
	var buffer bytes.Buffer

	buffer.WriteString(" (")
	buffer.WriteString("group ")
	buffer.WriteString(node.Expr.Print())
	buffer.WriteString(") ")

	return buffer.String()
}
