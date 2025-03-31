/* Licensed under GNU GPL v3. See LICENSE file for details. */
#ifndef VM_OPCODE_HANDLER_H
#define VM_OPCODE_HANDLER_H

#include "../common.h"
#include "value_handler.h"
#include "stack_manager.h"
#include "error_handler.h"
#include "../opcode/opcodes.h"

typedef struct VMContext VMContext;

typedef status_t (*OpcodeHandlerFunc)(VMContext* context);

typedef struct OpcodeHandler {
    void (*register_handler)(byte_t opcode, OpcodeHandlerFunc handler);

    status_t (*execute)(VMContext* context, byte_t opcode);

    const char* (*get_opcode_name)(byte_t opcode);
    bool (*is_opcode_valid)(byte_t opcode);

    void (*reset_handlers)(void);
} OpcodeHandler;

struct VMContext {
    byte_t* bytecode;
    size_t bytecode_len;
    size_t ip;  // Instruction pointer

    Value* constants;
    size_t constants_len;

    Value* variables;

    ValueHandler* value_handler;
    StackManager* stack_manager;
    ErrorHandler* error_handler;

    void* user_data;
};

OpcodeHandler* opcode_handler_new(void);

void opcode_handler_free(OpcodeHandler* handler);

status_t handle_load_const(VMContext* context);
status_t handle_print(VMContext* context);
status_t handle_halt(VMContext* context);
status_t handle_pop(VMContext* context);

status_t handle_add(VMContext* context);
status_t handle_sub(VMContext* context);
status_t handle_mul(VMContext* context);
status_t handle_div(VMContext* context);

status_t handle_store_var(VMContext* context);
status_t handle_load_var(VMContext* context);

status_t handle_jump(VMContext* context);
status_t handle_jump_if_false(VMContext* context);

status_t handle_gt(VMContext* context);
status_t handle_lt(VMContext* context);
status_t handle_eq(VMContext* context);

status_t handle_and(VMContext* context);
status_t handle_or(VMContext* context);

#endif /* VM_OPCODE_HANDLER_H */