#include "../includes/vm.h"
#include <stdio.h>

int main() {
    byte_t bytecode[] = {
        OP_LOAD_CONST, 0,
        OP_LOAD_CONST, 1,
        OP_ADD,
        OP_PRINT,
        OP_HALT
    };

    Value constants[2];
    constants[0].type = TYPE_INT;
    constants[0].value.int_val = 10;
    constants[1].type = TYPE_INT;
    constants[1].value.int_val = 20;

    VM* vm = vm_new();
    if (!vm) {
        fprintf(stderr, "Failed to create VM\n");
        return 1;
    }

    vm_set_debug_mode(vm, true);

    printf("--- Starting VM execution ---\n");
    status_t status = vm_execute(vm, bytecode, sizeof(bytecode), constants, 2);
    printf("--- VM execution completed with status: %d ---\n", status);

    vm_free(vm);

    return (status == STATUS_SUCCESS) ? 0 : 1;
}