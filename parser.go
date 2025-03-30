package main

import (
	"fmt"
	"strconv"
)

type Parser struct {
	tokens []Token
	pos    int
}

func NewParser(tokens []Token) *Parser {
	return &Parser{tokens: tokens}
}

func (p *Parser) peek() Token {
	if p.pos >= len(p.tokens) {
		return Token{Type: T_EOF, Value: ""}
	}
	return p.tokens[p.pos]
}

func (p *Parser) next() Token {
	tok := p.peek()
	if tok.Type != T_EOF {
		p.pos++
	}
	return tok
}

func (p *Parser) match(t TokenType) bool {
	if p.peek().Type == t {
		p.next()
		return true
	}
	return false
}

func (p *Parser) expect(t TokenType) (Token, error) {
	if p.peek().Type == t {
		return p.next(), nil
	}
	return Token{}, fmt.Errorf("Position %d: expected token %v, but found %v (%q)", p.pos, t, p.peek().Type, p.peek().Value)
}

func (p *Parser) Parse() ([]Stmt, error) {
	var stmts []Stmt
	for p.peek().Type != T_EOF {
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

func (p *Parser) parseStatement() (Stmt, error) {
	tok := p.peek()

	switch tok.Type {
	case T_DOLLAR:
		p.next() // $
		identToken, err := p.expect(T_IDENT)
		if err != nil {
			return nil, err
		}
		_, err = p.expect(T_EQ)
		if err != nil {
			return nil, err
		}
		expr, err := p.parseExpression()
		if err != nil {
			return nil, err
		}
		_, err = p.expect(T_SEMI)
		if err != nil {
			return nil, err
		}
		return &AssignStmt{Name: identToken.Value, Expr: expr}, nil

	case T_ECHO:
		p.next()
		expr, err := p.parseExpression()
		if err != nil {
			return nil, err
		}
		_, err = p.expect(T_SEMI)
		if err != nil {
			return nil, err
		}
		return &EchoStmt{Expr: expr}, nil

	case T_IF:
		p.next()
		_, err := p.expect(T_LPAREN)
		if err != nil {
			return nil, err
		}
		cond, err := p.parseExpression()
		if err != nil {
			return nil, err
		}
		_, err = p.expect(T_RPAREN)
		if err != nil {
			return nil, err
		}
		thenBlock, err := p.parseBlock()
		if err != nil {
			return nil, err
		}

		var elseBlock []Stmt
		if p.peek().Type == T_ELSE {
			p.next()
			elseBlock, err = p.parseBlock()
			if err != nil {
				return nil, err
			}
		}
		return &IfStmt{Cond: cond, Then: thenBlock, Else: elseBlock}, nil

	default:
		if tok.Type == T_ILLEGAL {
			return nil, fmt.Errorf("Lexer error in position %d: %s", p.pos, tok.Value)
		}
		return nil, fmt.Errorf("Position %d: unexpected token in the start of instruction: %v (%q)", p.pos, tok.Type, tok.Value)
	}
}

func (p *Parser) parseBlock() ([]Stmt, error) {
	_, err := p.expect(T_LBRACE)
	if err != nil {
		return nil, err
	}
	var stmts []Stmt
	for p.peek().Type != T_RBRACE && p.peek().Type != T_EOF {
		stmt, err := p.parseStatement()
		if err != nil {
			return nil, err
		}
		stmts = append(stmts, stmt)
	}
	_, err = p.expect(T_RBRACE)
	if err != nil {
		return nil, err
	}
	return stmts, nil
}

// --- Expressions ---

func (p *Parser) parseExpression() (Expr, error) {
	return p.parseOr()
}

func (p *Parser) parseOr() (Expr, error) {
	left, err := p.parseAnd()
	if err != nil {
		return nil, err
	}
	for p.peek().Type == T_OR {
		opTok := p.next()
		right, err := p.parseAnd()
		if err != nil {
			return nil, err
		}
		left = &BinaryExpr{Left: left, Op: opTok.Type, Right: right}
	}
	return left, nil
}

func (p *Parser) parseAnd() (Expr, error) {
	left, err := p.parseComparison()
	if err != nil {
		return nil, err
	}
	for p.peek().Type == T_AND {
		opTok := p.next()
		right, err := p.parseComparison()
		if err != nil {
			return nil, err
		}
		left = &BinaryExpr{Left: left, Op: opTok.Type, Right: right}
	}
	return left, nil
}

func (p *Parser) parseComparison() (Expr, error) {
	left, err := p.parseAddSub()
	if err != nil {
		return nil, err
	}
	for p.peek().Type == T_GT || p.peek().Type == T_LT || p.peek().Type == T_EQEQ {
		opTok := p.next()
		right, err := p.parseAddSub()
		if err != nil {
			return nil, err
		}
		left = &BinaryExpr{Left: left, Op: opTok.Type, Right: right}
	}
	return left, nil
}

func (p *Parser) parseAddSub() (Expr, error) {
	left, err := p.parseMulDiv()
	if err != nil {
		return nil, err
	}
	for p.peek().Type == T_PLUS || p.peek().Type == T_MINUS {
		opTok := p.next()
		right, err := p.parseMulDiv()
		if err != nil {
			return nil, err
		}
		left = &BinaryExpr{Left: left, Op: opTok.Type, Right: right}
	}
	return left, nil
}

func (p *Parser) parseMulDiv() (Expr, error) {
	left, err := p.parsePrimary()
	if err != nil {
		return nil, err
	}
	for p.peek().Type == T_STAR || p.peek().Type == T_SLASH {
		opTok := p.next()
		right, err := p.parsePrimary()
		if err != nil {
			return nil, err
		}
		left = &BinaryExpr{Left: left, Op: opTok.Type, Right: right}
	}
	return left, nil
}

func (p *Parser) parsePrimary() (Expr, error) {
	tok := p.peek()

	switch tok.Type {
	case T_NUMBER:
		p.next()
		val, err := strconv.Atoi(tok.Value)
		if err != nil {
			return nil, fmt.Errorf("Position %d: wrong number format: %s", p.pos-1, tok.Value)
		}
		return &NumberLiteral{Value: val}, nil
	case T_STRING:
		p.next()
		return &StringLiteral{Value: tok.Value}, nil
	case T_DOLLAR:
		p.next()
		identToken, err := p.expect(T_IDENT)
		if err != nil {
			return nil, err
		}
		return &VarExpr{Name: identToken.Value}, nil
	case T_LPAREN:
		p.next()
		expr, err := p.parseExpression()
		if err != nil {
			return nil, err
		}
		_, err = p.expect(T_RPAREN)
		if err != nil {
			return nil, err
		}
		return expr, nil
	case T_ILLEGAL:
		p.next()
		return nil, fmt.Errorf("Lexer error in position %d: %s", p.pos-1, tok.Value)
	default:
		p.next()
		return nil, fmt.Errorf("Position %d: expected expression (num, string, var, '('), but found token: %v (%q)", p.pos-1, tok.Type, tok.Value)
	}
}