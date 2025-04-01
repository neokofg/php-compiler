// Licensed under GNU GPL v3. See LICENSE file for details.
package expr

import (
	"fmt"
	"github.com/neokofg/php-compiler/internal/ast"
	"github.com/neokofg/php-compiler/internal/compiler/bytecode"
	"github.com/neokofg/php-compiler/internal/compiler/interfaces"
	"github.com/neokofg/php-compiler/internal/token"
)

type PostfixCompiler struct {
	context interfaces.CompilationContext
}

func NewPostfixCompiler(context interfaces.CompilationContext) *PostfixCompiler {
	return &PostfixCompiler{
		context: context,
	}
}

func (c *PostfixCompiler) Compile(expr *ast.PostfixExpr) error {
	varExpr, ok := expr.Expr.(*ast.VarExpr)
	if !ok {
		return fmt.Errorf("can only apply postfix operators to variables")
	}

	varIdx := c.context.GetVariableManager().GetIndex(varExpr.Name)

	c.context.GetBytecodeBuilder().Append(bytecode.OP_LOAD_VAR)
	c.context.GetBytecodeBuilder().Append(byte(varIdx))

	switch expr.Op {
	case token.T_INC:
		c.context.GetBytecodeBuilder().Append(bytecode.OP_POST_INC)
	case token.T_DEC:
		c.context.GetBytecodeBuilder().Append(bytecode.OP_POST_DEC)
	}

	c.context.GetBytecodeBuilder().Append(bytecode.OP_STORE_VAR)
	c.context.GetBytecodeBuilder().Append(byte(varIdx))

	return nil
}
