# WhistlerLang Changelog

---

## 26.0 Nightly 2 — Stability Release

**Released by TheDevin-labs**

This release focuses entirely on compiler stability and bug fixes.
No new language features — just making everything that exists work correctly.

### Bug fixes

**Android / Termux**
- Fixed `SIGSYS: bad system call` crash caused by `exec.LookPath` using `faccessat2` syscall which Android's seccomp filter blocks. Replaced with `os.Stat` based clang detection that uses `fstatat` instead which Android allows.
- Fixed Termux linker errors (`cannot open Scrt1.o`, `crti.o`, `crtbeginS.o`). Root cause was passing `-target aarch64-unknown-linux-gnu` to Termux's clang which expects its own native Bionic target. Fix: do not pass `-target` on Termux — clang already knows its own target.
- Fixed `-I /path` flag warning from clang. Changed to `-I/path` with no space which is the correct form.
- Fixed `cannot find library -lgcc` and `-lgcc_s` — Termux uses `compiler-rt` not `libgcc`. Added `-rtlib=compiler-rt -unwindlib=none` flags.

**Codegen**
- Fixed LLVM IR format string size mismatch. `%g\0A\00` is 4 bytes not 5, `%g\00` is 3 bytes not 4. All format string declarations and `getelementptr` references now use correct sizes.
- Fixed `undefined variable: f` — function parameters were not being registered in the variable map before the function body was generated. `genFunction` now saves outer scope, creates fresh inner scope, and registers all params before touching the body.
- Fixed `syntax error: unexpected keyword case, expected }` — `KnownUseStatement` case was being inserted outside the switch block by faulty Python patching. Entire `codegen.go` rewritten cleanly from scratch.
- Fixed `codegen.go` `BlockrockStatement` case misplacement from repeated patch operations.

**Parser**
- Fixed `Parse error: expected ASSIGN got ILLEGAL ":"` on type annotations. Root cause: type keywords like `int`, `float`, `byte` are tokenized as `TOKEN_TYPE_INT` etc, not `TOKEN_IDENT`. `parseLet` and `parseFn` now explicitly accept all nine type tokens after `:`.
- Fixed function parameter type annotations not parsing correctly in strict mode.

**Linker / Platform**
- Fixed triple mismatch between `codegen.go` hardcoded triple and `-target` flag passed by `llvm.go`. Triple is now determined once by `resolveTriple()` in `llvm.go` and passed into `NewCodeGen(triple)` so IR header and clang always agree.
- Added Windows amd64 support (`x86_64-pc-windows-msvc`, `-lmsvcrt`, `.exe` suffix auto-added).
- Added amd64 support for all platforms alongside existing arm64.

---

## 26.0 Nightly 1 — Initial Release

**Released by TheDevin-labs**

First release of WhistlerLang. Full compiler pipeline from `.wh` source to native binary via LLVM IR.

### Features
- Lexer, parser, AST, codegen, LLVM compilation pipeline
- Static typing with optional type annotations
- `say` keyword for output
- `fn` functions with implicit return
- `if / elif / else` conditionals
- `for i in range(n)` and `for item in array` loops
- `let x: type = value` variable declarations
- `blockrock { } panic { }` error handling
- `_knownuse { }` safety bypass block
- `--strict` compiler flag for critical systems
- Built-in math: `sin cos sqrt log exp pow abs ceil floor round`
- Built-in stats: `mean sum min max std variance median len`
- Built-in linalg: `dot cross norm transpose det inverse rank zeros ones identity`
- `csv.open()` and `csv.line()` with auto type detection
- `byte` and `bytes` types with hex literals `0xFF`
- `array` and `matrix` types with stack allocation in IR
- Android/Termux, Linux, macOS, Windows support
- Interactive shell with `help`, `llvm`, `llvm --strict`, `quit/exit`
- Makefile with interactive OS/architecture menus

---

*WhistlerLang — TheDevin-labs*
