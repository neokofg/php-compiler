package statement_types

import (
	"github.com/neokofg/php-compiler/internal/ast"
	"github.com/neokofg/php-compiler/internal/parser/expression_parts"
	"github.com/neokofg/php-compiler/internal/parser/parser_contract"
	"github.com/neokofg/php-compiler/internal/token"
)

func ParseAssignStatement(p parser_contract.ParserLike) (ast.Stmt, error) {
	p.Next() // $
	identToken, err := p.Expect(token.T_IDENT)
	if err != nil {
		return nil, err
	}
	_, err = p.Expect(token.T_EQ)
	if err != nil {
		return nil, err
	}
	expr, err := expression_parts.ParseExpression(p)
	if err != nil {
		return nil, err
	}
	_, err = p.Expect(token.T_SEMI)
	if err != nil {
		return nil, err
	}
	return &ast.AssignStmt{Name: identToken.Value, Expr: expr}, nil
}
