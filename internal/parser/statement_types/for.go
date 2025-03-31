package statement_types

import (
	"fmt"
	"github.com/neokofg/php-compiler/internal/ast"
	"github.com/neokofg/php-compiler/internal/parser/parser_contract"
	"github.com/neokofg/php-compiler/internal/token"
)

func ParseForStatement(p parser_contract.ParserLike) (*ast.ForStmt, error) {
	p.Next()
	_, err := p.Expect(token.T_LPAREN)
	if err != nil {
		return nil, err
	}

	initExpr, err := p.ParseOptionalExpression(token.T_SEMI)
	if err != nil {
		return nil, fmt.Errorf("error parsing for-loop initializer: %w", err)
	}
	_, err = p.Expect(token.T_SEMI)
	if err != nil {
		return nil, err
	}

	condExpr, err := p.ParseOptionalExpression(token.T_SEMI)
	if err != nil {
		return nil, fmt.Errorf("error parsing for-loop condition: %w", err)
	}
	_, err = p.Expect(token.T_SEMI)
	if err != nil {
		return nil, err
	}

	incrExpr, err := p.ParseOptionalExpression(token.T_RPAREN)
	if err != nil {
		return nil, fmt.Errorf("error parsing for-loop increment: %w", err)
	}
	_, err = p.Expect(token.T_RPAREN)
	if err != nil {
		return nil, err
	}

	bodyBlock, err := p.ParseBlock()
	if err != nil {
		return nil, err
	}

	return &ast.ForStmt{Init: initExpr, Cond: condExpr, Incr: incrExpr, Body: bodyBlock}, nil
}
