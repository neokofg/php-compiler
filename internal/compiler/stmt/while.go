package stmt

import (
	"github.com/neokofg/php-compiler/internal/ast"
	"github.com/neokofg/php-compiler/internal/compiler/bytecode"
	"github.com/neokofg/php-compiler/internal/compiler/interfaces"
)

type WhileCompiler struct {
	context      interfaces.CompilationContext
	exprCompiler interfaces.ExprCompiler
	stmtCompiler interfaces.StmtCompiler
}

func NewWhileCompiler(context interfaces.CompilationContext, exprCompiler interfaces.ExprCompiler, stmtCompiler interfaces.StmtCompiler) *WhileCompiler {
	return &WhileCompiler{
		context:      context,
		exprCompiler: exprCompiler,
		stmtCompiler: stmtCompiler,
	}
}

func (c *WhileCompiler) Compile(stmt *ast.WhileStmt) error {
	loopStartPos := c.context.GetBytecodeBuilder().CurrentPosition()

	if err := c.exprCompiler.CompileExpr(stmt.Cond); err != nil {
		return err
	}

	jumpFalsePos := c.context.GetBytecodeBuilder().CurrentPosition()
	c.context.GetBytecodeBuilder().Append(bytecode.OP_JUMP_IF_FALSE)
	c.context.GetBytecodeBuilder().AppendUint16(0xFFFF) // Placeholder

	for _, bodyStmt := range stmt.Body {
		if err := c.stmtCompiler.CompileStmt(bodyStmt); err != nil {
			return err
		}
	}

	jumpBackOffset := loopStartPos - (c.context.GetBytecodeBuilder().CurrentPosition() + 3)
	c.context.GetBytecodeBuilder().Append(bytecode.OP_JUMP)
	c.context.GetBytecodeBuilder().AppendInt16(int16(jumpBackOffset))

	loopEndPos := c.context.GetBytecodeBuilder().CurrentPosition()
	jumpFalseOffset := uint16(loopEndPos - (jumpFalsePos + 3))
	c.context.GetBytecodeBuilder().PatchUint16(jumpFalsePos+1, jumpFalseOffset)

	return nil
}
