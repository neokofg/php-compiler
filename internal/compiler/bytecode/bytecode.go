// Licensed under GNU GPL v3. See LICENSE file for details.
package bytecode

type BytecodeBuilder struct {
	code     []byte
	onUpdate func([]byte)
}

func NewBytecodeBuilder() *BytecodeBuilder {
	return &BytecodeBuilder{
		code: make([]byte, 0, 64),
	}
}

func (b *BytecodeBuilder) Append(value byte) {
	b.code = append(b.code, value)
	b.notifyUpdate()
}

func (b *BytecodeBuilder) AppendUint16(value uint16) {
	lowByte := byte(value & 0xFF)
	highByte := byte(value >> 8)
	b.code = append(b.code, lowByte, highByte)
	b.notifyUpdate()
}

func (b *BytecodeBuilder) AppendInt16(value int16) {
	b.AppendUint16(uint16(value))
}

func (b *BytecodeBuilder) PatchUint16(position int, value uint16) {
	if position+1 >= len(b.code) {
		panic("patch position out of bounds")
	}

	b.code[position] = byte(value & 0xFF)
	b.code[position+1] = byte(value >> 8)
	b.notifyUpdate()
}

func (b *BytecodeBuilder) CurrentPosition() int {
	return len(b.code)
}

func (b *BytecodeBuilder) Get() []byte {
	return b.code
}

func (b *BytecodeBuilder) SetSyncCallback(callback func([]byte)) {
	b.onUpdate = callback
}

func (b *BytecodeBuilder) notifyUpdate() {
	if b.onUpdate != nil {
		b.onUpdate(b.code)
	}
}
