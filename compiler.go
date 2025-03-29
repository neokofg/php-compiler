package main

import "fmt"

type Constant struct {
	Type  string
	Value string
}

const (
	OP_LOAD_CONST = 0x01
	OP_PRINT      = 0x02
	OP_ADD        = 0x03
	OP_SUB        = 0x04
	OP_MUL        = 0x05
	OP_DIV        = 0x06
	OP_STORE_VAR  = 0x10
	OP_LOAD_VAR   = 0x11
	OP_HALT       = 0xFF
)


var bytecode []byte
var constants []Constant

var varMap = map[string]int{}
var nextVarIndex = 0

func getVarIndex(name string) int {
	if idx, ok := varMap[name]; ok {
		return idx
	}
	varMap[name] = nextVarIndex
	nextVarIndex++
	return varMap[name]
}

func addConstant(c Constant) int {
	for i, existing := range constants {
		if existing == c {
			return i
		}
	}
	constants = append(constants, c)
	return len(constants) - 1
}

func CompileStmt(stmt Stmt) {
	switch s := stmt.(type) {
	case *AssignStmt:
		CompileExpr(s.Expr)
		varIdx := getVarIndex(s.Name)
		bytecode = append(bytecode, OP_STORE_VAR, byte(varIdx))

	case *EchoStmt:
		CompileExpr(s.Expr)
		bytecode = append(bytecode, OP_PRINT)

	default:
		panic("неподдерживаемый тип выражения в CompileStmt")
	}
}

func CompileExpr(expr Expr) {
	switch e := expr.(type) {
	case *NumberLiteral:
		idx := addConstant(Constant{Type: "int", Value: fmt.Sprint(e.Value)})
		bytecode = append(bytecode, OP_LOAD_CONST, byte(idx))

	case *StringLiteral:
		idx := addConstant(Constant{Type: "string", Value: e.Value})
		bytecode = append(bytecode, OP_LOAD_CONST, byte(idx))

	case *VarExpr:
		varIdx := getVarIndex(e.Name)
		bytecode = append(bytecode, OP_LOAD_VAR, byte(varIdx))

	case *BinaryExpr:
		CompileExpr(e.Left)
		CompileExpr(e.Right)
		switch e.Op {
		case T_PLUS:
			bytecode = append(bytecode, OP_ADD)
		case T_MINUS:
			bytecode = append(bytecode, OP_SUB)
		case T_STAR:
			bytecode = append(bytecode, OP_MUL)
		case T_SLASH:
			bytecode = append(bytecode, OP_DIV)
		default:
			panic("неподдерживаемый оператор в BinaryExpr")
		}

	default:
		panic("неподдерживаемое выражение в CompileExpr")
	}
}
