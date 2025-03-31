package expr

import (
	"github.com/neokofg/php-compiler/internal/ast"
	"github.com/neokofg/php-compiler/internal/compiler/bytecode"
	"github.com/neokofg/php-compiler/internal/compiler/constant"
	"github.com/neokofg/php-compiler/internal/compiler/interfaces"
)

type StringCompiler struct {
	context interfaces.CompilationContext
}

func NewStringCompiler(context interfaces.CompilationContext) *StringCompiler {
	return &StringCompiler{
		context: context,
	}
}

func (c *StringCompiler) Compile(expr *ast.StringLiteral) error {
	idx := c.context.GetConstantPool().Add(constant.Constant{
		Type:  "string",
		Value: expr.Value,
	})

	c.context.GetBytecodeBuilder().Append(bytecode.OP_LOAD_CONST)
	c.context.GetBytecodeBuilder().Append(byte(idx))
	return nil
}
