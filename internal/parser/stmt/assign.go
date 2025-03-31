package stmt

import (
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

	_, err = p.context.Expect(token.T_EQ)
	if err != nil {
		return nil, err
	}

	expr, err := p.exprParser.ParseExpression()
	if err != nil {
		return nil, err
	}

	_, err = p.context.Expect(token.T_SEMI)
	if err != nil {
		return nil, err
	}

	return &ast.AssignStmt{Name: identToken.Value, Expr: expr}, nil
}
