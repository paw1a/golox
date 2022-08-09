package runtime

import (
	"fmt"
	"github.com/paw1a/golox/internal/ast"
	"github.com/paw1a/golox/internal/lexing"
)

func (i *Interpreter) Evaluate(expr ast.Expr) interface{} {
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
	case ast.TernaryExpr:
		return i.evaluateTernaryExpr(expr.(ast.TernaryExpr))
	case ast.LogicalExpr:
		return i.evaluateLogicalExpr(expr.(ast.LogicalExpr))
	case ast.CallExpr:
		return i.evaluateCallExpr(expr.(ast.CallExpr))
	case ast.IndexExpr:
		return i.evaluateIndexExpr(expr.(ast.IndexExpr))
	case ast.ArrayExpr:
		return i.evaluateArrayExpr(expr.(ast.ArrayExpr))
	case ast.LambdaExpr:
		return i.evaluateLambdaExpr(expr.(ast.LambdaExpr))
	default:
		runtimeError(lexing.Token{}, "invalid ast type")
	}

	return nil
}

func (i *Interpreter) evaluateBinaryExpr(expr ast.BinaryExpr) interface{} {
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

func (i *Interpreter) evaluateUnaryExpr(expr ast.UnaryExpr) interface{} {
	value := i.Evaluate(expr.RightExpr)

	switch expr.Operator.TokenType {
	case lexing.Minus:
		return -value.(float64)
	case lexing.Bang:
		return !isTruthy(value)
	}

	return nil
}

func (i *Interpreter) evaluateLiteralExpr(expr ast.LiteralExpr) interface{} {
	return expr.LiteralValue
}

func (i *Interpreter) evaluateGroupingExpr(expr ast.GroupingExpr) interface{} {
	return i.Evaluate(expr.Expr)
}

func (i *Interpreter) evaluateVariableExpr(expr ast.VariableExpr) interface{} {
	return i.lookUpVariable(expr.Name, expr)
}

func (i *Interpreter) evaluateAssignExpr(expr ast.AssignExpr) interface{} {
	value := i.Evaluate(expr.Initializer)
	switch expr.Variable.(type) {
	case ast.IndexExpr:
		array := i.Evaluate(expr.Variable.(ast.IndexExpr).Array).([]interface{})
		indexValue := i.Evaluate(expr.Variable.(ast.IndexExpr).IndexExpr)
		var index int
		if isNumber(indexValue) {
			index = int(indexValue.(float64))
		} else {
			runtimeError(expr.Variable.(ast.IndexExpr).Bracket, "index must be integer number")
		}

		if index >= len(array) {
			runtimeError(expr.Variable.(ast.IndexExpr).Bracket,
				fmt.Sprintf("index %d out of range in array with len %d", index, len(array)))
		}
		array[index] = value
	case ast.VariableExpr:
		i.env.assign(expr.Variable.(ast.VariableExpr).Name, value)
	}
	return value
}

func (i *Interpreter) evaluateTernaryExpr(expr ast.TernaryExpr) interface{} {
	conditionValue := i.Evaluate(expr.Condition)

	var value interface{}
	if isTruthy(conditionValue) {
		value = i.Evaluate(expr.TrueExpr)
	} else {
		value = i.Evaluate(expr.FalseExpr)
	}

	return value
}

func (i *Interpreter) evaluateLogicalExpr(expr ast.LogicalExpr) interface{} {
	leftValue := i.Evaluate(expr.LeftExpr)

	if expr.Operator.TokenType == lexing.Or && isTruthy(leftValue) ||
		expr.Operator.TokenType == lexing.And && !isTruthy(leftValue) {
		return leftValue
	}

	return i.Evaluate(expr.RightExpr)
}

func (i *Interpreter) evaluateCallExpr(expr ast.CallExpr) interface{} {
	calleeValue := i.Evaluate(expr.Callee)

	argumentValues := make([]interface{}, 0)
	for _, argumentExpr := range expr.Arguments {
		argumentValues = append(argumentValues, i.Evaluate(argumentExpr))
	}

	switch calleeValue.(type) {
	case Caller:
		function := calleeValue.(Caller)
		if function.ParametersCount() != len(argumentValues) {
			runtimeError(expr.Paren,
				fmt.Sprintf("expect %d arguments, got %d",
					function.ParametersCount(), len(argumentValues)))
		}
		return function.Call(i, argumentValues)
	}

	runtimeError(expr.Paren, "invalid object to call")
	return nil
}

func (i *Interpreter) evaluateIndexExpr(expr ast.IndexExpr) interface{} {
	arrayValue := i.Evaluate(expr.Array)

	switch arrayValue.(type) {
	case []interface{}:
		array := arrayValue.([]interface{})

		indexValue := i.Evaluate(expr.IndexExpr)
		var index int
		if isNumber(indexValue) {
			index = int(indexValue.(float64))
		} else {
			runtimeError(expr.Bracket, "index must be integer number")
		}

		if index >= len(array) {
			runtimeError(expr.Bracket,
				fmt.Sprintf("index %d out of range in array with len %d", index, len(array)))
		}

		return array[index]
	}

	runtimeError(expr.Bracket, "invalid array object")
	return nil
}

func (i *Interpreter) evaluateArrayExpr(expr ast.ArrayExpr) interface{} {
	array := make([]interface{}, len(expr.Elements))

	for index, elemExpr := range expr.Elements {
		elemValue := i.Evaluate(elemExpr)
		array[index] = elemValue
	}

	return array
}

func (i *Interpreter) evaluateLambdaExpr(expr ast.LambdaExpr) interface{} {
	return LambdaFunction{
		LambdaExpr: expr,
		Closure:    i.env,
	}
}

func requireNumberOperand(operator lexing.Token, operand interface{}) {
	switch operand.(type) {
	case float64:
		return
	}

	runtimeError(operator, "number operand expected")
}
