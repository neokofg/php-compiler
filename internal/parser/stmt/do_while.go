// Licensed under GNU GPL v3. See LICENSE file for details.
package stmt

import (
	"github.com/neokofg/php-compiler/internal/ast"
	"github.com/neokofg/php-compiler/internal/parser/interfaces"
	"github.com/neokofg/php-compiler/internal/token"
)

type DoWhileParser struct {
	context     interfaces.TokenReader
	exprParser  interfaces.ExpressionParser
	blockParser *BlockParser
}

func NewDoWhileParser(context interfaces.TokenReader, exprParser interfaces.ExpressionParser, blockParser *BlockParser) *DoWhileParser {
	return &DoWhileParser{
		context:     context,
		exprParser:  exprParser,
		blockParser: blockParser,
	}
}

func (p *DoWhileParser) Parse() (ast.Stmt, error) {
	p.context.Next() // do

	bodyBlock, err := p.blockParser.Parse()
	if err != nil {
		return nil, err
	}

	_, err = p.context.Expect(token.T_WHILE)
	if err != nil {
		return nil, err
	}

	_, err = p.context.Expect(token.T_LPAREN)
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

	_, err = p.context.Expect(token.T_SEMI)
	if err != nil {
		return nil, err
	}

	return &ast.DoWhileStmt{Body: bodyBlock, Cond: condExpr}, nil
}
