#include "../../includes/interfaces/opcode_handler.h"

static uint16_t read_uint16(VMContext* context) {
    if (context->ip + 1 >= context->bytecode_len) {
        context->error_handler->runtime_error("Unexpected end of bytecode while reading uint16 at ip=%zu", context->ip);
        return 0;
    }

    byte_t low_byte = context->bytecode[context->ip++];
    byte_t high_byte = context->bytecode[context->ip++];

    return (uint16_t)((high_byte << 8) | low_byte);
}

static int16_t read_int16(VMContext* context) {
    return (int16_t)read_uint16(context);
}

status_t handle_jump(VMContext* context) {
    if (!context || !context->bytecode) {
        return STATUS_ERROR;
    }

    int16_t offset = read_int16(context);

    intptr_t current_ip_signed = (intptr_t)context->ip;
    intptr_t target_ip_signed = current_ip_signed + offset;

    if (target_ip_signed < 0 || (size_t)target_ip_signed >= context->bytecode_len) {
        context->error_handler->runtime_error("Jump target out of bounds (ip=%zu, offset=%d, target=%ld)",
                                             context->ip - 2, offset, target_ip_signed);
        return STATUS_ERROR;
    }

    context->ip = (size_t)target_ip_signed;

    return STATUS_SUCCESS;
}

status_t handle_jump_if_false(VMContext* context) {
    if (!context || !context->bytecode || !context->stack_manager) {
        return STATUS_ERROR;
    }

    uint16_t offset = read_uint16(context);

    if (context->stack_manager->is_empty()) {
        context->error_handler->runtime_error("Stack underflow in JUMP_IF_FALSE at ip=%zu", context->ip - 3);
        return STATUS_STACK_UNDERFLOW;
    }

    Value cond = context->stack_manager->pop();

    if (!context->value_handler->to_boolean(cond)) {
        size_t target_ip = context->ip + offset;

        if (target_ip >= context->bytecode_len) {
            context->error_handler->runtime_error("JUMP_IF_FALSE target out of bounds (ip=%zu, offset=%u, target=%zu)",
                                                 context->ip - 3, offset, target_ip);
            return STATUS_ERROR;
        }

        context->ip = target_ip;
    }

    return STATUS_SUCCESS;
}