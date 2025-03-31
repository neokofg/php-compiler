/* Licensed under GNU GPL v3. See LICENSE file for details. */
#include "../../includes/interfaces/memory_manager.h"
#include <stdlib.h>
#include <stdio.h>

static size_t allocated_bytes = 0;
static size_t allocation_count = 0;

static void* allocate_memory(size_t size) {
    void* memory = malloc(size);
    if (memory) {
        allocated_bytes += size;
        allocation_count++;
    }
    return memory;
}

static void* reallocate_memory(void* ptr, size_t old_size, size_t new_size) {
    void* memory = realloc(ptr, new_size);
    if (memory) {
        allocated_bytes = allocated_bytes - old_size + new_size;
    }
    return memory;
}

static void free_memory(void* ptr) {
    if (ptr) {
        free(ptr);
    }
}

static void collect_garbage(void) {
    printf("GC: No automatic garbage collection implemented yet\n");
}

static size_t get_allocated_bytes(void) {
    return allocated_bytes;
}

static size_t get_allocation_count(void) {
    return allocation_count;
}

static void print_stats(void) {
    printf("Memory stats: %zu bytes in %zu allocations\n",
           allocated_bytes, allocation_count);
}

MemoryManager* memory_manager_new(void) {
    MemoryManager* manager = (MemoryManager*)malloc(sizeof(MemoryManager));
    if (!manager) return NULL;

    manager->allocate = allocate_memory;
    manager->reallocate = reallocate_memory;
    manager->free = free_memory;
    manager->collect_garbage = collect_garbage;
    manager->get_allocated_bytes = get_allocated_bytes;
    manager->get_allocation_count = get_allocation_count;
    manager->print_stats = print_stats;

    return manager;
}

void memory_manager_free(MemoryManager* manager) {
    free(manager);
}