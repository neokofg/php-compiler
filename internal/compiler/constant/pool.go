package constant

type Constant struct {
	Type  string
	Value string
}

type Pool struct {
	constants []Constant
	onUpdate  func([]Constant)
}

func NewPool() *Pool {
	return &Pool{
		constants: make([]Constant, 0, 32),
	}
}

func (p *Pool) Add(constant Constant) int {
	for i, existing := range p.constants {
		if existing.Type == constant.Type && existing.Value == constant.Value {
			return i
		}
	}

	p.constants = append(p.constants, constant)
	index := len(p.constants) - 1

	p.notifyUpdate()
	return index
}

func (p *Pool) GetAll() []Constant {
	return p.constants
}

func (p *Pool) SetSyncCallback(callback func([]Constant)) {
	p.onUpdate = callback
}

func (p *Pool) notifyUpdate() {
	if p.onUpdate != nil {
		p.onUpdate(p.constants)
	}
}
