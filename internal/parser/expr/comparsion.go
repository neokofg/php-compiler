// Licensed under GNU GPL v3. See LICENSE file for details.
package expr

import (
	"github.com/neokofg/php-compiler/internal/ast"
	"github.com/neokofg/php-compiler/internal/parser/interfaces"
	"github.com/neokofg/php-compiler/internal/token"
)

type ComparisonParser struct {
	context      interfaces.TokenReader
	addSubParser *AddSubParser
}

func NewComparisonParser(context interfaces.TokenReader, addSubParser *AddSubParser) *ComparisonParser {
	return &ComparisonParser{
		context:      context,
		addSubParser: addSubParser,
	}
}

func (p *ComparisonParser) Parse() (ast.Expr, error) {
	left, err := p.addSubParser.Parse()
	if err != nil {
		return nil, err
	}

	for p.context.Peek().Type == token.T_GT ||
		p.context.Peek().Type == token.T_LT ||
		p.context.Peek().Type == token.T_EQEQ ||
		p.context.Peek().Type == token.T_NOTEQ {

		opTok := p.context.Next()
		right, err := p.addSubParser.Parse()
		if err != nil {
			return nil, err
		}
		left = &ast.BinaryExpr{Left: left, Op: opTok.Type, Right: right}
	}

	return left, nil
}
