/* Licensed under GNU GPL v3. See LICENSE file for details. */
#include "../../includes/interfaces/opcode_handler.h"

status_t handle_load_const(VMContext* context) {
    if (!context || !context->bytecode || !context->constants) {
        return STATUS_ERROR;
    }

    if (context->ip >= context->bytecode_len) {
        context->error_handler->runtime_error("Unexpected end of bytecode at ip=%zu", context->ip);
        return STATUS_ERROR;
    }

    byte_t const_idx = context->bytecode[context->ip++];

    if (const_idx >= context->constants_len) {
        context->error_handler->runtime_error("Invalid constant index %d at ip=%zu, max allowed: %zu",
                                              const_idx, context->ip - 1, context->constants_len - 1);
        return STATUS_ERROR;
    }

    context->stack_manager->push(context->constants[const_idx]);

    return STATUS_SUCCESS;
}

status_t handle_print(VMContext* context) {
    if (!context || !context->stack_manager) {
        return STATUS_ERROR;
    }

    if (context->stack_manager->is_empty()) {
        context->error_handler->runtime_error("Stack underflow in PRINT at ip=%zu", context->ip - 1);
        return STATUS_STACK_UNDERFLOW;
    }

    Value value = context->stack_manager->pop();
    switch (value.type) {
        case TYPE_INT:
            printf("%d", value.value.int_val);
            break;
        case TYPE_STRING:
            if (value.value.str_val) {
                printf("%s", value.value.str_val);
            }
            break;
        case TYPE_BOOLEAN:
            printf("%s", value.value.bool_val ? "true" : "false");
            break;
        case TYPE_NULL:
            printf("null");
            break;
        default:
            printf("unknown");
            break;
    }

    fflush(stdout);

    return STATUS_SUCCESS;
}

status_t handle_halt(VMContext* context) {
    if (!context) {
        return STATUS_ERROR;
    }

    // Set IP to the end of the bytecode to stop execution
    context->ip = context->bytecode_len;

    return STATUS_SUCCESS;
}

status_t handle_pop(VMContext* context) {
    if (!context || !context->stack_manager) {
        return STATUS_ERROR;
    }

    if (context->stack_manager->is_empty()) {
        context->error_handler->runtime_error("Stack underflow in POP at ip=%zu", context->ip - 1);
        return STATUS_STACK_UNDERFLOW;
    }

    context->stack_manager->pop();

    return STATUS_SUCCESS;
}

status_t handle_store_var(VMContext* context) {
    if (!context || !context->bytecode || !context->variables) {
        return STATUS_ERROR;
    }

    if (context->ip >= context->bytecode_len) {
        context->error_handler->runtime_error("Unexpected end of bytecode at ip=%zu", context->ip);
        return STATUS_ERROR;
    }

    byte_t var_idx = context->bytecode[context->ip++];

    if (var_idx >= VAR_COUNT) {
        context->error_handler->runtime_error("Invalid variable index %d at ip=%zu, max allowed: %d",
                                              var_idx, context->ip - 1, VAR_COUNT - 1);
        return STATUS_ERROR;
    }

    if (context->stack_manager->is_empty()) {
        context->error_handler->runtime_error("Stack underflow in STORE_VAR at ip=%zu", context->ip - 1);
        return STATUS_STACK_UNDERFLOW;
    }

    Value value = context->stack_manager->pop();
    context->variables[var_idx] = value;

    return STATUS_SUCCESS;
}

status_t handle_load_var(VMContext* context) {
    if (!context || !context->bytecode || !context->variables) {
        return STATUS_ERROR;
    }

    if (context->ip >= context->bytecode_len) {
        context->error_handler->runtime_error("Unexpected end of bytecode at ip=%zu", context->ip);
        return STATUS_ERROR;
    }

    byte_t var_idx = context->bytecode[context->ip++];

    if (var_idx >= VAR_COUNT) {
        context->error_handler->runtime_error("Invalid variable index %d at ip=%zu, max allowed: %d",
                                              var_idx, context->ip - 1, VAR_COUNT - 1);
        return STATUS_ERROR;
    }

    context->stack_manager->push(context->variables[var_idx]);

    return STATUS_SUCCESS;
}