package stmt

import (
	"github.com/neokofg/php-compiler/internal/ast"
	"github.com/neokofg/php-compiler/internal/compiler/bytecode"
	"github.com/neokofg/php-compiler/internal/compiler/interfaces"
)

type IfCompiler struct {
	context      interfaces.CompilationContext
	exprCompiler interfaces.ExprCompiler
	stmtCompiler interfaces.StmtCompiler
}

func NewIfCompiler(context interfaces.CompilationContext, exprCompiler interfaces.ExprCompiler, stmtCompiler interfaces.StmtCompiler) *IfCompiler {
	return &IfCompiler{
		context:      context,
		exprCompiler: exprCompiler,
		stmtCompiler: stmtCompiler,
	}
}

func (c *IfCompiler) Compile(stmt *ast.IfStmt) error {
	if err := c.exprCompiler.CompileExpr(stmt.Cond); err != nil {
		return err
	}

	jumpIfFalsePos := c.context.GetBytecodeBuilder().CurrentPosition()
	c.context.GetBytecodeBuilder().Append(bytecode.OP_JUMP_IF_FALSE)
	c.context.GetBytecodeBuilder().AppendUint16(0xFFFF) // Placeholder

	// Compile THEN block
	for _, thenStmt := range stmt.Then {
		if err := c.stmtCompiler.CompileStmt(thenStmt); err != nil {
			return err
		}
	}

	if len(stmt.Else) > 0 {
		jumpOverElsePos := c.context.GetBytecodeBuilder().CurrentPosition()
		c.context.GetBytecodeBuilder().Append(bytecode.OP_JUMP)
		c.context.GetBytecodeBuilder().AppendUint16(0xFFFF) // Placeholder

		elseStart := c.context.GetBytecodeBuilder().CurrentPosition()

		// Compile ELSE block
		for _, elseStmt := range stmt.Else {
			if err := c.stmtCompiler.CompileStmt(elseStmt); err != nil {
				return err
			}
		}

		endElse := c.context.GetBytecodeBuilder().CurrentPosition()

		// Patch offsets
		offsetJumpIfFalse := uint16(elseStart - (jumpIfFalsePos + 3))
		c.context.GetBytecodeBuilder().PatchUint16(jumpIfFalsePos+1, offsetJumpIfFalse)

		offsetJumpOverElse := uint16(endElse - (jumpOverElsePos + 3))
		c.context.GetBytecodeBuilder().PatchUint16(jumpOverElsePos+1, offsetJumpOverElse)
	} else {
		afterThen := c.context.GetBytecodeBuilder().CurrentPosition()
		offsetJumpIfFalse := uint16(afterThen - (jumpIfFalsePos + 3))
		c.context.GetBytecodeBuilder().PatchUint16(jumpIfFalsePos+1, offsetJumpIfFalse)
	}

	return nil
}
