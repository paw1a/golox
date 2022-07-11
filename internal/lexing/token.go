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
	TokenType TokenType
	Lexeme    string
	Literal   interface{}
	Line      int
}

func (t *Token) String() string {
	return fmt.Sprintf("token %v, Lexeme = %s, Literal = %v, Line = %d",
		t.TokenType, t.Lexeme, t.Literal, t.Line)
}

func NewToken(tokenType TokenType, lexeme string, literal interface{}, line int) Token {
	return Token{
		TokenType: tokenType,
		Lexeme:    lexeme,
		Literal:   literal,
		Line:      line,
	}
}
