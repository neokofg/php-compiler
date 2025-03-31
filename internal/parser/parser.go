package parser

import (
	"fmt"
	"github.com/neokofg/php-compiler/internal/ast"
	"github.com/neokofg/php-compiler/internal/token"
)

type Parser struct {
	tokens []token.Token
	pos    int
}

func (p *Parser) GetPos() int {
	return p.pos
}

func (p *Parser) SetPos(pos int) {
	p.pos = pos
}

func NewParser(tokens []token.Token) *Parser {
	return &Parser{tokens: tokens}
}

func (p *Parser) Peek() token.Token {
	if p.pos >= len(p.tokens) {
		return token.Token{Type: token.T_EOF, Value: ""}
	}
	return p.tokens[p.pos]
}

func (p *Parser) Next() token.Token {
	peekedToken := p.Peek()
	if peekedToken.Type != token.T_EOF {
		p.pos++
	}
	return peekedToken
}

func (p *Parser) match(t token.TokenType) bool {
	if p.Peek().Type == t {
		p.Next()
		return true
	}
	return false
}

func (p *Parser) Expect(t token.TokenType) (token.Token, error) {
	if p.Peek().Type == t {
		return p.Next(), nil
	}
	return token.Token{}, fmt.Errorf("Position %d: expected token %v, but found %v (%q)", p.pos, t, p.Peek().Type, p.Peek().Value)
}

func (p *Parser) Parse() ([]ast.Stmt, error) {
	var stmts []ast.Stmt
	for p.Peek().Type != token.T_EOF {
		stmt, err := p.parseStatement()
		if err != nil {
			return nil, err
		}
		if stmt != nil {
			stmts = append(stmts, stmt)
		}
	}
	return stmts, nil
}
