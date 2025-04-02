// Licensed under GNU GPL v3. See LICENSE file for details.
package context

import (
	"fmt"
	"github.com/neokofg/php-compiler/internal/token"
)

type ParserContext struct {
	tokens []token.Token
	pos    int
}

func NewParserContext(tokens []token.Token) *ParserContext {
	return &ParserContext{
		tokens: tokens,
		pos:    0,
	}
}

func (c *ParserContext) GetPos() int {
	return c.pos
}

func (c *ParserContext) SetPos(pos int) {
	c.pos = pos
}

func (c *ParserContext) Peek() token.Token {
	if c.pos >= len(c.tokens) {
		return token.Token{Type: token.T_EOF, Value: ""}
	}
	return c.tokens[c.pos]
}

func (c *ParserContext) PeekNext() token.Token {
	if c.pos >= len(c.tokens) {
		return token.Token{Type: token.T_EOF, Value: ""}
	}
	return c.tokens[c.pos+1]
}

func (c *ParserContext) Next() token.Token {
	peekedToken := c.Peek()
	if peekedToken.Type != token.T_EOF {
		c.pos++
	}
	return peekedToken
}

func (c *ParserContext) Match(t token.TokenType) bool {
	if c.Peek().Type == t {
		c.Next()
		return true
	}
	return false
}

func (c *ParserContext) Expect(t token.TokenType) (token.Token, error) {
	if c.Peek().Type == t {
		return c.Next(), nil
	}
	return token.Token{}, fmt.Errorf("Position %d: expected token %v, but found %v (%q)", c.pos, t, c.Peek().Type, c.Peek().Value)
}
