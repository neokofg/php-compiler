// Licensed under GNU GPL v3. See LICENSE file for details.
package ast

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

type BreakStmt struct{}

type ContinueStmt struct{}

type DoWhileStmt struct {
	Body []Stmt
	Cond Expr
}

type SwitchStmt struct {
	Expr  Expr
	Cases []CaseStmt
}

type CaseStmt struct {
	Expr  Expr
	Stmts []Stmt
}

type FunctionDecl struct {
	Name      string
	Params    []string
	Body      []Stmt
	StartAddr int
}

type ReturnStmt struct {
	Expr Expr
}

type FunctionCallStmt struct {
	Call *FunctionCall
}
