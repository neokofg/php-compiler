package main

import "fmt"

type Constant struct {
	Type  string
	Value string
}

const (
	OP_LOAD_CONST    = 0x01
	OP_PRINT         = 0x02
	OP_ADD           = 0x03
	OP_SUB           = 0x04
	OP_MUL           = 0x05
	OP_DIV           = 0x06
	OP_STORE_VAR     = 0x10
	OP_LOAD_VAR      = 0x11
	OP_HALT          = 0xFF
	OP_JUMP_IF_FALSE = 0x20
	OP_JUMP          = 0x21
	OP_GT 			 = 0x07
	OP_LT  			 = 0x08
	OP_AND 			 = 0x09
	OP_OR  			 = 0x0A
	OP_EQ 			 = 0x0B
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

	case *IfStmt:
		CompileExpr(s.Cond)

		jumpIfFalsePos := len(bytecode)
		bytecode = append(bytecode, OP_JUMP_IF_FALSE, 0xFF) // временно

		// THEN блок
		for _, stmt := range s.Then {
			CompileStmt(stmt)
		}
		endThen := len(bytecode)

		if len(s.Else) > 0 {
			jumpOverElsePos := len(bytecode)
			bytecode = append(bytecode, OP_JUMP, 0xFF) // временно

			// ELSE блок
			for _, stmt := range s.Else {
				CompileStmt(stmt)
			}
			endElse := len(bytecode)

			// правим оффсеты
			bytecode[jumpIfFalsePos+1] = byte(jumpOverElsePos - jumpIfFalsePos)
			bytecode[jumpOverElsePos+1] = byte(endElse - (jumpOverElsePos + 2))
		} else {
			bytecode[jumpIfFalsePos+1] = byte(endThen - (jumpIfFalsePos + 2))
		}

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
		case T_GT:
			bytecode = append(bytecode, OP_GT)
		case T_LT:
			bytecode = append(bytecode, OP_LT)
		case T_EQEQ:
			bytecode = append(bytecode, OP_EQ)
		case T_AND:
			bytecode = append(bytecode, OP_AND)
		case T_OR:
			bytecode = append(bytecode, OP_OR)
		default:
			panic("неподдерживаемый оператор в BinaryExpr")
		}
	default:
		panic("неподдерживаемое выражение в CompileExpr")
	}
}
