// Licensed under GNU GPL v3. See LICENSE file for details.
package stmt

import (
	"github.com/neokofg/php-compiler/internal/ast"
	"github.com/neokofg/php-compiler/internal/parser/interfaces"
	"github.com/neokofg/php-compiler/internal/token"
)

type EchoParser struct {
	context    interfaces.TokenReader
	exprParser interfaces.ExpressionParser
}

func NewEchoParser(context interfaces.TokenReader, exprParser interfaces.ExpressionParser) *EchoParser {
	return &EchoParser{
		context:    context,
		exprParser: exprParser,
	}
}

func (p *EchoParser) Parse() (ast.Stmt, error) {
	p.context.Next()

	expr, err := p.exprParser.ParseExpression()
	if err != nil {
		return nil, err
	}

	_, err = p.context.Expect(token.T_SEMI)
	if err != nil {
		return nil, err
	}

	return &ast.EchoStmt{Expr: expr}, nil
}
