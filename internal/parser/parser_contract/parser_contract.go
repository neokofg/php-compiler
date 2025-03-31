package parser_contract

import (
	"github.com/neokofg/php-compiler/internal/ast"
	"github.com/neokofg/php-compiler/internal/token"
)

type ParserLike interface {
	Next() token.Token
	Expect(t token.TokenType) (token.Token, error)
	Peek() token.Token
	ParseBlock() ([]ast.Stmt, error)
	ParseOptionalExpression(terminator token.TokenType) (ast.Expr, error)

	GetPos() int
	SetPos(int)
}
