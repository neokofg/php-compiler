// Licensed under GNU GPL v3. See LICENSE file for details.
package expr

import (
	"github.com/neokofg/php-compiler/internal/ast"
	"github.com/neokofg/php-compiler/internal/parser/interfaces"
	"github.com/neokofg/php-compiler/internal/token"
)

type MulDivParser struct {
	context       interfaces.TokenReader
	primaryParser *PrimaryParser
}

func NewMulDivParser(context interfaces.TokenReader, primaryParser *PrimaryParser) *MulDivParser {
	return &MulDivParser{
		context:       context,
		primaryParser: primaryParser,
	}
}

func (p *MulDivParser) Parse() (ast.Expr, error) {
	left, err := p.primaryParser.Parse()
	if err != nil {
		return nil, err
	}

	for p.context.Peek().Type == token.T_STAR || p.context.Peek().Type == token.T_SLASH {
		opTok := p.context.Next()
		right, err := p.primaryParser.Parse()
		if err != nil {
			return nil, err
		}
		left = &ast.BinaryExpr{Left: left, Op: opTok.Type, Right: right}
	}

	return left, nil
}
