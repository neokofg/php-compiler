/*
 * PHPС VM - virtual machine for language PHP Compiler on C
 * Copyright (C) 2025 Andrey Vasilev (neokofg)
 *
 * This program is free software: you can redistribute it and/or modify
 * it under the terms of the GNU General Public License as published by
 * the Free Software Foundation, either version 3 of the License, or
 * (at your option) any later version.
 *
 * This program is distributed in the hope that it will be useful,
 * but WITHOUT ANY WARRANTY; without even the implied warranty of
 * MERCHANTABILITY or FITNESS FOR A PARTICULAR PURPOSE.  See the
 * GNU General Public License for more details.
 *
 * You should have received a copy of the GNU General Public License
 * along with this program.  If not, see <https://www.gnu.org/licenses/>.
 */

#include "includes/vm.h"
#include <stdio.h>
#include <stdlib.h>

int main(int argc, char** argv) {
    if (argc < 2) {
        fprintf(stderr, "Usage: %s <bytecode_file> [--debug]\n", argv[0]);
        return 1;
    }

    bool debug_mode = false;
    if (argc > 2 && strcmp(argv[2], "--debug") == 0) {
        debug_mode = true;
    }

    FILE* file = fopen(argv[1], "rb");
    if (!file) {
        fprintf(stderr, "Error: Cannot open file %s\n", argv[1]);
        return 1;
    }

    fseek(file, 0, SEEK_END);
    long file_length = ftell(file);
    fseek(file, 0, SEEK_SET);

    if (file_length <= 0) {
        fprintf(stderr, "Error: Empty bytecode file\n");
        fclose(file);
        return 1;
    }

    // Header format:
    // - uint32_t: bytecode_len
    // - uint32_t: constants_len
    // - bytecode (bytecode_len bytes)
    // - constants (constants_len * sizeof(Value) bytes)

    uint32_t bytecode_len, constants_len;
    if (fread(&bytecode_len, sizeof(bytecode_len), 1, file) != 1 ||
        fread(&constants_len, sizeof(constants_len), 1, file) != 1) {
        fprintf(stderr, "Error: Invalid bytecode file format\n");
        fclose(file);
        return 1;
    }

    byte_t* bytecode = (byte_t*)malloc(bytecode_len);
    Value* constants = (Value*)malloc(constants_len * sizeof(Value));

    if (!bytecode || !constants) {
        fprintf(stderr, "Error: Memory allocation failed\n");
        fclose(file);
        free(bytecode);
        free(constants);
        return 1;
    }

    if (fread(bytecode, 1, bytecode_len, file) != bytecode_len ||
        fread(constants, sizeof(Value), constants_len, file) != constants_len) {
        fprintf(stderr, "Error: Failed to read bytecode or constants\n");
        fclose(file);
        free(bytecode);
        free(constants);
        return 1;
    }

    fclose(file);

    VM* vm = vm_new();
    if (!vm) {
        fprintf(stderr, "Error: Failed to create VM\n");
        free(bytecode);
        free(constants);
        return 1;
    }

    vm_set_debug_mode(vm, debug_mode);

    printf("--- Starting VM execution ---\n");
    status_t status = vm_execute(vm, bytecode, bytecode_len, constants, constants_len);
    printf("--- VM execution completed with status: %d ---\n", status);

    vm_free(vm);
    free(bytecode);

    for (size_t i = 0; i < constants_len; i++) {
        if (constants[i].type == TYPE_STRING && constants[i].value.str_val) {
            free(constants[i].value.str_val);
        }
    }
    free(constants);

    return (status == STATUS_SUCCESS) ? 0 : 1;
}