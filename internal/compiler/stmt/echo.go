package stmt

import (
	"github.com/neokofg/php-compiler/internal/ast"
	"github.com/neokofg/php-compiler/internal/compiler/bytecode"
	"github.com/neokofg/php-compiler/internal/compiler/interfaces"
)

type EchoCompiler struct {
	context      interfaces.CompilationContext
	exprCompiler interfaces.ExprCompiler
}

func NewEchoCompiler(context interfaces.CompilationContext, exprCompiler interfaces.ExprCompiler) *EchoCompiler {
	return &EchoCompiler{
		context:      context,
		exprCompiler: exprCompiler,
	}
}

func (c *EchoCompiler) Compile(stmt *ast.EchoStmt) error {
	if err := c.exprCompiler.CompileExpr(stmt.Expr); err != nil {
		return err
	}

	c.context.GetBytecodeBuilder().Append(bytecode.OP_PRINT)
	return nil
}
