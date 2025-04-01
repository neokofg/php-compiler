/* Licensed under GNU GPL v3. See LICENSE file for details. */
#ifndef VM_OPCODES_H
#define VM_OPCODES_H

#define OP_LOAD_CONST       0x01
#define OP_PRINT            0x02
#define OP_HALT             0xFF
#define OP_POP              0x0C

#define OP_ADD              0x03
#define OP_SUB              0x04
#define OP_MUL              0x05
#define OP_DIV              0x06

#define OP_CONCAT           0x0F

#define OP_STORE_VAR        0x10
#define OP_LOAD_VAR         0x11

#define OP_JUMP             0x21
#define OP_JUMP_IF_FALSE    0x20

#define OP_GT               0x07
#define OP_LT               0x08
#define OP_EQ               0x0B
#define OP_NOT              0x0D

#define OP_AND              0x09
#define OP_OR               0x0A

#define OP_INC              0x30
#define OP_DEC              0x31
#define OP_POST_INC         0x32
#define OP_POST_DEC         0x33

#define OP_MOD              0x34

#define OP_BIT_AND          0x40
#define OP_BIT_OR           0x41
#define OP_BIT_XOR          0x42
#define OP_BIT_NOT          0x43
#define OP_LSHIFT           0x44
#define OP_RSHIFT           0x45

#define OP_GTE              0x50
#define OP_LTE              0x51
#define OP_IDENTITY_EQ      0x52
#define OP_IDENTITY_NE      0x53

#define OP_ASSIGN_ADD       0x60
#define OP_ASSIGN_SUB       0x61
#define OP_ASSIGN_MUL       0x62
#define OP_ASSIGN_DIV       0x63
#define OP_ASSIGN_MOD       0x64
#define OP_ASSIGN_CONCAT    0x65

#define OP_BREAK            0x70
#define OP_CONTINUE         0x71

#endif /* VM_OPCODES_H */