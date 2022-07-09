package lexing

import "fmt"

type TokenType int

const (
	Eof TokenType = iota

	LeftParen
	RightParen
	LeftBrace
	RightBrace

	Comma
	Dot
	Minus
	Plus
	Semicolon
	Slash
	Star

	Bang
	BangEqual
	Equal
	EqualEqual
	Greater
	GreaterEqual
	Less
	LessEqual

	Identifier
	String
	Number

	And
	Class
	Else
	False
	Fun
	For
	If
	Nil
	Or
	Print
	Return
	Super
	This
	True
	Var
	While
)

type Token struct {
	tokenType TokenType
	lexeme    string
	literal   interface{}
	line      int
}

func (t *Token) String() string {
	return fmt.Sprintf("token %v, lexeme = %s, literal = %v, line = %d",
		t.tokenType, t.lexeme, t.literal, t.line)
}

func NewToken(tokenType TokenType, lexeme string, literal interface{}, line int) Token {
	return Token{
		tokenType: tokenType,
		lexeme:    lexeme,
		literal:   literal,
		line:      line,
	}
}
