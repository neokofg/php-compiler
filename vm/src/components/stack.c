/* Licensed under GNU GPL v3. See LICENSE file for details. */
#include "../../includes/interfaces/stack_manager.h"
#include <stdio.h>
#include <stdlib.h>

static Value stack[STACK_SIZE];
static int sp = -1;
static ValueHandler* value_handler_ref = NULL;

static void push(Value value) {
    if (sp >= STACK_SIZE - 1) {
        fprintf(stderr, "Stack overflow\n");
        exit(EXIT_FAILURE);
    }
    stack[++sp] = value;
}

static Value pop() {
    if (sp < 0) {
        fprintf(stderr, "Stack underflow\n");
        exit(EXIT_FAILURE);
    }
    return stack[sp--];
}

static Value peek(int offset) {
    int index = sp - offset;
    if (index < 0 || index > sp) {
        fprintf(stderr, "Stack peek out of bounds\n");
        exit(EXIT_FAILURE);
    }
    return stack[index];
}

static void swap(int a, int b) {
    int idx_a = sp - a;
    int idx_b = sp - b;

    if (idx_a < 0 || idx_a > sp || idx_b < 0 || idx_b > sp) {
        fprintf(stderr, "Stack swap out of bounds\n");
        exit(EXIT_FAILURE);
    }

    Value temp = stack[idx_a];
    stack[idx_a] = stack[idx_b];
    stack[idx_b] = temp;
}

static void dup() {
    if (sp < 0) {
        fprintf(stderr, "Stack underflow in dup\n");
        exit(EXIT_FAILURE);
    }
    if (sp >= STACK_SIZE - 1) {
        fprintf(stderr, "Stack overflow in dup\n");
        exit(EXIT_FAILURE);
    }

    stack[sp + 1] = stack[sp];
    sp++;
}

static void rotate(int n) {
    if (n <= 1 || sp < n - 1) {
        return;
    }

    Value top = stack[sp];
    for (int i = sp; i > sp - n + 1; i--) {
        stack[i] = stack[i - 1];
    }
    stack[sp - n + 1] = top;
}

static int size() {
    return sp + 1;
}

static bool is_empty() {
    return sp < 0;
}

static bool is_full() {
    return sp >= STACK_SIZE - 1;
}

static void reset() {
    sp = -1;
}

static void print() {
    printf("Stack (bottom to top): ");
    if (sp < 0) {
        printf("<empty>\n");
        return;
    }

    for (int i = 0; i <= sp; i++) {
        Value value = stack[i];

        if (i > 0) {
            printf(", ");
        }

        switch (value.type) {
            case TYPE_INT:
                printf("%d", value.value.int_val);
                break;
            case TYPE_STRING:
                printf("\"%s\"", value.value.str_val ? value.value.str_val : "");
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
    printf("\n");
}

StackManager* stack_manager_new(ValueHandler* value_handler) {
    StackManager* manager = (StackManager*)malloc(sizeof(StackManager));
    if (!manager) return NULL;

    value_handler_ref = value_handler;
    sp = -1;

    manager->push = push;
    manager->pop = pop;
    manager->peek = peek;
    manager->swap = swap;
    manager->dup = dup;
    manager->rotate = rotate;
    manager->size = size;
    manager->is_empty = is_empty;
    manager->is_full = is_full;
    manager->reset = reset;
    manager->print = print;

    return manager;
}

void stack_manager_free(StackManager* manager) {
    free(manager);
}