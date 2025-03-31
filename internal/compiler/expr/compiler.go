package expr

import (
	"fmt"
	"github.com/neokofg/php-compiler/internal/ast"
	"github.com/neokofg/php-compiler/internal/compiler/interfaces"
)

type exprCompiler struct {
	context        interfaces.CompilationContext
	numberCompiler *NumberCompiler
	stringCompiler *StringCompiler
	varCompiler    *VarCompiler
	binaryCompiler *BinaryCompiler
}

func NewCompiler(context interfaces.CompilationContext) interfaces.ExprCompiler {
	compiler := &exprCompiler{
		context: context,
	}

	// Инициализируем отдельные компиляторы для каждого типа выражений
	compiler.numberCompiler = NewNumberCompiler(context)
	compiler.stringCompiler = NewStringCompiler(context)
	compiler.varCompiler = NewVarCompiler(context)

	// Для BinaryCompiler нужна циклическая ссылка на сам exprCompiler
	compiler.binaryCompiler = NewBinaryCompiler(context, compiler)

	return compiler
}

func (c *exprCompiler) CompileExpr(expr ast.Expr) error {
	switch e := expr.(type) {
	case *ast.NumberLiteral:
		return c.numberCompiler.Compile(e)
	case *ast.StringLiteral:
		return c.stringCompiler.Compile(e)
	case *ast.VarExpr:
		return c.varCompiler.Compile(e)
	case *ast.BinaryExpr:
		return c.binaryCompiler.Compile(e)
	default:
		return fmt.Errorf("unsupported expression type: %T", expr)
	}
}
