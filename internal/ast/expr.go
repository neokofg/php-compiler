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
