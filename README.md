## It's still PHP. It just compiles. It's just async. It's just fast.

### Hello everyone!

This is php compiler project. This project tries to compile interpreted language — PHP — into native code.
Ultimate goal of this project is to be able to compile an entire Laravel application.

---

## 🧭 Roadmap for PHP Compiler (phpc)

> Core idea: Compile PHP source code into IR bytecode and execute it via a custom C-based VM.

---

### ✅ Phase 1: Minimal working compiler (DONE)
- [x] Lexer and parser for basic expressions
- [x] Support for:
  - Integer and string literals
  - Arithmetic operators (`+`, `-`, `*`, `/`)
  - Variables (`$x = ...`)
  - `echo` statements
  - `if/else` with conditionals
- [x] Bytecode compiler in Go
- [x] Virtual machine in C
- [x] CLI tool with `--out` mode to generate executable

---

### 🔄 Phase 2: Expression and control flow extensions (IN PROGRESS)
- [x] Logical operators: `&&`, `||`
- [x] Comparison operators: `==`, `>`, `<`
- [ ] `while`, `for`, `foreach`
- [ ] `switch`, `break`, `continue`

---

### 🧩 Phase 3: Functions & scopes
- [ ] Function declaration and calls
- [ ] Return values
- [ ] Local variables and scoping rules
- [ ] Recursion support

---

### 🧱 Phase 4: Data structures
- [ ] Arrays (indexed and associative)
- [ ] Basic array operations (`[]`, `$arr['key']`)
- [ ] Iteration over arrays (`foreach`)
- [ ] String interpolation

---

### 🏗 Phase 5: Object-oriented features
- [ ] Class declarations
- [ ] Properties and methods
- [ ] Constructors (`__construct`)
- [ ] `this` keyword
- [ ] Inheritance and interfaces
- [ ] Traits
- [ ] Private classes

---

### 📚 Phase 6: Built-in functions and standard library
- [ ] Add native C implementations for core functions (`strlen`, `isset`, etc.)
- [ ] Bridge between compiled PHP and C-implemented internals
- [ ] Add support for include/require
- [ ] Async runtime
- [ ] Multithreaded runtime

---

### ⚙️ Phase 7: Overall PHP support
- [ ] Namespace support (`use`, `namespace`)
- [ ] Composer-style autoloading
- [ ] Exception handling (`try/catch`)
- [ ] Type hints and attributes
- [ ] Anonymous functions / closures
- [ ] Reflection (if possible)
- [ ] Runtime environment setup (DB, routes, config)

---

## 🚀 Ultimate goal
- Compile and run a **real Laravel app** without relying on the PHP interpreter.

---

Feel free to contribute, suggest ideas, or just follow along!  
This is a fun and educational journey into the world of compilers and virtual machines 🚀
