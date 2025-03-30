#include <stdio.h>
#include <stdint.h>
#include <stdlib.h>
#include <string.h>
#include <stdbool.h> 
#include "vm.h"

#define VM_ERR_EXIT(msg, ...) do { fprintf(stderr, "VM Runtime Error: " msg "\n", ##__VA_ARGS__); exit(EXIT_FAILURE); } while(0)
#define CHECK_STACK_MIN(min_required) do { \
    if (sp < (min_required - 1)) { \
        VM_ERR_EXIT("Stack underflow. Need %d elements, have %d.", min_required, sp + 1); \
    } \
} while(0)

Value stack[STACK_SIZE];
int sp = -1;

Value variables[VAR_COUNT];

void push(Value val) {
    if (sp >= STACK_SIZE - 1) {
        VM_ERR_EXIT("Stack overflow.");
    }
    stack[++sp] = val;
}

Value pop() {
    CHECK_STACK_MIN(1);
    return stack[sp--];
}

bool is_truthy(Value v) {
    if (v.type == TYPE_INT) {
        return v.value.int_val != 0;
    }
    if (v.type == TYPE_STRING) {
        if (strcmp(v.value.str_val, "0") == 0 || strlen(v.value.str_val) == 0) {
            return false;
        }
        return true;
    }
    return false;
}

int value_to_int(Value v) {
    if (v.type == TYPE_INT) {
        return v.value.int_val;
    }
    if (v.type == TYPE_STRING) {
        return atoi(v.value.str_val);
    }
    // fprintf(stderr, "Warning: Trying to convert unsupported type %d to int (ip=%zu).\n", v.type, ip-1);
    return 0;
}


void run(uint8_t* bytecode, size_t length, Value constants[], size_t constants_len) {
    size_t ip = 0; // instruction pointer

    for (int i = 0; i < VAR_COUNT; i++) {
        variables[i].type = TYPE_INT;
        variables[i].value.int_val = 0;
        // EITHER TODO: variables[i].type = TYPE_NULL;
    }
    sp = -1;

    while (ip < length) {
        uint8_t instr = bytecode[ip++];
        //printf("DEBUG: ip=%zu, instr=0x%02X, sp=%d\n", ip-1, instr, sp);

        switch (instr) {
            case OP_LOAD_CONST: {
                uint8_t const_index = bytecode[ip++];
                if (const_index >= constants_len) {
                    VM_ERR_EXIT("Constant index %u out of bounds (ip=%zu, constants_len=%zu).", const_index, ip-2, constants_len);
                }
                push(constants[const_index]);
                break;
            }
            case OP_STORE_VAR: {
                uint8_t var_index = bytecode[ip++];
                if (var_index >= VAR_COUNT) {
                    VM_ERR_EXIT("Variable index %u out of bounds (ip=%zu, VAR_COUNT=%d).", var_index, ip-2, VAR_COUNT);
                }
                CHECK_STACK_MIN(1);
                Value val = pop();
                variables[var_index] = val;
                break;
            }
            case OP_LOAD_VAR: {
                uint8_t var_index = bytecode[ip++];
                 if (var_index >= VAR_COUNT) {
                     VM_ERR_EXIT("Variable index %u out of bounds (ip=%zu, VAR_COUNT=%d).", var_index, ip-2, VAR_COUNT);
                 }
                push(variables[var_index]);
                break;
            }
            case OP_PRINT: {
                CHECK_STACK_MIN(1);
                Value val = pop();
                if (val.type == TYPE_STRING) {
                    printf("%s\n", val.value.str_val);
                } else if (val.type == TYPE_INT) {
                    printf("%d\n", val.value.int_val);
                } else {
                     printf("null\n");
                }
                break;
            }
            case OP_ADD:
            case OP_SUB:
            case OP_MUL:
            case OP_DIV: {
                CHECK_STACK_MIN(2);
                Value b = pop();
                Value a = pop();

                int int_a = value_to_int(a);
                int int_b = value_to_int(b);
                int result = 0;

                switch (instr) {
                    case OP_ADD: result = int_a + int_b; break;
                    case OP_SUB: result = int_a - int_b; break;
                    case OP_MUL: result = int_a * int_b; break;
                    case OP_DIV:
                        if (int_b == 0) {
                            fprintf(stderr, "Warning: Division by zero on line 0\n"); // Line tracking not available
                            result = 0; // PHP > 8 returns INF/NAN/false, let's return 0 (like false)
                        } else {
                            // TODO: float
                            if (int_a % int_b != 0) {
                                fprintf(stderr, "Warning: Integer division resulting in truncation (%d / %d) on line 0\n", int_a, int_b);
                            }
                            result = int_a / int_b;
                        }
                        break;
                    default:
                        VM_ERR_EXIT("Internal error: Unexpected arithmetic op %d (ip=%zu)", instr, ip-1);
                }
                Value res = {.type = TYPE_INT, .value.int_val = result};
                push(res);
                break;
            }
            case OP_HALT: {
                // printf("DEBUG: Halting VM.\n");
                return;
            }
            case OP_JUMP_IF_FALSE: {
                if (ip + 1 >= length) { VM_ERR_EXIT("Unexpected EOF reading JUMP_IF_FALSE offset (ip=%zu)", ip-1); }
                uint8_t lowByte = bytecode[ip++];
                uint8_t highByte = bytecode[ip++];
                uint16_t offset = (uint16_t)(highByte << 8) | lowByte;

                CHECK_STACK_MIN(1);
                Value cond = pop();
                if (!is_truthy(cond)) {
                    size_t target_ip = ip + offset;
                     if (target_ip > length) {
                        VM_ERR_EXIT("JUMP_IF_FALSE target out of bounds (ip=%zu, offset=%u, target=%zu, length=%zu)", ip-3, offset, target_ip, length);
                     }
                    ip = target_ip;
                }
                break;
            }                        
            case OP_JUMP: {
                if (ip + 1 >= length) { VM_ERR_EXIT("Unexpected EOF reading JUMP offset (ip=%zu)", ip-1); }
                uint8_t lowByte = bytecode[ip++];
                uint8_t highByte = bytecode[ip++];
                int16_t offset = (int16_t)(((uint16_t)highByte << 8) | lowByte);

                intptr_t current_ip_signed = (intptr_t)ip;
                intptr_t target_ip_signed = current_ip_signed + offset;

                if (target_ip_signed < 0 || (size_t)target_ip_signed >= length) {
                    VM_ERR_EXIT("JUMP target out of bounds (ip=%zu, offset=%d, target=%ld, length=%zu)",
                                ip - 3,
                                offset,
                                target_ip_signed,
                                length);
                }

                ip = (size_t)target_ip_signed;
                break;
            }            
            case OP_GT:
            case OP_LT:
            case OP_EQ: {
                CHECK_STACK_MIN(2);
                Value b = pop();
                Value a = pop();
                int result = 0;

                int int_a = value_to_int(a);
                int int_b = value_to_int(b);

                 switch(instr) {
                    case OP_EQ: result = (int_a == int_b); break;
                    case OP_GT: result = (int_a > int_b); break;
                    case OP_LT: result = (int_a < int_b); break;
                    default:
                         VM_ERR_EXIT("Internal error: Unexpected comparison op %d (ip=%zu)", instr, ip-1);
                 }

                Value res = {.type = TYPE_INT, .value.int_val = result};
                push(res);
                break;
            }
            case OP_AND:
            case OP_OR: {
                 CHECK_STACK_MIN(2);
                Value b = pop();
                Value a = pop();
                int result = 0;
                 if (instr == OP_AND) {
                     result = is_truthy(a) && is_truthy(b);
                 } else { // OP_OR
                     result = is_truthy(a) || is_truthy(b);
                 }
                Value res = {.type = TYPE_INT, .value.int_val = result};
                push(res);
                break;
            }
            case OP_POP: {
                CHECK_STACK_MIN(1);
                pop();
                break;
            }                     
            default: {
                VM_ERR_EXIT("Unknown instruction: 0x%02X (ip=%zu)", instr, ip-1);
            }
        } // end switch
    } // end while

    // printf("DEBUG: Reached end of bytecode without HALT.\n");
}

// --- main as COMPILE_AS_EXECUTABLE ---
#ifdef COMPILE_AS_EXECUTABLE
extern uint8_t bytecode[];
extern size_t bytecode_len;
extern Value constants[];
extern size_t constants_len;

int main() {
    if (constants == NULL || bytecode == NULL) {
        fprintf(stderr, "Error: Bytecode or constants not linked correctly.\n");
        return 1;
    }
    run(bytecode, bytecode_len, constants, constants_len);
    return 0;
}
#endif
