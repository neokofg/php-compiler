/* Licensed under GNU GPL v3. See LICENSE file for details. */
#ifndef VM_H
#define VM_H

#include "common.h"
#include "interfaces/value_handler.h"
#include "interfaces/memory_manager.h"
#include "interfaces/stack_manager.h"
#include "interfaces/error_handler.h"
#include "interfaces/opcode_handler.h"

typedef struct VM {
    VMContext* context;

    ValueHandler* value_handler;
    StackManager* stack_manager;
    ErrorHandler* error_handler;
    OpcodeHandler* opcode_handler;
    MemoryManager* memory_manager;

    bool running;
    status_t last_status;

    bool debug_mode;
} VM;

VM* vm_new(void);
void vm_free(VM* vm);

status_t vm_execute(VM* vm, byte_t* bytecode, size_t bytecode_len, Value* constants, size_t constants_len);
status_t vm_execute_instruction(VM* vm);
void vm_reset(VM* vm);

bool vm_is_running(VM* vm);
status_t vm_get_status(VM* vm);

void vm_set_debug_mode(VM* vm, bool debug_mode);

void vm_register_opcode_handler(VM* vm, byte_t opcode, OpcodeHandlerFunc handler);
void vm_set_user_data(VM* vm, void* user_data);
void* vm_get_user_data(VM* vm);

#endif /* VM_H */