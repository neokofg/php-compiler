// Licensed under GNU GPL v3. See LICENSE file for details.
package stmt

import (
	"github.com/neokofg/php-compiler/internal/ast"
	"github.com/neokofg/php-compiler/internal/compiler/bytecode"
	"github.com/neokofg/php-compiler/internal/compiler/interfaces"
)

type FunctionCompiler struct {
	context      interfaces.CompilationContext
	stmtCompiler interfaces.StmtCompiler
}

func NewFunctionCompiler(context interfaces.CompilationContext, stmtCompiler interfaces.StmtCompiler) *FunctionCompiler {
	return &FunctionCompiler{
		context:      context,
		stmtCompiler: stmtCompiler,
	}
}

func (c *FunctionCompiler) Compile(stmt *ast.FunctionDecl) error {
	jumpPos := c.context.GetBytecodeBuilder().CurrentPosition()
	c.context.GetBytecodeBuilder().Append(bytecode.OP_JUMP)
	c.context.GetBytecodeBuilder().AppendUint16(0)

	funcStartAddr := c.context.GetBytecodeBuilder().CurrentPosition()
	stmt.StartAddr = funcStartAddr

	err := c.context.GetFunctionManager().AddFunction(stmt.Name, len(stmt.Params), funcStartAddr)
	if err != nil {
		return err
	}

	c.context.GetBytecodeBuilder().Append(bytecode.OP_FUNC_DECL)
	c.context.GetBytecodeBuilder().Append(byte(len(stmt.Params)))

	for _, param := range stmt.Params {
		varIdx := c.context.GetVariableManager().GetIndex(param)
		c.context.GetBytecodeBuilder().Append(byte(varIdx))
	}

	for _, bodyStmt := range stmt.Body {
		if err := c.stmtCompiler.CompileStmt(bodyStmt); err != nil {
			return err
		}
	}

	c.context.GetBytecodeBuilder().Append(bytecode.OP_EXIT_FUNC)

	endFuncPos := c.context.GetBytecodeBuilder().CurrentPosition()
	offset := uint16(endFuncPos - (jumpPos + 3))
	c.context.GetBytecodeBuilder().PatchUint16(jumpPos+1, offset)

	return nil
}
