// Licensed under GNU GPL v3. See LICENSE file for details.
package expr

import (
	"fmt"
	"github.com/neokofg/php-compiler/internal/ast"
	"github.com/neokofg/php-compiler/internal/compiler/bytecode"
	"github.com/neokofg/php-compiler/internal/compiler/interfaces"
	"github.com/neokofg/php-compiler/internal/token"
)

type PrefixCompiler struct {
	context interfaces.CompilationContext
}

func NewPrefixCompiler(context interfaces.CompilationContext) *PrefixCompiler {
	return &PrefixCompiler{
		context: context,
	}
}

func (c *PrefixCompiler) Compile(expr *ast.PrefixExpr) error {
	varExpr, ok := expr.Expr.(*ast.VarExpr)
	if !ok {
		return fmt.Errorf("can only apply prefix operators to variables")
	}

	varIdx := c.context.GetVariableManager().GetIndex(varExpr.Name)

	c.context.GetBytecodeBuilder().Append(bytecode.OP_LOAD_VAR)
	c.context.GetBytecodeBuilder().Append(byte(varIdx))

	switch expr.Op {
	case token.T_INC:
		c.context.GetBytecodeBuilder().Append(bytecode.OP_INC)
	case token.T_DEC:
		c.context.GetBytecodeBuilder().Append(bytecode.OP_DEC)
	}

	c.context.GetBytecodeBuilder().Append(bytecode.OP_STORE_VAR)
	c.context.GetBytecodeBuilder().Append(byte(varIdx))

	return nil
}
