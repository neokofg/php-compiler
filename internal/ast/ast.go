package ast

import (
	"github.com/neokofg/php-compiler/internal/token"
)

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

type Stmt interface{}

type AssignStmt struct {
	Name string
	Expr Expr
}

type EchoStmt struct {
	Expr Expr
}

type IfStmt struct {
	Cond Expr
	Then []Stmt
	Else []Stmt
}

type WhileStmt struct {
	Cond Expr
	Body []Stmt
}

type ForStmt struct {
	Init Expr
	Cond Expr
	Incr Expr
	Body []Stmt
}