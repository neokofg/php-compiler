package expr

import (
	"fmt"
	"github.com/neokofg/php-compiler/internal/ast"
	"github.com/neokofg/php-compiler/internal/compiler/bytecode"
	"github.com/neokofg/php-compiler/internal/compiler/constant"
	"github.com/neokofg/php-compiler/internal/compiler/interfaces"
)

type NumberCompiler struct {
	context interfaces.CompilationContext
}

func NewNumberCompiler(context interfaces.CompilationContext) *NumberCompiler {
	return &NumberCompiler{
		context: context,
	}
}

func (c *NumberCompiler) Compile(expr *ast.NumberLiteral) error {
	idx := c.context.GetConstantPool().Add(constant.Constant{
		Type:  "int",
		Value: fmt.Sprint(expr.Value),
	})

	c.context.GetBytecodeBuilder().Append(bytecode.OP_LOAD_CONST)
	c.context.GetBytecodeBuilder().Append(byte(idx))
	return nil
}
