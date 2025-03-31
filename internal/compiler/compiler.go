package compiler

import (
	"github.com/neokofg/php-compiler/internal/ast"
	"github.com/neokofg/php-compiler/internal/compiler/bytecode"
	"github.com/neokofg/php-compiler/internal/compiler/constant"
	"github.com/neokofg/php-compiler/internal/compiler/expr"
	"github.com/neokofg/php-compiler/internal/compiler/interfaces"
	"github.com/neokofg/php-compiler/internal/compiler/stmt"
)

type Compiler struct {
	context      *interfaces.Context
	stmtCompiler interfaces.StmtCompiler
	exprCompiler interfaces.ExprCompiler
}

func New() *Compiler {
	// Создаем контекст
	context := interfaces.NewContext()

	// Создаем компиляторы
	exprCompiler := expr.NewCompiler(context)
	stmtCompiler := stmt.NewCompiler(context, exprCompiler)

	return &Compiler{
		context:      context,
		stmtCompiler: stmtCompiler,
		exprCompiler: exprCompiler,
	}
}

func (c *Compiler) CompileProgram(stmts []ast.Stmt) error {
	for _, statement := range stmts {
		if err := c.stmtCompiler.CompileStmt(statement); err != nil {
			return err
		}
	}

	c.context.BytecodeBuilder.Append(bytecode.OP_HALT)

	return nil
}

func (c *Compiler) GetBytecode() []byte {
	return c.context.BytecodeBuilder.Get()
}

func (c *Compiler) GetConstants() []constant.Constant {
	return c.context.ConstantPool.GetAll()
}

// Глобальные переменные для обратной совместимости
var (
	Bytecode  []byte
	Constants []constant.Constant
)

// CompileStmt для обратной совместимости
func CompileStmt(stmt ast.Stmt) {
	if Bytecode == nil {
		compiler := New()

		Bytecode = compiler.context.BytecodeBuilder.Get()
		Constants = compiler.context.ConstantPool.GetAll()

		compiler.context.BytecodeBuilder.SetSyncCallback(func(code []byte) {
			Bytecode = code
		})

		compiler.context.ConstantPool.SetSyncCallback(func(consts []constant.Constant) {
			Constants = consts
		})

		compiler.stmtCompiler.CompileStmt(stmt)
	} else {
		compiler := New()
		compiler.stmtCompiler.CompileStmt(stmt)

		newBytecode := compiler.context.BytecodeBuilder.Get()
		Bytecode = append(Bytecode, newBytecode...)

		newConstants := compiler.context.ConstantPool.GetAll()
		existingConstMap := make(map[string]bool)

		for _, c := range Constants {
			existingConstMap[c.Type+":"+c.Value] = true
		}

		for _, c := range newConstants {
			key := c.Type + ":" + c.Value
			if !existingConstMap[key] {
				Constants = append(Constants, c)
				existingConstMap[key] = true
			}
		}
	}
}
