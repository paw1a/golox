package parsing

import (
	"fmt"
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

func (expr BinaryExpr) Evaluate() interface{} {
	leftValue := expr.LeftExpr.Evaluate()
	rightValue := expr.RightExpr.Evaluate()

	switch expr.Operator.TokenType {
	case lexing.Plus:
		switch {
		case isNumber(leftValue) && isNumber(rightValue):
			return leftValue.(float64) + rightValue.(float64)
		case isString(leftValue) && isString(rightValue):
			return leftValue.(string) + rightValue.(string)
		default:
			runtimeError(expr.Operator, "number or string operands expected")
		}
	case lexing.Minus:
		requireNumberOperand(expr.Operator, leftValue)
		requireNumberOperand(expr.Operator, rightValue)
		return leftValue.(float64) - rightValue.(float64)
	case lexing.Star:
		requireNumberOperand(expr.Operator, leftValue)
		requireNumberOperand(expr.Operator, rightValue)
		return leftValue.(float64) * rightValue.(float64)
	case lexing.Slash:
		requireNumberOperand(expr.Operator, leftValue)
		requireNumberOperand(expr.Operator, rightValue)

		if rightValue.(float64) == 0 {
			runtimeError(expr.Operator, "zero division")
		}

		return leftValue.(float64) / rightValue.(float64)
	case lexing.Less:
		switch {
		case isNumber(leftValue) && isNumber(rightValue):
			return leftValue.(float64) < rightValue.(float64)
		case isString(leftValue) && isString(rightValue):
			return leftValue.(string) < rightValue.(string)
		default:
			runtimeError(expr.Operator, "number or string operands expected")
		}
	case lexing.Greater:
		switch {
		case isNumber(leftValue) && isNumber(rightValue):
			return leftValue.(float64) > rightValue.(float64)
		case isString(leftValue) && isString(rightValue):
			return leftValue.(string) > rightValue.(string)
		default:
			runtimeError(expr.Operator, "number or string operands expected")
		}
	case lexing.LessEqual:
		switch {
		case isNumber(leftValue) && isNumber(rightValue):
			return leftValue.(float64) <= rightValue.(float64)
		case isString(leftValue) && isString(rightValue):
			return leftValue.(string) <= rightValue.(string)
		default:
			runtimeError(expr.Operator, "number or string operands expected")
		}
	case lexing.GreaterEqual:
		switch {
		case isNumber(leftValue) && isNumber(rightValue):
			return leftValue.(float64) >= rightValue.(float64)
		case isString(leftValue) && isString(rightValue):
			return leftValue.(string) >= rightValue.(string)
		default:
			runtimeError(expr.Operator, "number or string operands expected")
		}
	case lexing.EqualEqual:
		switch {
		case leftValue == nil && rightValue == nil:
			return true
		case leftValue == nil || rightValue == nil:
			return false
		case isNumber(leftValue) && isNumber(rightValue):
			return leftValue.(float64) == rightValue.(float64)
		case isString(leftValue) && isString(rightValue):
			return leftValue.(string) == rightValue.(string)
		default:
			runtimeError(expr.Operator, "number, string or nil operands expected")
		}
	case lexing.BangEqual:
		switch {
		case leftValue == nil && rightValue == nil:
			return false
		case leftValue == nil || rightValue == nil:
			return true
		case isNumber(leftValue) && isNumber(rightValue):
			return leftValue.(float64) != rightValue.(float64)
		case isString(leftValue) && isString(rightValue):
			return leftValue.(string) != rightValue.(string)
		default:
			runtimeError(expr.Operator, "number, string or nil operands expected")
		}
	}

	return nil
}

type UnaryExpr struct {
	Operator  lexing.Token
	RightExpr Expr
}

func (expr UnaryExpr) Evaluate() interface{} {
	value := expr.RightExpr.Evaluate()

	switch expr.Operator.TokenType {
	case lexing.Minus:
		return -value.(float64)
	case lexing.Bang:
		return !isTruthy(value)
	}

	return nil
}

type LiteralExpr struct {
	LiteralValue interface{}
}

func (expr LiteralExpr) Evaluate() interface{} {
	return expr.LiteralValue
}

type GroupingExpr struct {
	Expr Expr
}

func (expr GroupingExpr) Evaluate() interface{} {
	return expr.Expr.Evaluate()
}

type Stmt interface {
	Execute()
}

type ExpressionStmt struct {
	Expr Expr
}

func (stmt ExpressionStmt) Execute() {
	stmt.Expr.Evaluate()
}

type PrintStmt struct {
	Expr Expr
}

func (stmt PrintStmt) Execute() {
	value := stmt.Expr.Evaluate()
	fmt.Printf("%v\n", value)
}

func requireNumberOperand(operator lexing.Token, operand interface{}) {
	switch operand.(type) {
	case float64:
		return
	}

	runtimeError(operator, "number operand expected")
}
