// Licensed under GNU GPL v3. See LICENSE file for details.
package reader

import (
	"unicode"
)

type SourceReader struct {
	input []rune
	pos   int
}

func NewSourceReader(input string) *SourceReader {
	return &SourceReader{
		input: []rune(input),
		pos:   0,
	}
}

func (r *SourceReader) GetPos() int {
	return r.pos
}

func (r *SourceReader) SetPos(pos int) {
	r.pos = pos
}

func (r *SourceReader) Next() rune {
	if r.pos >= len(r.input) {
		return 0
	}
	ch := r.input[r.pos]
	r.pos++
	return ch
}

func (r *SourceReader) Peek() rune {
	if r.pos >= len(r.input) {
		return 0
	}
	return r.input[r.pos]
}

func (r *SourceReader) PeekNext() rune {
	if r.pos+1 >= len(r.input) {
		return 0
	}
	return r.input[r.pos+1]
}

func (r *SourceReader) ReadWhile(cond func(rune) bool) string {
	start := r.pos
	for cond(r.Peek()) {
		r.Next()
	}
	return string(r.input[start:r.pos])
}

func (r *SourceReader) skipWhitespace() {
	for unicode.IsSpace(r.Peek()) {
		r.Next()
	}
}

func (r *SourceReader) skipSingleLineComment() {
	for {
		ch := r.Peek()
		if ch == 0 || ch == '\n' {
			return
		}
		r.Next()
	}
}

func (r *SourceReader) skipMultiLineComment() {
	for {
		ch := r.Peek()
		if ch == 0 {
			return
		}
		if ch == '*' && r.PeekNext() == '/' {
			r.Next() // Skip *
			r.Next() // Skip /
			return
		}
		r.Next()
	}
}

func (r *SourceReader) SkipWhitespaceAndComments() {
	for {
		r.skipWhitespace()

		ch := r.Peek()
		if ch == '/' {
			nextCh := r.PeekNext()
			if nextCh == '/' {
				r.Next()
				r.Next()
				r.skipSingleLineComment()
				continue
			} else if nextCh == '*' {
				r.Next()
				r.Next()
				r.skipMultiLineComment()
				continue
			}
		} else if ch == '#' {
			r.Next()
			r.skipSingleLineComment()
			continue
		}

		break
	}
}
