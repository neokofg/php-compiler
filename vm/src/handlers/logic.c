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

status_t handle_gt(VMContext* context) {
    status_t status = check_stack_size(context, 2);
    if (status != STATUS_SUCCESS) {
        return status;
    }

    Value b = context->stack_manager->pop();
    Value a = context->stack_manager->pop();

    bool result = context->value_handler->greater_than(a, b);

    context->stack_manager->push(context->value_handler->create_boolean(result));

    return STATUS_SUCCESS;
}

status_t handle_lt(VMContext* context) {
    status_t status = check_stack_size(context, 2);
    if (status != STATUS_SUCCESS) {
        return status;
    }

    Value b = context->stack_manager->pop();
    Value a = context->stack_manager->pop();

    bool result = context->value_handler->less_than(a, b);

    context->stack_manager->push(context->value_handler->create_boolean(result));

    return STATUS_SUCCESS;
}

status_t handle_eq(VMContext* context) {
    status_t status = check_stack_size(context, 2);
    if (status != STATUS_SUCCESS) {
        return status;
    }

    Value b = context->stack_manager->pop();
    Value a = context->stack_manager->pop();

    bool result = context->value_handler->equals(a, b);

    context->stack_manager->push(context->value_handler->create_boolean(result));

    return STATUS_SUCCESS;
}

status_t handle_and(VMContext* context) {
    status_t status = check_stack_size(context, 2);
    if (status != STATUS_SUCCESS) {
        return status;
    }

    Value b = context->stack_manager->pop();
    Value a = context->stack_manager->pop();

    bool result = context->value_handler->to_boolean(a) && context->value_handler->to_boolean(b);

    context->stack_manager->push(context->value_handler->create_boolean(result));

    return STATUS_SUCCESS;
}

status_t handle_or(VMContext* context) {
    status_t status = check_stack_size(context, 2);
    if (status != STATUS_SUCCESS) {
        return status;
    }

    Value b = context->stack_manager->pop();
    Value a = context->stack_manager->pop();

    bool result = context->value_handler->to_boolean(a) || context->value_handler->to_boolean(b);

    context->stack_manager->push(context->value_handler->create_boolean(result));

    return STATUS_SUCCESS;
}