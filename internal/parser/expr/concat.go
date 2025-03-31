// Licensed under GNU GPL v3. See LICENSE file for details.
package expr

import (
	"github.com/neokofg/php-compiler/internal/ast"
	"github.com/neokofg/php-compiler/internal/parser/interfaces"
	"github.com/neokofg/php-compiler/internal/token"
)

type ConcatParser struct {
	context      interfaces.TokenReader
	addSubParser *AddSubParser
}

func NewConcatParser(context interfaces.TokenReader, addSubParser *AddSubParser) *ConcatParser {
	return &ConcatParser{
		context:      context,
		addSubParser: addSubParser,
	}
}

func (p *ConcatParser) Parse() (ast.Expr, error) {
	left, err := p.addSubParser.Parse()
	if err != nil {
		return nil, err
	}

	for p.context.Peek().Type == token.T_DOT {
		opTok := p.context.Next()
		right, err := p.addSubParser.Parse()
		if err != nil {
			return nil, err
		}
		left = &ast.BinaryExpr{Left: left, Op: opTok.Type, Right: right}
	}

	return left, nil
}
