// Licensed under GNU GPL v3. See LICENSE file for details.
package stmt

import (
	"github.com/neokofg/php-compiler/internal/ast"
	"github.com/neokofg/php-compiler/internal/parser/interfaces"
	"github.com/neokofg/php-compiler/internal/token"
)

type IfParser struct {
	context     interfaces.TokenReader
	exprParser  interfaces.ExpressionParser
	blockParser *BlockParser
}

func NewIfParser(context interfaces.TokenReader, exprParser interfaces.ExpressionParser, blockParser *BlockParser) *IfParser {
	return &IfParser{
		context:     context,
		exprParser:  exprParser,
		blockParser: blockParser,
	}
}

func (p *IfParser) Parse() (ast.Stmt, error) {
	p.context.Next() // if

	_, err := p.context.Expect(token.T_LPAREN)
	if err != nil {
		return nil, err
	}

	cond, err := p.exprParser.ParseExpression()
	if err != nil {
		return nil, err
	}

	_, err = p.context.Expect(token.T_RPAREN)
	if err != nil {
		return nil, err
	}

	thenBlock, err := p.blockParser.Parse()
	if err != nil {
		return nil, err
	}

	var elseBlock []ast.Stmt
	if p.context.Peek().Type == token.T_ELSE {
		p.context.Next() // else
		elseBlock, err = p.blockParser.Parse()
		if err != nil {
			return nil, err
		}
	}

	return &ast.IfStmt{Cond: cond, Then: thenBlock, Else: elseBlock}, nil
}
