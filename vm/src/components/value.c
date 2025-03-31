/* Licensed under GNU GPL v3. See LICENSE file for details. */
#include "../../includes/interfaces/value_handler.h"
#include <stdlib.h>
#include <string.h>
#include <stdio.h>

static Value create_int(int_t value) {
    Value val;
    val.type = TYPE_INT;
    val.value.int_val = value;
    return val;
}

static Value create_string(const char* value) {
    Value val;
    val.type = TYPE_STRING;
    val.value.str_val = strdup(value);
    return val;
}

static Value create_boolean(bool value) {
    Value val;
    val.type = TYPE_BOOLEAN;
    val.value.bool_val = value;
    return val;
}

static Value create_null(void) {
    Value val;
    val.type = TYPE_NULL;
    return val;
}

static int_t to_int(Value value) {
    switch (value.type) {
        case TYPE_INT:
            return value.value.int_val;
        case TYPE_STRING:
            return value.value.str_val ? atoi(value.value.str_val) : 0;
        case TYPE_BOOLEAN:
            return value.value.bool_val ? 1 : 0;
        case TYPE_NULL:
            return 0;
        default:
            return 0;
    }
}

static char* to_string(Value value) {
    char buffer[64];

    switch (value.type) {
        case TYPE_INT:
            snprintf(buffer, sizeof(buffer), "%d", value.value.int_val);
            return strdup(buffer);
        case TYPE_STRING:
            return value.value.str_val ? strdup(value.value.str_val) : strdup("");
        case TYPE_BOOLEAN:
            return strdup(value.value.bool_val ? "true" : "false");
        case TYPE_NULL:
            return strdup("null");
        default:
            return strdup("unknown");
    }
}

static bool to_boolean(Value value) {
    switch (value.type) {
        case TYPE_INT:
            return value.value.int_val != 0;
        case TYPE_STRING:
            return value.value.str_val &&
                  (strcmp(value.value.str_val, "") != 0) &&
                  (strcmp(value.value.str_val, "0") != 0);
        case TYPE_BOOLEAN:
            return value.value.bool_val;
        case TYPE_NULL:
            return false;
        default:
            return false;
    }
}

static bool equals(Value a, Value b) {
    if (a.type != b.type) {
        if ((a.type == TYPE_INT && b.type == TYPE_STRING) ||
            (a.type == TYPE_STRING && b.type == TYPE_INT)) {
            return to_int(a) == to_int(b);
        }

        if ((a.type == TYPE_BOOLEAN && (b.type == TYPE_INT || b.type == TYPE_STRING)) ||
            (b.type == TYPE_BOOLEAN && (a.type == TYPE_INT || a.type == TYPE_STRING))) {
            return to_boolean(a) == to_boolean(b);
        }

        return false;
    }

    switch (a.type) {
        case TYPE_INT:
            return a.value.int_val == b.value.int_val;
        case TYPE_STRING:
            if (a.value.str_val == NULL && b.value.str_val == NULL) {
                return true;
            }
            if (a.value.str_val == NULL || b.value.str_val == NULL) {
                return false;
            }
            return strcmp(a.value.str_val, b.value.str_val) == 0;
        case TYPE_BOOLEAN:
            return a.value.bool_val == b.value.bool_val;
        case TYPE_NULL:
            return true;
        default:
            return false;
    }
}

static bool less_than(Value a, Value b) {
    if (a.type != b.type) {
        if ((a.type == TYPE_INT && b.type == TYPE_STRING) ||
            (a.type == TYPE_STRING && b.type == TYPE_INT)) {
            return to_int(a) < to_int(b);
        }

        if (a.type == TYPE_BOOLEAN || b.type == TYPE_BOOLEAN) {
            return to_boolean(a) < to_boolean(b);
        }
    }

    switch (a.type) {
        case TYPE_INT:
            return a.value.int_val < b.value.int_val;
        case TYPE_STRING:
            if (a.value.str_val == NULL && b.value.str_val == NULL) {
                return false; // Equal, not less than
            }
            if (a.value.str_val == NULL) {
                return true; // NULL is less than any string
            }
            if (b.value.str_val == NULL) {
                return false; // Any string is not less than NULL
            }
            return strcmp(a.value.str_val, b.value.str_val) < 0;
        case TYPE_BOOLEAN:
            return a.value.bool_val < b.value.bool_val;
        case TYPE_NULL:
            return false; // Null is not less than null
        default:
            return false;
    }
}

static bool greater_than(Value a, Value b) {
    return !equals(a, b) && !less_than(a, b);
}

static void print(Value value) {
    switch (value.type) {
        case TYPE_INT:
            printf("%d", value.value.int_val);
            break;
        case TYPE_STRING:
            printf("%s", value.value.str_val ? value.value.str_val : "");
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
}

static void free_value(Value value) {
    if (value.type == TYPE_STRING && value.value.str_val != NULL) {
        free(value.value.str_val);
    }
}

ValueHandler* value_handler_new(void) {
    ValueHandler* handler = (ValueHandler*)malloc(sizeof(ValueHandler));
    if (!handler) return NULL;

    handler->create_int = create_int;
    handler->create_string = create_string;
    handler->create_boolean = create_boolean;
    handler->create_null = create_null;
    handler->to_int = to_int;
    handler->to_string = to_string;
    handler->to_boolean = to_boolean;
    handler->equals = equals;
    handler->less_than = less_than;
    handler->greater_than = greater_than;
    handler->print = print;
    handler->free = free_value;

    return handler;
}

void value_handler_free(ValueHandler* handler) {
    free(handler);
}