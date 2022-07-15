package parsing

import "github.com/paw1a/golox/internal/lexing"

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
