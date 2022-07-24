package resolving

import (
	"github.com/paw1a/golox/internal/ast"
	"github.com/paw1a/golox/internal/lexing"
)

func (r *Resolver) ResolveExpr(expr ast.Expr) {
	switch expr.(type) {
	case ast.BinaryExpr:
		r.binaryExpr(expr.(ast.BinaryExpr))
	case ast.UnaryExpr:
		r.unaryExpr(expr.(ast.UnaryExpr))
	case ast.LiteralExpr:
		r.literalExpr()
	case ast.GroupingExpr:
		r.groupingExpr(expr.(ast.GroupingExpr))
	case ast.VariableExpr:
		r.variableExpr(expr.(ast.VariableExpr))
	case ast.AssignExpr:
		r.assignExpr(expr.(ast.AssignExpr))
	case ast.TernaryExpr:
		r.ternaryExpr(expr.(ast.TernaryExpr))
	case ast.LogicalExpr:
		r.logicalExpr(expr.(ast.LogicalExpr))
	case ast.CallExpr:
		r.callExpr(expr.(ast.CallExpr))
	case ast.LambdaExpr:
		r.lambdaExpr(expr.(ast.LambdaExpr))
	default:
		resolveError(lexing.Token{}, "resolver: invalid ast type")
	}
}

func (r *Resolver) binaryExpr(expr ast.BinaryExpr) {
	r.ResolveExpr(expr.LeftExpr)
	r.ResolveExpr(expr.RightExpr)
}

func (r *Resolver) unaryExpr(expr ast.UnaryExpr) {
	r.ResolveExpr(expr.RightExpr)
}

func (r *Resolver) literalExpr() {
}

func (r *Resolver) groupingExpr(expr ast.GroupingExpr) {
	r.ResolveExpr(expr.Expr)
}

func (r *Resolver) variableExpr(expr ast.VariableExpr) {
	if scope, ok := r.scopes.Peek(); ok && scope[expr.Name.Lexeme] {
		r.resolveLocal(expr, expr.Name)
	} else {
		resolveError(expr.Name, "can't read local variable in its own initializer.")
	}
}

func (r *Resolver) assignExpr(expr ast.AssignExpr) {
	r.ResolveExpr(expr.Initializer)
	r.resolveLocal(expr, expr.Name)
}

func (r *Resolver) ternaryExpr(expr ast.TernaryExpr) {
	r.ResolveExpr(expr.Condition)
	r.ResolveExpr(expr.TrueExpr)
	r.ResolveExpr(expr.FalseExpr)
}

func (r *Resolver) logicalExpr(expr ast.LogicalExpr) {
	r.ResolveExpr(expr.LeftExpr)
	r.ResolveExpr(expr.RightExpr)
}

func (r *Resolver) callExpr(expr ast.CallExpr) {
	r.ResolveExpr(expr.Callee)
	for _, arg := range expr.Arguments {
		r.ResolveExpr(arg)
	}
}

func (r *Resolver) lambdaExpr(expr ast.LambdaExpr) {
	r.ResolveStmt(expr.Statement)
}
