package lexing

import (
	"github.com/neokofg/php-compiler/internal/lexer/lexer_contract"
	"github.com/neokofg/php-compiler/internal/token"
	"unicode"
)

func ReadNumber(l lexer_contract.LexerLike) token.Token {
	val := l.ReadWhile(unicode.IsDigit)
	return token.Token{token.T_NUMBER, val}
}
