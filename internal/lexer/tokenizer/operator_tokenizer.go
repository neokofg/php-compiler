// Licensed under GNU GPL v3. See LICENSE file for details.
package tokenizer

import (
	"github.com/neokofg/php-compiler/internal/lexer/interfaces"
	"github.com/neokofg/php-compiler/internal/token"
)

type OperatorTokenizer struct{}

func NewOperatorTokenizer() *OperatorTokenizer {
	return &OperatorTokenizer{}
}

func (t *OperatorTokenizer) CanTokenize(r rune) bool {
	switch r {
	case '+', '-', '*', '/', '=', ';', '$', '(', ')', '{', '}', '>', '<', '&', '|', '!', '.', '%', '^', '~', ':':
		return true
	default:
		return false
	}
}

func (t *OperatorTokenizer) Tokenize(reader interfaces.Reader) token.Token {
	switch reader.Peek() {
	case '+':
		reader.Next()
		if reader.Peek() == '+' {
			reader.Next()
			return token.Token{Type: token.T_INC, Value: "++"}
		} else if reader.Peek() == '=' {
			reader.Next()
			return token.Token{Type: token.T_PLUS_EQ, Value: "+="}
		}
		return token.Token{Type: token.T_PLUS, Value: "+"}
	case '-':
		reader.Next()
		if reader.Peek() == '-' {
			reader.Next()
			return token.Token{Type: token.T_DEC, Value: "--"}
		} else if reader.Peek() == '=' {
			reader.Next()
			return token.Token{Type: token.T_MINUS_EQ, Value: "-="}
		}
		return token.Token{Type: token.T_MINUS, Value: "-"}
	case '*':
		reader.Next()
		if reader.Peek() == '=' {
			reader.Next()
			return token.Token{Type: token.T_MUL_EQ, Value: "*="}
		}
		return token.Token{Type: token.T_STAR, Value: "*"}
	case '/':
		reader.Next()
		if reader.Peek() == '=' {
			reader.Next()
			return token.Token{Type: token.T_DIV_EQ, Value: "/="}
		}
		return token.Token{Type: token.T_SLASH, Value: "/"}
	case '%':
		reader.Next()
		if reader.Peek() == '=' {
			reader.Next()
			return token.Token{Type: token.T_MOD_EQ, Value: "%="}
		}
		return token.Token{Type: token.T_MOD, Value: "%"}
	case '=':
		reader.Next()
		if reader.Peek() == '=' {
			reader.Next()
			if reader.Peek() == '=' {
				reader.Next()
				return token.Token{Type: token.T_EQEQEQ, Value: "==="}
			}
			return token.Token{Type: token.T_EQEQ, Value: "=="}
		}
		return token.Token{Type: token.T_EQ, Value: "="}
	case ';':
		reader.Next()
		return token.Token{Type: token.T_SEMI, Value: ";"}
	case '$':
		reader.Next()
		return token.Token{Type: token.T_DOLLAR, Value: "$"}
	case '(':
		reader.Next()
		return token.Token{Type: token.T_LPAREN, Value: "("}
	case ')':
		reader.Next()
		return token.Token{Type: token.T_RPAREN, Value: ")"}
	case '{':
		reader.Next()
		return token.Token{Type: token.T_LBRACE, Value: "{"}
	case '}':
		reader.Next()
		return token.Token{Type: token.T_RBRACE, Value: "}"}
	case '>':
		reader.Next()
		if reader.Peek() == '=' {
			reader.Next()
			return token.Token{Type: token.T_GTE, Value: ">="}
		} else if reader.Peek() == '>' {
			reader.Next()
			return token.Token{Type: token.T_RSHIFT, Value: ">>"}
		}
		return token.Token{Type: token.T_GT, Value: ">"}
	case '<':
		reader.Next()
		if reader.Peek() == '=' {
			reader.Next()
			return token.Token{Type: token.T_LTE, Value: "<="}
		} else if reader.Peek() == '<' {
			reader.Next()
			return token.Token{Type: token.T_LSHIFT, Value: "<<"}
		}
		return token.Token{Type: token.T_LT, Value: "<"}
	case '&':
		reader.Next()
		if reader.Peek() == '&' {
			reader.Next()
			return token.Token{Type: token.T_AND, Value: "&&"}
		}
		return token.Token{Type: token.T_BIT_AND, Value: "&"}
	case '|':
		reader.Next()
		if reader.Peek() == '|' {
			reader.Next()
			return token.Token{Type: token.T_OR, Value: "||"}
		}
		return token.Token{Type: token.T_BIT_OR, Value: "|"}
	case '!':
		reader.Next()
		if reader.Peek() == '=' {
			reader.Next()
			if reader.Peek() == '=' {
				reader.Next()
				return token.Token{Type: token.T_NOTEQEQ, Value: "!=="}
			}
			return token.Token{Type: token.T_NOTEQ, Value: "!="}
		}
		return token.Token{Type: token.T_NOT, Value: "!"}
	case '^':
		reader.Next()
		return token.Token{Type: token.T_BIT_XOR, Value: "^"}
	case '~':
		reader.Next()
		return token.Token{Type: token.T_BIT_NOT, Value: "~"}
	case '.':
		reader.Next()
		if reader.Peek() == '=' {
			reader.Next()
			return token.Token{Type: token.T_DOT_EQ, Value: ".="}
		}
		return token.Token{Type: token.T_DOT, Value: "."}
	case ':':
		reader.Next()
		return token.Token{Type: token.T_COLON, Value: ":"}
	default:
		return token.Token{Type: token.T_ILLEGAL, Value: ""}
	}
}
