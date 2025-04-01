// Licensed under GNU GPL v3. See LICENSE file for details.
package stmt

import (
	"github.com/neokofg/php-compiler/internal/ast"
	"github.com/neokofg/php-compiler/internal/parser/interfaces"
	"github.com/neokofg/php-compiler/internal/token"
)

type SwitchParser struct {
	context    interfaces.TokenReader
	exprParser interfaces.ExpressionParser
	stmtParser interfaces.StatementParser
}

func NewSwitchParser(context interfaces.TokenReader, exprParser interfaces.ExpressionParser, stmtParser interfaces.StatementParser) *SwitchParser {
	return &SwitchParser{
		context:    context,
		exprParser: exprParser,
		stmtParser: stmtParser,
	}
}

func (p *SwitchParser) Parse() (ast.Stmt, error) {
	p.context.Next()

	_, err := p.context.Expect(token.T_LPAREN)
	if err != nil {
		return nil, err
	}

	expr, err := p.exprParser.ParseExpression()
	if err != nil {
		return nil, err
	}

	_, err = p.context.Expect(token.T_RPAREN)
	if err != nil {
		return nil, err
	}

	_, err = p.context.Expect(token.T_LBRACE)
	if err != nil {
		return nil, err
	}

	var cases []ast.CaseStmt
	for p.context.Peek().Type != token.T_RBRACE && p.context.Peek().Type != token.T_EOF {
		var caseExpr ast.Expr
		var caseStmts []ast.Stmt

		if p.context.Peek().Type == token.T_CASE {
			p.context.Next()

			caseExpr, err = p.exprParser.ParseExpression()
			if err != nil {
				return nil, err
			}

			_, err = p.context.Expect(token.T_COLON)
			if err != nil {
				return nil, err
			}
		} else if p.context.Peek().Type == token.T_DEFAULT {
			p.context.Next()

			_, err = p.context.Expect(token.T_COLON)
			if err != nil {
				return nil, err
			}

		} else {
			return nil, &token.UnexpectedTokenError{
				Expected: "case или default",
				Found:    p.context.Peek(),
				Pos:      p.context.GetPos(),
			}
		}

		for p.context.Peek().Type != token.T_CASE &&
			p.context.Peek().Type != token.T_DEFAULT &&
			p.context.Peek().Type != token.T_RBRACE &&
			p.context.Peek().Type != token.T_EOF {

			stmt, err := p.stmtParser.ParseStatement()
			if err != nil {
				return nil, err
			}

			caseStmts = append(caseStmts, stmt)
		}

		cases = append(cases, ast.CaseStmt{
			Expr:  caseExpr,
			Stmts: caseStmts,
		})
	}

	_, err = p.context.Expect(token.T_RBRACE)
	if err != nil {
		return nil, err
	}

	return &ast.SwitchStmt{
		Expr:  expr,
		Cases: cases,
	}, nil
}
