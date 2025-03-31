// Licensed under GNU GPL v3. See LICENSE file for details.
package tokenizer

import (
	"github.com/neokofg/php-compiler/internal/lexer/interfaces"
	"github.com/neokofg/php-compiler/internal/token"
	"strings"
	"unicode"
)

type KeywordTokenizer struct{}

func NewKeywordTokenizer() *KeywordTokenizer {
	return &KeywordTokenizer{}
}

func (t *KeywordTokenizer) CanTokenize(r rune) bool {
	return unicode.IsLetter(r) || r == '_'
}

func (t *KeywordTokenizer) Tokenize(reader interfaces.Reader) token.Token {
	val := reader.ReadWhile(func(r rune) bool {
		return unicode.IsLetter(r) || unicode.IsDigit(r) || r == '_'
	})

	switch strings.ToLower(val) {
	case "echo":
		return token.Token{Type: token.T_ECHO, Value: val}
	case "if":
		return token.Token{Type: token.T_IF, Value: val}
	case "else":
		return token.Token{Type: token.T_ELSE, Value: val}
	case "while":
		return token.Token{Type: token.T_WHILE, Value: val}
	case "for":
		return token.Token{Type: token.T_FOR, Value: val}
	case "true":
		return token.Token{Type: token.T_TRUE, Value: val}
	case "false":
		return token.Token{Type: token.T_FALSE, Value: val}
	default:
		return token.Token{Type: token.T_IDENT, Value: val}
	}
}
