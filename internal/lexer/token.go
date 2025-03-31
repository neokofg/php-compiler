package lexer

import (
	"fmt"
	"github.com/neokofg/php-compiler/internal/lexer/lexing"
	"github.com/neokofg/php-compiler/internal/token"
	"unicode"
)

func (l *Lexer) NextToken() token.Token {
	l.skipWhitespace()

	ch := l.Peek()
	if ch == 0 {
		return token.Token{Type: token.T_EOF, Value: ""}
	}

	if tok := lexing.ReadOperator(l); tok.Type != token.T_ILLEGAL {
		return tok
	}

	switch {
	case unicode.IsDigit(ch):
		return lexing.ReadNumber(l)
	case unicode.IsLetter(ch) || ch == '_':
		return lexing.ReadIdentOrKeyword(l)
	case ch == '"':
		return lexing.ReadString(l)
	default:
		charStr := string(ch)
		l.Next()
		return token.Token{token.T_ILLEGAL, fmt.Sprintf("Undefined symbol: '%s'", charStr)}
	}
}
