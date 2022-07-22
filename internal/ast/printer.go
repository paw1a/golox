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

func (expr TernaryExpr) Print() string {
	var buffer bytes.Buffer

	buffer.WriteString(" (")
	buffer.WriteString("? ")
	buffer.WriteString(expr.Condition.Print())
	buffer.WriteString(expr.TrueExpr.Print())
	buffer.WriteString(" : ")
	buffer.WriteString(expr.FalseExpr.Print())
	buffer.WriteString(") ")

	return buffer.String()
}

func (expr LogicalExpr) Print() string {
	var buffer bytes.Buffer

	buffer.WriteString(" (")
	buffer.WriteString(expr.Operator.Lexeme)
	buffer.WriteString(expr.LeftExpr.Print())
	buffer.WriteString(expr.RightExpr.Print())
	buffer.WriteString(") ")

	return buffer.String()
}

func (expr CallExpr) Print() string {
	return ""
}

func (stmt ExpressionStmt) Print() string {
	var buffer bytes.Buffer

	buffer.WriteString(" (")
	buffer.WriteString("expr ")
	buffer.WriteString(stmt.Expr.Print())
	buffer.WriteString(") ")

	return buffer.String()
}

func (stmt PrintStmt) Print() string {
	var buffer bytes.Buffer

	buffer.WriteString(" (")
	buffer.WriteString("print ")
	buffer.WriteString(stmt.Expr.Print())
	buffer.WriteString(") ")

	return buffer.String()
}

func (stmt VarDeclarationStmt) Print() string {
	var buffer bytes.Buffer

	buffer.WriteString(" (")
	buffer.WriteString("var ")
	buffer.WriteString(stmt.Name.Lexeme)
	if stmt.Initializer != nil {
		buffer.WriteString(" = ")
		buffer.WriteString(stmt.Initializer.Print())
	}
	buffer.WriteString(") ")

	return buffer.String()
}

func (stmt BlockStmt) Print() string {
	var buffer bytes.Buffer

	buffer.WriteString(" {")
	for _, st := range stmt.Stmts {
		buffer.WriteString(" ")
		buffer.WriteString(st.Print())
		buffer.WriteString(";")
	}
	buffer.WriteString("} ")

	return buffer.String()
}

func (stmt IfStmt) Print() string {
	var buffer bytes.Buffer

	buffer.WriteString(" (")
	buffer.WriteString("if ")
	buffer.WriteString(stmt.ConditionExpr.Print())
	buffer.WriteString(" then ")
	buffer.WriteString(stmt.IfStatement.Print())
	if stmt.ElseStatement != nil {
		buffer.WriteString(" else ")
		buffer.WriteString(stmt.ElseStatement.Print())
	}
	buffer.WriteString(") ")

	return buffer.String()
}

func (stmt ForStmt) Print() string {
	return ""
}

func (stmt BreakStmt) Print() string {
	return ""
}

func (stmt ContinueStmt) Print() string {
	return ""
}

func (stmt FunDeclarationStmt) Print() string {
	return ""
}
