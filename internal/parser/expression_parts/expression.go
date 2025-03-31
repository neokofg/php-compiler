package expression_parts

import (
	"github.com/neokofg/php-compiler/internal/ast"
	"github.com/neokofg/php-compiler/internal/parser/parser_contract"
)

func ParseExpression(p parser_contract.ParserLike) (ast.Expr, error) {
	return ParseOr(p)
}
