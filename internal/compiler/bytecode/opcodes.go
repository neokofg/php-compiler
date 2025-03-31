// Licensed under GNU GPL v3. See LICENSE file for details.
package bytecode

const (
	OP_LOAD_CONST    = 0x01
	OP_PRINT         = 0x02
	OP_ADD           = 0x03
	OP_SUB           = 0x04
	OP_MUL           = 0x05
	OP_DIV           = 0x06
	OP_STORE_VAR     = 0x10
	OP_LOAD_VAR      = 0x11
	OP_HALT          = 0xFF
	OP_JUMP_IF_FALSE = 0x20
	OP_JUMP          = 0x21
	OP_GT            = 0x07
	OP_LT            = 0x08
	OP_AND           = 0x09
	OP_OR            = 0x0A
	OP_EQ            = 0x0B
	OP_POP           = 0x0C
)
