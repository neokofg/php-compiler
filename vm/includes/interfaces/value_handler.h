#ifndef VM_VALUE_HANDLER_H
#define VM_VALUE_HANDLER_H

#include "../common.h"

typedef enum {
    TYPE_INT = 0,
    TYPE_STRING = 1,
    TYPE_BOOLEAN = 2,
    TYPE_NULL = 3
} ValueType;

typedef struct {
    ValueType type;
    union {
        int_t int_val;
        char* str_val;
        bool bool_val;
    } value;
} Value;

typedef struct ValueHandler {
    Value (*create_int)(int_t value);
    Value (*create_string)(const char* value);
    Value (*create_boolean)(bool value);
    Value (*create_null)(void);

    int_t (*to_int)(Value value);
    char* (*to_string)(Value value);
    bool (*to_boolean)(Value value);

    bool (*equals)(Value a, Value b);
    bool (*less_than)(Value a, Value b);
    bool (*greater_than)(Value a, Value b);

    void (*print)(Value value);

    void (*free)(Value value);
} ValueHandler;

ValueHandler* value_handler_new(void);

void value_handler_free(ValueHandler* handler);

#endif /* VM_VALUE_HANDLER_H */