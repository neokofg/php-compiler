#include "../../includes/interfaces/opcode_handler.h"
#include <stdlib.h>
#include <string.h>

VMContext* vm_context_new() {
    VMContext* context = (VMContext*)malloc(sizeof(VMContext));
    if (!context) return NULL;

    context->bytecode = NULL;
    context->bytecode_len = 0;
    context->ip = 0;
    context->constants = NULL;
    context->constants_len = 0;
    context->variables = (Value*)calloc(VAR_COUNT, sizeof(Value));
    context->value_handler = NULL;
    context->stack_manager = NULL;
    context->error_handler = NULL;
    context->user_data = NULL;

    return context;
}

void vm_context_free(VMContext* context) {
    if (!context) return;

    if (context->variables) {
        free(context->variables);
    }

    free(context);
}

void vm_context_set_bytecode(VMContext* context, byte_t* bytecode, size_t bytecode_len) {
    if (!context) return;

    context->bytecode = bytecode;
    context->bytecode_len = bytecode_len;
    context->ip = 0;
}

void vm_context_set_constants(VMContext* context, Value* constants, size_t constants_len) {
    if (!context) return;

    context->constants = constants;
    context->constants_len = constants_len;
}

void vm_context_set_handlers(VMContext* context,
                             ValueHandler* value_handler,
                             StackManager* stack_manager,
                             ErrorHandler* error_handler) {
    if (!context) return;

    context->value_handler = value_handler;
    context->stack_manager = stack_manager;
    context->error_handler = error_handler;
}

void vm_context_set_user_data(VMContext* context, void* user_data) {
    if (!context) return;

    context->user_data = user_data;
}

size_t vm_context_get_ip(VMContext* context) {
    if (!context) return 0;

    return context->ip;
}

void vm_context_set_ip(VMContext* context, size_t ip) {
    if (!context) return;

    context->ip = ip;
}

void vm_context_reset(VMContext* context) {
    if (!context) return;

    context->bytecode = NULL;
    context->bytecode_len = 0;
    context->ip = 0;
    context->constants = NULL;
    context->constants_len = 0;

    if (context->variables) {
        memset(context->variables, 0, VAR_COUNT * sizeof(Value));
    }
}