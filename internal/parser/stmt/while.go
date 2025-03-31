// Licensed under GNU GPL v3. See LICENSE file for details.
package stmt

import (
	"github.com/neokofg/php-compiler/internal/ast"
	"github.com/neokofg/php-compiler/internal/parser/interfaces"
	"github.com/neokofg/php-compiler/internal/token"
)

type WhileParser struct {
	context     interfaces.TokenReader
	exprParser  interfaces.ExpressionParser
	blockParser *BlockParser
}

func NewWhileParser(context interfaces.TokenReader, exprParser interfaces.ExpressionParser, blockParser *BlockParser) *WhileParser {
	return &WhileParser{
		context:     context,
		exprParser:  exprParser,
		blockParser: blockParser,
	}
}

func (p *WhileParser) Parse() (ast.Stmt, error) {
	p.context.Next() // while

	_, err := p.context.Expect(token.T_LPAREN)
	if err != nil {
		return nil, err
	}

	condExpr, err := p.exprParser.ParseExpression()
	if err != nil {
		return nil, err
	}

	_, err = p.context.Expect(token.T_RPAREN)
	if err != nil {
		return nil, err
	}

	bodyBlock, err := p.blockParser.Parse()
	if err != nil {
		return nil, err
	}

	return &ast.WhileStmt{Cond: condExpr, Body: bodyBlock}, nil
}
