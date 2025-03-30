package parser

import (
	"fmt"
	"strconv"
	"github.com/neokofg/php-compiler/internal/token"
	"github.com/neokofg/php-compiler/internal/ast"
)

type Parser struct {
	tokens []token.Token
	pos    int
}

func NewParser(tokens []token.Token) *Parser {
	return &Parser{tokens: tokens}
}

func (p *Parser) peek() token.Token {
	if p.pos >= len(p.tokens) {
		return token.Token{Type: token.T_EOF, Value: ""}
	}
	return p.tokens[p.pos]
}

func (p *Parser) next() token.Token {
	tok := p.peek()
	if tok.Type != token.T_EOF {
		p.pos++
	}
	return tok
}

func (p *Parser) match(t token.TokenType) bool {
	if p.peek().Type == t {
		p.next()
		return true
	}
	return false
}

func (p *Parser) expect(t token.TokenType) (token.Token, error) {
	if p.peek().Type == t {
		return p.next(), nil
	}
	return token.Token{}, fmt.Errorf("Position %d: expected token %v, but found %v (%q)", p.pos, t, p.peek().Type, p.peek().Value)
}

func (p *Parser) Parse() ([]ast.Stmt, error) {
	var stmts []ast.Stmt
	for p.peek().Type != token.T_EOF {
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

func (p *Parser) parseStatement() (ast.Stmt, error) {
	tok := p.peek()

	switch tok.Type {
	case token.T_DOLLAR:
		p.next() // $
		identToken, err := p.expect(token.T_IDENT)
		if err != nil {
			return nil, err
		}
		_, err = p.expect(token.T_EQ)
		if err != nil {
			return nil, err
		}
		expr, err := p.parseExpression()
		if err != nil {
			return nil, err
		}
		_, err = p.expect(token.T_SEMI)
		if err != nil {
			return nil, err
		}
		return &ast.AssignStmt{Name: identToken.Value, Expr: expr}, nil

	case token.T_ECHO:
		p.next()
		expr, err := p.parseExpression()
		if err != nil {
			return nil, err
		}
		_, err = p.expect(token.T_SEMI)
		if err != nil {
			return nil, err
		}
		return &ast.EchoStmt{Expr: expr}, nil

	case token.T_IF:
		p.next()
		_, err := p.expect(token.T_LPAREN)
		if err != nil {
			return nil, err
		}
		cond, err := p.parseExpression()
		if err != nil {
			return nil, err
		}
		_, err = p.expect(token.T_RPAREN)
		if err != nil {
			return nil, err
		}
		thenBlock, err := p.parseBlock()
		if err != nil {
			return nil, err
		}

		var elseBlock []ast.Stmt
		if p.peek().Type == token.T_ELSE {
			p.next()
			elseBlock, err = p.parseBlock()
			if err != nil {
				return nil, err
			}
		}
		return &ast.IfStmt{Cond: cond, Then: thenBlock, Else: elseBlock}, nil
	case token.T_WHILE:
		p.next()
		_, err := p.expect(token.T_LPAREN)
		if err != nil {
			return nil, err
		}

		condExpr, err := p.parseExpression()
		if err != nil {
			return nil, err
		}

		_, err = p.expect(token.T_RPAREN)
		if err != nil {
			return nil, err
		}

		bodyBlock, err := p.parseBlock()
		if err != nil {
			return nil, err
		}

		return &ast.WhileStmt{Cond: condExpr, Body: bodyBlock}, nil
	case token.T_FOR:
		p.next()
		_, err := p.expect(token.T_LPAREN)
		if err != nil {
			return nil, err
		}
	
		initExpr, err := p.parseOptionalExpression(token.T_SEMI)
		if err != nil {
			return nil, fmt.Errorf("error parsing for-loop initializer: %w", err)
		}
		_, err = p.expect(token.T_SEMI)
		if err != nil {
			return nil, err
		}
	
		condExpr, err := p.parseOptionalExpression(token.T_SEMI)
		if err != nil {
			return nil, fmt.Errorf("error parsing for-loop condition: %w", err)
		}
		_, err = p.expect(token.T_SEMI)
		if err != nil {
			return nil, err
		}
	
		incrExpr, err := p.parseOptionalExpression(token.T_RPAREN)
		if err != nil {
			return nil, fmt.Errorf("error parsing for-loop increment: %w", err)
		}
		_, err = p.expect(token.T_RPAREN)
		if err != nil {
			return nil, err
		}
	
		bodyBlock, err := p.parseBlock()
		if err != nil {
			return nil, err
		}
	
		return &ast.ForStmt{Init: initExpr, Cond: condExpr, Incr: incrExpr, Body: bodyBlock}, nil
	default:
		if tok.Type == token.T_ILLEGAL {
			return nil, fmt.Errorf("Lexer error in position %d: %s", p.pos, tok.Value)
		}
		return nil, fmt.Errorf("Position %d: unexpected token in the start of instruction: %v (%q)", p.pos, tok.Type, tok.Value)
	}
}

func (p *Parser) parseOptionalExpression(terminator token.TokenType) (ast.Expr, error) {
	if p.peek().Type == terminator {
		return nil, nil
	}
	expr, err := p.parseExpression()
	if err != nil {
		return nil, fmt.Errorf("Failed to parse optional expression before %v: %w", terminator, err)
	}
	return expr, nil
}

func (p *Parser) parseBlock() ([]ast.Stmt, error) {
	_, err := p.expect(token.T_LBRACE)
	if err != nil {
		return nil, err
	}
	var stmts []ast.Stmt
	for p.peek().Type != token.T_RBRACE && p.peek().Type != token.T_EOF {
		stmt, err := p.parseStatement()
		if err != nil {
			return nil, err
		}
		stmts = append(stmts, stmt)
	}
	_, err = p.expect(token.T_RBRACE)
	if err != nil {
		return nil, err
	}
	return stmts, nil
}

// --- Expressions ---

func (p *Parser) parseExpression() (ast.Expr, error) {
	return p.parseOr()
}

func (p *Parser) parseOr() (ast.Expr, error) {
	left, err := p.parseAnd()
	if err != nil {
		return nil, err
	}
	for p.peek().Type == token.T_OR {
		opTok := p.next()
		right, err := p.parseAnd()
		if err != nil {
			return nil, err
		}
		left = &ast.BinaryExpr{Left: left, Op: opTok.Type, Right: right}
	}
	return left, nil
}

func (p *Parser) parseAnd() (ast.Expr, error) {
	left, err := p.parseComparison()
	if err != nil {
		return nil, err
	}
	for p.peek().Type == token.T_AND {
		opTok := p.next()
		right, err := p.parseComparison()
		if err != nil {
			return nil, err
		}
		left = &ast.BinaryExpr{Left: left, Op: opTok.Type, Right: right}
	}
	return left, nil
}

func (p *Parser) parseComparison() (ast.Expr, error) {
	left, err := p.parseAddSub()
	if err != nil {
		return nil, err
	}
	for p.peek().Type == token.T_GT || p.peek().Type == token.T_LT || p.peek().Type == token.T_EQEQ {
		opTok := p.next()
		right, err := p.parseAddSub()
		if err != nil {
			return nil, err
		}
		left = &ast.BinaryExpr{Left: left, Op: opTok.Type, Right: right}
	}
	return left, nil
}

func (p *Parser) parseAddSub() (ast.Expr, error) {
	left, err := p.parseMulDiv()
	if err != nil {
		return nil, err
	}
	for p.peek().Type == token.T_PLUS || p.peek().Type == token.T_MINUS {
		opTok := p.next()
		right, err := p.parseMulDiv()
		if err != nil {
			return nil, err
		}
		left = &ast.BinaryExpr{Left: left, Op: opTok.Type, Right: right}
	}
	return left, nil
}

func (p *Parser) parseMulDiv() (ast.Expr, error) {
	left, err := p.parsePrimary()
	if err != nil {
		return nil, err
	}
	for p.peek().Type == token.T_STAR || p.peek().Type == token.T_SLASH {
		opTok := p.next()
		right, err := p.parsePrimary()
		if err != nil {
			return nil, err
		}
		left = &ast.BinaryExpr{Left: left, Op: opTok.Type, Right: right}
	}
	return left, nil
}

func (p *Parser) parsePrimary() (ast.Expr, error) {
	tok := p.peek()

	switch tok.Type {
	case token.T_NUMBER:
		p.next()
		val, err := strconv.Atoi(tok.Value)
		if err != nil {
			return nil, fmt.Errorf("Position %d: wrong number format: %s", p.pos-1, tok.Value)
		}
		return &ast.NumberLiteral{Value: val}, nil
	case token.T_STRING:
		p.next()
		return &ast.StringLiteral{Value: tok.Value}, nil
	case token.T_DOLLAR:
		p.next()
		identToken, err := p.expect(token.T_IDENT)
		if err != nil {
			return nil, err
		}
		return &ast.VarExpr{Name: identToken.Value}, nil
	case token.T_LPAREN:
		p.next()
		expr, err := p.parseExpression()
		if err != nil {
			return nil, err
		}
		_, err = p.expect(token.T_RPAREN)
		if err != nil {
			return nil, err
		}
		return expr, nil
	case token.T_ILLEGAL:
		p.next()
		return nil, fmt.Errorf("Lexer error in position %d: %s", p.pos-1, tok.Value)
	default:
		p.next()
		return nil, fmt.Errorf("Position %d: expected expression (num, string, var, '('), but found token: %v (%q)", p.pos-1, tok.Type, tok.Value)
	}
}