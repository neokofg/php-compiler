package interfaces

import (
	"github.com/neokofg/php-compiler/internal/ast"
	"github.com/neokofg/php-compiler/internal/token"
)

type TokenReader interface {
	Next() token.Token
	Peek() token.Token
	Expect(t token.TokenType) (token.Token, error)
	GetPos() int
	SetPos(int)
}

type ExpressionParser interface {
	ParseExpression() (ast.Expr, error)
}

type StatementParser interface {
	ParseStatement() (ast.Stmt, error)
	ParseBlock() ([]ast.Stmt, error)
	ParseOptionalExpression(terminator token.TokenType) (ast.Expr, error)
}

type Parser interface {
	TokenReader
	Parse() ([]ast.Stmt, error)
}
