// Licensed under GNU GPL v3. See LICENSE file for details.
package expr

import (
	"github.com/neokofg/php-compiler/internal/ast"
	"github.com/neokofg/php-compiler/internal/parser/interfaces"
)

type Parser struct {
	context interfaces.TokenReader

	primaryParser    *PrimaryParser
	mulDivParser     *MulDivParser
	addSubParser     *AddSubParser
	concatParser     *ConcatParser
	comparisonParser *ComparisonParser
	andParser        *AndParser
	orParser         *OrParser
}

func NewParser(context interfaces.TokenReader) interfaces.ExpressionParser {
	parser := &Parser{
		context: context,
	}

	parser.primaryParser = NewPrimaryParser(context)
	parser.mulDivParser = NewMulDivParser(context, parser.primaryParser)
	parser.addSubParser = NewAddSubParser(context, parser.mulDivParser)
	parser.concatParser = NewConcatParser(context, parser.addSubParser)
	parser.comparisonParser = NewComparisonParser(context, parser.concatParser)
	parser.andParser = NewAndParser(context, parser.comparisonParser)
	parser.orParser = NewOrParser(context, parser.andParser)

	parser.primaryParser.SetExprParser(parser)

	return parser
}

func (p *Parser) ParseExpression() (ast.Expr, error) {
	return p.orParser.Parse()
}
