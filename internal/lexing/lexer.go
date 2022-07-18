package lexing

import (
	"bytes"
	"fmt"
	"strconv"
	"strings"
)

type Lexer struct {
	Tokens []Token

	start   int
	current int
	line    int

	Errors []error

	source    string
	lineStart int
	Lines     []string
}

func (l *Lexer) ScanTokens() []Token {
	for !l.isEOF() {
		l.start = l.current
		l.ScanToken()
	}

	if l.source[len(l.source)-2] != '\n' {
		l.nextLine()
	}

	l.Tokens = append(l.Tokens, NewToken(Eof, "", nil, l.line, l.start-l.lineStart+1))

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
		l.nextLine()
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
			l.nextLine()
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
					l.nextLine()
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
			l.nextLine()
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
	position := l.start - l.lineStart
	l.Tokens = append(l.Tokens, NewToken(tokenType, lexeme, literal, l.line, position))
}

func (l *Lexer) error(message string) {
	var buffer bytes.Buffer
	buffer.WriteString(fmt.Sprintf("[ %d:%d ]: error: %s\n",
		l.line, l.current-l.lineStart-1, message))

	lineStr := strconv.Itoa(l.line)
	buffer.WriteString(fmt.Sprintf("      %d |         %s\n", l.line, l.source[l.lineStart:l.current]))
	buffer.WriteString(fmt.Sprintf("      "))
	buffer.WriteString(strings.Repeat(" ", len(lineStr)))
	buffer.WriteString(" |         ")
	buffer.WriteString(fmt.Sprintf("%s^\n", strings.Repeat(" ", l.current-l.lineStart-1)))

	l.Errors = append(l.Errors, fmt.Errorf(buffer.String()))
}

func NewLexer(source string) *Lexer {
	return &Lexer{
		source: source,
		line:   1,
	}
}
