/* Licensed under GNU GPL v3. See LICENSE file for details. */
#ifndef VM_MEMORY_MANAGER_H
#define VM_MEMORY_MANAGER_H

#include "../common.h"

typedef struct MemoryManager {
    void* (*allocate)(size_t size);

    void* (*reallocate)(void* ptr, size_t old_size, size_t new_size);

    void (*free)(void* ptr);

    void (*collect_garbage)(void);

    size_t (*get_allocated_bytes)(void);
    size_t (*get_allocation_count)(void);

    void (*print_stats)(void);
} MemoryManager;

MemoryManager* memory_manager_new(void);

void memory_manager_free(MemoryManager* manager);

#endif /* VM_MEMORY_MANAGER_H */