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

	var expr ast.Expr
	var err error

	if tok.Type == token.T_INC || tok.Type == token.T_DEC {
		op := p.context.Next().Type

		if p.context.Peek().Type != token.T_DOLLAR {
			return nil, fmt.Errorf("Position %d: expected variable after increment/decrement", p.context.GetPos())
		}

		p.context.Next()
		identToken, err := p.context.Expect(token.T_IDENT)
		if err != nil {
			return nil, err
		}

		return &ast.PrefixExpr{
			Op:   op,
			Expr: &ast.VarExpr{Name: identToken.Value},
		}, nil
	}

	switch tok.Type {
	case token.T_NUMBER:
		p.context.Next()
		val, err := strconv.Atoi(tok.Value)
		if err != nil {
			return nil, fmt.Errorf("Position %d: wrong number format: %s", p.context.GetPos()-1, tok.Value)
		}
		expr = &ast.NumberLiteral{Value: val}

	case token.T_STRING:
		p.context.Next()
		expr = &ast.StringLiteral{Value: tok.Value}

	case token.T_DOLLAR:
		p.context.Next()
		identToken, err := p.context.Expect(token.T_IDENT)
		if err != nil {
			return nil, err
		}
		expr = &ast.VarExpr{Name: identToken.Value}

	case token.T_LPAREN:
		p.context.Next()

		exprParser := NewParser(p.context)
		innerExpr, err := exprParser.ParseExpression()
		if err != nil {
			return nil, err
		}

		_, err = p.context.Expect(token.T_RPAREN)
		if err != nil {
			return nil, err
		}
		expr = innerExpr

	case token.T_ILLEGAL:
		p.context.Next()
		return nil, fmt.Errorf("Lexer error in position %d: %s", p.context.GetPos()-1, tok.Value)

	case token.T_NOT:
		p.context.Next()
		innerExpr, err := p.Parse()
		if err != nil {
			return nil, err
		}
		expr = &ast.UnaryExpr{Op: token.T_NOT, Expr: innerExpr}

	case token.T_TRUE:
		p.context.Next()
		expr = &ast.BooleanLiteral{Value: true}

	case token.T_FALSE:
		p.context.Next()
		expr = &ast.BooleanLiteral{Value: false}

	case token.T_INC, token.T_DEC:
		op := p.context.Next().Type
		exprValue, err := p.Parse()
		if err != nil {
			return nil, err
		}

		varExpr, ok := exprValue.(*ast.VarExpr)
		if !ok {
			return nil, fmt.Errorf("Position %d: can only increment/decrement variables", p.context.GetPos()-1)
		}

		return &ast.PrefixExpr{Op: op, Expr: varExpr}, nil
	default:
		p.context.Next()
		return nil, fmt.Errorf("Position %d: expected expression (num, string, var, '('), but found token: %v (%q)", p.context.GetPos()-1, tok.Type, tok.Value)
	}

	if expr != nil {
		if p.context.Peek().Type == token.T_INC || p.context.Peek().Type == token.T_DEC {
			// Проверяем, что инкремент/декремент применяется к переменной
			_, ok := expr.(*ast.VarExpr)
			if !ok {
				return nil, fmt.Errorf("Position %d: can only increment/decrement variables", p.context.GetPos())
			}

			op := p.context.Next().Type
			return &ast.PostfixExpr{Expr: expr, Op: op}, nil
		}
	}
	
	return expr, err
}
