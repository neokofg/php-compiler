package interfaces

import (
	"github.com/neokofg/php-compiler/internal/token"
)

type Reader interface {
	Next() rune
	Peek() rune
	PeekNext() rune
	ReadWhile(cond func(rune) bool) string
	GetPos() int
	SetPos(pos int)
}

type Tokenizer interface {
	CanTokenize(r rune) bool
	Tokenize(reader Reader) token.Token
}

type Lexer interface {
	NextToken() token.Token
}
