// Licensed under GNU GPL v3. See LICENSE file for details.
package expr

import (
	"github.com/neokofg/php-compiler/internal/ast"
	"github.com/neokofg/php-compiler/internal/parser/interfaces"
	"github.com/neokofg/php-compiler/internal/token"
)

type FunctionCallParser struct {
	context    interfaces.TokenReader
	exprParser interfaces.ExpressionParser
}

func NewFunctionCallParser(context interfaces.TokenReader, exprParser interfaces.ExpressionParser) *FunctionCallParser {
	return &FunctionCallParser{
		context:    context,
		exprParser: exprParser,
	}
}

func (p *FunctionCallParser) Parse(funcName string) (ast.Expr, error) {
	_, err := p.context.Expect(token.T_LPAREN)
	if err != nil {
		return nil, err
	}

	var args []ast.Expr

	if p.context.Peek().Type != token.T_RPAREN {
		for {
			arg, err := p.exprParser.ParseExpression()
			if err != nil {
				return nil, err
			}

			args = append(args, arg)

			if p.context.Peek().Type != token.T_COMMA {
				break
			}
			p.context.Next() // ,
		}
	}

	_, err = p.context.Expect(token.T_RPAREN)
	if err != nil {
		return nil, err
	}

	return &ast.FunctionCall{
		Name: funcName,
		Args: args,
	}, nil
}
