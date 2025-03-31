#include "../../includes/interfaces/opcode_handler.h"

static status_t check_stack_size(VMContext* context, int required_size) {
    if (!context || !context->stack_manager) {
        return STATUS_ERROR;
    }

    if (context->stack_manager->size() < required_size) {
        context->error_handler->runtime_error("Stack underflow at ip=%zu, need %d elements, have %d",
                                             context->ip - 1, required_size, context->stack_manager->size());
        return STATUS_STACK_UNDERFLOW;
    }

    return STATUS_SUCCESS;
}

status_t handle_add(VMContext* context) {
    status_t status = check_stack_size(context, 2);
    if (status != STATUS_SUCCESS) {
        return status;
    }

    Value b = context->stack_manager->pop();
    Value a = context->stack_manager->pop();

    int_t a_int = context->value_handler->to_int(a);
    int_t b_int = context->value_handler->to_int(b);

    int_t result = a_int + b_int;

    context->stack_manager->push(context->value_handler->create_int(result));

    return STATUS_SUCCESS;
}

status_t handle_sub(VMContext* context) {
    status_t status = check_stack_size(context, 2);
    if (status != STATUS_SUCCESS) {
        return status;
    }

    Value b = context->stack_manager->pop();
    Value a = context->stack_manager->pop();

    int_t a_int = context->value_handler->to_int(a);
    int_t b_int = context->value_handler->to_int(b);

    int_t result = a_int - b_int;

    context->stack_manager->push(context->value_handler->create_int(result));

    return STATUS_SUCCESS;
}

status_t handle_mul(VMContext* context) {
    status_t status = check_stack_size(context, 2);
    if (status != STATUS_SUCCESS) {
        return status;
    }

    Value b = context->stack_manager->pop();
    Value a = context->stack_manager->pop();

    int_t a_int = context->value_handler->to_int(a);
    int_t b_int = context->value_handler->to_int(b);

    int_t result = a_int * b_int;

    context->stack_manager->push(context->value_handler->create_int(result));

    return STATUS_SUCCESS;
}

status_t handle_div(VMContext* context) {
    status_t status = check_stack_size(context, 2);
    if (status != STATUS_SUCCESS) {
        return status;
    }

    Value b = context->stack_manager->pop();
    Value a = context->stack_manager->pop();

    int_t a_int = context->value_handler->to_int(a);
    int_t b_int = context->value_handler->to_int(b);

    if (b_int == 0) {
        context->error_handler->warning("Division by zero at ip=%zu", context->ip - 1);
        context->stack_manager->push(context->value_handler->create_int(0));
        return STATUS_DIVISION_BY_ZERO;
    }

    if (a_int % b_int != 0) {
        context->error_handler->warning("Integer division resulting in truncation (%d / %d) at ip=%zu",
                                       a_int, b_int, context->ip - 1);
    }

    int_t result = a_int / b_int;

    context->stack_manager->push(context->value_handler->create_int(result));

    return STATUS_SUCCESS;
}