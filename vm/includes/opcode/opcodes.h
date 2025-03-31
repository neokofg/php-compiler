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

#define OP_STORE_VAR        0x10
#define OP_LOAD_VAR         0x11

#define OP_JUMP             0x21
#define OP_JUMP_IF_FALSE    0x20

#define OP_GT               0x07
#define OP_LT               0x08
#define OP_EQ               0x0B

#define OP_AND              0x09
#define OP_OR               0x0A

#endif /* VM_OPCODES_H */