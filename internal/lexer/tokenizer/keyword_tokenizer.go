package tokenizer

import (
	"github.com/neokofg/php-compiler/internal/lexer/interfaces"
	token2 "github.com/neokofg/php-compiler/internal/token"
	"unicode"
)

type KeywordTokenizer struct{}

func NewKeywordTokenizer() *KeywordTokenizer {
	return &KeywordTokenizer{}
}

func (t *KeywordTokenizer) CanTokenize(r rune) bool {
	return unicode.IsLetter(r) || r == '_'
}

func (t *KeywordTokenizer) Tokenize(reader interfaces.Reader) token2.Token {
	val := reader.ReadWhile(func(r rune) bool {
		return unicode.IsLetter(r) || unicode.IsDigit(r) || r == '_'
	})

	switch val {
	case "echo":
		return token2.Token{Type: token2.T_ECHO, Value: val}
	case "if":
		return token2.Token{Type: token2.T_IF, Value: val}
	case "else":
		return token2.Token{Type: token2.T_ELSE, Value: val}
	case "while":
		return token2.Token{Type: token2.T_WHILE, Value: val}
	case "for":
		return token2.Token{Type: token2.T_FOR, Value: val}
	default:
		return token2.Token{Type: token2.T_IDENT, Value: val}
	}
}
