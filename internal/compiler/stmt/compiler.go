package stmt

import (
	"fmt"
	"github.com/neokofg/php-compiler/internal/ast"
	"github.com/neokofg/php-compiler/internal/compiler/interfaces"
)

type stmtCompiler struct {
	context        interfaces.CompilationContext
	exprCompiler   interfaces.ExprCompiler
	assignCompiler *AssignCompiler
	echoCompiler   *EchoCompiler
	ifCompiler     *IfCompiler
	whileCompiler  *WhileCompiler
	forCompiler    *ForCompiler
}

func NewCompiler(context interfaces.CompilationContext, exprCompiler interfaces.ExprCompiler) interfaces.StmtCompiler {
	compiler := &stmtCompiler{
		context:      context,
		exprCompiler: exprCompiler,
	}

	// Инициализируем отдельные компиляторы для каждого типа утверждений
	compiler.assignCompiler = NewAssignCompiler(context, exprCompiler)
	compiler.echoCompiler = NewEchoCompiler(context, exprCompiler)

	// Эти компиляторы требуют циклической ссылки на сам stmtCompiler
	compiler.ifCompiler = NewIfCompiler(context, exprCompiler, compiler)
	compiler.whileCompiler = NewWhileCompiler(context, exprCompiler, compiler)
	compiler.forCompiler = NewForCompiler(context, exprCompiler, compiler)

	return compiler
}

func (c *stmtCompiler) CompileStmt(stmt ast.Stmt) error {
	switch s := stmt.(type) {
	case *ast.AssignStmt:
		return c.assignCompiler.Compile(s)
	case *ast.EchoStmt:
		return c.echoCompiler.Compile(s)
	case *ast.IfStmt:
		return c.ifCompiler.Compile(s)
	case *ast.WhileStmt:
		return c.whileCompiler.Compile(s)
	case *ast.ForStmt:
		return c.forCompiler.Compile(s)
	default:
		return fmt.Errorf("unsupported statement type: %T", stmt)
	}
}
