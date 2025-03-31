#ifndef VM_CONTEXT_H
#define VM_CONTEXT_H

#include "common.h"
#include "interfaces/opcode_handler.h"

struct VMContext;

VMContext* vm_context_new(void);
void vm_context_free(VMContext* context);
void vm_context_set_bytecode(VMContext* context, byte_t* bytecode, size_t bytecode_len);
void vm_context_set_constants(VMContext* context, Value* constants, size_t constants_len);
void vm_context_set_handlers(VMContext* context,
                             ValueHandler* value_handler,
                             StackManager* stack_manager,
                             ErrorHandler* error_handler);
void vm_context_set_user_data(VMContext* context, void* user_data);
size_t vm_context_get_ip(VMContext* context);
void vm_context_set_ip(VMContext* context, size_t ip);
void vm_context_reset(VMContext* context);

#endif /* VM_CONTEXT_H */