// Licensed under GNU GPL v3. See LICENSE file for details.
package interfaces

import (
	"github.com/neokofg/php-compiler/internal/ast"
	"github.com/neokofg/php-compiler/internal/compiler/bytecode"
	"github.com/neokofg/php-compiler/internal/compiler/constant"
	"github.com/neokofg/php-compiler/internal/compiler/function"
	"github.com/neokofg/php-compiler/internal/compiler/variable"
)

type ExprCompiler interface {
	CompileExpr(expr ast.Expr) error
}

type StmtCompiler interface {
	CompileStmt(stmt ast.Stmt) error
}

type CompilationContext interface {
	GetBytecodeBuilder() *bytecode.BytecodeBuilder
	GetConstantPool() *constant.Pool
	GetVariableManager() *variable.Manager

	EnterLoop() *LoopContext
	ExitLoop()
	GetCurrentLoop() *LoopContext

	AddPendingJump(position int, isBreak bool)
	ApplyPendingJumps()

	GetFunctionManager() *function.Manager
}

type JumpPatch struct {
	Position int
	IsBreak  bool
}

type LoopContext struct {
	StartPos     int
	ConditionPos int
	EndPos       int
	Parent       *LoopContext
	PendingJumps []JumpPatch
}

type Context struct {
	BytecodeBuilder *bytecode.BytecodeBuilder
	ConstantPool    *constant.Pool
	VariableManager *variable.Manager
	CurrentLoop     *LoopContext
	FunctionManager *function.Manager
}

func NewContext() *Context {
	return &Context{
		BytecodeBuilder: bytecode.NewBytecodeBuilder(),
		ConstantPool:    constant.NewPool(),
		VariableManager: variable.NewManager(),
		CurrentLoop:     nil,
		FunctionManager: function.NewManager(),
	}
}

func (c *Context) GetBytecodeBuilder() *bytecode.BytecodeBuilder {
	return c.BytecodeBuilder
}

func (c *Context) GetConstantPool() *constant.Pool {
	return c.ConstantPool
}

func (c *Context) GetVariableManager() *variable.Manager {
	return c.VariableManager
}

func (c *Context) GetCurrentLoop() *LoopContext {
	return c.CurrentLoop
}

func (c *Context) GetFunctionManager() *function.Manager {
	return c.FunctionManager
}

func (c *Context) EnterLoop() *LoopContext {
	loop := &LoopContext{
		Parent:       c.CurrentLoop,
		PendingJumps: make([]JumpPatch, 0),
	}
	c.CurrentLoop = loop
	return loop
}

func (c *Context) ExitLoop() {
	if c.CurrentLoop != nil {
		c.CurrentLoop = c.CurrentLoop.Parent
	}
}

func (c *Context) AddPendingJump(position int, isBreak bool) {
	if c.CurrentLoop != nil {
		c.CurrentLoop.PendingJumps = append(c.CurrentLoop.PendingJumps, JumpPatch{
			Position: position,
			IsBreak:  isBreak,
		})
	}
}

func (c *Context) ApplyPendingJumps() {
	if c.CurrentLoop == nil {
		return
	}

	for _, patch := range c.CurrentLoop.PendingJumps {
		targetPos := c.CurrentLoop.ConditionPos
		if patch.IsBreak {
			targetPos = c.CurrentLoop.EndPos
		}

		if targetPos > 0 {
			offset := uint16(targetPos - (patch.Position + 3))
			c.BytecodeBuilder.PatchUint16(patch.Position+1, offset)
		}
	}
}
