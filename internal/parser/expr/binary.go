// Licensed under GNU GPL v3. See LICENSE file for details.
package expr

import (
	"github.com/neokofg/php-compiler/internal/ast"
	"github.com/neokofg/php-compiler/internal/parser/interfaces"
	"github.com/neokofg/php-compiler/internal/token"
)

type BinaryExpressionParser struct {
	context     interfaces.TokenReader
	leftParser  interfaces.ExpressionParser
	validTokens []token.TokenType
}

func NewBinaryExpressionParser(context interfaces.TokenReader, leftParser interfaces.ExpressionParser, validTokens []token.TokenType) *BinaryExpressionParser {
	return &BinaryExpressionParser{
		context:     context,
		leftParser:  leftParser,
		validTokens: validTokens,
	}
}

func (p *BinaryExpressionParser) IsValidOperator(tokenType token.TokenType) bool {
	for _, t := range p.validTokens {
		if t == tokenType {
			return true
		}
	}
	return false
}

func (p *BinaryExpressionParser) Parse() (ast.Expr, error) {
	left, err := p.leftParser.ParseExpression()
	if err != nil {
		return nil, err
	}

	for p.IsValidOperator(p.context.Peek().Type) {
		opTok := p.context.Next()
		right, err := p.leftParser.ParseExpression()
		if err != nil {
			return nil, err
		}
		left = &ast.BinaryExpr{Left: left, Op: opTok.Type, Right: right}
	}

	return left, nil
}
