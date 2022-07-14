package parsing

import (
	"fmt"
	"github.com/paw1a/golox/internal/lexing"
)

type Parser struct {
	tokens  []lexing.Token
	current int

	Errors []error
}

func (p *Parser) Parse() Expr {
	return p.expression()
}

func (p *Parser) expression() Expr {
	return p.equality()
}

func (p *Parser) equality() Expr {
	var expr Expr

	expr = p.comparison()
	for p.match(lexing.BangEqual, lexing.EqualEqual) {
		operator := p.advance()
		rightExpr := p.comparison()
		expr = BinaryExpr{
			LeftExpr:  expr,
			Operator:  operator,
			RightExpr: rightExpr,
		}
	}

	return expr
}

func (p *Parser) comparison() Expr {
	var expr Expr

	expr = p.term()
	for p.match(lexing.Less, lexing.LessEqual, lexing.Greater, lexing.GreaterEqual) {
		operator := p.advance()
		rightExpr := p.term()
		expr = BinaryExpr{
			LeftExpr:  expr,
			Operator:  operator,
			RightExpr: rightExpr,
		}
	}

	return expr
}

func (p *Parser) term() Expr {
	var expr Expr

	expr = p.factor()
	for p.match(lexing.Minus, lexing.Plus) {
		operator := p.advance()
		rightExpr := p.factor()
		expr = BinaryExpr{
			LeftExpr:  expr,
			Operator:  operator,
			RightExpr: rightExpr,
		}
	}

	return expr
}

func (p *Parser) factor() Expr {
	var expr Expr

	expr = p.unary()
	for p.match(lexing.Star, lexing.Slash) {
		operator := p.advance()
		rightExpr := p.unary()
		expr = BinaryExpr{
			LeftExpr:  expr,
			Operator:  operator,
			RightExpr: rightExpr,
		}
	}

	return expr
}

func (p *Parser) unary() Expr {
	for p.match(lexing.Bang, lexing.Minus) {
		operator := p.advance()
		rightExpr := p.unary()
		return UnaryExpr{
			Operator:  operator,
			RightExpr: rightExpr,
		}
	}

	return p.primary()
}

func (p *Parser) primary() Expr {
	switch {
	case p.match(lexing.False):
		p.advance()
		return LiteralExpr{LiteralValue: false}
	case p.match(lexing.True):
		p.advance()
		return LiteralExpr{LiteralValue: true}
	case p.match(lexing.Nil):
		p.advance()
		return LiteralExpr{LiteralValue: nil}
	case p.match(lexing.Number, lexing.String):
		astNode := LiteralExpr{LiteralValue: p.peek().Literal}
		p.advance()
		return astNode
	case p.match(lexing.LeftParen):
		p.advance()
		expr := p.expression()
		p.require(lexing.RightParen, "expect ')' token after expression")
		return GroupingExpr{Expr: expr}
	default:
		runtimeError(p.peek(), "expect expression")
		return nil
	}
}

func (p *Parser) require(tokenType lexing.TokenType, message string) {
	if p.match(tokenType) {
		p.advance()
	} else {
		runtimeError(p.peek(), message)
	}
}

func runtimeError(token lexing.Token, message string) {
	var errorMessage string
	if token.TokenType == lexing.Eof {
		errorMessage = fmt.Sprintf("line %d | at end of input: %s", token.Line, message)
	} else {
		errorMessage = fmt.Sprintf("line %d | at '%s': %s", token.Line, token.Lexeme, message)
	}
	panic(errorMessage)
}

func NewParser(tokens []lexing.Token) *Parser {
	return &Parser{tokens: tokens}
}
