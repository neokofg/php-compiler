/* Licensed under GNU GPL v3. See LICENSE file for details. */
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

status_t handle_not(VMContext* context) {
    if (!context || !context->stack_manager) {
        return STATUS_ERROR;
    }

    if (context->stack_manager->is_empty()) {
        context->error_handler->runtime_error("Stack underflow in NOT at ip=%zu", context->ip - 1);
        return STATUS_STACK_UNDERFLOW;
    }

    Value val = context->stack_manager->pop();
    bool result = !context->value_handler->to_boolean(val);

    context->stack_manager->push(context->value_handler->create_boolean(result));

    return STATUS_SUCCESS;
}

status_t handle_gte(VMContext* context) {
    status_t status = check_stack_size(context, 2);
    if (status != STATUS_SUCCESS) {
        return status;
    }

    Value b = context->stack_manager->pop();
    Value a = context->stack_manager->pop();

    bool result = context->value_handler->greater_than(a, b) ||
                 context->value_handler->equals(a, b);

    context->stack_manager->push(context->value_handler->create_boolean(result));
    return STATUS_SUCCESS;
}

status_t handle_lte(VMContext* context) {
    status_t status = check_stack_size(context, 2);
    if (status != STATUS_SUCCESS) {
        return status;
    }

    Value b = context->stack_manager->pop();
    Value a = context->stack_manager->pop();

    bool result = context->value_handler->less_than(a, b) ||
                 context->value_handler->equals(a, b);

    context->stack_manager->push(context->value_handler->create_boolean(result));
    return STATUS_SUCCESS;
}

status_t handle_identity_eq(VMContext* context) {
    status_t status = check_stack_size(context, 2);
    if (status != STATUS_SUCCESS) {
        return status;
    }

    Value b = context->stack_manager->pop();
    Value a = context->stack_manager->pop();

    bool result = a.type == b.type && context->value_handler->equals(a, b);
    context->stack_manager->push(context->value_handler->create_boolean(result));
    return STATUS_SUCCESS;
}

status_t handle_identity_ne(VMContext* context) {
    status_t status = check_stack_size(context, 2);
    if (status != STATUS_SUCCESS) {
        return status;
    }

    Value b = context->stack_manager->pop();
    Value a = context->stack_manager->pop();

    bool result = a.type != b.type || !context->value_handler->equals(a, b);
    context->stack_manager->push(context->value_handler->create_boolean(result));
    return STATUS_SUCCESS;
}

status_t handle_bit_and(VMContext* context) {
    status_t status = check_stack_size(context, 2);
    if (status != STATUS_SUCCESS) {
        return status;
    }

    Value b = context->stack_manager->pop();
    Value a = context->stack_manager->pop();

    int_t a_int = context->value_handler->to_int(a);
    int_t b_int = context->value_handler->to_int(b);
    int_t result = a_int & b_int;

    context->stack_manager->push(context->value_handler->create_int(result));
    return STATUS_SUCCESS;
}

status_t handle_bit_or(VMContext* context) {
    status_t status = check_stack_size(context, 2);
    if (status != STATUS_SUCCESS) {
        return status;
    }

    Value b = context->stack_manager->pop();
    Value a = context->stack_manager->pop();

    int_t a_int = context->value_handler->to_int(a);
    int_t b_int = context->value_handler->to_int(b);
    int_t result = a_int | b_int;

    context->stack_manager->push(context->value_handler->create_int(result));
    return STATUS_SUCCESS;
}

status_t handle_bit_xor(VMContext* context) {
    status_t status = check_stack_size(context, 2);
    if (status != STATUS_SUCCESS) {
        return status;
    }

    Value b = context->stack_manager->pop();
    Value a = context->stack_manager->pop();

    int_t a_int = context->value_handler->to_int(a);
    int_t b_int = context->value_handler->to_int(b);
    int_t result = a_int ^ b_int;

    context->stack_manager->push(context->value_handler->create_int(result));
    return STATUS_SUCCESS;
}

status_t handle_bit_not(VMContext* context) {
    if (!context || !context->stack_manager) {
        return STATUS_ERROR;
    }

    if (context->stack_manager->is_empty()) {
        context->error_handler->runtime_error("Stack underflow in BIT_NOT at ip=%zu", context->ip - 1);
        return STATUS_STACK_UNDERFLOW;
    }

    Value val = context->stack_manager->pop();
    int_t int_val = context->value_handler->to_int(val);
    int_t result = ~int_val;

    context->stack_manager->push(context->value_handler->create_int(result));
    return STATUS_SUCCESS;
}

status_t handle_lshift(VMContext* context) {
    status_t status = check_stack_size(context, 2);
    if (status != STATUS_SUCCESS) {
        return status;
    }

    Value b = context->stack_manager->pop();
    Value a = context->stack_manager->pop();

    int_t a_int = context->value_handler->to_int(a);
    int_t b_int = context->value_handler->to_int(b);

    if (b_int < 0) {
        context->error_handler->warning("Negative shift amount at ip=%zu", context->ip - 1);
        b_int = 0;
    }

    int_t result = a_int << b_int;
    context->stack_manager->push(context->value_handler->create_int(result));
    return STATUS_SUCCESS;
}

status_t handle_rshift(VMContext* context) {
    status_t status = check_stack_size(context, 2);
    if (status != STATUS_SUCCESS) {
        return status;
    }

    Value b = context->stack_manager->pop();
    Value a = context->stack_manager->pop();

    int_t a_int = context->value_handler->to_int(a);
    int_t b_int = context->value_handler->to_int(b);

    if (b_int < 0) {
        context->error_handler->warning("Negative shift amount at ip=%zu", context->ip - 1);
        b_int = 0;
    }

    int_t result = a_int >> b_int;
    context->stack_manager->push(context->value_handler->create_int(result));
    return STATUS_SUCCESS;
}

status_t handle_assign_add(VMContext* context) {
    status_t status = check_stack_size(context, 2);
    if (status != STATUS_SUCCESS) {
        return status;
    }

    Value value = context->stack_manager->pop();
    Value var = context->stack_manager->pop();

    int_t var_int = context->value_handler->to_int(var);
    int_t value_int = context->value_handler->to_int(value);
    int_t result = var_int + value_int;

    context->stack_manager->push(context->value_handler->create_int(result));
    return STATUS_SUCCESS;
}

status_t handle_assign_sub(VMContext* context) {
    status_t status = check_stack_size(context, 2);
    if (status != STATUS_SUCCESS) {
        return status;
    }

    Value value = context->stack_manager->pop();
    Value var = context->stack_manager->pop();

    int_t var_int = context->value_handler->to_int(var);
    int_t value_int = context->value_handler->to_int(value);
    int_t result = var_int - value_int;

    context->stack_manager->push(context->value_handler->create_int(result));
    return STATUS_SUCCESS;
}

status_t handle_assign_mul(VMContext* context) {
    status_t status = check_stack_size(context, 2);
    if (status != STATUS_SUCCESS) {
        return status;
    }

    Value value = context->stack_manager->pop();
    Value var = context->stack_manager->pop();

    int_t var_int = context->value_handler->to_int(var);
    int_t value_int = context->value_handler->to_int(value);
    int_t result = var_int * value_int;

    context->stack_manager->push(context->value_handler->create_int(result));
    return STATUS_SUCCESS;
}

status_t handle_assign_div(VMContext* context) {
    status_t status = check_stack_size(context, 2);
    if (status != STATUS_SUCCESS) {
        return status;
    }

    Value value = context->stack_manager->pop();
    Value var = context->stack_manager->pop();

    int_t var_int = context->value_handler->to_int(var);
    int_t value_int = context->value_handler->to_int(value);

    if (value_int == 0) {
        context->error_handler->warning("Division by zero in /= at ip=%zu", context->ip - 1);
        context->stack_manager->push(context->value_handler->create_int(0));
        return STATUS_DIVISION_BY_ZERO;
    }

    int_t result = var_int / value_int;
    context->stack_manager->push(context->value_handler->create_int(result));
    return STATUS_SUCCESS;
}

status_t handle_assign_mod(VMContext* context) {
    status_t status = check_stack_size(context, 2);
    if (status != STATUS_SUCCESS) {
        return status;
    }

    Value value = context->stack_manager->pop();
    Value var = context->stack_manager->pop();

    int_t var_int = context->value_handler->to_int(var);
    int_t value_int = context->value_handler->to_int(value);

    if (value_int == 0) {
        context->error_handler->warning("Modulo by zero in %= at ip=%zu", context->ip - 1);
        context->stack_manager->push(context->value_handler->create_int(0));
        return STATUS_DIVISION_BY_ZERO;
    }

    int_t result = var_int % value_int;
    context->stack_manager->push(context->value_handler->create_int(result));
    return STATUS_SUCCESS;
}

status_t handle_assign_concat(VMContext* context) {
    status_t status = check_stack_size(context, 2);
    if (status != STATUS_SUCCESS) {
        return status;
    }

    Value b = context->stack_manager->peek(0);
    Value a = context->stack_manager->peek(1);

    char* str_a = context->value_handler->to_string(a);
    char* str_b = context->value_handler->to_string(b);

    if (!str_a || !str_b) {
        if (str_a) free(str_a);
        if (str_b) free(str_b);
        context->error_handler->runtime_error("Failed to convert values to strings for .=");
        return STATUS_ERROR;
    }

    size_t len_a = strlen(str_a);
    size_t len_b = strlen(str_b);
    char* result = (char*)malloc(len_a + len_b + 1);

    if (!result) {
        free(str_a);
        free(str_b);
        context->error_handler->runtime_error("Memory allocation failed for .=");
        return STATUS_OUT_OF_MEMORY;
    }

    memcpy(result, str_a, len_a);
    memcpy(result + len_a, str_b, len_b);
    result[len_a + len_b] = '\0';

    context->stack_manager->pop(); // b
    context->stack_manager->pop(); // a

    Value concat_value;
    concat_value.type = TYPE_STRING;
    concat_value.value.str_val = result;

    context->stack_manager->push(concat_value);

    free(str_a);
    free(str_b);

    return STATUS_SUCCESS;
}