package main

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
	Op    TokenType
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
