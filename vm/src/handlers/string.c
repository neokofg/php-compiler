/* Licensed under GNU GPL v3. See LICENSE file for details. */
#include "../../includes/interfaces/opcode_handler.h"
#include <string.h>

status_t handle_concat(VMContext* context) {
    if (!context || !context->stack_manager) {
        return STATUS_ERROR;
    }

    if (context->stack_manager->size() < 2) {
        context->error_handler->runtime_error("Stack underflow in CONCAT at ip=%zu", context->ip - 1);
        return STATUS_STACK_UNDERFLOW;
    }

    // Get values but don't pop them yet
    Value b = context->stack_manager->peek(0);
    Value a = context->stack_manager->peek(1);

    // Convert to strings
    char* str_a = context->value_handler->to_string(a);
    char* str_b = context->value_handler->to_string(b);

    if (!str_a || !str_b) {
        if (str_a) free(str_a);
        if (str_b) free(str_b);
        context->error_handler->runtime_error("Failed to convert values to strings for concatenation");
        return STATUS_ERROR;
    }

    // Create result string
    size_t len_a = strlen(str_a);
    size_t len_b = strlen(str_b);
    char* result = (char*)malloc(len_a + len_b + 1);

    if (!result) {
        free(str_a);
        free(str_b);
        context->error_handler->runtime_error("Memory allocation failed for concatenation");
        return STATUS_OUT_OF_MEMORY;
    }

    // Copy strings
    memcpy(result, str_a, len_a);
    memcpy(result + len_a, str_b, len_b);
    result[len_a + len_b] = '\0';

    // Now pop the original values
    context->stack_manager->pop(); // b
    context->stack_manager->pop(); // a

    // Create result value and push to stack
    Value concat_value;
    concat_value.type = TYPE_STRING;
    concat_value.value.str_val = result;

    context->stack_manager->push(concat_value);

    // Free the temporary strings from to_string conversions
    free(str_a);
    free(str_b);

    return STATUS_SUCCESS;
}