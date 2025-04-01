// Licensed under GNU GPL v3. See LICENSE file for details.
package stmt

import (
	"fmt"
	"github.com/neokofg/php-compiler/internal/ast"
	"github.com/neokofg/php-compiler/internal/parser/interfaces"
	"github.com/neokofg/php-compiler/internal/token"
)

type AssignParser struct {
	context    interfaces.TokenReader
	exprParser interfaces.ExpressionParser
}

func NewAssignParser(context interfaces.TokenReader, exprParser interfaces.ExpressionParser) *AssignParser {
	return &AssignParser{
		context:    context,
		exprParser: exprParser,
	}
}

func (p *AssignParser) Parse() (ast.Stmt, error) {
	p.context.Next() // $
	identToken, err := p.context.Expect(token.T_IDENT)
	if err != nil {
		return nil, err
	}

	next := p.context.Peek().Type
	switch next {
	case token.T_EQ:
		p.context.Next() // =
		expr, err := p.exprParser.ParseExpression()
		if err != nil {
			return nil, err
		}

		_, err = p.context.Expect(token.T_SEMI)
		if err != nil {
			return nil, err
		}

		return &ast.AssignStmt{Name: identToken.Value, Expr: expr}, nil

	case token.T_PLUS_EQ, token.T_MINUS_EQ, token.T_MUL_EQ, token.T_DIV_EQ, token.T_MOD_EQ, token.T_DOT_EQ:
		op := p.context.Next().Type
		expr, err := p.exprParser.ParseExpression()
		if err != nil {
			return nil, err
		}

		_, err = p.context.Expect(token.T_SEMI)
		if err != nil {
			return nil, err
		}

		return &ast.CompoundAssignStmt{Name: identToken.Value, Op: op, Expr: expr}, nil

	case token.T_INC, token.T_DEC:
		op := p.context.Next().Type

		_, err = p.context.Expect(token.T_SEMI)
		if err != nil {
			return nil, err
		}

		return &ast.PostfixExpr{
			Expr: &ast.VarExpr{Name: identToken.Value},
			Op:   op,
		}, nil

	default:
		return nil, fmt.Errorf("Position %d: expected assignment operator after variable", p.context.GetPos())
	}
}
