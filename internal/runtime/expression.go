package runtime

import (
	"github.com/paw1a/golox/internal/ast"
	"github.com/paw1a/golox/internal/lexing"
)

func (i Interpreter) Evaluate(expr ast.Expr) interface{} {
	switch expr.(type) {
	case ast.BinaryExpr:
		return i.evaluateBinaryExpr(expr.(ast.BinaryExpr))
	case ast.UnaryExpr:
		return i.evaluateUnaryExpr(expr.(ast.UnaryExpr))
	case ast.LiteralExpr:
		return i.evaluateLiteralExpr(expr.(ast.LiteralExpr))
	case ast.GroupingExpr:
		return i.evaluateGroupingExpr(expr.(ast.GroupingExpr))
	case ast.VariableExpr:
		return i.evaluateVariableExpr(expr.(ast.VariableExpr))
	case ast.AssignExpr:
		return i.evaluateAssignExpr(expr.(ast.AssignExpr))
	default:
		runtimeError(lexing.Token{}, "invalid ast type")
	}

	return nil
}

func (i Interpreter) evaluateBinaryExpr(expr ast.BinaryExpr) interface{} {
	leftValue := i.Evaluate(expr.LeftExpr)
	rightValue := i.Evaluate(expr.RightExpr)

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
	case lexing.Comma:
		return rightValue
	}

	return nil
}

func (i Interpreter) evaluateUnaryExpr(expr ast.UnaryExpr) interface{} {
	value := i.Evaluate(expr.RightExpr)

	switch expr.Operator.TokenType {
	case lexing.Minus:
		return -value.(float64)
	case lexing.Bang:
		return !isTruthy(value)
	}

	return nil
}

func (i Interpreter) evaluateLiteralExpr(expr ast.LiteralExpr) interface{} {
	return expr.LiteralValue
}

func (i Interpreter) evaluateGroupingExpr(expr ast.GroupingExpr) interface{} {
	return i.Evaluate(expr.Expr)
}

func (i Interpreter) evaluateVariableExpr(expr ast.VariableExpr) interface{} {
	return i.env.get(expr.Name)
}

func (i Interpreter) evaluateAssignExpr(expr ast.AssignExpr) interface{} {
	value := i.Evaluate(expr.Initializer)
	i.env.assign(expr.Name, value)
	return value
}

func requireNumberOperand(operator lexing.Token, operand interface{}) {
	switch operand.(type) {
	case float64:
		return
	}

	runtimeError(operator, "number operand expected")
}
