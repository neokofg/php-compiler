#ifndef VM_COMMON_H
#define VM_COMMON_H

#include <stdint.h>
#include <stddef.h>
#include <stdbool.h>
#include <stdio.h>
#include <stdlib.h>
#include <string.h>

typedef uint8_t byte_t;
typedef int32_t int_t;
typedef int_t status_t;

#define STATUS_SUCCESS 0
#define STATUS_ERROR -1
#define STATUS_RUNTIME_ERROR -2
#define STATUS_STACK_OVERFLOW -3
#define STATUS_STACK_UNDERFLOW -4
#define STATUS_INVALID_OPCODE -5
#define STATUS_OUT_OF_MEMORY -6
#define STATUS_DIVISION_BY_ZERO -7

#include "config.h"

#endif /* VM_COMMON_H */