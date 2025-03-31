package statement_types

import (
	"github.com/neokofg/php-compiler/internal/ast"
	"github.com/neokofg/php-compiler/internal/parser/expression_parts"
	"github.com/neokofg/php-compiler/internal/parser/parser_contract"
	"github.com/neokofg/php-compiler/internal/token"
)

func ParseIfStatement(p parser_contract.ParserLike) (*ast.IfStmt, error) {
	p.Next()
	_, err := p.Expect(token.T_LPAREN)
	if err != nil {
		return nil, err
	}
	cond, err := expression_parts.ParseExpression(p)
	if err != nil {
		return nil, err
	}
	_, err = p.Expect(token.T_RPAREN)
	if err != nil {
		return nil, err
	}
	thenBlock, err := p.ParseBlock()
	if err != nil {
		return nil, err
	}

	var elseBlock []ast.Stmt
	if p.Peek().Type == token.T_ELSE {
		p.Next()
		elseBlock, err = p.ParseBlock()
		if err != nil {
			return nil, err
		}
	}
	return &ast.IfStmt{Cond: cond, Then: thenBlock, Else: elseBlock}, nil
}
