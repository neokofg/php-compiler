// Licensed under GNU GPL v3. See LICENSE file for details.
package tokenizer

import (
	"github.com/neokofg/php-compiler/internal/lexer/interfaces"
	"github.com/neokofg/php-compiler/internal/token"
)

type Tokenizer interface {
	CanTokenize(r rune) bool
	Tokenize(reader interfaces.Reader) token.Token
}

type TokenizerRegistry struct {
	tokenizers []Tokenizer
}

func NewTokenizerRegistry() *TokenizerRegistry {
	return &TokenizerRegistry{
		tokenizers: make([]Tokenizer, 0),
	}
}

func (r *TokenizerRegistry) RegisterTokenizer(tokenizer Tokenizer) {
	r.tokenizers = append(r.tokenizers, tokenizer)
}

func (r *TokenizerRegistry) FindTokenizer(ch rune) Tokenizer {
	for _, tokenizer := range r.tokenizers {
		if tokenizer.CanTokenize(ch) {
			return tokenizer
		}
	}
	return nil
}
