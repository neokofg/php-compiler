/* Licensed under GNU GPL v3. See LICENSE file for details. */
#ifndef VM_ERROR_HANDLER_H
#define VM_ERROR_HANDLER_H

#include "../common.h"

typedef enum {
    ERROR_FATAL,      // Fatal error, execution stops
    ERROR_RUNTIME,    // Runtime error, potentially recoverable
    ERROR_WARNING,    // Warning, execution continues
    ERROR_INFO        // Informational message
} ErrorSeverity;

typedef struct ErrorHandler {
    void (*report)(ErrorSeverity severity, const char* format, ...);

    status_t (*runtime_error)(const char* format, ...);

    void (*fatal_error)(const char* format, ...);

    void (*warning)(const char* format, ...);

    bool (*has_error)(void);
    void (*clear_error)(void);

    void (*set_source_info)(const char* filename, int line);

    void (*set_error_callback)(void (*callback)(ErrorSeverity severity, const char* message));
} ErrorHandler;

ErrorHandler* error_handler_new(void);

void error_handler_free(ErrorHandler* handler);

#endif /* VM_ERROR_HANDLER_H */