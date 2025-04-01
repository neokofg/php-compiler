// Licensed under GNU GPL v3. See LICENSE file for details.
package expr

import (
	"fmt"
	"github.com/neokofg/php-compiler/internal/ast"
	"github.com/neokofg/php-compiler/internal/compiler/interfaces"
)

type exprCompiler struct {
	context         interfaces.CompilationContext
	numberCompiler  *NumberCompiler
	stringCompiler  *StringCompiler
	booleanCompiler *BooleanCompiler
	varCompiler     *VarCompiler
	postfixCompiler *PostfixCompiler
	prefixCompiler  *PrefixCompiler
	binaryCompiler  *BinaryCompiler
	unaryCompiler   *UnaryCompiler
}

func NewCompiler(context interfaces.CompilationContext) interfaces.ExprCompiler {
	compiler := &exprCompiler{
		context: context,
	}

	compiler.numberCompiler = NewNumberCompiler(context)
	compiler.stringCompiler = NewStringCompiler(context)
	compiler.booleanCompiler = NewBooleanCompiler(context)
	compiler.varCompiler = NewVarCompiler(context)
	compiler.postfixCompiler = NewPostfixCompiler(context)
	compiler.prefixCompiler = NewPrefixCompiler(context)

	compiler.unaryCompiler = NewUnaryCompiler(context, compiler)
	compiler.binaryCompiler = NewBinaryCompiler(context, compiler)

	return compiler
}

func (c *exprCompiler) CompileExpr(expr ast.Expr) error {
	switch e := expr.(type) {
	case *ast.NumberLiteral:
		return c.numberCompiler.Compile(e)
	case *ast.StringLiteral:
		return c.stringCompiler.Compile(e)
	case *ast.BooleanLiteral:
		return c.booleanCompiler.Compile(e)
	case *ast.VarExpr:
		return c.varCompiler.Compile(e)
	case *ast.UnaryExpr:
		return c.unaryCompiler.Compile(e)
	case *ast.BinaryExpr:
		return c.binaryCompiler.Compile(e)
	case *ast.PostfixExpr:
		return c.postfixCompiler.Compile(e)
	case *ast.PrefixExpr:
		return c.prefixCompiler.Compile(e)
	default:
		return fmt.Errorf("unsupported expression type: %T", expr)
	}
}
