package lexing

import "fmt"

type TokenType int

const (
	Eof TokenType = iota

	LeftParen
	RightParen
	LeftBrace
	RightBrace
	LeftBracket
	RightBracket

	Comma
	Dot
	Minus
	Plus
	Semicolon
	Slash
	Star
	Question
	Colon

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
	Break
	Continue
)

type Token struct {
	TokenType TokenType
	Lexeme    string
	Literal   interface{}
	Line      int
	Position  int
}

func (t *Token) String() string {
	return fmt.Sprintf("token %v, Lexeme = %s, Literal = %v, Line = %d",
		t.TokenType, t.Lexeme, t.Literal, t.Line)
}

func NewToken(tokenType TokenType, lexeme string, literal interface{}, line int, position int) Token {
	return Token{
		TokenType: tokenType,
		Lexeme:    lexeme,
		Literal:   literal,
		Line:      line,
		Position:  position,
	}
}
