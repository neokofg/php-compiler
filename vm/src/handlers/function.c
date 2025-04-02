/* Licensed under GNU GPL v3. See LICENSE file for details. */
#include "../../includes/interfaces/opcode_handler.h"
#include "../../includes/vm.h"

static int_t return_address = -1;

status_t handle_func_call(VMContext* context) {
    byte_t param_count = context->bytecode[context->ip++];

    byte_t low_byte = context->bytecode[context->ip++];
    byte_t high_byte = context->bytecode[context->ip++];
    uint16_t func_addr = (uint16_t)((high_byte << 8) | low_byte);

    if (func_addr >= context->bytecode_len) {
        context->error_handler->runtime_error("Invalid function address %u at ip=%zu, bytecode_len=%zu",
                                             func_addr, context->ip - 2, context->bytecode_len);
        return STATUS_ERROR;
    }

    byte_t func_opcode = context->bytecode[func_addr];

    if (func_opcode != 0x80) {
        context->error_handler->runtime_error("Invalid function opcode 0x%02X at address %u",
                                             func_opcode, func_addr);
        return STATUS_ERROR;
    }

    return_address = context->ip;
    context->ip = func_addr;

    return STATUS_SUCCESS;
}

status_t handle_func_decl(VMContext* context) {
    byte_t param_count = context->bytecode[context->ip++];

    for (int i = 0; i < param_count; i++) {
        byte_t var_idx = context->bytecode[context->ip++];

        if (!context->stack_manager->is_empty()) {
            Value param_val = context->stack_manager->pop();
            context->variables[var_idx] = param_val;
        } else {
            context->error_handler->warning("Not enough parameters for function at ip=%zu", context->ip);

            Value default_val;
            default_val.type = TYPE_INT;
            default_val.value.int_val = 0;
            context->variables[var_idx] = default_val;
        }
    }

    return STATUS_SUCCESS;
}


status_t handle_enter_func(VMContext* context) {
    return STATUS_SUCCESS;
}

status_t handle_return(VMContext* context) {
    if (return_address < 0) {
        context->error_handler->runtime_error("No return address set in RETURN at ip=%zu", context->ip);
        return STATUS_ERROR;
    }

    if (context->stack_manager->is_empty()) {
        context->error_handler->runtime_error("Stack empty in RETURN at ip=%zu", context->ip);
        return STATUS_STACK_UNDERFLOW;
    }

    Value return_value = context->stack_manager->pop();

    if (return_address < 0 || return_address >= (int)context->bytecode_len) {
        context->error_handler->runtime_error("Invalid return address %d at ip=%zu, bytecode_len=%zu",
                                             return_address, context->ip, context->bytecode_len);
        return STATUS_ERROR;
    }

    context->ip = return_address;
    context->stack_manager->push(return_value);

    return_address = -1;

    return STATUS_SUCCESS;
}


status_t handle_exit_func(VMContext* context) {
    Value null_val;
    null_val.type = TYPE_NULL;

    if (return_address < 0) {
        context->error_handler->runtime_error("No return address set in EXIT_FUNC at ip=%zu", context->ip);
        return STATUS_ERROR;
    }

    if (return_address < 0 || return_address >= (int)context->bytecode_len) {
        context->error_handler->runtime_error("Invalid return address %d at ip=%zu, bytecode_len=%zu",
                                             return_address, context->ip, context->bytecode_len);
        return STATUS_ERROR;
    }

    context->ip = return_address;
    context->stack_manager->push(null_val);

    return_address = -1;

    return STATUS_SUCCESS;
}

