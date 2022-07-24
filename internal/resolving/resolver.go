package resolving

import (
	"fmt"
	"github.com/paw1a/golox/internal/ast"
	"github.com/paw1a/golox/internal/lexing"
	"github.com/paw1a/golox/internal/runtime"
	"github.com/paw1a/golox/internal/stack"
)

type Resolver struct {
	interpreter *runtime.Interpreter
	scopes      stack.Stack[map[string]bool]
}

func (r *Resolver) beginScope() {
	r.scopes.Push(make(map[string]bool))
}

func (r *Resolver) endScope() {
	r.scopes.Pop()
}

func (r *Resolver) declare(token lexing.Token) {
	if scope, ok := r.scopes.Peek(); ok {
		scope[token.Lexeme] = false
	}
}

func (r *Resolver) define(token lexing.Token) {
	if scope, ok := r.scopes.Peek(); ok {
		scope[token.Lexeme] = true
	}
}

func (r *Resolver) resolveLocal(expr ast.Expr, token lexing.Token) {
	for i := r.scopes.Size() - 1; i >= 0; i-- {
		scope, _ := r.scopes.Get(i)
		if _, ok := scope[token.Lexeme]; ok {
			r.interpreter.Resolve(expr, r.scopes.Size()-1-i)
			return
		}
	}
}

func resolveError(token lexing.Token, message string) {
	var errorMessage string
	if token.TokenType == lexing.Eof {
		errorMessage = fmt.Sprintf("line %d | at end of input: %s", token.Line, message)
	} else {
		errorMessage = fmt.Sprintf("line %d | at '%s': %s", token.Line, token.Lexeme, message)
	}
	panic(errorMessage)
}

func NewResolver(interpreter *runtime.Interpreter) *Resolver {
	return &Resolver{
		interpreter: interpreter,
		scopes:      stack.New[map[string]bool](),
	}
}
