/* Licensed under GNU GPL v3. See LICENSE file for details. */
#ifndef VM_STACK_MANAGER_H
#define VM_STACK_MANAGER_H

#include "../common.h"
#include "value_handler.h"

typedef struct StackManager {
    void (*push)(Value value);
    Value (*pop)(void);
    Value (*peek)(int offset);

    void (*swap)(int a, int b);
    void (*dup)(void);
    void (*rotate)(int n);

    int (*size)(void);
    bool (*is_empty)(void);
    bool (*is_full)(void);

    void (*reset)(void);

    void (*print)(void);
} StackManager;

StackManager* stack_manager_new(ValueHandler* value_handler);

void stack_manager_free(StackManager* manager);

#endif /* VM_STACK_MANAGER_H */