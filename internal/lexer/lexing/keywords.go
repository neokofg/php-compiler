package lexing

import (
	"github.com/neokofg/php-compiler/internal/lexer/lexer_contract"
	"github.com/neokofg/php-compiler/internal/token"
	"unicode"
)

func ReadIdentOrKeyword(l lexer_contract.LexerLike) token.Token {
	val := l.ReadWhile(func(r rune) bool {
		return unicode.IsLetter(r) || unicode.IsDigit(r) || r == '_'
	})

	switch val {
	case "echo":
		return token.Token{token.T_ECHO, val}
	case "if":
		return token.Token{token.T_IF, val}
	case "else":
		return token.Token{token.T_ELSE, val}
	case "while":
		return token.Token{token.T_WHILE, val}
	case "for":
		return token.Token{token.T_FOR, val}
	default:
		return token.Token{token.T_IDENT, val}
	}
}
