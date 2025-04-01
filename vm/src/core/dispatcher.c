/* Licensed under GNU GPL v3. See LICENSE file for details. */
#include "../../includes/interfaces/opcode_handler.h"

#define MAX_OPCODES 256

typedef struct {
    OpcodeHandlerFunc handlers[MAX_OPCODES];
    const char* opcode_names[MAX_OPCODES];
} OpcodeHandlerImpl;

static OpcodeHandlerImpl impl;

static void register_handler(byte_t opcode, OpcodeHandlerFunc handler) {
    impl.handlers[opcode] = handler;
}

static status_t execute(VMContext* context, byte_t opcode) {
    if (!context) return STATUS_ERROR;

    if (!impl.handlers[opcode]) {
        if (context->error_handler) {
            context->error_handler->runtime_error("Invalid opcode: 0x%02X at position %zu", opcode, context->ip - 1);
        }
        return STATUS_INVALID_OPCODE;
    }

    return impl.handlers[opcode](context);
}

static const char* get_opcode_name(byte_t opcode) {
    return impl.opcode_names[opcode] ? impl.opcode_names[opcode] : "UNKNOWN";
}

static bool is_opcode_valid(byte_t opcode) {
    return impl.handlers[opcode] != NULL;
}

static void reset_handlers(void) {
    memset(&impl, 0, sizeof(impl));
}

static void init_opcode_names(void) {
    impl.opcode_names[OP_LOAD_CONST] = "LOAD_CONST";
    impl.opcode_names[OP_PRINT] = "PRINT";
    impl.opcode_names[OP_HALT] = "HALT";
    impl.opcode_names[OP_POP] = "POP";

    impl.opcode_names[OP_ADD] = "ADD";
    impl.opcode_names[OP_SUB] = "SUB";
    impl.opcode_names[OP_MUL] = "MUL";
    impl.opcode_names[OP_DIV] = "DIV";

    impl.opcode_names[OP_CONCAT] = "CONCAT";

    impl.opcode_names[OP_STORE_VAR] = "STORE_VAR";
    impl.opcode_names[OP_LOAD_VAR] = "LOAD_VAR";

    impl.opcode_names[OP_JUMP] = "JUMP";
    impl.opcode_names[OP_JUMP_IF_FALSE] = "JUMP_IF_FALSE";

    impl.opcode_names[OP_GT] = "GT";
    impl.opcode_names[OP_LT] = "LT";
    impl.opcode_names[OP_EQ] = "EQ";
    impl.opcode_names[OP_NOT] = "NOT";

    impl.opcode_names[OP_AND] = "AND";
    impl.opcode_names[OP_OR] = "OR";

    impl.opcode_names[OP_INC] = "INC";
    impl.opcode_names[OP_DEC] = "DEC";
    impl.opcode_names[OP_POST_INC] = "POST_INC";
    impl.opcode_names[OP_POST_DEC] = "POST_DEC";

    impl.opcode_names[OP_MOD] = "MOD";

    impl.opcode_names[OP_GTE] = "GTE";
    impl.opcode_names[OP_LTE] = "LTE";
    impl.opcode_names[OP_IDENTITY_EQ] = "IDENTITY_EQ";
    impl.opcode_names[OP_IDENTITY_NE] = "IDENTITY_NE";

    impl.opcode_names[OP_BIT_AND] = "BIT_AND";
    impl.opcode_names[OP_BIT_OR] = "BIT_OR";
    impl.opcode_names[OP_BIT_XOR] = "BIT_XOR";
    impl.opcode_names[OP_BIT_NOT] = "BIT_NOT";
    impl.opcode_names[OP_LSHIFT] = "LSHIFT";
    impl.opcode_names[OP_RSHIFT] = "RSHIFT";

    impl.opcode_names[OP_ASSIGN_ADD] = "ASSIGN_ADD";
    impl.opcode_names[OP_ASSIGN_SUB] = "ASSIGN_SUB";
    impl.opcode_names[OP_ASSIGN_MUL] = "ASSIGN_MUL";
    impl.opcode_names[OP_ASSIGN_DIV] = "ASSIGN_DIV";
    impl.opcode_names[OP_ASSIGN_MOD] = "ASSIGN_MOD";
    impl.opcode_names[OP_ASSIGN_CONCAT] = "ASSIGN_CONCAT";

    impl.opcode_names[OP_BREAK] = "BREAK";
    impl.opcode_names[OP_CONTINUE] = "CONTINUE";
}

OpcodeHandler* opcode_handler_new(void) {
    static bool initialized = false;
    if (!initialized) {
        memset(&impl, 0, sizeof(impl));
        init_opcode_names();
        initialized = true;
    }

    OpcodeHandler* handler = (OpcodeHandler*)malloc(sizeof(OpcodeHandler));
    if (!handler) return NULL;

    handler->register_handler = register_handler;
    handler->execute = execute;
    handler->get_opcode_name = get_opcode_name;
    handler->is_opcode_valid = is_opcode_valid;
    handler->reset_handlers = reset_handlers;

    return handler;
}

void opcode_handler_free(OpcodeHandler* handler) {
    free(handler);
}