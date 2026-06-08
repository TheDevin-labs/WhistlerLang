
# WhistlerLang
A statically typed, high-performance scientific language built for data, math, embedded systems, and critical computing. Compiles natively via LLVM IR on Linux, macOS, Windows, and Android.
```
let measurements: array = [98.6, 97.9, 99.1, 98.3, 100.2]
say mean(measurements)   -- 98.82
say std(measurements)    -- 0.81
```
---
## Table of contents
- [Getting started](#getting-started)
- [Variables](#variables)
- [Data types](#data-types)
- [Output](#output)
- [Operators](#operators)
- [Arrays](#arrays)
- [Matrices](#matrices)
- [Bytes](#bytes)
- [Functions](#functions)
- [Conditionals](#conditionals)
- [Loops](#loops)
- [Comments](#comments)
- [Error handling (Blockrock)](#error-handling-blockrock)
- [Safety Bypass (_knownuse)](#safety-bypass-_knownuse)
- [CSV Data](#csv-data)
- [Built-in math](#built-in-math)
- [Built-in statistics](#built-in-statistics)
- [Built-in linear algebra](#built-in-linear-algebra)
- [Strict mode](#strict-mode)
- [Full example](#full-example)
---
## Getting started
WhistlerLang files use the `.wh` extension. Compiles natively via LLVM IR. Run your code using the compiler:
```bash
llvm science_demo.wh
```
For critical systems, you can enforce compile-time safety checks using:
```bash
llvm --strict science_demo.wh
```
## Variables
Declare variables with let. WhistlerLang supports optional static type annotations. In --strict mode, type annotations are strictly required on every variable.
```
-- with type annotation (required in strict mode)
let x: int    = 10
let pi: float  = 3.14159
let name: string = "WhistlerLang"
let flag: bool  = true
-- without annotation (normal mode only)
let y = 42
let z = 9.81
```
## Data types
WhistlerLang has nine native types covering scientific, systems, and data use cases:

| Type | Example | Description |
| :--- | :--- | :--- |
| int | 42 | Whole numbers |
| float | 9.81 | Decimal numbers |
| complex | 3+2i | Complex numbers (real + imaginary) |
| bool | true / false | Boolean values |
| string | "hello" | Text in double quotes |
| byte | 0xFF / 255 | Single byte (0–255) |
| bytes | [0xFF, 0x00] | Byte array / raw buffer |
| array | [1.0, 2.0, 3.0] | Ordered sequences of values |
| matrix | [[1,2],[3,4]] | 2D array for linear algebra |

## Output
Use the say keyword to print to the console. No parentheses needed — say is a keyword, not a function. Works seamlessly with scalars, strings, arrays, and matrices.
```
say "Hello World"    -- Hello World
say x                -- prints value of x
say 2 + 2            -- evaluates then prints -> 4
say nums             -- prints: [1, 2, 3]
say mat              -- prints: [[1, 2], [3, 4]]
```
## Operators
**Arithmetic**

| Operator | Description | Example | Result |
| :--- | :--- | :--- | :--- |
| + | Addition | 10 + 5 | 15 |
| - | Subtraction | 10 - 5 | 5 |
| * | Multiplication | 10 * 5 | 50 |
| / | Division | 10 / 5 | 2 |
| % | Modulo (remainder) | 10 % 3 | 1 |
| ^ | Exponentiation | 2 ^ 8 | 256 | <br> **Comparison** — all return true or false
| Operator | Description | Example | Result |
| :--- | :--- | :--- | :--- |
| == | Equal to | 5 == 5 | true |
| != | Not equal to | 5 != 4 | true |
| < | Less than | 3 < 5 | true |
| > | Greater than | 5 > 3 | true |
| <= | Less than or equal | 3 <= 3 | true |
| > | Greater than or equal | 5 >= 5 | true | <br> **Logical**
| Operator | Description | Example | Result |
| :--- | :--- | :--- | :--- |
| and | Logical AND — both must be true | true and false | false |
| or | Logical OR — one must be true | true or false | true |
| not | Logical NOT — flips the value | not true | false |

## Arrays
Ordered sequences of values declared with square brackets. Indexing is zero-indexed. You can run statistical and linear algebra functions directly on arrays.
```
let nums: array = [1.0, 2.0, 3.0, 4.0, 5.0]
let first = nums[0]   -- 1.0
let third = nums[2]   -- 3.0
say mean(nums)        -- 3.0
```
## Matrices
Matrices are first-class two-dimensional arrays. Access elements using mat[row][col]. All built-in linear algebra functions operate directly on this type.
```
let mat: matrix = [[1.0, 2.0, 3.0],
                   [4.0, 5.0, 6.0],
                   [7.0, 8.0, 9.0]]
let center = mat[1][1]   -- row 1, col 1 -> 5.0
let corner = mat[0][2]   -- row 0, col 2 -> 3.0
```
## Bytes
WhistlerLang provides explicit low-level system types for byte and bytes buffers. Byte values can be declared via hex literals or decimals.
```
-- single byte
let b1: byte = 0xFF    -- hex literal
let b2: byte = 255     -- same value, decimal
-- byte array / buffer
let buf: bytes = [0x48, 0x65, 0x6C, 0x6C, 0x6F]
say buf   -- [72, 101, 108, 108, 111]
```
> **Note:** An array containing only hex literals is automatically inferred as bytes by the compiler.
> 
## Functions
Declare with fn, parameters, and ->. The last expression in the body is implicitly returned without a return keyword. In strict mode, all parameter types and the return type must be explicitly annotated.
```
-- normal mode (no annotations required)
fn add(a, b) -> {
    a + b
}
-- strict mode (full annotations required)
fn multiply(a: float, b: float) -> float {
    a * b
}
fn celsius(f: float) -> float {
    (f - 32.0) * 5.0 / 9.0
}
let result = add(3, 4)
say multiply(6.0, 7.0)
```
## Conditionals
Use if, elif, and else for branching logic. You can chain multiple elif blocks.
```
let score: int = 85
if score >= 90 {
    say "Grade: A"
} elif score >= 75 {
    say "Grade: B"
} elif score >= 60 {
    say "Grade: C"
} else {
    say "Grade: F"
}
```
## Loops
**Range loop** — iterates n times from 0 through n-1:
```
for i in range(5) {
    say i       -- 0, 1, 2, 3, 4
}
```
**Array loop** — iterates directly over elements:
```
let data: array = [10.0, 20.0, 30.0]
for item in data {
    say item    -- 10.0, 20.0, 30.0
}
```
## Comments
WhistlerLang supports single-line comments only. A comment begins with -- and extends to the end of the line.
```
-- This is a full-line comment
let x: int = 10   -- inline comment
```
> Multi-line or block comments are not supported. Prefix each line with --.
> 
## Error handling (Blockrock)
WhistlerLang utilizes the blockrock system for handling runtime errors. Wrap risky operations inside a blockrock block. If a failure occurs, execution instantly jumps to the panic block.
```
blockrock {
    let data = csv.open("sensors.csv")
    say data
} panic {
    say "Failed to read sensor data"
}
```
> **Strict mode rule:** Every blockrock statement must include a panic handler block, otherwise it will trigger a compile-time error.
> 
## Safety Bypass (_knownuse)
An escape hatch designed for expert low-level code. Wrapping code within a _knownuse block instructs the compiler to completely bypass strict type and error checking.
```
_knownuse {
    let raw = 0xFF
    let unsafe_val = 42
    say raw
}
```
## CSV Data
Built-in native CSV parsing with automatic cell type detection (checks for int first, then float, then string). CSV operations must always be safely wrapped inside a blockrock block.
```
blockrock {
    let table = csv.open("data.csv")   -- Reads entire file into a matrix
    say table
} panic {
    say "Could not open data.csv"
}
```

| Function | Returns | Description |
| :--- | :--- | :--- |
| csv.open(path) | matrix | Reads the entire file where each row is an array of auto-detected values. |
| csv.line(path) | array of arrays | Reads line-by-line; yields one array per row for custom iteration loops. | <br> ## Built-in math <br> Always available globally with zero imports required.
| Function | Description | Example | Result |
| :--- | :--- | :--- | :--- |
| sin(x) | Sine (radians) | sin(3.14) | ≈ 0.0 |
| cos(x) | Cosine (radians) | cos(0.0) | 1.0 |
| tan(x) | Tangent (radians) | tan(0.785) | ≈ 1.0 |
| sqrt(x) | Square root | sqrt(16.0) | 4.0 |
| log(x) | Natural logarithm | log(2.718) | ≈ 1.0 |
| exp(x) | e raised to power x | exp(1.0) | 2.718 |
| pow(x, y) | x to the power of y | pow(2.0, 10.0) | 1024.0 |
| abs(x) | Absolute value | abs(-42.0) | 42.0 |
| ceil(x) | Round up | ceil(3.2) | 4.0 |
| floor(x) | Round down | floor(3.8) | 3.0 |
| round(x) | Round to nearest | round(3.5) | 4.0 |

## Built-in statistics
Statistics functions operate natively on arrays of numbers.
```
let values: array = [4.0, 8.0, 15.0, 16.0, 23.0, 42.0]
```

| Function | Description | Result |
| :--- | :--- | :--- |
| mean(arr) | Average of all values | 18.0 |
| sum(arr) | Sum of all values | 108.0 |
| min(arr) | Smallest value | 4.0 |
| max(arr) | Largest value | 42.0 |
| median(arr) | Middle value when sorted | 15.5 |
| variance(arr) | Spread of the values (Variance) | — |
| std(arr) | Standard deviation | — |
| len(arr) | Count of elements | 6 |

## Built-in linear algebra
Matrix and vector operations for scientific operations. No imports required.
```
let v1: array  = [1.0, 2.0, 3.0]
let v2: array  = [4.0, 5.0, 6.0]
let m1: matrix = [[1.0, 2.0], [3.0, 4.0]]
let m2: matrix = [[5.0, 6.0], [7.0, 8.0]]
```

| Function | Description | Notes / Results |
| :--- | :--- | :--- |
| dot(v1, v2) | Dot product of two vectors | 32.0 |
| cross(v1, v2) | Cross product (3D vectors) | Returns array |
| norm(v) | Magnitude of a vector | 3.74 |
| dot(m1, m2) | Matrix multiplication | Also works on matrices |
| transpose(m) | Flip rows and columns | Returns matrix |
| det(m) | Determinant of a square matrix | -2.0 |
| inverse(m) | Inverse of a matrix | Returns matrix |
| rank(m) | Rank of a matrix | — |
| zeros(r, c) | r \times c matrix of zeros | zeros(3,3) |
| ones(r, c) | r \times c matrix of ones | ones(2,4) |
| identity(n) | n \times n identity matrix | identity(3) |

## Strict mode
Running your code with llvm --strict targets critical safety systems (embedded, kernels, aviation). In strict mode, all warnings are upgraded to explicit **compile errors**:
 1. **Explicit Annotations:** Every variable declaration must feature an explicit type block: let x: int = 10.
 2. **Strict Signatures:** Every function must declare its parameter types and return type explicitly.
 3. **Guaranteed Blockrock:** Every blockrock instance requires an accompanying panic handler block.
 4. **Controlled Bypass:** The _knownuse block acts as the singular escape hatch for raw operations.
## Full example
```
say "=== WhistlerLang Science Demo ==="
-- temperature measurements
let temps: array = [98.6, 97.9, 99.1, 98.3, 100.2]
say "Measurements:"
say temps
let avg: float = mean(temps)
let dev: float = std(temps)
say "Average (F):"
say avg
say "Std deviation:"
say dev
-- convert to Celsius
fn celsius(f: float) -> float {
    (f - 32.0) * 5.0 / 9.0
}
say "Average in Celsius:"
say celsius(avg)
-- load extra data from CSV
blockrock {
    let extra = csv.open("extra.csv")
    say extra
} panic {
    say "No extra data found"
}
-- status check
if avg > 99.5 {
    say "Status: Fever"
} elif avg > 98.9 {
    say "Status: Slightly elevated"
} else {
    say "Status: Normal"
}
```
*WhistlerLang — TheDevin-labs*
```
```
