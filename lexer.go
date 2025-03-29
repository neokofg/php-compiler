package main

import (
	"strings"
	"unicode"
)

type Lexer struct {
	input []rune
	pos   int
}

func NewLexer(input string) *Lexer {
	return &Lexer{input: []rune(input)}
}

func (l *Lexer) next() rune {
	if l.pos >= len(l.input) {
		return 0
	}
	ch := l.input[l.pos]
	l.pos++
	return ch
}

func (l *Lexer) peek() rune {
	if l.pos >= len(l.input) {
		return 0
	}
	return l.input[l.pos]
}

func (l *Lexer) skipWhitespace() {
	for unicode.IsSpace(l.peek()) {
		l.next()
	}
}

func (l *Lexer) readWhile(cond func(rune) bool) string {
	var sb strings.Builder
	for cond(l.peek()) {
		sb.WriteRune(l.next())
	}
	return sb.String()
}

func (l *Lexer) NextToken() Token {
	l.skipWhitespace()

	ch := l.peek()
	if ch == 0 {
		return Token{T_EOF, ""}
	}

	switch ch {
	case '+':
		l.next()
		return Token{T_PLUS, "+"}
	case '-':
		l.next()
		return Token{T_MINUS, "-"}
	case '*':
		l.next()
		return Token{T_STAR, "*"}
	case '/':
		l.next()
		return Token{T_SLASH, "/"}
	case '=':
		l.next()
		if l.peek() == '=' {
			l.next()
			return Token{T_EQEQ, "=="}
		}
		return Token{T_EQ, "="}	
	case ';':
		l.next()
		return Token{T_SEMI, ";"}
	case '$':
		l.next()
		return Token{T_DOLLAR, "$"}
	case '(':
		l.next()
		return Token{T_LPAREN, "("}
	case ')':
		l.next()
		return Token{T_RPAREN, ")"}
	case '"':
		l.next()
		var sb strings.Builder
		for {
			ch := l.peek()
			if ch == '"' || ch == 0 {
				break
			}
			sb.WriteRune(l.next())
		}
		l.next()
		return Token{T_STRING, sb.String()}	
	case '{':
		l.next()
		return Token{T_LBRACE, "{"}
	case '}':
		l.next()
		return Token{T_RBRACE, "}"}	
	case '>':
		l.next()
		return Token{T_GT, ">"}	
	case '<':
		l.next()
		return Token{T_LT, "<"}
	case '&':
		l.next()
		if l.peek() == '&' {
			l.next()
			return Token{T_AND, "&&"}
		} else {
			return Token{T_EOF, ""}
		}
	case '|':
		l.next()
		if l.peek() == '|' {
			l.next()
			return Token{T_OR, "||"}
		} else {
			return Token{T_EOF, ""}
		}
	default:
		if unicode.IsDigit(ch) {
			val := l.readWhile(unicode.IsDigit)
			return Token{T_NUMBER, val}
		}
	
		if unicode.IsLetter(ch) {
			val := l.readWhile(func(r rune) bool {
				return unicode.IsLetter(r) || unicode.IsDigit(r) || r == '_'
			})
			if val == "echo" {
				return Token{T_ECHO, val}
			}
			if val == "if" {
				return Token{T_IF, val}
			}
			if val == "else" {
				return Token{T_ELSE, val}
			}			
			return Token{T_IDENT, val}
		}

		if ch == '"' {
			l.next()
			val := l.readWhile(func(r rune) bool { return r != '"' })
			l.next()
			return Token{T_STRING, val}
		}
		l.next()
		return Token{T_EOF, ""}
	}
}
