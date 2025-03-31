// Licensed under GNU GPL v3. See LICENSE file for details.
package stmt

import (
	"github.com/neokofg/php-compiler/internal/ast"
	"github.com/neokofg/php-compiler/internal/compiler/bytecode"
	"github.com/neokofg/php-compiler/internal/compiler/constant"
	"github.com/neokofg/php-compiler/internal/compiler/interfaces"
)

type ForCompiler struct {
	context      interfaces.CompilationContext
	exprCompiler interfaces.ExprCompiler
	stmtCompiler interfaces.StmtCompiler
}

func NewForCompiler(context interfaces.CompilationContext, exprCompiler interfaces.ExprCompiler, stmtCompiler interfaces.StmtCompiler) *ForCompiler {
	return &ForCompiler{
		context:      context,
		exprCompiler: exprCompiler,
		stmtCompiler: stmtCompiler,
	}
}

func (c *ForCompiler) Compile(stmt *ast.ForStmt) error {
	// Initialize
	if stmt.Init != nil {
		if err := c.exprCompiler.CompileExpr(stmt.Init); err != nil {
			return err
		}
		c.context.GetBytecodeBuilder().Append(bytecode.OP_POP)
	}

	conditionStartPos := c.context.GetBytecodeBuilder().CurrentPosition()

	// Condition
	if stmt.Cond != nil {
		if err := c.exprCompiler.CompileExpr(stmt.Cond); err != nil {
			return err
		}
	} else {
		// Default condition is true
		trueConstIdx := c.context.GetConstantPool().Add(constant.Constant{
			Type:  "int",
			Value: "1",
		})
		c.context.GetBytecodeBuilder().Append(bytecode.OP_LOAD_CONST)
		c.context.GetBytecodeBuilder().Append(byte(trueConstIdx))
	}

	jumpFalsePos := c.context.GetBytecodeBuilder().CurrentPosition()
	c.context.GetBytecodeBuilder().Append(bytecode.OP_JUMP_IF_FALSE)
	c.context.GetBytecodeBuilder().AppendUint16(0xFFFF) // Placeholder

	// Body
	for _, bodyStmt := range stmt.Body {
		if err := c.stmtCompiler.CompileStmt(bodyStmt); err != nil {
			return err
		}
	}

	// Increment
	if stmt.Incr != nil {
		if err := c.exprCompiler.CompileExpr(stmt.Incr); err != nil {
			return err
		}
		c.context.GetBytecodeBuilder().Append(bytecode.OP_POP)
	}

	// Jump back to condition
	jumpBackOffset := conditionStartPos - (c.context.GetBytecodeBuilder().CurrentPosition() + 3)
	c.context.GetBytecodeBuilder().Append(bytecode.OP_JUMP)
	c.context.GetBytecodeBuilder().AppendInt16(int16(jumpBackOffset))

	// Patch the jump if false
	loopEndPos := c.context.GetBytecodeBuilder().CurrentPosition()
	jumpFalseOffset := uint16(loopEndPos - (jumpFalsePos + 3))
	c.context.GetBytecodeBuilder().PatchUint16(jumpFalsePos+1, jumpFalseOffset)

	return nil
}
