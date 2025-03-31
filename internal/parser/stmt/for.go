// Licensed under GNU GPL v3. See LICENSE file for details.
package stmt

import (
	"fmt"
	"github.com/neokofg/php-compiler/internal/ast"
	"github.com/neokofg/php-compiler/internal/parser/interfaces"
	"github.com/neokofg/php-compiler/internal/token"
)

type ForParser struct {
	context     interfaces.TokenReader
	exprParser  interfaces.ExpressionParser
	blockParser *BlockParser
}

func NewForParser(context interfaces.TokenReader, exprParser interfaces.ExpressionParser, blockParser *BlockParser) *ForParser {
	return &ForParser{
		context:     context,
		exprParser:  exprParser,
		blockParser: blockParser,
	}
}

func (p *ForParser) Parse() (ast.Stmt, error) {
	p.context.Next()

	_, err := p.context.Expect(token.T_LPAREN)
	if err != nil {
		return nil, err
	}

	var initExpr ast.Expr
	if p.context.Peek().Type != token.T_SEMI {
		initExpr, err = p.exprParser.ParseExpression()
		if err != nil {
			return nil, fmt.Errorf("error parsing for-loop initializer: %w", err)
		}
	}

	_, err = p.context.Expect(token.T_SEMI)
	if err != nil {
		return nil, err
	}

	var condExpr ast.Expr
	if p.context.Peek().Type != token.T_SEMI {
		condExpr, err = p.exprParser.ParseExpression()
		if err != nil {
			return nil, fmt.Errorf("error parsing for-loop condition: %w", err)
		}
	}

	_, err = p.context.Expect(token.T_SEMI)
	if err != nil {
		return nil, err
	}

	var incrExpr ast.Expr
	if p.context.Peek().Type != token.T_RPAREN {
		incrExpr, err = p.exprParser.ParseExpression()
		if err != nil {
			return nil, fmt.Errorf("error parsing for-loop increment: %w", err)
		}
	}

	_, err = p.context.Expect(token.T_RPAREN)
	if err != nil {
		return nil, err
	}

	bodyBlock, err := p.blockParser.Parse()
	if err != nil {
		return nil, err
	}

	return &ast.ForStmt{Init: initExpr, Cond: condExpr, Incr: incrExpr, Body: bodyBlock}, nil
}
