package statement_types

import (
	"github.com/neokofg/php-compiler/internal/ast"
	"github.com/neokofg/php-compiler/internal/parser/expression_parts"
	"github.com/neokofg/php-compiler/internal/parser/parser_contract"
	"github.com/neokofg/php-compiler/internal/token"
)

func ParseEchoStatement(p parser_contract.ParserLike) (*ast.EchoStmt, error) {
	p.Next()
	expr, err := expression_parts.ParseExpression(p)
	if err != nil {
		return nil, err
	}
	_, err = p.Expect(token.T_SEMI)
	if err != nil {
		return nil, err
	}
	return &ast.EchoStmt{Expr: expr}, nil
}
