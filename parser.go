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
		return Token{Type: T_EOF}
	}
	return p.tokens[p.pos]
}

func (p *Parser) next() Token {
	tok := p.peek()
	p.pos++
	return tok
}

func (p *Parser) match(t TokenType) bool {
	if p.peek().Type == t {
		p.next()
		return true
	}
	return false
}

func (p *Parser) expect(t TokenType) Token {
	if !p.match(t) {
		panic(fmt.Sprintf("Ожидался токен %v, но найден %v", t, p.peek()))
	}
	return p.tokens[p.pos-1]
}

func (p *Parser) Parse() []Stmt {
	var stmts []Stmt
	for p.peek().Type != T_EOF {
		stmts = append(stmts, p.parseStatement())
	}
	return stmts
}

func (p *Parser) parseStatement() Stmt {
	tok := p.peek()

	if tok.Type == T_DOLLAR {
		p.next() // $
		ident := p.expect(T_IDENT).Value
		p.expect(T_EQ)
		expr := p.parseExpression()
		p.expect(T_SEMI)
		return &AssignStmt{Name: ident, Expr: expr}
	}

	if tok.Type == T_ECHO {
		p.next()
		expr := p.parseExpression()
		p.expect(T_SEMI)
		return &EchoStmt{Expr: expr}
	}

	if tok.Type == T_IF {
		p.next()
		p.expect(T_LPAREN)
		cond := p.parseExpression()
		p.expect(T_RPAREN)
		thenBlock := p.parseBlock()
		var elseBlock []Stmt
		if p.peek().Type == T_ELSE {
			p.next()
			elseBlock = p.parseBlock()
		}
		return &IfStmt{Cond: cond, Then: thenBlock, Else: elseBlock}
	}	

	panic(fmt.Sprintf("Неожиданный токен в начале строки: %v", tok))
}

func (p *Parser) parseBlock() []Stmt {
	p.expect(T_LBRACE)
	var stmts []Stmt
	for p.peek().Type != T_RBRACE && p.peek().Type != T_EOF {
		stmts = append(stmts, p.parseStatement())
	}
	p.expect(T_RBRACE)
	return stmts
}

// --------------------
// Выражения с приоритетом
// --------------------

func (p *Parser) parseExpression() Expr {
	return p.parseOr()
}

func (p *Parser) parseAddSub() Expr {
	left := p.parseMulDiv()
	for {
		op := p.peek()
		if op.Type == T_PLUS || op.Type == T_MINUS {
			p.next()
			right := p.parseMulDiv()
			left = &BinaryExpr{Left: left, Op: op.Type, Right: right}
		} else {
			break
		}
	}
	return left
}

func (p *Parser) parseMulDiv() Expr {
	left := p.parsePrimary()
	for {
		op := p.peek()
		if op.Type == T_STAR || op.Type == T_SLASH {
			p.next()
			right := p.parsePrimary()
			left = &BinaryExpr{Left: left, Op: op.Type, Right: right}
		} else {
			break
		}
	}
	return left
}

func (p *Parser) parsePrimary() Expr {
	tok := p.next()

	switch tok.Type {
	case T_NUMBER:
		val, _ := strconv.Atoi(tok.Value)
		return &NumberLiteral{Value: val}
	case T_STRING:
		return &StringLiteral{Value: tok.Value}
	case T_DOLLAR:
		ident := p.expect(T_IDENT).Value
		return &VarExpr{Name: ident}
	case T_LPAREN:
		expr := p.parseExpression()
		p.expect(T_RPAREN)
		return expr
	default:
		panic(fmt.Sprintf("Ожидалось выражение, найден токен: %v", tok))
	}
}

func (p *Parser) parseOr() Expr {
	left := p.parseAnd()
	for {
		if p.peek().Type == T_OR {
			p.next()
			right := p.parseAnd()
			left = &BinaryExpr{Left: left, Op: T_OR, Right: right}
		} else {
			break
		}
	}
	return left
}

func (p *Parser) parseAnd() Expr {
	left := p.parseComparison()
	for {
		if p.peek().Type == T_AND {
			p.next()
			right := p.parseComparison()
			left = &BinaryExpr{Left: left, Op: T_AND, Right: right}
		} else {
			break
		}
	}
	return left
}

func (p *Parser) parseComparison() Expr {
	left := p.parseAddSub()
	for {
		op := p.peek()
		if op.Type == T_GT || op.Type == T_LT || op.Type == T_EQEQ {
			p.next()
			right := p.parseAddSub()
			left = &BinaryExpr{Left: left, Op: op.Type, Right: right}
		} else {
			break
		}
	}
	return left
}

