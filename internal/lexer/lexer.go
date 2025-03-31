package lexer

import (
	"unicode"
)

type Lexer struct {
	input []rune
	pos   int
}

func NewLexer(input string) *Lexer {
	return &Lexer{input: []rune(input)}
}

func (l *Lexer) GetPos() int {
	return l.pos
}

func (l *Lexer) SetPos(pos int) {
	l.pos = pos
}

func (l *Lexer) Next() rune {
	if l.pos >= len(l.input) {
		return 0
	}
	ch := l.input[l.pos]
	l.pos++
	return ch
}

func (l *Lexer) Peek() rune {
	if l.pos >= len(l.input) {
		return 0
	}
	return l.input[l.pos]
}

func (l *Lexer) PeekNext() rune {
	if l.pos+1 >= len(l.input) {
		return 0
	}
	return l.input[l.pos+1]
}

func (l *Lexer) skipWhitespace() {
	for unicode.IsSpace(l.Peek()) {
		l.Next()
	}
}

func (l *Lexer) ReadWhile(cond func(rune) bool) string {
	start := l.pos
	for cond(l.Peek()) {
		l.Next()
	}
	return string(l.input[start:l.pos])
}
