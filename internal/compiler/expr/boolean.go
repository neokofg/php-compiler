// Licensed under GNU GPL v3. See LICENSE file for details.
package expr

import (
	"github.com/neokofg/php-compiler/internal/ast"
	"github.com/neokofg/php-compiler/internal/compiler/bytecode"
	"github.com/neokofg/php-compiler/internal/compiler/constant"
	"github.com/neokofg/php-compiler/internal/compiler/interfaces"
)

type BooleanCompiler struct {
	context interfaces.CompilationContext
}

func NewBooleanCompiler(context interfaces.CompilationContext) *BooleanCompiler {
	return &BooleanCompiler{
		context: context,
	}
}

func (c *BooleanCompiler) Compile(expr *ast.BooleanLiteral) error {
	value := "0"
	if expr.Value {
		value = "1"
	}

	idx := c.context.GetConstantPool().Add(constant.Constant{
		Type:  "int",
		Value: value,
	})

	c.context.GetBytecodeBuilder().Append(bytecode.OP_LOAD_CONST)
	c.context.GetBytecodeBuilder().Append(byte(idx))
	return nil
}
