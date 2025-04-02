// Licensed under GNU GPL v3. See LICENSE file for details.
package stmt

import (
	"fmt"
	"github.com/neokofg/php-compiler/internal/ast"
	"github.com/neokofg/php-compiler/internal/compiler/bytecode"
	"github.com/neokofg/php-compiler/internal/compiler/interfaces"
)

type FunctionCallStmtCompiler struct {
	context      interfaces.CompilationContext
	exprCompiler interfaces.ExprCompiler
}

func NewFunctionCallStmtCompiler(context interfaces.CompilationContext, exprCompiler interfaces.ExprCompiler) *FunctionCallStmtCompiler {
	return &FunctionCallStmtCompiler{
		context:      context,
		exprCompiler: exprCompiler,
	}
}

func (c *FunctionCallStmtCompiler) Compile(stmt *ast.FunctionCallStmt) error {
	function, exists := c.context.GetFunctionManager().GetFunction(stmt.Call.Name)
	if !exists {
		return fmt.Errorf("undefined function: %s", stmt.Call.Name)
	}

	if len(stmt.Call.Args) != function.ParamCount {
		return fmt.Errorf("function %s requires %d arguments, %d given",
			stmt.Call.Name, function.ParamCount, len(stmt.Call.Args))
	}

	for i := len(stmt.Call.Args) - 1; i >= 0; i-- {
		if err := c.exprCompiler.CompileExpr(stmt.Call.Args[i]); err != nil {
			return err
		}
	}

	c.context.GetBytecodeBuilder().Append(bytecode.OP_FUNC_CALL)
	c.context.GetBytecodeBuilder().Append(byte(len(stmt.Call.Args)))

	c.context.GetBytecodeBuilder().AppendUint16(uint16(function.Address))

	c.context.GetBytecodeBuilder().Append(bytecode.OP_POP)

	return nil
}
