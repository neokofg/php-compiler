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
	l.reader.SkipWhitespace()

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
