// Licensed under GNU GPL v3. See LICENSE file for details.
package expr

import (
	"github.com/neokofg/php-compiler/internal/ast"
	"github.com/neokofg/php-compiler/internal/parser/interfaces"
	"github.com/neokofg/php-compiler/internal/token"
)

type OrParser struct {
	context   interfaces.TokenReader
	andParser *AndParser
}

func NewOrParser(context interfaces.TokenReader, andParser *AndParser) *OrParser {
	return &OrParser{
		context:   context,
		andParser: andParser,
	}
}

func (p *OrParser) Parse() (ast.Expr, error) {
	left, err := p.andParser.Parse()
	if err != nil {
		return nil, err
	}

	for p.context.Peek().Type == token.T_OR {
		opTok := p.context.Next()
		right, err := p.andParser.Parse()
		if err != nil {
			return nil, err
		}
		left = &ast.BinaryExpr{Left: left, Op: opTok.Type, Right: right}
	}

	return left, nil
}
