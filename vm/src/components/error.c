/* Licensed under GNU GPL v3. See LICENSE file for details. */
#include "../../includes/interfaces/error_handler.h"
#include <stdio.h>
#include <stdlib.h>
#include <string.h>
#include <stdarg.h>

static bool has_error_flag = false;
static char error_message[1024] = {0};
static char source_filename[256] = {0};
static int source_line = 0;
static void (*error_callback)(ErrorSeverity, const char*) = NULL;

static void report_error(ErrorSeverity severity, const char* format, ...) {
    va_list args;
    va_start(args, format);

    char message[1024];
    vsnprintf(message, sizeof(message), format, args);

    has_error_flag = true;
    strncpy(error_message, message, sizeof(error_message) - 1);

    const char* prefix = "INFO";
    switch (severity) {
        case ERROR_FATAL:
            prefix = "FATAL ERROR";
            break;
        case ERROR_RUNTIME:
            prefix = "RUNTIME ERROR";
            break;
        case ERROR_WARNING:
            prefix = "WARNING";
            break;
        case ERROR_INFO:
            prefix = "INFO";
            break;
    }

    if (source_filename[0] != '\0') {
        fprintf(stderr, "%s at %s:%d: %s\n", prefix, source_filename, source_line, message);
    } else {
        fprintf(stderr, "%s: %s\n", prefix, message);
    }

    if (error_callback) {
        error_callback(severity, message);
    }

    va_end(args);

    if (severity == ERROR_FATAL) {
        exit(EXIT_FAILURE);
    }
}

static status_t runtime_error(const char* format, ...) {
    va_list args;
    va_start(args, format);

    char message[1024];
    vsnprintf(message, sizeof(message), format, args);

    report_error(ERROR_RUNTIME, "%s", message);

    va_end(args);

    return STATUS_RUNTIME_ERROR;
}

static void fatal_error(const char* format, ...) {
    va_list args;
    va_start(args, format);

    char message[1024];
    vsnprintf(message, sizeof(message), format, args);

    report_error(ERROR_FATAL, "%s", message);

    va_end(args);

    // This should never be reached due to exit() in report_error
}

static void warning(const char* format, ...) {
    va_list args;
    va_start(args, format);

    char message[1024];
    vsnprintf(message, sizeof(message), format, args);

    report_error(ERROR_WARNING, "%s", message);

    va_end(args);
}

static bool has_error() {
    return has_error_flag;
}

static void clear_error() {
    has_error_flag = false;
    error_message[0] = '\0';
}

static void set_source_info(const char* filename, int line) {
    if (filename) {
        strncpy(source_filename, filename, sizeof(source_filename) - 1);
    } else {
        source_filename[0] = '\0';
    }
    source_line = line;
}

static void set_error_callback(void (*callback)(ErrorSeverity severity, const char* message)) {
    error_callback = callback;
}

ErrorHandler* error_handler_new() {
    ErrorHandler* handler = (ErrorHandler*)malloc(sizeof(ErrorHandler));
    if (!handler) return NULL;

    handler->report = report_error;
    handler->runtime_error = runtime_error;
    handler->fatal_error = fatal_error;
    handler->warning = warning;
    handler->has_error = has_error;
    handler->clear_error = clear_error;
    handler->set_source_info = set_source_info;
    handler->set_error_callback = set_error_callback;

    has_error_flag = false;
    error_message[0] = '\0';
    source_filename[0] = '\0';
    source_line = 0;
    error_callback = NULL;

    return handler;
}

void error_handler_free(ErrorHandler* handler) {
    free(handler);
}