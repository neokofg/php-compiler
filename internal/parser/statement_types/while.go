package statement_types

import (
	"github.com/neokofg/php-compiler/internal/ast"
	"github.com/neokofg/php-compiler/internal/parser/expression_parts"
	"github.com/neokofg/php-compiler/internal/parser/parser_contract"
	"github.com/neokofg/php-compiler/internal/token"
)

func ParseWhileStatement(p parser_contract.ParserLike) (*ast.WhileStmt, error) {
	p.Next()
	_, err := p.Expect(token.T_LPAREN)
	if err != nil {
		return nil, err
	}

	condExpr, err := expression_parts.ParseExpression(p)
	if err != nil {
		return nil, err
	}

	_, err = p.Expect(token.T_RPAREN)
	if err != nil {
		return nil, err
	}

	bodyBlock, err := p.ParseBlock()
	if err != nil {
		return nil, err
	}

	return &ast.WhileStmt{Cond: condExpr, Body: bodyBlock}, nil
}
