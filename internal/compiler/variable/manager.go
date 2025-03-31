// Licensed under GNU GPL v3. See LICENSE file for details.
package variable

type Manager struct {
	variableMap  map[string]int
	nextVarIndex int
}

func NewManager() *Manager {
	return &Manager{
		variableMap:  make(map[string]int),
		nextVarIndex: 0,
	}
}

func (m *Manager) GetIndex(name string) int {
	if index, exists := m.variableMap[name]; exists {
		return index
	}

	index := m.nextVarIndex
	m.variableMap[name] = index
	m.nextVarIndex++
	return index
}

func (m *Manager) GetAllVariables() map[string]int {
	return m.variableMap
}
