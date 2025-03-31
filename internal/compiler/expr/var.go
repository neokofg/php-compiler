// Licensed under GNU GPL v3. See LICENSE file for details.
package expr

import (
	"github.com/neokofg/php-compiler/internal/ast"
	"github.com/neokofg/php-compiler/internal/compiler/bytecode"
	"github.com/neokofg/php-compiler/internal/compiler/interfaces"
)

type VarCompiler struct {
	context interfaces.CompilationContext
}

func NewVarCompiler(context interfaces.CompilationContext) *VarCompiler {
	return &VarCompiler{
		context: context,
	}
}

func (c *VarCompiler) Compile(expr *ast.VarExpr) error {
	varIdx := c.context.GetVariableManager().GetIndex(expr.Name)

	c.context.GetBytecodeBuilder().Append(bytecode.OP_LOAD_VAR)
	c.context.GetBytecodeBuilder().Append(byte(varIdx))
	return nil
}
