package lexing

import "unicode"

func (l *Lexer) advance() uint8 {
	c := l.source[l.current]
	l.current++
	return c
}

func (l *Lexer) peek() uint8 {
	return l.source[l.current]
}

func (l *Lexer) peekNext() uint8 {
	return l.source[l.current+1]
}

func (l *Lexer) isDigit(c uint8) bool {
	return unicode.IsDigit(rune(c))
}

func (l *Lexer) isAlpha(c uint8) bool {
	return c >= 'a' && c <= 'z' || c >= 'A' && c <= 'Z' || c == '_'
}

func (l *Lexer) isEOF() bool {
	return l.current >= len(l.source)
}
