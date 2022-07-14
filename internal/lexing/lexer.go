package lexing

import (
	"fmt"
	"strconv"
)

type Lexer struct {
	Tokens []Token
	source string

	start   int
	current int
	line    int

	Errors []error
}

func (l *Lexer) ScanTokens() []Token {
	for !l.isEOF() {
		l.start = l.current
		l.ScanToken()
	}

	l.Tokens = append(l.Tokens, NewToken(Eof, "", nil, l.line))

	return l.Tokens
}

func (l *Lexer) ScanToken() {
	c := l.advance()
	switch c {
	case '(':
		l.addToken(LeftParen)
	case ')':
		l.addToken(RightParen)
	case '{':
		l.addToken(LeftBrace)
	case '}':
		l.addToken(RightBrace)
	case ',':
		l.addToken(Comma)
	case '.':
		l.addToken(Dot)
	case '-':
		l.addToken(Minus)
	case '+':
		l.addToken(Plus)
	case ';':
		l.addToken(Semicolon)
	case '*':
		l.addToken(Star)
	case '!':
		if l.peek() == '=' {
			l.advance()
			l.addToken(BangEqual)
		} else {
			l.addToken(Bang)
		}
	case '<':
		if l.peek() == '=' {
			l.advance()
			l.addToken(LessEqual)
		} else {
			l.addToken(Less)
		}
	case '>':
		if l.peek() == '=' {
			l.advance()
			l.addToken(GreaterEqual)
		} else {
			l.addToken(Greater)
		}
	case '=':
		if l.peek() == '=' {
			l.advance()
			l.addToken(EqualEqual)
		} else {
			l.addToken(Equal)
		}
	case '/':
		if l.peek() == '/' {
			l.advance()
			for !l.isEOF() && l.peek() != '\n' {
				l.advance()
			}
		} else if l.peek() == '*' {
			l.advance()
			l.blockComment()
		} else {
			l.addToken(Slash)
		}
	case '"':
		l.string()
	case '\n':
		l.line++
	case ' ', '\t', '\r':
	default:
		if l.isDigit(c) {
			l.number()
		} else if l.isAlpha(c) {
			l.identifier()
		} else {
			l.error("unexpected token")
		}
	}
}

func (l *Lexer) string() {
	for !l.isEOF() && l.peek() != '"' {
		if l.peek() == '\n' {
			l.line++
		}
		l.advance()
	}

	if l.isEOF() {
		l.error("no closing \" quote")
		return
	}

	l.advance()

	literalValue := l.source[l.start+1 : l.current-1]
	l.addTokenWithLiteral(String, literalValue)
}

func (l *Lexer) number() {
	for !l.isEOF() && l.isDigit(l.peek()) {
		l.advance()
	}

	if !l.isEOF() && l.peek() == '.' && l.isDigit(l.peekNext()) {
		l.advance()
	} else {
		numberValue, _ := strconv.ParseFloat(l.source[l.start:l.current], 64)
		l.addTokenWithLiteral(Number, numberValue)
		return
	}

	for !l.isEOF() && l.isDigit(l.peek()) {
		l.advance()
	}

	numberValue, _ := strconv.ParseFloat(l.source[l.start:l.current], 64)
	l.addTokenWithLiteral(Number, numberValue)
}

func (l *Lexer) blockComment() {
	for {
		if l.peek() == '*' {
			l.advance()
			if l.peek() == '/' {
				l.advance()
				if l.peek() == '\n' {
					l.line++
					l.advance()
				}
				return
			}
		}

		if l.peek() == '/' {
			l.advance()
			if l.peek() == '*' {
				l.advance()
				l.blockComment()
			}
		}

		if l.peek() == '\n' {
			l.line++
		}

		if l.isEOF() {
			l.error("unterminated block comment")
			return
		}
		l.advance()
	}
}

var keywords = map[string]TokenType{
	"and":    And,
	"class":  Class,
	"else":   Else,
	"false":  False,
	"fun":    Fun,
	"for":    For,
	"if":     If,
	"nil":    Nil,
	"or":     Or,
	"print":  Print,
	"return": Return,
	"super":  Super,
	"this":   This,
	"true":   True,
	"var":    Var,
	"while":  While,
}

func (l *Lexer) identifier() {
	for !l.isEOF() && (l.isAlpha(l.peek()) || l.isDigit(l.peek())) {
		l.advance()
	}

	identifier := l.source[l.start:l.current]
	identifierType, ok := keywords[identifier]
	if ok {
		l.addToken(identifierType)
	} else {
		l.addToken(Identifier)
	}
}

func (l *Lexer) addToken(tokenType TokenType) {
	l.addTokenWithLiteral(tokenType, nil)
}

func (l *Lexer) addTokenWithLiteral(tokenType TokenType, literal interface{}) {
	lexeme := l.source[l.start:l.current]
	l.Tokens = append(l.Tokens, NewToken(tokenType, lexeme, literal, l.line))
}

func (l *Lexer) error(message string) {
	err := fmt.Errorf("line %d | error: %s", l.line, message)
	l.Errors = append(l.Errors, err)
}

func NewLexer(source string) *Lexer {
	return &Lexer{
		source: source,
		line:   1,
	}
}
