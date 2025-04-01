// Licensed under GNU GPL v3. See LICENSE file for details.
package bytecode

const (
	OP_LOAD_CONST = 0x01
	OP_PRINT      = 0x02
	OP_HALT       = 0xFF
	OP_POP        = 0x0C

	OP_ADD = 0x03
	OP_SUB = 0x04
	OP_MUL = 0x05
	OP_DIV = 0x06

	OP_CONCAT = 0x0F

	OP_STORE_VAR = 0x10
	OP_LOAD_VAR  = 0x11

	OP_JUMP          = 0x21
	OP_JUMP_IF_FALSE = 0x20

	OP_GT  = 0x07
	OP_LT  = 0x08
	OP_EQ  = 0x0B
	OP_NOT = 0x0D

	OP_AND = 0x09
	OP_OR  = 0x0A

	OP_INC      = 0x30
	OP_DEC      = 0x31
	OP_POST_INC = 0x32
	OP_POST_DEC = 0x33

	OP_MOD = 0x34

	OP_BIT_AND = 0x40
	OP_BIT_OR  = 0x41
	OP_BIT_XOR = 0x42
	OP_BIT_NOT = 0x43
	OP_LSHIFT  = 0x44
	OP_RSHIFT  = 0x45

	OP_GTE         = 0x50
	OP_LTE         = 0x51
	OP_IDENTITY_EQ = 0x52
	OP_IDENTITY_NE = 0x53

	OP_ASSIGN_ADD    = 0x60
	OP_ASSIGN_SUB    = 0x61
	OP_ASSIGN_MUL    = 0x62
	OP_ASSIGN_DIV    = 0x63
	OP_ASSIGN_MOD    = 0x64
	OP_ASSIGN_CONCAT = 0x65
)
