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
		str := l.readWhile(func(r rune) bool { return r != '"' })
		l.next() // consume closing "
		return Token{T_STRING, str}
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
