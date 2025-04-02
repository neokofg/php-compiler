// Licensed under GNU GPL v3. See LICENSE file for details.
package stmt

import (
	"github.com/neokofg/php-compiler/internal/ast"
	"github.com/neokofg/php-compiler/internal/compiler/bytecode"
	"github.com/neokofg/php-compiler/internal/compiler/constant"
	"github.com/neokofg/php-compiler/internal/compiler/interfaces"
)

type ReturnCompiler struct {
	context      interfaces.CompilationContext
	exprCompiler interfaces.ExprCompiler
}

func NewReturnCompiler(context interfaces.CompilationContext, exprCompiler interfaces.ExprCompiler) *ReturnCompiler {
	return &ReturnCompiler{
		context:      context,
		exprCompiler: exprCompiler,
	}
}

func (c *ReturnCompiler) Compile(stmt *ast.ReturnStmt) error {
	if stmt.Expr != nil {
		if err := c.exprCompiler.CompileExpr(stmt.Expr); err != nil {
			return err
		}
	} else {
		nullIdx := c.context.GetConstantPool().Add(constant.Constant{
			Type:  "int",
			Value: "0",
		})
		c.context.GetBytecodeBuilder().Append(bytecode.OP_LOAD_CONST)
		c.context.GetBytecodeBuilder().Append(byte(nullIdx))
	}

	c.context.GetBytecodeBuilder().Append(bytecode.OP_RETURN)

	return nil
}
