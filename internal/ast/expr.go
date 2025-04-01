// Licensed under GNU GPL v3. See LICENSE file for details.
package ast

import "github.com/neokofg/php-compiler/internal/token"

type Expr interface{}

type NumberLiteral struct {
	Value int
}

type StringLiteral struct {
	Value string
}

type VarExpr struct {
	Name string
}

type BinaryExpr struct {
	Left  Expr
	Op    token.TokenType
	Right Expr
}

type UnaryExpr struct {
	Op   token.TokenType
	Expr Expr
}

type BooleanLiteral struct {
	Value bool
}

type PostfixExpr struct {
	Expr Expr
	Op   token.TokenType
}

type PrefixExpr struct {
	Op   token.TokenType
	Expr Expr
}

type CompoundAssignStmt struct {
	Name string
	Op   token.TokenType
	Expr Expr
}
