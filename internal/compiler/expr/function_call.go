// Licensed under GNU GPL v3. See LICENSE file for details.
package expr

import (
	"fmt"
	"github.com/neokofg/php-compiler/internal/ast"
	"github.com/neokofg/php-compiler/internal/compiler/bytecode"
	"github.com/neokofg/php-compiler/internal/compiler/interfaces"
)

type FunctionCallCompiler struct {
	context      interfaces.CompilationContext
	exprCompiler interfaces.ExprCompiler
}

func NewFunctionCallCompiler(context interfaces.CompilationContext, exprCompiler interfaces.ExprCompiler) *FunctionCallCompiler {
	return &FunctionCallCompiler{
		context:      context,
		exprCompiler: exprCompiler,
	}
}

func (c *FunctionCallCompiler) Compile(expr *ast.FunctionCall) error {
	function, exists := c.context.GetFunctionManager().GetFunction(expr.Name)
	if !exists {
		return fmt.Errorf("undefined function: %s", expr.Name)
	}

	if len(expr.Args) != function.ParamCount {
		return fmt.Errorf("function %s requires %d arguments, %d given",
			expr.Name, function.ParamCount, len(expr.Args))
	}

	for i := len(expr.Args) - 1; i >= 0; i-- {
		if err := c.exprCompiler.CompileExpr(expr.Args[i]); err != nil {
			return err
		}
	}

	c.context.GetBytecodeBuilder().Append(bytecode.OP_FUNC_CALL)
	c.context.GetBytecodeBuilder().Append(byte(len(expr.Args)))

	c.context.GetBytecodeBuilder().AppendUint16(uint16(function.Address))

	return nil
}
