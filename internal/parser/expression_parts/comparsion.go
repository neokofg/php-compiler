package expression_parts

import (
	"github.com/neokofg/php-compiler/internal/ast"
	"github.com/neokofg/php-compiler/internal/parser/parser_contract"
	"github.com/neokofg/php-compiler/internal/token"
)

func ParseComparison(p parser_contract.ParserLike) (ast.Expr, error) {
	left, err := ParseAddSub(p)
	if err != nil {
		return nil, err
	}
	for p.Peek().Type == token.T_GT || p.Peek().Type == token.T_LT || p.Peek().Type == token.T_EQEQ {
		opTok := p.Next()
		right, err := ParseAddSub(p)
		if err != nil {
			return nil, err
		}
		left = &ast.BinaryExpr{Left: left, Op: opTok.Type, Right: right}
	}
	return left, nil
}
