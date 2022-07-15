package runtime

import (
	"github.com/paw1a/golox/internal/ast"
	"github.com/paw1a/golox/internal/lexing"
)

func Evaluate(expr ast.Expr) interface{} {
	switch expr.(type) {
	case ast.BinaryExpr:
		return evaluateBinaryExpr(expr.(ast.BinaryExpr))
	case ast.UnaryExpr:
		return evaluateUnaryExpr(expr.(ast.UnaryExpr))
	case ast.LiteralExpr:
		return evaluateLiteralExpr(expr.(ast.LiteralExpr))
	case ast.GroupingExpr:
		return evaluateGroupingExpr(expr.(ast.GroupingExpr))
	case ast.VariableExpr:
		return evaluateVariableExpr(expr.(ast.VariableExpr))
	}

	return nil
}

func evaluateBinaryExpr(expr ast.BinaryExpr) interface{} {
	leftValue := Evaluate(expr.LeftExpr)
	rightValue := Evaluate(expr.RightExpr)

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

func evaluateUnaryExpr(expr ast.UnaryExpr) interface{} {
	value := Evaluate(expr.RightExpr)

	switch expr.Operator.TokenType {
	case lexing.Minus:
		return -value.(float64)
	case lexing.Bang:
		return !isTruthy(value)
	}

	return nil
}

func evaluateLiteralExpr(expr ast.LiteralExpr) interface{} {
	return expr.LiteralValue
}

func evaluateGroupingExpr(expr ast.GroupingExpr) interface{} {
	return Evaluate(expr.Expr)
}

func evaluateVariableExpr(expr ast.VariableExpr) interface{} {
	return nil
}

func requireNumberOperand(operator lexing.Token, operand interface{}) {
	switch operand.(type) {
	case float64:
		return
	}

	runtimeError(operator, "number operand expected")
}
