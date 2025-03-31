package expression_parts

import (
	"github.com/neokofg/php-compiler/internal/ast"
	"github.com/neokofg/php-compiler/internal/parser/parser_contract"
	"github.com/neokofg/php-compiler/internal/token"
)

func ParseAnd(p parser_contract.ParserLike) (ast.Expr, error) {
	left, err := ParseComparison(p)
	if err != nil {
		return nil, err
	}
	for p.Peek().Type == token.T_AND {
		opTok := p.Next()
		right, err := ParseComparison(p)
		if err != nil {
			return nil, err
		}
		left = &ast.BinaryExpr{Left: left, Op: opTok.Type, Right: right}
	}
	return left, nil
}
