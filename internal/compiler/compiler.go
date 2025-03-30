package compiler

import (
	"fmt"
	"github.com/neokofg/php-compiler/internal/token"
	"github.com/neokofg/php-compiler/internal/ast"
)


type Constant struct {
	Type  string
	Value string
}


var Bytecode []byte
var Constants []Constant

var varMap = map[string]int{}
var nextVarIndex = 0

func writeUint16(Bytecode *[]byte, offsetValue uint16) {
	lowByte := byte(offsetValue & 0xFF)
	highByte := byte(offsetValue >> 8)
	*Bytecode = append(*Bytecode, lowByte, highByte)
}

func patchUint16(Bytecode []byte, patchPos int, offsetValue uint16) {
	if patchPos+1 >= len(Bytecode) {
		panic("Patch error: position is too far from Bytecode")
	}
	lowByte := byte(offsetValue & 0xFF)
	highByte := byte(offsetValue >> 8)
	Bytecode[patchPos] = lowByte
	Bytecode[patchPos+1] = highByte
}

func getVarIndex(name string) int {
	if idx, ok := varMap[name]; ok {
		return idx
	}
	varMap[name] = nextVarIndex
	nextVarIndex++
	return varMap[name]
}

func addConstant(c Constant) int {
	for i, existing := range Constants {
		if existing == c {
			return i
		}
	}
	Constants = append(Constants, c)
	return len(Constants) - 1
}

func CompileStmt(stmt ast.Stmt) {
	switch s := stmt.(type) {
	case *ast.AssignStmt:
		CompileExpr(s.Expr)
		varIdx := getVarIndex(s.Name)
		Bytecode = append(Bytecode, OP_STORE_VAR, byte(varIdx))

	case *ast.EchoStmt:
		CompileExpr(s.Expr)
		Bytecode = append(Bytecode, OP_PRINT)

	case *ast.IfStmt:
		CompileExpr(s.Cond)

		jumpIfFalsePos := len(Bytecode)
		Bytecode = append(Bytecode, OP_JUMP_IF_FALSE, 0xFF, 0xFF) // Placeholder

		// THEN
		// thenStart := len(Bytecode) // thenStart
		for _, stmt := range s.Then {
			CompileStmt(stmt)
		}
		// endThen := len(Bytecode) // endThen

		if len(s.Else) > 0 {
			jumpOverElsePos := len(Bytecode)
			Bytecode = append(Bytecode, OP_JUMP, 0xFF, 0xFF) // placeholder
			
			elseStart := len(Bytecode)
			// ELSE
			for _, stmt := range s.Else {
				CompileStmt(stmt)
			}
			endElse := len(Bytecode)

			// offsets
			offsetJumpIfFalse := uint16(elseStart - (jumpIfFalsePos + 3))
			patchUint16(Bytecode, jumpIfFalsePos+1, offsetJumpIfFalse)

			offsetJumpOverElse := uint16(endElse - (jumpOverElsePos + 3))
			patchUint16(Bytecode, jumpOverElsePos+1, offsetJumpOverElse)
		} else {
			afterThen := len(Bytecode)
			offsetJumpIfFalse := uint16(afterThen - (jumpIfFalsePos + 3))
			patchUint16(Bytecode, jumpIfFalsePos+1, offsetJumpIfFalse)
		}
	case *ast.WhileStmt:
		loopStartPos := len(Bytecode)

		CompileExpr(s.Cond)

		jumpFalsePos := len(Bytecode)
		Bytecode = append(Bytecode, OP_JUMP_IF_FALSE, 0xFF, 0xFF) // Placeholder

		for _, bodyStmt := range s.Body {
			CompileStmt(bodyStmt)
		}

		jumpBackOffset := loopStartPos - (len(Bytecode) + 3)
		Bytecode = append(Bytecode, OP_JUMP)
		writeUint16(&Bytecode, uint16(int16(jumpBackOffset)))

		loopEndPos := len(Bytecode)
		jumpFalseOffset := uint16(loopEndPos - (jumpFalsePos + 3))
		patchUint16(Bytecode, jumpFalsePos+1, jumpFalseOffset)

	case *ast.ForStmt:
		if s.Init != nil {
			CompileExpr(s.Init)
			Bytecode = append(Bytecode, OP_POP)
		}

		conditionStartPos := len(Bytecode)

		if s.Cond != nil {
			CompileExpr(s.Cond)
		} else {
			trueConstIdx := addConstant(Constant{Type: "int", Value: "1"})
			Bytecode = append(Bytecode, OP_LOAD_CONST, byte(trueConstIdx))
		}

		jumpFalsePos := len(Bytecode)
		Bytecode = append(Bytecode, OP_JUMP_IF_FALSE, 0xFF, 0xFF) // Placeholder

		for _, bodyStmt := range s.Body {
			CompileStmt(bodyStmt)
		}

		incrementStartPos := len(Bytecode)
		_ = incrementStartPos
		if s.Incr != nil {
			CompileExpr(s.Incr)
			Bytecode = append(Bytecode, OP_POP)
		}

		jumpBackOffset := conditionStartPos - (len(Bytecode) + 3)
		Bytecode = append(Bytecode, OP_JUMP)
		writeUint16(&Bytecode, uint16(int16(jumpBackOffset)))

		loopEndPos := len(Bytecode)
		jumpFalseOffset := uint16(loopEndPos - (jumpFalsePos + 3))
		patchUint16(Bytecode, jumpFalsePos+1, jumpFalseOffset)
	default:
		panic(fmt.Sprintf("Unsupported type of statement: %T", stmt))
	}
}

func CompileExpr(expr ast.Expr) {
	switch e := expr.(type) {
	case *ast.NumberLiteral:
		idx := addConstant(Constant{Type: "int", Value: fmt.Sprint(e.Value)})
		Bytecode = append(Bytecode, OP_LOAD_CONST, byte(idx))

	case *ast.StringLiteral:
		idx := addConstant(Constant{Type: "string", Value: e.Value})
		Bytecode = append(Bytecode, OP_LOAD_CONST, byte(idx))

	case *ast.VarExpr:
		varIdx := getVarIndex(e.Name)
		Bytecode = append(Bytecode, OP_LOAD_VAR, byte(varIdx))

	case *ast.BinaryExpr:
		CompileExpr(e.Left)
		CompileExpr(e.Right)
		switch e.Op {
		case token.T_PLUS:  Bytecode = append(Bytecode, OP_ADD)
		case token.T_MINUS: Bytecode = append(Bytecode, OP_SUB)
		case token.T_STAR:  Bytecode = append(Bytecode, OP_MUL)
		case token.T_SLASH: Bytecode = append(Bytecode, OP_DIV)
		case token.T_GT:    Bytecode = append(Bytecode, OP_GT)
		case token.T_LT:    Bytecode = append(Bytecode, OP_LT)
		case token.T_EQEQ:  Bytecode = append(Bytecode, OP_EQ)
		case token.T_AND:   Bytecode = append(Bytecode, OP_AND)
		case token.T_OR:    Bytecode = append(Bytecode, OP_OR)
		default:
            panic(fmt.Sprintf("Unsupported binary operator: %v", e.Op))
		}
	default:
        panic(fmt.Sprintf("Unsupported type of expression: %T", expr))
	}
}
