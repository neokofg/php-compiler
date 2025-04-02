// Licensed under GNU GPL v3. See LICENSE file for details.
package stmt

import (
	"fmt"
	"github.com/neokofg/php-compiler/internal/ast"
	"github.com/neokofg/php-compiler/internal/parser/interfaces"
	"github.com/neokofg/php-compiler/internal/token"
)

type FunctionParser struct {
	context     interfaces.TokenReader
	exprParser  interfaces.ExpressionParser
	blockParser *BlockParser
}

func NewFunctionParser(context interfaces.TokenReader, exprParser interfaces.ExpressionParser, blockParser *BlockParser) *FunctionParser {
	return &FunctionParser{
		context:     context,
		exprParser:  exprParser,
		blockParser: blockParser,
	}
}

func (p *FunctionParser) Parse() (ast.Stmt, error) {
	p.context.Next()

	if p.context.Peek().Type != token.T_IDENT {
		return nil, fmt.Errorf("Position %d: expected function name after 'function' keyword, got: %v (%s)",
			p.context.GetPos(), p.context.Peek().Type, p.context.Peek().Value)
	}

	name := p.context.Next()

	if p.context.Peek().Type != token.T_LPAREN {
		return nil, fmt.Errorf("Position %d: expected '(' after function name, got: %v (%s)",
			p.context.GetPos(), p.context.Peek().Type, p.context.Peek().Value)
	}

	p.context.Next()

	var params []string
	if p.context.Peek().Type != token.T_RPAREN {
		for {
			if p.context.Peek().Type != token.T_DOLLAR {
				break
			}

			p.context.Next() // $
			param, err := p.context.Expect(token.T_IDENT)
			if err != nil {
				return nil, err
			}

			params = append(params, param.Value)

			if p.context.Peek().Type != token.T_COMMA {
				break
			}
			p.context.Next() // ,
		}
	}

	if p.context.Peek().Type != token.T_RPAREN {
		return nil, fmt.Errorf("Position %d: expected ')' after function parameters, got: %v (%s)",
			p.context.GetPos(), p.context.Peek().Type, p.context.Peek().Value)
	}

	p.context.Next()

	body, err := p.blockParser.Parse()
	if err != nil {
		return nil, err
	}

	return &ast.FunctionDecl{
		Name:   name.Value,
		Params: params,
		Body:   body,
	}, nil
}
