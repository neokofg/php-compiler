# Licensed under GNU GPL v3. See LICENSE file for details.
CC = gcc
CFLAGS = -Wall -Wextra -g
INCLUDES = -I./includes

# Source directories
CORE_DIR = src/core
HANDLERS_DIR = src/handlers
COMPONENTS_DIR = src/components
EXAMPLES_DIR = examples

# Object files
CORE_OBJS = $(CORE_DIR)/vm.o $(CORE_DIR)/context.o $(CORE_DIR)/dispatcher.o
HANDLERS_OBJS = $(HANDLERS_DIR)/arithmetic.o $(HANDLERS_DIR)/core.o $(HANDLERS_DIR)/flow.o $(HANDLERS_DIR)/logic.o $(HANDLERS_DIR)/string.o $(HANDLERS_DIR)/function.o
COMPONENTS_OBJS = $(COMPONENTS_DIR)/value.o $(COMPONENTS_DIR)/memory.o $(COMPONENTS_DIR)/stack.o $(COMPONENTS_DIR)/error.o
COMMON_OBJS = $(CORE_OBJS) $(HANDLERS_OBJS) $(COMPONENTS_OBJS)
MAIN_OBJS = main.o
EXAMPLE_OBJS = $(EXAMPLES_DIR)/simple.o

# Targets
TARGET = vm_exec
EXAMPLE_TARGET = vm_example

all: directories $(TARGET) $(EXAMPLE_TARGET)

# Library for VM core
libvm.a: $(COMMON_OBJS)
	ar rcs $@ $^

# Link the main executable
$(TARGET): $(MAIN_OBJS) libvm.a
	$(CC) $(CFLAGS) -o $@ $(MAIN_OBJS) -L. -lvm

# Link the example executable
$(EXAMPLE_TARGET): $(EXAMPLE_OBJS) libvm.a
	$(CC) $(CFLAGS) -o $@ $(EXAMPLE_OBJS) -L. -lvm

# Compile source files to object files
%.o: %.c
	$(CC) $(CFLAGS) $(INCLUDES) -c $< -o $@

# Create directories if they don't exist
directories:
	mkdir -p $(CORE_DIR) $(HANDLERS_DIR) $(COMPONENTS_DIR) $(EXAMPLES_DIR)

# Clean up
clean:
	rm -f $(COMMON_OBJS) $(MAIN_OBJS) $(EXAMPLE_OBJS) libvm.a $(TARGET) $(EXAMPLE_TARGET)

# Phony targets
.PHONY: all clean directories