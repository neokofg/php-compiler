// Licensed under GNU GPL v3. See LICENSE file for details.
package tokenizer

import (
	"github.com/neokofg/php-compiler/internal/lexer/interfaces"
	token2 "github.com/neokofg/php-compiler/internal/token"
)

type OperatorTokenizer struct{}

func NewOperatorTokenizer() *OperatorTokenizer {
	return &OperatorTokenizer{}
}

func (t *OperatorTokenizer) CanTokenize(r rune) bool {
	switch r {
	case '+', '-', '*', '/', '=', ';', '$', '(', ')', '{', '}', '>', '<', '&', '|':
		return true
	default:
		return false
	}
}

func (t *OperatorTokenizer) Tokenize(reader interfaces.Reader) token2.Token {
	switch reader.Peek() {
	case '+':
		reader.Next()
		return token2.Token{Type: token2.T_PLUS, Value: "+"}
	case '-':
		reader.Next()
		return token2.Token{Type: token2.T_MINUS, Value: "-"}
	case '*':
		reader.Next()
		return token2.Token{Type: token2.T_STAR, Value: "*"}
	case '/':
		reader.Next()
		return token2.Token{Type: token2.T_SLASH, Value: "/"}
	case '=':
		reader.Next()
		if reader.Peek() == '=' {
			reader.Next()
			return token2.Token{Type: token2.T_EQEQ, Value: "=="}
		}
		return token2.Token{Type: token2.T_EQ, Value: "="}
	case ';':
		reader.Next()
		return token2.Token{Type: token2.T_SEMI, Value: ";"}
	case '$':
		reader.Next()
		return token2.Token{Type: token2.T_DOLLAR, Value: "$"}
	case '(':
		reader.Next()
		return token2.Token{Type: token2.T_LPAREN, Value: "("}
	case ')':
		reader.Next()
		return token2.Token{Type: token2.T_RPAREN, Value: ")"}
	case '{':
		reader.Next()
		return token2.Token{Type: token2.T_LBRACE, Value: "{"}
	case '}':
		reader.Next()
		return token2.Token{Type: token2.T_RBRACE, Value: "}"}
	case '>':
		reader.Next()
		return token2.Token{Type: token2.T_GT, Value: ">"}
	case '<':
		reader.Next()
		return token2.Token{Type: token2.T_LT, Value: "<"}
	case '&':
		reader.Next()
		if reader.Peek() == '&' {
			reader.Next()
			return token2.Token{Type: token2.T_AND, Value: "&&"}
		}
		return token2.Token{Type: token2.T_ILLEGAL, Value: "&"}
	case '|':
		reader.Next()
		if reader.Peek() == '|' {
			reader.Next()
			return token2.Token{Type: token2.T_OR, Value: "||"}
		}
		return token2.Token{Type: token2.T_ILLEGAL, Value: "|"}
	default:
		return token2.Token{Type: token2.T_ILLEGAL, Value: ""}
	}
}
