package interfaces

import (
	"github.com/neokofg/php-compiler/internal/ast"
	"github.com/neokofg/php-compiler/internal/compiler/bytecode"
	"github.com/neokofg/php-compiler/internal/compiler/constant"
	"github.com/neokofg/php-compiler/internal/compiler/variable"
)

// Интерфейсы компиляторов
type ExprCompiler interface {
	CompileExpr(expr ast.Expr) error
}

type StmtCompiler interface {
	CompileStmt(stmt ast.Stmt) error
}

// CompilationContext предоставляет контекст для компиляции
type CompilationContext interface {
	GetBytecodeBuilder() *bytecode.BytecodeBuilder
	GetConstantPool() *constant.Pool
	GetVariableManager() *variable.Manager
}

// Базовый контекст компиляции
type Context struct {
	BytecodeBuilder *bytecode.BytecodeBuilder
	ConstantPool    *constant.Pool
	VariableManager *variable.Manager
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

// Функция-фабрика для создания контекста
func NewContext() *Context {
	return &Context{
		BytecodeBuilder: bytecode.NewBytecodeBuilder(),
		ConstantPool:    constant.NewPool(),
		VariableManager: variable.NewManager(),
	}
}
