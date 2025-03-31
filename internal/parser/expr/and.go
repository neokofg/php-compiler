package expr

import (
	"github.com/neokofg/php-compiler/internal/ast"
	"github.com/neokofg/php-compiler/internal/parser/interfaces"
	"github.com/neokofg/php-compiler/internal/token"
)

type AndParser struct {
	context          interfaces.TokenReader
	comparisonParser *ComparisonParser
}

func NewAndParser(context interfaces.TokenReader, comparisonParser *ComparisonParser) *AndParser {
	return &AndParser{
		context:          context,
		comparisonParser: comparisonParser,
	}
}

func (p *AndParser) Parse() (ast.Expr, error) {
	left, err := p.comparisonParser.Parse()
	if err != nil {
		return nil, err
	}

	for p.context.Peek().Type == token.T_AND {
		opTok := p.context.Next()
		right, err := p.comparisonParser.Parse()
		if err != nil {
			return nil, err
		}
		left = &ast.BinaryExpr{Left: left, Op: opTok.Type, Right: right}
	}

	return left, nil
}
