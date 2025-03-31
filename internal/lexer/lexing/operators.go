package lexing

import (
	"github.com/neokofg/php-compiler/internal/lexer/lexer_contract"
	"github.com/neokofg/php-compiler/internal/token"
)

func ReadOperator(l lexer_contract.LexerLike) token.Token {
	switch l.Peek() {
	case '+':
		l.Next()
		return token.Token{token.T_PLUS, "+"}
	case '-':
		l.Next()
		return token.Token{token.T_MINUS, "-"}
	case '*':
		l.Next()
		return token.Token{token.T_STAR, "*"}
	case '/':
		l.Next()
		return token.Token{token.T_SLASH, "/"} // потом сюда же можно добавить // и /*
	case '=':
		l.Next()
		if l.Peek() == '=' {
			l.Next()
			return token.Token{token.T_EQEQ, "=="}
		}
		return token.Token{token.T_EQ, "="}
	case ';':
		l.Next()
		return token.Token{token.T_SEMI, ";"}
	case '$':
		l.Next()
		return token.Token{token.T_DOLLAR, "$"}
	case '(':
		l.Next()
		return token.Token{token.T_LPAREN, "("}
	case ')':
		l.Next()
		return token.Token{token.T_RPAREN, ")"}
	case '{':
		l.Next()
		return token.Token{token.T_LBRACE, "{"}
	case '}':
		l.Next()
		return token.Token{token.T_RBRACE, "}"}
	case '>':
		l.Next()
		return token.Token{token.T_GT, ">"}
	case '<':
		l.Next()
		return token.Token{token.T_LT, "<"}
	case '&':
		l.Next()
		if l.Peek() == '&' {
			l.Next()
			return token.Token{token.T_AND, "&&"}
		}
		return token.Token{token.T_ILLEGAL, "&"}
	case '|':
		l.Next()
		if l.Peek() == '|' {
			l.Next()
			return token.Token{token.T_OR, "||"}
		}
		return token.Token{token.T_ILLEGAL, "|"}
	default:
		return token.Token{token.T_ILLEGAL, ""}
	}
}
