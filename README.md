# WhistlerLang

A strongly typed scientific scripting language designed for clarity, expressiveness, and fast iteration. Built for data, math, and linear algebra â€” with zero imports required.

```
let measurements = [98.6, 97.9, 99.1, 98.3, 100.2]
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
- [Functions](#functions)
- [Conditionals](#conditionals)
- [Loops](#loops)
- [Comments](#comments)
- [Built-in math](#built-in-math)
- [Built-in statistics](#built-in-statistics)
- [Built-in linear algebra](#built-in-linear-algebra)
- [Full example](#full-example)

---

## Getting started

WhistlerLang files use the `.wl` extension. Run a program with:

```
whistler run myscript.wl
```

---

## Variables

Declare variables with `let`. No type annotation required â€” types are inferred at runtime.

```
let x    = 10
let pi   = 3.14159
let name = "Whistler"
let flag = true
```

---

## Data types

| Type | Example | Description |
|------|---------|-------------|
| `int` | `42` | Whole numbers |
| `float` | `9.81` | Decimal numbers |
| `complex` | `3+2i` | Complex numbers (real + imaginary) |
| `bool` | `true` / `false` | Boolean values |
| `string` | `"hello"` | Text in double quotes |

---

## Output

Use the `say` keyword to print to the console. No parentheses needed.

```
say "Hello World"    -- Hello World
say x                -- 10
say 2 + 2            -- 4
```

---

## Operators

**Arithmetic**

| Operator | Description | Example | Result |
|----------|-------------|---------|--------|
| `+` | Addition | `10 + 5` | `15` |
| `-` | Subtraction | `10 - 5` | `5` |
| `*` | Multiplication | `10 * 5` | `50` |
| `/` | Division | `10 / 5` | `2` |
| `%` | Modulo | `10 % 3` | `1` |
| `^` | Exponentiation | `2 ^ 8` | `256` |

**Comparison** â€” all return `true` or `false`

| Operator | Description |
|----------|-------------|
| `==` | Equal to |
| `!=` | Not equal to |
| `<` | Less than |
| `>` | Greater than |
| `<=` | Less than or equal |
| `>=` | Greater than or equal |

**Logical**

| Operator | Description | Example | Result |
|----------|-------------|---------|--------|
| `and` | Both must be true | `true and false` | `false` |
| `or` | One must be true | `true or false` | `true` |
| `not` | Flips the value | `not true` | `false` |

---

## Arrays

Ordered sequences of values. Declared with square brackets, zero-indexed.

```
let nums  = [1, 2, 3, 4, 5]
let words = ["hi", "world"]

let first = nums[0]   -- 1
let third = nums[2]   -- 3
```

---

## Matrices

Two-dimensional arrays. Access elements with `mat[row][col]`.

```
let mat = [[1, 2, 3],
           [4, 5, 6],
           [7, 8, 9]]

let center = mat[1][1]   -- 5
let corner = mat[0][2]   -- 3
```

---

## Functions

Declare with `fn`, parameters in `()`, and `->`. The last expression in the body is automatically returned â€” no `return` keyword needed.

```
fn add(a, b) -> {
    a + b
}

fn celsius(f) -> {
    (f - 32.0) * 5.0 / 9.0
}

fn greet(name) -> {
    say "Hello"
    say name
}

-- Calling functions
let result = add(3, 4)   -- 7
say celsius(98.6)         -- 37.0
greet("World")            -- Hello \n World
```

---

## Conditionals

Use `if`, `elif`, and `else` for branching. `elif` is short for "else if".

```
let score = 85

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

---

## Loops

**Range loop** â€” iterates `i` from `0` to `n-1`:

```
for i in range(5) {
    say i       -- 0, 1, 2, 3, 4
}
```

**Array loop** â€” iterates directly over elements:

```
let data = [10, 20, 30]
for item in data {
    say item    -- 10, 20, 30
}
```

---

## Comments

WhistlerLang supports single-line comments only. Start a comment with `--`.

```
-- This is a full-line comment
let x = 10  -- This is an inline comment
```

> Multi-line comments are not supported. Prefix each line with `--`.

---

## Built-in math

No imports required. All math functions accept numeric arguments in radians where applicable.

| Function | Description | Example | Result |
|----------|-------------|---------|--------|
| `sin(x)` | Sine | `sin(3.14)` | `â‰ˆ 0.0` |
| `cos(x)` | Cosine | `cos(0.0)` | `1.0` |
| `tan(x)` | Tangent | `tan(0.785)` | `â‰ˆ 1.0` |
| `sqrt(x)` | Square root | `sqrt(16.0)` | `4.0` |
| `log(x)` | Natural log (base e) | `log(2.718)` | `â‰ˆ 1.0` |
| `exp(x)` | e to the power x | `exp(1.0)` | `2.718` |
| `pow(x, y)` | x to the power y | `pow(2.0, 10.0)` | `1024.0` |
| `abs(x)` | Absolute value | `abs(-42.0)` | `42.0` |
| `ceil(x)` | Round up | `ceil(3.2)` | `4.0` |
| `floor(x)` | Round down | `floor(3.8)` | `3.0` |
| `round(x)` | Round to nearest | `round(3.5)` | `4.0` |

---

## Built-in statistics

All statistics functions operate on arrays of numbers.

```
let values = [4.0, 8.0, 15.0, 16.0, 23.0, 42.0]
```

| Function | Description | Result |
|----------|-------------|--------|
| `mean(arr)` | Average of all values | `18.0` |
| `sum(arr)` | Sum of all values | `108.0` |
| `min(arr)` | Smallest value | `4.0` |
| `max(arr)` | Largest value | `42.0` |
| `median(arr)` | Middle value when sorted | `15.5` |
| `variance(arr)` | Variance of the values | â€” |
| `std(arr)` | Standard deviation | â€” |
| `len(arr)` | Number of elements | `6` |

---

## Built-in linear algebra

Vector and matrix operations for scientific computing.

```
let v1 = [1.0, 2.0, 3.0]
let v2 = [4.0, 5.0, 6.0]

let m1 = [[1.0, 2.0], [3.0, 4.0]]
let m2 = [[5.0, 6.0], [7.0, 8.0]]
```

**Vectors**

| Function | Description | Result |
|----------|-------------|--------|
| `dot(v1, v2)` | Dot product | `32.0` |
| `cross(v1, v2)` | Cross product (3D only) | vector |
| `norm(v)` | Magnitude/length | `3.74` |

**Matrices**

| Function | Description | Result |
|----------|-------------|--------|
| `dot(m1, m2)` | Matrix multiplication | matrix |
| `transpose(m)` | Flip rows and columns | matrix |
| `det(m)` | Determinant | `-2.0` |
| `inverse(m)` | Inverse matrix | matrix |
| `zeros(r, c)` | rÃ—c matrix of zeros | matrix |
| `ones(r, c)` | rÃ—c matrix of ones | matrix |
| `identity(n)` | nÃ—n identity matrix | matrix |

---

## Full example

```
say "=== WhistlerLang Science Demo ==="

let measurements = [98.6, 97.9, 99.1, 98.3, 100.2]

say "Measurements:"
say measurements

let avg_temp  = mean(measurements)
let deviation = std(measurements)

say "Average temperature (F):"
say avg_temp

say "Standard deviation:"
say deviation

fn celsius(f) -> {
    (f - 32.0) * 5.0 / 9.0
}

say "Average in Celsius:"
say celsius(avg_temp)

if avg_temp > 99.5 {
    say "Status: Fever detected"
} elif avg_temp > 98.9 {
    say "Status: Slightly elevated"
} else {
    say "Status: Normal"
}
```

---

*WhistlerLang 26.0 â€” TheDevin-labs*
