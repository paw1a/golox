package resolving

import (
	"github.com/paw1a/golox/internal/ast"
	"github.com/paw1a/golox/internal/lexing"
)

func (r *Resolver) ResolveStmt(stmt ast.Stmt) {
	switch stmt.(type) {
	case ast.ExpressionStmt:
		r.exprStmt(stmt.(ast.ExpressionStmt))
	case ast.PrintStmt:
		r.printStmt(stmt.(ast.PrintStmt))
	case ast.VarDeclarationStmt:
		r.varDeclarationStmt(stmt.(ast.VarDeclarationStmt))
	case ast.BlockStmt:
		r.blockStmt(stmt.(ast.BlockStmt))
	case ast.IfStmt:
		r.ifStmt(stmt.(ast.IfStmt))
	case ast.ForStmt:
		r.forStmt(stmt.(ast.ForStmt))
	case ast.BreakStmt:
		r.breakStmt()
	case ast.ContinueStmt:
		r.continueStmt()
	case ast.FunDeclarationStmt:
		r.funDeclarationStmt(stmt.(ast.FunDeclarationStmt))
	case ast.ReturnStmt:
		r.returnStmt(stmt.(ast.ReturnStmt))
	default:
		resolveError(lexing.Token{}, "resolver: invalid ast type")
	}
}

func (r *Resolver) exprStmt(stmt ast.ExpressionStmt) {
	r.ResolveExpr(stmt.Expr)
}

func (r *Resolver) printStmt(stmt ast.PrintStmt) {
	r.ResolveExpr(stmt.Expr)
}

func (r *Resolver) varDeclarationStmt(stmt ast.VarDeclarationStmt) {
	r.declare(stmt.Name)
	if stmt.Initializer != nil {
		r.ResolveExpr(stmt.Initializer)
	}
	r.define(stmt.Name)
}

func (r *Resolver) funDeclarationStmt(stmt ast.FunDeclarationStmt) {
	r.define(stmt.Name)

	r.beginScope()
	for _, param := range stmt.Params {
		r.define(param)
	}
	r.ResolveStmt(stmt.Statement)
	r.endScope()
}

func (r *Resolver) blockStmt(blockStmt ast.BlockStmt) {
	r.beginScope()
	for _, stmt := range blockStmt.Stmts {
		r.ResolveStmt(stmt)
	}
	r.endScope()
}

func (r *Resolver) ifStmt(stmt ast.IfStmt) {
	r.ResolveExpr(stmt.ConditionExpr)
	r.ResolveStmt(stmt.IfStatement)

	if stmt.ElseStatement != nil {
		r.ResolveStmt(stmt.ElseStatement)
	}
}

func (r *Resolver) forStmt(stmt ast.ForStmt) {
	r.beginScope()
	if stmt.InitializerStmt != nil {
		r.ResolveStmt(stmt.InitializerStmt)
	}
	if stmt.ConditionExpr != nil {
		r.ResolveExpr(stmt.ConditionExpr)
	}
	if stmt.IncrementExpr != nil {
		r.ResolveExpr(stmt.IncrementExpr)
	}
	r.ResolveStmt(stmt.Statement)
	r.endScope()
}

func (r *Resolver) breakStmt() {
}

func (r *Resolver) continueStmt() {
}

func (r *Resolver) returnStmt(stmt ast.ReturnStmt) {
	if stmt.Expr != nil {
		r.ResolveExpr(stmt.Expr)
	}
}
