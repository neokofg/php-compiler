// Licensed under GNU GPL v3. See LICENSE file for details.
package expr

import (
	"github.com/neokofg/php-compiler/internal/ast"
	"github.com/neokofg/php-compiler/internal/parser/interfaces"
)

type Expression struct {
	orParser *OrParser
}

func NewExpression(context interfaces.TokenReader) *Expression {
	primaryParser := NewPrimaryParser(context)
	mulDivParser := NewMulDivParser(context, primaryParser)
	addSubParser := NewAddSubParser(context, mulDivParser)
	comparisonParser := NewComparisonParser(context, addSubParser)
	andParser := NewAndParser(context, comparisonParser)
	orParser := NewOrParser(context, andParser)

	return &Expression{
		orParser: orParser,
	}
}

func (e *Expression) ParseExpression() (ast.Expr, error) {
	return e.orParser.Parse()
}
