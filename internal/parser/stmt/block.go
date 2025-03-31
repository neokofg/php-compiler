// Licensed under GNU GPL v3. See LICENSE file for details.
package stmt

import (
	"github.com/neokofg/php-compiler/internal/ast"
	"github.com/neokofg/php-compiler/internal/parser/interfaces"
	"github.com/neokofg/php-compiler/internal/token"
)

type BlockParser struct {
	context    interfaces.TokenReader
	stmtParser interfaces.StatementParser
}

func NewBlockParser(context interfaces.TokenReader, stmtParser interfaces.StatementParser) *BlockParser {
	return &BlockParser{
		context:    context,
		stmtParser: stmtParser,
	}
}

func (p *BlockParser) Parse() ([]ast.Stmt, error) {
	_, err := p.context.Expect(token.T_LBRACE)
	if err != nil {
		return nil, err
	}

	var stmts []ast.Stmt
	for p.context.Peek().Type != token.T_RBRACE && p.context.Peek().Type != token.T_EOF {
		stmt, err := p.stmtParser.ParseStatement()
		if err != nil {
			return nil, err
		}
		stmts = append(stmts, stmt)
	}

	_, err = p.context.Expect(token.T_RBRACE)
	if err != nil {
		return nil, err
	}

	return stmts, nil
}
