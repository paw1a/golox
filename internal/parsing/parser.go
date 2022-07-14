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
		p.require(lexing.RightParen, "no right paren")
		return GroupingExpr{Expr: expr}
	default:
		p.error(p.peek(), "no primary tokens")
		return nil
	}
}

func (p *Parser) require(tokenType lexing.TokenType, message string) error {
	if p.match(tokenType) {
		p.advance()
		return nil
	} else {
		return nil
	}
}

func (p *Parser) match(tokenTypes ...lexing.TokenType) bool {
	for _, tokenType := range tokenTypes {
		if p.peek().TokenType == tokenType {
			return true
		}
	}
	return false
}

func (p *Parser) advance() lexing.Token {
	if !p.isEof() {
		token := p.tokens[p.current]
		p.current++
		return token
	}

	return p.tokens[p.current]
}

func (p *Parser) peek() lexing.Token {
	return p.tokens[p.current]
}

func (p *Parser) isEof() bool {
	return p.peek().TokenType == lexing.Eof
}

func (p *Parser) error(token lexing.Token, message string) {
	err := fmt.Errorf("line %d: '%s' | error: %s", token.Line, token.Lexeme, message)
	p.Errors = append(p.Errors, err)
}

func NewParser(tokens []lexing.Token) *Parser {
	return &Parser{tokens: tokens}
}
