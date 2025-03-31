package lexing

import (
	"github.com/neokofg/php-compiler/internal/lexer/lexer_contract"
	"github.com/neokofg/php-compiler/internal/token"
	"strings"
)

func ReadString(l lexer_contract.LexerLike) token.Token {
	l.Next() // skip opening quote
	var sb strings.Builder
	startPos := l.GetPos()

	for {
		ch := l.Peek()
		if ch == '"' {
			l.Next()
			break
		}
		if ch == 0 {
			l.SetPos(startPos)
			return token.Token{token.T_ILLEGAL, "Unexpected end of string"}
		}
		if ch == '\\' {
			l.Next()
			switch l.Peek() {
			case 'n':
				sb.WriteRune('\n')
				l.Next()
			case 't':
				sb.WriteRune('\t')
				l.Next()
			case 'r':
				sb.WriteRune('\r')
				l.Next()
			case '"':
				sb.WriteRune('"')
				l.Next()
			case '\\':
				sb.WriteRune('\\')
				l.Next()
			default:
				sb.WriteRune('\\')
				sb.WriteRune(l.Next())
			}
		} else {
			sb.WriteRune(l.Next())
		}
	}

	return token.Token{token.T_STRING, sb.String()}
}
