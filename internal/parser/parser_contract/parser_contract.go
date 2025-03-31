package parser_contract

import (
	"github.com/neokofg/php-compiler/internal/ast"
	token2 "github.com/neokofg/php-compiler/internal/token"
)

type ParserLike interface {
	Next() token2.Token
	Expect(t token2.TokenType) (token2.Token, error)
	Peek() token2.Token
	ParseBlock() ([]ast.Stmt, error)
	ParseOptionalExpression(terminator token2.TokenType) (ast.Expr, error)

	GetPos() int
	SetPos(int)
}
