// PHP Compiler - compiles php code to IR and then running it on PHPC VM
// Copyright (C) 2025  Andrey Vasilev (neokofg)
//
// This program is free software: you can redistribute it and/or modify
// it under the terms of the GNU General Public License as published by
// the Free Software Foundation, either version 3 of the License, or
// (at your option) any later version.
//
// This program is distributed in the hope that it will be useful,
// but WITHOUT ANY WARRANTY; without even the implied warranty of
// MERCHANTABILITY or FITNESS FOR A PARTICULAR PURPOSE.  See the
// GNU General Public License for more details.
//
// You should have received a copy of the GNU General Public License
// along with this program.  If not, see <https://www.gnu.org/licenses/>.
package lexer

import (
	"fmt"
	"github.com/neokofg/php-compiler/internal/lexer/reader"
	"github.com/neokofg/php-compiler/internal/lexer/tokenizer"
	token2 "github.com/neokofg/php-compiler/internal/token"
)

type Lexer struct {
	reader     *reader.SourceReader
	tokenizers *tokenizer.TokenizerRegistry
}

func NewLexer(input string) *Lexer {
	sourceReader := reader.NewSourceReader(input)
	tokenizerRegistry := tokenizer.NewTokenizerRegistry()

	tokenizerRegistry.RegisterTokenizer(tokenizer.NewKeywordTokenizer())
	tokenizerRegistry.RegisterTokenizer(tokenizer.NewNumberTokenizer())
	tokenizerRegistry.RegisterTokenizer(tokenizer.NewOperatorTokenizer())
	tokenizerRegistry.RegisterTokenizer(tokenizer.NewStringTokenizer())

	return &Lexer{
		reader:     sourceReader,
		tokenizers: tokenizerRegistry,
	}
}

func (l *Lexer) NextToken() token2.Token {
	l.reader.SkipWhitespaceAndComments()

	ch := l.reader.Peek()
	if ch == 0 {
		return token2.Token{Type: token2.T_EOF, Value: ""}
	}

	tokenizerInstance := l.tokenizers.FindTokenizer(ch)
	if tokenizerInstance != nil {
		return tokenizerInstance.Tokenize(l.reader)
	}

	charStr := string(ch)
	l.reader.Next()
	return token2.Token{Type: token2.T_ILLEGAL, Value: fmt.Sprintf("Undefined symbol: '%s'", charStr)}
}
