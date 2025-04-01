# 💡 phpc — The PHP Compiler

> **It’s still PHP. It just compiles. It's just fast.**

---

## 👋 What is this?

**`phpc`** is a next-gen compiler for PHP. It compiles your PHP code to bytecode, runs it in a blazing-fast virtual machine (written in C), and aims for full compatibility with existing PHP applications — including Laravel.

This isn’t a reimagining.  
This is still **PHP** — just without the interpreter.

---

## 🚀 Why?

Because we want:

- **Real performance gains** — without rewriting apps in Go, Rust, or C++
- **Portable, single-file binaries** — no FPM, no PHP runtime, no weird setup
- **Developer joy** — write PHP as always, run it like a system language

And maybe — just maybe — make PHP cool again 😎

---

## 🧠 How does it work?

`phpc` uses a **hybrid compilation model**:

- Static parts of your PHP are compiled to **bytecode**
- Dynamic constructs (like `eval`, `__call`, etc) are handled by a **built-in mini interpreter**
- The bytecode is executed in a **custom virtual machine** written in C
- Result: native-speed execution with PHP-level flexibility

---

## 🧭 Roadmap (Short & Sweet)

| Phase | Status | Highlights |
|-------|--------|-----------|
| 🛠 Core Compiler & VM | ✅ Done | Lexing, parsing, bytecode, C-VM |
| 🔁 Control flow & logic | ✅ Done | if/else, loops, logical ops |
| 🔣 Functions & scopes | 🔄 In progress | user-defined functions, recursion |
| 📦 Arrays & structures | 🟡 Soon | arrays, foreach, string handling |
| 🧱 OOP | ⏳ Planned | classes, inheritance, magic methods |
| ⚙️ Runtime features | ⏳ Planned | async, multithreading, autoload |
| 🎯 Laravel-ready | 🚀 Ultimate goal | run Laravel without PHP installed |

---

## ✨ Why should you follow this repo?

Because this is **not a toy project.**  
This is a serious attempt to:

- Turn PHP into a **systems-capable compiled language**
- Deliver **drop-in speed** for legacy and modern codebases
- Build a **dev toolchain** for a language that deserves better

Whether you're curious, nostalgic, or want to contribute —  
**hit that ⭐️, press `watch`, and drop by whenever.**

Your support means everything.

---

## 🧪 Try it out

```bash
git clone https://github.com/neokofg/php-compiler
cd php-compiler
go build -o phpc cmd/phpc/main.go
./phpc examples/hello.php --out hello.bin
```

---
## 🧑‍💻 Want to contribute?

Yes, please.  
Ideas, bug reports, pull requests, code reviews — all welcome.

---
## 📜 License
Licensed under GPLv3.  
Open forever. Free forever. PHP forever.