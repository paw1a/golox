package runtime

import (
	"fmt"
	"github.com/paw1a/golox/internal/ast"
	"github.com/paw1a/golox/internal/lexing"
)

func (i *Interpreter) Execute(stmt ast.Stmt) {
	switch stmt.(type) {
	case ast.ExpressionStmt:
		i.executeExprStmt(stmt.(ast.ExpressionStmt))
	case ast.PrintStmt:
		i.executePrintStmt(stmt.(ast.PrintStmt))
	case ast.VarDeclarationStmt:
		i.executeVarDeclarationStmt(stmt.(ast.VarDeclarationStmt))
	case ast.BlockStmt:
		i.executeBlockStmt(stmt.(ast.BlockStmt))
	case ast.IfStmt:
		i.executeIfStmt(stmt.(ast.IfStmt))
	case ast.ForStmt:
		i.executeForStmt(stmt.(ast.ForStmt))
	case ast.BreakStmt:
		i.executeBreakStmt()
	case ast.ContinueStmt:
		i.executeContinueStmt()
	case ast.FunDeclarationStmt:
		i.executeFunDeclarationStmt(stmt.(ast.FunDeclarationStmt))
	case ast.ReturnStmt:
		i.executeReturnStmt(stmt.(ast.ReturnStmt))
	default:
		runtimeError(lexing.Token{}, "invalid ast type")
	}
}

func (i *Interpreter) executeExprStmt(stmt ast.ExpressionStmt) {
	i.Evaluate(stmt.Expr)
}

func (i *Interpreter) executePrintStmt(stmt ast.PrintStmt) {
	value := i.Evaluate(stmt.Expr)
	fmt.Printf("%v\n", value)
}

func (i *Interpreter) executeVarDeclarationStmt(stmt ast.VarDeclarationStmt) {
	var value interface{}
	if stmt.Initializer != nil {
		value = i.Evaluate(stmt.Initializer)
	}

	i.env.define(stmt.Name.Lexeme, value)
}

func (i *Interpreter) executeFunDeclarationStmt(stmt ast.FunDeclarationStmt) {
	function := Function{Declaration: stmt}
	i.global.define(stmt.Name.Lexeme, function)
}

func (i *Interpreter) executeBlockStmt(blockStmt ast.BlockStmt) {
	enclosingEnv := i.env
	i.env = NewEnvironment(enclosingEnv)
	defer func() {
		i.env = enclosingEnv
	}()

	for _, stmt := range blockStmt.Stmts {
		if i.loopContext.breakFlag || i.loopContext.continueFlag ||
			i.returnContext.returnFlag {
			break
		}
		i.Execute(stmt)
	}
}

func (i *Interpreter) executeIfStmt(stmt ast.IfStmt) {
	conditionValue := i.Evaluate(stmt.ConditionExpr)

	if isTruthy(conditionValue) {
		i.Execute(stmt.IfStatement)
		return
	}

	if stmt.ElseStatement != nil {
		i.Execute(stmt.ElseStatement)
	}
}

func (i *Interpreter) executeForStmt(stmt ast.ForStmt) {
	enclosingEnv := i.env
	i.env = NewEnvironment(enclosingEnv)
	defer func() {
		i.env = enclosingEnv
	}()

	if stmt.InitializerStmt != nil {
		i.Execute(stmt.InitializerStmt)
	}

	for isTruthy(i.Evaluate(stmt.ConditionExpr)) {
		i.Execute(stmt.Statement)
		if i.loopContext.breakFlag {
			i.loopContext.breakFlag = false
			break
		}
		if i.loopContext.continueFlag {
			i.loopContext.continueFlag = false
		}

		if stmt.IncrementExpr != nil {
			i.Evaluate(stmt.IncrementExpr)
		}
	}
}

func (i *Interpreter) executeBreakStmt() {
	i.loopContext.breakFlag = true
}

func (i *Interpreter) executeContinueStmt() {
	i.loopContext.continueFlag = true
}

func (i *Interpreter) executeReturnStmt(stmt ast.ReturnStmt) {
	if stmt.Expr != nil {
		i.returnContext.returnValue = i.Evaluate(stmt.Expr)
	} else {
		i.returnContext.returnValue = nil
	}
	i.returnContext.returnFlag = true
}
