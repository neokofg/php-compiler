// Licensed under GNU GPL v3. See LICENSE file for details.
package stmt

import (
	"github.com/neokofg/php-compiler/internal/ast"
	"github.com/neokofg/php-compiler/internal/compiler/bytecode"
	"github.com/neokofg/php-compiler/internal/compiler/interfaces"
)

type DoWhileCompiler struct {
	context      interfaces.CompilationContext
	exprCompiler interfaces.ExprCompiler
	stmtCompiler interfaces.StmtCompiler
}

func NewDoWhileCompiler(context interfaces.CompilationContext, exprCompiler interfaces.ExprCompiler, stmtCompiler interfaces.StmtCompiler) *DoWhileCompiler {
	return &DoWhileCompiler{
		context:      context,
		exprCompiler: exprCompiler,
		stmtCompiler: stmtCompiler,
	}
}

func (c *DoWhileCompiler) Compile(stmt *ast.DoWhileStmt) error {
	loop := c.context.EnterLoop()
	defer c.context.ExitLoop()

	loopBodyStart := c.context.GetBytecodeBuilder().CurrentPosition()
	loop.StartPos = loopBodyStart

	for _, bodyStmt := range stmt.Body {
		if err := c.stmtCompiler.CompileStmt(bodyStmt); err != nil {
			return err
		}
	}

	conditionPos := c.context.GetBytecodeBuilder().CurrentPosition()
	loop.ConditionPos = conditionPos

	if err := c.exprCompiler.CompileExpr(stmt.Cond); err != nil {
		return err
	}

	c.context.GetBytecodeBuilder().Append(bytecode.OP_JUMP)
	jumpBackOffset := int16(loopBodyStart - (c.context.GetBytecodeBuilder().CurrentPosition() + 3))
	c.context.GetBytecodeBuilder().AppendInt16(jumpBackOffset)

	loop.EndPos = c.context.GetBytecodeBuilder().CurrentPosition()

	return nil
}
