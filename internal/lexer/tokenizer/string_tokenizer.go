// Licensed under GNU GPL v3. See LICENSE file for details.
package tokenizer

import (
	"github.com/neokofg/php-compiler/internal/lexer/interfaces"
	token2 "github.com/neokofg/php-compiler/internal/token"
	"strings"
)

type StringTokenizer struct{}

func NewStringTokenizer() *StringTokenizer {
	return &StringTokenizer{}
}

func (t *StringTokenizer) CanTokenize(r rune) bool {
	return r == '"'
}

func (t *StringTokenizer) Tokenize(reader interfaces.Reader) token2.Token {
	reader.Next()
	var sb strings.Builder
	startPos := reader.GetPos()

	for {
		ch := reader.Peek()
		if ch == '"' {
			reader.Next()
			break
		}
		if ch == 0 {
			reader.SetPos(startPos)
			return token2.Token{Type: token2.T_ILLEGAL, Value: "Unexpected end of string"}
		}
		if ch == '\\' {
			reader.Next()
			switch reader.Peek() {
			case 'n':
				sb.WriteRune('\n')
				reader.Next()
			case 't':
				sb.WriteRune('\t')
				reader.Next()
			case 'r':
				sb.WriteRune('\r')
				reader.Next()
			case '"':
				sb.WriteRune('"')
				reader.Next()
			case '\\':
				sb.WriteRune('\\')
				reader.Next()
			default:
				sb.WriteRune('\\')
				sb.WriteRune(reader.Next())
			}
		} else {
			sb.WriteRune(reader.Next())
		}
	}

	return token2.Token{Type: token2.T_STRING, Value: sb.String()}
}
