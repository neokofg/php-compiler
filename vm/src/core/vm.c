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
#include "../../includes/vm.h"
#include "../../includes/context.h"
#include <stdlib.h>
#include <stdio.h>
#include <string.h>

VM* vm_new(void) {
    VM* vm = (VM*)malloc(sizeof(VM));
    if (!vm) return NULL;

    vm->memory_manager = memory_manager_new();
    vm->value_handler = value_handler_new();
    vm->stack_manager = stack_manager_new(vm->value_handler);
    vm->error_handler = error_handler_new();
    vm->opcode_handler = opcode_handler_new();

    vm->context = vm_context_new();
    if (!vm->context) {
        vm_free(vm);
        return NULL;
    }

    vm_context_set_handlers(vm->context, vm->value_handler, vm->stack_manager, vm->error_handler);

    vm->running = false;
    vm->last_status = STATUS_SUCCESS;
    vm->debug_mode = false;

    vm_register_opcode_handler(vm, OP_LOAD_CONST, handle_load_const);
    vm_register_opcode_handler(vm, OP_PRINT, handle_print);
    vm_register_opcode_handler(vm, OP_HALT, handle_halt);
    vm_register_opcode_handler(vm, OP_POP, handle_pop);

    vm_register_opcode_handler(vm, OP_ADD, handle_add);
    vm_register_opcode_handler(vm, OP_SUB, handle_sub);
    vm_register_opcode_handler(vm, OP_MUL, handle_mul);
    vm_register_opcode_handler(vm, OP_DIV, handle_div);

    vm_register_opcode_handler(vm, OP_CONCAT, handle_concat);

    vm_register_opcode_handler(vm, OP_STORE_VAR, handle_store_var);
    vm_register_opcode_handler(vm, OP_LOAD_VAR, handle_load_var);

    vm_register_opcode_handler(vm, OP_JUMP, handle_jump);
    vm_register_opcode_handler(vm, OP_JUMP_IF_FALSE, handle_jump_if_false);

    vm_register_opcode_handler(vm, OP_GT, handle_gt);
    vm_register_opcode_handler(vm, OP_LT, handle_lt);
    vm_register_opcode_handler(vm, OP_EQ, handle_eq);
    vm_register_opcode_handler(vm, OP_NOT, handle_not);

    vm_register_opcode_handler(vm, OP_AND, handle_and);
    vm_register_opcode_handler(vm, OP_OR, handle_or);

    vm_register_opcode_handler(vm, OP_INC, handle_inc);
    vm_register_opcode_handler(vm, OP_DEC, handle_dec);
    vm_register_opcode_handler(vm, OP_POST_INC, handle_post_inc);
    vm_register_opcode_handler(vm, OP_POST_DEC, handle_post_dec);

    vm_register_opcode_handler(vm, OP_MOD, handle_mod);

    vm_register_opcode_handler(vm, OP_GTE, handle_gte);
    vm_register_opcode_handler(vm, OP_LTE, handle_lte);
    vm_register_opcode_handler(vm, OP_IDENTITY_EQ, handle_identity_eq);
    vm_register_opcode_handler(vm, OP_IDENTITY_NE, handle_identity_ne);

    vm_register_opcode_handler(vm, OP_BIT_AND, handle_bit_and);
    vm_register_opcode_handler(vm, OP_BIT_OR, handle_bit_or);
    vm_register_opcode_handler(vm, OP_BIT_XOR, handle_bit_xor);
    vm_register_opcode_handler(vm, OP_BIT_NOT, handle_bit_not);

    vm_register_opcode_handler(vm, OP_LSHIFT, handle_lshift);
    vm_register_opcode_handler(vm, OP_RSHIFT, handle_rshift);

    vm_register_opcode_handler(vm, OP_ASSIGN_ADD, handle_assign_add);
    vm_register_opcode_handler(vm, OP_ASSIGN_SUB, handle_assign_sub);
    vm_register_opcode_handler(vm, OP_ASSIGN_MUL, handle_assign_mul);
    vm_register_opcode_handler(vm, OP_ASSIGN_DIV, handle_assign_div);
    vm_register_opcode_handler(vm, OP_ASSIGN_MOD, handle_assign_mod);
    vm_register_opcode_handler(vm, OP_ASSIGN_CONCAT, handle_assign_concat);

    vm_register_opcode_handler(vm, OP_BREAK, handle_break);
    vm_register_opcode_handler(vm, OP_CONTINUE, handle_continue);

    return vm;
}

void vm_free(VM* vm) {
    if (!vm) return;

    if (vm->context) {
        vm_context_free(vm->context);
    }

    if (vm->opcode_handler) opcode_handler_free(vm->opcode_handler);
    if (vm->error_handler) error_handler_free(vm->error_handler);
    if (vm->stack_manager) stack_manager_free(vm->stack_manager);
    if (vm->value_handler) value_handler_free(vm->value_handler);
    if (vm->memory_manager) memory_manager_free(vm->memory_manager);

    free(vm);
}

status_t vm_execute(VM* vm, byte_t* bytecode, size_t bytecode_len, Value* constants, size_t constants_len) {
    if (!vm || !bytecode) return STATUS_ERROR;

    vm_context_set_bytecode(vm->context, bytecode, bytecode_len);
    vm_context_set_constants(vm->context, constants, constants_len);

    vm->running = true;
    vm->last_status = STATUS_SUCCESS;

    while (vm->running && vm->context->ip < vm->context->bytecode_len) {
        status_t status = vm_execute_instruction(vm);
        if (status != STATUS_SUCCESS) {
            vm->last_status = status;
            vm->running = false;
            return status;
        }
    }

    return STATUS_SUCCESS;
}

status_t vm_execute_instruction(VM* vm) {
    if (!vm || !vm->running) return STATUS_ERROR;
    if (vm->context->ip >= vm->context->bytecode_len) {
        vm->running = false;
        return STATUS_SUCCESS;
    }

    byte_t opcode = vm->context->bytecode[vm->context->ip++];

    if (vm->debug_mode) {
        printf("DEBUG: ip=%zu, opcode=0x%02X (%s)\n",
               vm->context->ip-1, opcode,
               vm->opcode_handler->get_opcode_name(opcode));
    }

    return vm->opcode_handler->execute(vm->context, opcode);
}

void vm_reset(VM* vm) {
    if (!vm) return;

    if (vm->context) {
        vm_context_reset(vm->context);
    }

    if (vm->stack_manager) {
        vm->stack_manager->reset();
    }

    vm->running = false;
    vm->last_status = STATUS_SUCCESS;
}

bool vm_is_running(VM* vm) {
    return vm && vm->running;
}

status_t vm_get_status(VM* vm) {
    return vm ? vm->last_status : STATUS_ERROR;
}

void vm_set_debug_mode(VM* vm, bool debug_mode) {
    if (vm) vm->debug_mode = debug_mode;
}

void vm_register_opcode_handler(VM* vm, byte_t opcode, OpcodeHandlerFunc handler) {
    if (vm && vm->opcode_handler) {
        vm->opcode_handler->register_handler(opcode, handler);
    }
}

void vm_set_user_data(VM* vm, void* user_data) {
    if (vm && vm->context) {
        vm_context_set_user_data(vm->context, user_data);
    }
}

void* vm_get_user_data(VM* vm) {
    return (vm && vm->context) ? vm->context->user_data : NULL;
}