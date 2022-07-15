package parsing

import (
	"fmt"
	"github.com/paw1a/golox/internal/ast"
	"github.com/paw1a/golox/internal/lexing"
)

type Parser struct {
	tokens  []lexing.Token
	current int

	Errors []error
}

func (p *Parser) Parse() []ast.Stmt {
	var statements []ast.Stmt

	for !p.isEof() {
		statements = append(statements, p.declaration())
	}

	return statements
}

func (p *Parser) declaration() ast.Stmt {
	defer p.parseRecoverFunc()

	if p.match(lexing.Var) {
		p.advance()
		return p.varDeclaration()
	}

	return p.statement()
}

func (p *Parser) varDeclaration() ast.Stmt {
	varName := p.requireToken(lexing.Identifier, "variable name expected")

	var initializer ast.Expr
	if p.peek().TokenType == lexing.Equal {
		p.advance()
		initializer = p.expression()
	}

	p.requireToken(lexing.Semicolon, "';' expected")
	return ast.VarDeclarationStmt{
		Name:        varName,
		Initializer: initializer,
	}
}

func (p *Parser) statement() ast.Stmt {
	if p.match(lexing.Print) {
		p.advance()
		return p.printStatement()
	}

	if p.match(lexing.LeftBrace) {
		p.advance()
		return p.blockStatement()
	}

	return p.expressionStatement()
}

func (p *Parser) blockStatement() ast.Stmt {
	var stmts []ast.Stmt

	for !p.match(lexing.RightBrace) && !p.isEof() {
		stmts = append(stmts, p.declaration())
	}

	p.requireToken(lexing.RightBrace, "'}' end of block expected")

	return ast.BlockStmt{Stmts: stmts}
}

func (p *Parser) printStatement() ast.Stmt {
	expr := p.expression()
	p.requireToken(lexing.Semicolon, "';' expected")
	return ast.PrintStmt{Expr: expr}
}

func (p *Parser) expressionStatement() ast.Stmt {
	expr := p.expression()
	p.requireToken(lexing.Semicolon, "';' expected")
	return ast.ExpressionStmt{Expr: expr}
}

func (p *Parser) expression() ast.Expr {
	return p.assignment()
}

func (p *Parser) assignment() ast.Expr {
	expr := p.equality()

	if p.match(lexing.Equal) {
		equalToken := p.advance()
		value := p.expression()

		switch expr.(type) {
		case ast.VariableExpr:
			return ast.AssignExpr{
				Name:        expr.(ast.VariableExpr).Name,
				Initializer: value,
			}
		}

		parseError(equalToken, "invalid assignment target")
	}

	return expr
}

func (p *Parser) equality() ast.Expr {
	var expr ast.Expr

	expr = p.comparison()
	for p.match(lexing.BangEqual, lexing.EqualEqual) {
		operator := p.advance()
		rightExpr := p.comparison()
		expr = ast.BinaryExpr{
			LeftExpr:  expr,
			Operator:  operator,
			RightExpr: rightExpr,
		}
	}

	return expr
}

func (p *Parser) comparison() ast.Expr {
	var expr ast.Expr

	expr = p.term()
	for p.match(lexing.Less, lexing.LessEqual, lexing.Greater, lexing.GreaterEqual) {
		operator := p.advance()
		rightExpr := p.term()
		expr = ast.BinaryExpr{
			LeftExpr:  expr,
			Operator:  operator,
			RightExpr: rightExpr,
		}
	}

	return expr
}

func (p *Parser) term() ast.Expr {
	var expr ast.Expr

	expr = p.factor()
	for p.match(lexing.Minus, lexing.Plus) {
		operator := p.advance()
		rightExpr := p.factor()
		expr = ast.BinaryExpr{
			LeftExpr:  expr,
			Operator:  operator,
			RightExpr: rightExpr,
		}
	}

	return expr
}

func (p *Parser) factor() ast.Expr {
	var expr ast.Expr

	expr = p.unary()
	for p.match(lexing.Star, lexing.Slash) {
		operator := p.advance()
		rightExpr := p.unary()
		expr = ast.BinaryExpr{
			LeftExpr:  expr,
			Operator:  operator,
			RightExpr: rightExpr,
		}
	}

	return expr
}

func (p *Parser) unary() ast.Expr {
	for p.match(lexing.Bang, lexing.Minus) {
		operator := p.advance()
		rightExpr := p.unary()
		return ast.UnaryExpr{
			Operator:  operator,
			RightExpr: rightExpr,
		}
	}

	return p.primary()
}

func (p *Parser) primary() ast.Expr {
	switch {
	case p.match(lexing.False):
		p.advance()
		return ast.LiteralExpr{LiteralValue: false}
	case p.match(lexing.True):
		p.advance()
		return ast.LiteralExpr{LiteralValue: true}
	case p.match(lexing.Nil):
		p.advance()
		return ast.LiteralExpr{LiteralValue: nil}
	case p.match(lexing.Number, lexing.String):
		astNode := ast.LiteralExpr{LiteralValue: p.peek().Literal}
		p.advance()
		return astNode
	case p.match(lexing.LeftParen):
		p.advance()
		expr := p.expression()
		p.requireToken(lexing.RightParen, "expect ')' token after expression")
		return ast.GroupingExpr{Expr: expr}
	case p.match(lexing.Identifier):
		return ast.VariableExpr{Name: p.advance()}
	default:
		parseError(p.peek(), "expect expression")
		return nil
	}
}

func (p *Parser) requireToken(tokenType lexing.TokenType, message string) lexing.Token {
	if p.match(tokenType) {
		return p.advance()
	} else {
		parseError(p.peek(), message)
		return lexing.Token{}
	}
}

func parseError(token lexing.Token, message string) {
	var errorMessage string
	if token.TokenType == lexing.Eof {
		errorMessage = fmt.Sprintf("line %d | at end of input: %s", token.Line, message)
	} else {
		errorMessage = fmt.Sprintf("line %d | at '%s': %s", token.Line, token.Lexeme, message)
	}
	panic(errorMessage)
}

func (p *Parser) parseRecoverFunc() {
	if err := recover(); err != nil {
		p.Errors = append(p.Errors, fmt.Errorf("%v", err))
	}
}

func NewParser(tokens []lexing.Token) *Parser {
	return &Parser{
		tokens: tokens,
	}
}
