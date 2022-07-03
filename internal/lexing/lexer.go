package lexing

type Lexer struct {
	tokens []Token
	source string

	start   int
	current int
	line    int
}

func (l *Lexer) ScanTokens() []Token {
	return l.tokens
}

func NewLexer(source string) *Lexer {
	return &Lexer{
		source: source,
		line:   1,
	}
}
