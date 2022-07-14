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
	leftValue := node.LeftExpr.Evaluate()
	rightValue := node.RightExpr.Evaluate()

	switch node.Operator.TokenType {
	case lexing.Plus:
		switch {
		case isNumber(leftValue) && isNumber(rightValue):
			return leftValue.(float64) + rightValue.(float64)
		case isString(leftValue) && isString(rightValue):
			return leftValue.(string) + rightValue.(string)
		default:
			runtimeError(node.Operator, "number or string operands expected")
		}
	case lexing.Minus:
		requireNumberOperand(node.Operator, leftValue)
		requireNumberOperand(node.Operator, rightValue)
		return leftValue.(float64) - rightValue.(float64)
	case lexing.Star:
		requireNumberOperand(node.Operator, leftValue)
		requireNumberOperand(node.Operator, rightValue)
		return leftValue.(float64) * rightValue.(float64)
	case lexing.Slash:
		requireNumberOperand(node.Operator, leftValue)
		requireNumberOperand(node.Operator, rightValue)

		if rightValue.(float64) == 0 {
			runtimeError(node.Operator, "zero division")
		}

		return leftValue.(float64) / rightValue.(float64)
	case lexing.Less:
		switch {
		case isNumber(leftValue) && isNumber(rightValue):
			return leftValue.(float64) < rightValue.(float64)
		case isString(leftValue) && isString(rightValue):
			return leftValue.(string) < rightValue.(string)
		default:
			runtimeError(node.Operator, "number or string operands expected")
		}
	case lexing.Greater:
		switch {
		case isNumber(leftValue) && isNumber(rightValue):
			return leftValue.(float64) > rightValue.(float64)
		case isString(leftValue) && isString(rightValue):
			return leftValue.(string) > rightValue.(string)
		default:
			runtimeError(node.Operator, "number or string operands expected")
		}
	case lexing.LessEqual:
		switch {
		case isNumber(leftValue) && isNumber(rightValue):
			return leftValue.(float64) <= rightValue.(float64)
		case isString(leftValue) && isString(rightValue):
			return leftValue.(string) <= rightValue.(string)
		default:
			runtimeError(node.Operator, "number or string operands expected")
		}
	case lexing.GreaterEqual:
		switch {
		case isNumber(leftValue) && isNumber(rightValue):
			return leftValue.(float64) >= rightValue.(float64)
		case isString(leftValue) && isString(rightValue):
			return leftValue.(string) >= rightValue.(string)
		default:
			runtimeError(node.Operator, "number or string operands expected")
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
			runtimeError(node.Operator, "number, string or nil operands expected")
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
			runtimeError(node.Operator, "number, string or nil operands expected")
		}
	}

	return nil
}

type UnaryExpr struct {
	Operator  lexing.Token
	RightExpr Expr
}

func (node UnaryExpr) Evaluate() interface{} {
	value := node.RightExpr.Evaluate()

	switch node.Operator.TokenType {
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

func (node LiteralExpr) Evaluate() interface{} {
	return node.LiteralValue
}

type GroupingExpr struct {
	Expr Expr
}

func (node GroupingExpr) Evaluate() interface{} {
	return node.Expr.Evaluate()
}

func requireNumberOperand(operator lexing.Token, operand interface{}) {
	switch operand.(type) {
	case float64:
		return
	}

	runtimeError(operator, "number operand expected")
}
