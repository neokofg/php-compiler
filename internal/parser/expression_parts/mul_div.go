package expression_parts

import (
	"github.com/neokofg/php-compiler/internal/ast"
	"github.com/neokofg/php-compiler/internal/parser/parser_contract"
	"github.com/neokofg/php-compiler/internal/token"
)

func ParseMulDiv(p parser_contract.ParserLike) (ast.Expr, error) {
	left, err := ParsePrimary(p)
	if err != nil {
		return nil, err
	}
	for p.Peek().Type == token.T_STAR || p.Peek().Type == token.T_SLASH {
		opTok := p.Next()
		right, err := ParsePrimary(p)
		if err != nil {
			return nil, err
		}
		left = &ast.BinaryExpr{Left: left, Op: opTok.Type, Right: right}
	}
	return left, nil
}
