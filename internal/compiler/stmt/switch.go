// Licensed under GNU GPL v3. See LICENSE file for details.
package stmt

import (
	"github.com/neokofg/php-compiler/internal/ast"
	"github.com/neokofg/php-compiler/internal/compiler/bytecode"
	"github.com/neokofg/php-compiler/internal/compiler/interfaces"
)

type SwitchCompiler struct {
	context      interfaces.CompilationContext
	exprCompiler interfaces.ExprCompiler
	stmtCompiler interfaces.StmtCompiler
}

func NewSwitchCompiler(context interfaces.CompilationContext, exprCompiler interfaces.ExprCompiler, stmtCompiler interfaces.StmtCompiler) *SwitchCompiler {
	return &SwitchCompiler{
		context:      context,
		exprCompiler: exprCompiler,
		stmtCompiler: stmtCompiler,
	}
}

func (c *SwitchCompiler) Compile(stmt *ast.SwitchStmt) error {
	loop := c.context.EnterLoop()
	defer c.context.ExitLoop()

	if err := c.exprCompiler.CompileExpr(stmt.Expr); err != nil {
		return err
	}

	switchVarIdx := c.context.GetVariableManager().GetIndex("__switch_value")
	c.context.GetBytecodeBuilder().Append(bytecode.OP_STORE_VAR)
	c.context.GetBytecodeBuilder().Append(byte(switchVarIdx))

	defaultCaseIndex := -1
	for i, caseStmt := range stmt.Cases {
		if caseStmt.Expr == nil {
			defaultCaseIndex = i
			break
		}
	}

	endSwitchJumps := make([]int, 0, len(stmt.Cases))

	for i, caseStmt := range stmt.Cases {
		if caseStmt.Expr == nil {
			continue
		}

		c.context.GetBytecodeBuilder().Append(bytecode.OP_LOAD_VAR)
		c.context.GetBytecodeBuilder().Append(byte(switchVarIdx))

		if err := c.exprCompiler.CompileExpr(caseStmt.Expr); err != nil {
			return err
		}

		c.context.GetBytecodeBuilder().Append(bytecode.OP_EQ)

		c.context.GetBytecodeBuilder().Append(bytecode.OP_JUMP_IF_FALSE)

		ifFalseJumpPos := c.context.GetBytecodeBuilder().CurrentPosition()
		c.context.GetBytecodeBuilder().AppendUint16(5)

		for _, s := range caseStmt.Stmts {
			if err := c.stmtCompiler.CompileStmt(s); err != nil {
				return err
			}
		}

		if i < len(stmt.Cases)-1 || !endsWithBreak(caseStmt.Stmts) {
			c.context.GetBytecodeBuilder().Append(bytecode.OP_JUMP)
			endJumpPos := c.context.GetBytecodeBuilder().CurrentPosition()
			c.context.GetBytecodeBuilder().AppendUint16(0)
			endSwitchJumps = append(endSwitchJumps, endJumpPos)
		}

		nextCasePos := c.context.GetBytecodeBuilder().CurrentPosition()
		jumpOffset := nextCasePos - (ifFalseJumpPos + 2)

		if jumpOffset <= 0 || jumpOffset > 255 {
			c.context.GetBytecodeBuilder().PatchUint16(ifFalseJumpPos, uint16(jumpOffset))
		} else {
			c.context.GetBytecodeBuilder().PatchUint16(ifFalseJumpPos, uint16(jumpOffset))
		}
	}

	if defaultCaseIndex >= 0 {
		defaultCase := stmt.Cases[defaultCaseIndex]
		for _, s := range defaultCase.Stmts {
			if err := c.stmtCompiler.CompileStmt(s); err != nil {
				return err
			}
		}
	}

	endSwitchPos := c.context.GetBytecodeBuilder().CurrentPosition()
	loop.EndPos = endSwitchPos

	for _, jumpPos := range endSwitchJumps {
		offset := endSwitchPos - (jumpPos + 2)
		if offset <= 0 || offset > 65535 {
			offset = 2
		}
		c.context.GetBytecodeBuilder().PatchUint16(jumpPos, uint16(offset))
	}

	c.context.ApplyPendingJumps()

	return nil
}

func endsWithBreak(stmts []ast.Stmt) bool {
	if len(stmts) == 0 {
		return false
	}

	_, isBreak := stmts[len(stmts)-1].(*ast.BreakStmt)
	return isBreak
}
