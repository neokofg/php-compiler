// Licensed under GNU GPL v3. See LICENSE file for details.
package expr

import (
	"fmt"
	"github.com/neokofg/php-compiler/internal/ast"
	"github.com/neokofg/php-compiler/internal/compiler/bytecode"
	"github.com/neokofg/php-compiler/internal/compiler/interfaces"
	"github.com/neokofg/php-compiler/internal/token"
)

type BinaryCompiler struct {
	context      interfaces.CompilationContext
	exprCompiler interfaces.ExprCompiler
}

func NewBinaryCompiler(context interfaces.CompilationContext, exprCompiler interfaces.ExprCompiler) *BinaryCompiler {
	return &BinaryCompiler{
		context:      context,
		exprCompiler: exprCompiler,
	}
}

func (c *BinaryCompiler) Compile(expr *ast.BinaryExpr) error {
	if err := c.exprCompiler.CompileExpr(expr.Left); err != nil {
		return err
	}

	if err := c.exprCompiler.CompileExpr(expr.Right); err != nil {
		return err
	}

	switch expr.Op {
	case token.T_PLUS:
		c.context.GetBytecodeBuilder().Append(bytecode.OP_ADD)
	case token.T_MINUS:
		c.context.GetBytecodeBuilder().Append(bytecode.OP_SUB)
	case token.T_STAR:
		c.context.GetBytecodeBuilder().Append(bytecode.OP_MUL)
	case token.T_SLASH:
		c.context.GetBytecodeBuilder().Append(bytecode.OP_DIV)
	case token.T_MOD:
		c.context.GetBytecodeBuilder().Append(bytecode.OP_MOD)
	case token.T_GT:
		c.context.GetBytecodeBuilder().Append(bytecode.OP_GT)
	case token.T_LT:
		c.context.GetBytecodeBuilder().Append(bytecode.OP_LT)
	case token.T_GTE:
		c.context.GetBytecodeBuilder().Append(bytecode.OP_GTE)
	case token.T_LTE:
		c.context.GetBytecodeBuilder().Append(bytecode.OP_LTE)
	case token.T_EQEQ:
		c.context.GetBytecodeBuilder().Append(bytecode.OP_EQ)
	case token.T_EQEQEQ:
		c.context.GetBytecodeBuilder().Append(bytecode.OP_IDENTITY_EQ)
	case token.T_NOTEQ:
		c.context.GetBytecodeBuilder().Append(bytecode.OP_EQ)
		c.context.GetBytecodeBuilder().Append(bytecode.OP_NOT)
	case token.T_NOTEQEQ:
		c.context.GetBytecodeBuilder().Append(bytecode.OP_IDENTITY_NE)
	case token.T_AND:
		c.context.GetBytecodeBuilder().Append(bytecode.OP_AND)
	case token.T_OR:
		c.context.GetBytecodeBuilder().Append(bytecode.OP_OR)
	case token.T_BIT_AND:
		c.context.GetBytecodeBuilder().Append(bytecode.OP_BIT_AND)
	case token.T_BIT_OR:
		c.context.GetBytecodeBuilder().Append(bytecode.OP_BIT_OR)
	case token.T_BIT_XOR:
		c.context.GetBytecodeBuilder().Append(bytecode.OP_BIT_XOR)
	case token.T_LSHIFT:
		c.context.GetBytecodeBuilder().Append(bytecode.OP_LSHIFT)
	case token.T_RSHIFT:
		c.context.GetBytecodeBuilder().Append(bytecode.OP_RSHIFT)
	case token.T_DOT:
		c.context.GetBytecodeBuilder().Append(bytecode.OP_CONCAT)
	default:
		return fmt.Errorf("unsupported binary operator: %v", expr.Op)
	}

	return nil
}
