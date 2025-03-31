package stmt

import (
	"github.com/neokofg/php-compiler/internal/ast"
	"github.com/neokofg/php-compiler/internal/compiler/bytecode"
	"github.com/neokofg/php-compiler/internal/compiler/interfaces"
)

type AssignCompiler struct {
	context      interfaces.CompilationContext
	exprCompiler interfaces.ExprCompiler
}

func NewAssignCompiler(context interfaces.CompilationContext, exprCompiler interfaces.ExprCompiler) *AssignCompiler {
	return &AssignCompiler{
		context:      context,
		exprCompiler: exprCompiler,
	}
}

func (a *AssignCompiler) Compile(stmt *ast.AssignStmt) error {
	if err := a.exprCompiler.CompileExpr(stmt.Expr); err != nil {
		return err
	}

	varIdx := a.context.GetVariableManager().GetIndex(stmt.Name)
	a.context.GetBytecodeBuilder().Append(bytecode.OP_STORE_VAR)
	a.context.GetBytecodeBuilder().Append(byte(varIdx))

	return nil
}
