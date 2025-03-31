package parser

import (
	"fmt"
	"github.com/neokofg/php-compiler/internal/ast"
	"github.com/neokofg/php-compiler/internal/parser/expression_parts"
	"github.com/neokofg/php-compiler/internal/parser/statement_types"
	"github.com/neokofg/php-compiler/internal/token"
)

func (p *Parser) parseStatement() (ast.Stmt, error) {
	peekedToken := p.Peek()

	switch peekedToken.Type {
	case token.T_DOLLAR:
		return statement_types.ParseAssignStatement(p)

	case token.T_ECHO:
		return statement_types.ParseEchoStatement(p)

	case token.T_IF:
		return statement_types.ParseIfStatement(p)

	case token.T_WHILE:
		return statement_types.ParseWhileStatement(p)

	case token.T_FOR:
		return statement_types.ParseForStatement(p)

	default:
		if peekedToken.Type == token.T_ILLEGAL {
			return nil, fmt.Errorf("Lexer error in position %d: %s", p.pos, peekedToken.Value)
		}
		return nil, fmt.Errorf("Position %d: unexpected token in the start of instruction: %v (%q)", p.pos, peekedToken.Type, peekedToken.Value)
	}
}

func (p *Parser) ParseOptionalExpression(terminator token.TokenType) (ast.Expr, error) {
	if p.Peek().Type == terminator {
		return nil, nil
	}
	expr, err := expression_parts.ParseExpression(p)
	if err != nil {
		return nil, fmt.Errorf("Failed to parse optional expression before %v: %w", terminator, err)
	}
	return expr, nil
}

func (p *Parser) ParseBlock() ([]ast.Stmt, error) {
	_, err := p.Expect(token.T_LBRACE)
	if err != nil {
		return nil, err
	}
	var stmts []ast.Stmt
	for p.Peek().Type != token.T_RBRACE && p.Peek().Type != token.T_EOF {
		stmt, err := p.parseStatement()
		if err != nil {
			return nil, err
		}
		stmts = append(stmts, stmt)
	}
	_, err = p.Expect(token.T_RBRACE)
	if err != nil {
		return nil, err
	}
	return stmts, nil
}
