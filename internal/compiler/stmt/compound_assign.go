// Licensed under GNU GPL v3. See LICENSE file for details.
package stmt

import (
	"github.com/neokofg/php-compiler/internal/ast"
	"github.com/neokofg/php-compiler/internal/compiler/bytecode"
	"github.com/neokofg/php-compiler/internal/compiler/interfaces"
	"github.com/neokofg/php-compiler/internal/token"
)

type CompoundAssignCompiler struct {
	context      interfaces.CompilationContext
	exprCompiler interfaces.ExprCompiler
}

func NewCompoundAssignCompiler(context interfaces.CompilationContext, exprCompiler interfaces.ExprCompiler) *CompoundAssignCompiler {
	return &CompoundAssignCompiler{
		context:      context,
		exprCompiler: exprCompiler,
	}
}

func (c *CompoundAssignCompiler) Compile(stmt *ast.CompoundAssignStmt) error {
	varIdx := c.context.GetVariableManager().GetIndex(stmt.Name)

	c.context.GetBytecodeBuilder().Append(bytecode.OP_LOAD_VAR)
	c.context.GetBytecodeBuilder().Append(byte(varIdx))

	if err := c.exprCompiler.CompileExpr(stmt.Expr); err != nil {
		return err
	}

	switch stmt.Op {
	case token.T_PLUS_EQ:
		c.context.GetBytecodeBuilder().Append(bytecode.OP_ASSIGN_ADD)
	case token.T_MINUS_EQ:
		c.context.GetBytecodeBuilder().Append(bytecode.OP_ASSIGN_SUB)
	case token.T_MUL_EQ:
		c.context.GetBytecodeBuilder().Append(bytecode.OP_ASSIGN_MUL)
	case token.T_DIV_EQ:
		c.context.GetBytecodeBuilder().Append(bytecode.OP_ASSIGN_DIV)
	case token.T_MOD_EQ:
		c.context.GetBytecodeBuilder().Append(bytecode.OP_ASSIGN_MOD)
	case token.T_DOT_EQ:
		c.context.GetBytecodeBuilder().Append(bytecode.OP_ASSIGN_CONCAT)
	}

	c.context.GetBytecodeBuilder().Append(bytecode.OP_STORE_VAR)
	c.context.GetBytecodeBuilder().Append(byte(varIdx))

	return nil
}
