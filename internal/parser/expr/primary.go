// Licensed under GNU GPL v3. See LICENSE file for details.
package expr

import (
	"fmt"
	"github.com/neokofg/php-compiler/internal/ast"
	"github.com/neokofg/php-compiler/internal/parser/interfaces"
	"github.com/neokofg/php-compiler/internal/token"
	"strconv"
)

type PrimaryParser struct {
	context interfaces.TokenReader
}

func NewPrimaryParser(context interfaces.TokenReader) *PrimaryParser {
	return &PrimaryParser{
		context: context,
	}
}

func (p *PrimaryParser) Parse() (ast.Expr, error) {
	tok := p.context.Peek()

	switch tok.Type {
	case token.T_NUMBER:
		p.context.Next()
		val, err := strconv.Atoi(tok.Value)
		if err != nil {
			return nil, fmt.Errorf("Position %d: wrong number format: %s", p.context.GetPos()-1, tok.Value)
		}
		return &ast.NumberLiteral{Value: val}, nil

	case token.T_STRING:
		p.context.Next()
		return &ast.StringLiteral{Value: tok.Value}, nil

	case token.T_DOLLAR:
		p.context.Next()
		identToken, err := p.context.Expect(token.T_IDENT)
		if err != nil {
			return nil, err
		}
		return &ast.VarExpr{Name: identToken.Value}, nil

	case token.T_LPAREN:
		p.context.Next()

		exprParser := NewParser(p.context)
		expr, err := exprParser.ParseExpression()
		if err != nil {
			return nil, err
		}

		_, err = p.context.Expect(token.T_RPAREN)
		if err != nil {
			return nil, err
		}
		return expr, nil

	case token.T_ILLEGAL:
		p.context.Next()
		return nil, fmt.Errorf("Lexer error in position %d: %s", p.context.GetPos()-1, tok.Value)

	default:
		p.context.Next()
		return nil, fmt.Errorf("Position %d: expected expression (num, string, var, '('), but found token: %v (%q)", p.context.GetPos()-1, tok.Type, tok.Value)
	}
}
