// Licensed under GNU GPL v3. See LICENSE file for details.
package expr

import (
	"fmt"
	"github.com/neokofg/php-compiler/internal/ast"
	"github.com/neokofg/php-compiler/internal/compiler/bytecode"
	"github.com/neokofg/php-compiler/internal/compiler/interfaces"
	"github.com/neokofg/php-compiler/internal/token"
)

type UnaryCompiler struct {
	context      interfaces.CompilationContext
	exprCompiler interfaces.ExprCompiler
}

func NewUnaryCompiler(context interfaces.CompilationContext, exprCompiler interfaces.ExprCompiler) *UnaryCompiler {
	return &UnaryCompiler{
		context:      context,
		exprCompiler: exprCompiler,
	}
}

func (c *UnaryCompiler) Compile(expr *ast.UnaryExpr) error {
	if err := c.exprCompiler.CompileExpr(expr.Expr); err != nil {
		return err
	}

	switch expr.Op {
	case token.T_NOT:
		c.context.GetBytecodeBuilder().Append(bytecode.OP_NOT)
	default:
		return fmt.Errorf("unsupported unary operator: %v", expr.Op)
	}

	return nil
}
