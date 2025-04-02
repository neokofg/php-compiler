// Licensed under GNU GPL v3. See LICENSE file for details.
package function

import (
	"fmt"
)

type Function struct {
	Name       string
	ParamCount int
	Address    int
}

type Manager struct {
	functions map[string]Function
}

func NewManager() *Manager {
	return &Manager{
		functions: make(map[string]Function),
	}
}

func (m *Manager) AddFunction(name string, paramCount int, address int) error {
	if _, exists := m.functions[name]; exists {
		return fmt.Errorf("function '%s' already defined", name)
	}

	m.functions[name] = Function{
		Name:       name,
		ParamCount: paramCount,
		Address:    address,
	}

	return nil
}

func (m *Manager) GetFunction(name string) (Function, bool) {
	function, exists := m.functions[name]
	return function, exists
}
