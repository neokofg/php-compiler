package tokenizer

import (
	"github.com/neokofg/php-compiler/internal/lexer/interfaces"
	token2 "github.com/neokofg/php-compiler/internal/token"
	"unicode"
)

type NumberTokenizer struct{}

func NewNumberTokenizer() *NumberTokenizer {
	return &NumberTokenizer{}
}

func (t *NumberTokenizer) CanTokenize(r rune) bool {
	return unicode.IsDigit(r)
}

func (t *NumberTokenizer) Tokenize(reader interfaces.Reader) token2.Token {
	val := reader.ReadWhile(unicode.IsDigit)
	return token2.Token{Type: token2.T_NUMBER, Value: val}
}
