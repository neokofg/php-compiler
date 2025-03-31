package expr

import (
	"github.com/neokofg/php-compiler/internal/ast"
	"github.com/neokofg/php-compiler/internal/parser/interfaces"
	"github.com/neokofg/php-compiler/internal/token"
)

type AddSubParser struct {
	context      interfaces.TokenReader
	mulDivParser *MulDivParser
}

func NewAddSubParser(context interfaces.TokenReader, mulDivParser *MulDivParser) *AddSubParser {
	return &AddSubParser{
		context:      context,
		mulDivParser: mulDivParser,
	}
}

func (p *AddSubParser) Parse() (ast.Expr, error) {
	left, err := p.mulDivParser.Parse()
	if err != nil {
		return nil, err
	}

	for p.context.Peek().Type == token.T_PLUS || p.context.Peek().Type == token.T_MINUS {
		opTok := p.context.Next()
		right, err := p.mulDivParser.Parse()
		if err != nil {
			return nil, err
		}
		left = &ast.BinaryExpr{Left: left, Op: opTok.Type, Right: right}
	}

	return left, nil
}
