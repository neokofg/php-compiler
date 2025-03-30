#ifndef VM_H
#define VM_H

#include <stdint.h>
#include <stdbool.h> 
#include <stddef.h>

#define STACK_SIZE 256
#define CONST_POOL_SIZE 256
#define VAR_COUNT 256

#define OP_LOAD_CONST       0x01
#define OP_PRINT            0x02
#define OP_ADD              0x03
#define OP_SUB              0x04
#define OP_MUL              0x05
#define OP_DIV              0x06
#define OP_HALT             0xFF
#define OP_STORE_VAR        0x10
#define OP_LOAD_VAR         0x11
#define OP_JUMP             0x21
#define OP_JUMP_IF_FALSE    0x20
#define OP_GT               0x07
#define OP_LT               0x08
#define OP_AND              0x09
#define OP_OR               0x0A
#define OP_EQ               0x0B
#define OP_POP              0x0C

typedef enum {
    TYPE_INT = 0,
    TYPE_STRING = 1
} ValueType;

typedef struct {
    ValueType type;
    union {
        char* str_val;
        int int_val;
    } value;
} Value;

bool is_truthy(Value v);

Value stack_pop();

void run(uint8_t* bytecode, size_t length, Value constants[], size_t constants_len);

#endif