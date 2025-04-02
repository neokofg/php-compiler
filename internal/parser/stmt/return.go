// Licensed under GNU GPL v3. See LICENSE file for details.
package stmt

import (
	"github.com/neokofg/php-compiler/internal/ast"
	"github.com/neokofg/php-compiler/internal/parser/interfaces"
	"github.com/neokofg/php-compiler/internal/token"
)

type ReturnParser struct {
	context    interfaces.TokenReader
	exprParser interfaces.ExpressionParser
}

func NewReturnParser(context interfaces.TokenReader, exprParser interfaces.ExpressionParser) *ReturnParser {
	return &ReturnParser{
		context:    context,
		exprParser: exprParser,
	}
}

func (p *ReturnParser) Parse() (ast.Stmt, error) {
	p.context.Next() // return keyword

	var expr ast.Expr
	var err error

	if p.context.Peek().Type != token.T_SEMI {
		expr, err = p.exprParser.ParseExpression()
		if err != nil {
			return nil, err
		}
	}

	_, err = p.context.Expect(token.T_SEMI)
	if err != nil {
		return nil, err
	}

	return &ast.ReturnStmt{Expr: expr}, nil
}
