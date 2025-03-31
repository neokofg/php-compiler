package expression_parts

import (
	"fmt"
	"github.com/neokofg/php-compiler/internal/ast"
	"github.com/neokofg/php-compiler/internal/parser/parser_contract"
	"github.com/neokofg/php-compiler/internal/token"
	"strconv"
)

func ParsePrimary(p parser_contract.ParserLike) (ast.Expr, error) {
	tok := p.Peek()

	switch tok.Type {
	case token.T_NUMBER:
		p.Next()
		val, err := strconv.Atoi(tok.Value)
		if err != nil {
			return nil, fmt.Errorf("Position %d: wrong number format: %s", p.GetPos()-1, tok.Value)
		}
		return &ast.NumberLiteral{Value: val}, nil
	case token.T_STRING:
		p.Next()
		return &ast.StringLiteral{Value: tok.Value}, nil
	case token.T_DOLLAR:
		p.Next()
		identToken, err := p.Expect(token.T_IDENT)
		if err != nil {
			return nil, err
		}
		return &ast.VarExpr{Name: identToken.Value}, nil
	case token.T_LPAREN:
		p.Next()
		expr, err := ParseExpression(p)
		if err != nil {
			return nil, err
		}
		_, err = p.Expect(token.T_RPAREN)
		if err != nil {
			return nil, err
		}
		return expr, nil
	case token.T_ILLEGAL:
		p.Next()
		return nil, fmt.Errorf("Lexer error in position %d: %s", p.GetPos()-1, tok.Value)
	default:
		p.Next()
		return nil, fmt.Errorf("Position %d: expected expression (num, string, var, '('), but found token: %v (%q)", p.GetPos()-1, tok.Type, tok.Value)
	}
}
