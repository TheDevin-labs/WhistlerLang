# WhistlerLang

**WhistlerLang** is a custom programming language for **learning, prototyping, and creating scripts**.  
—MATE!

---

## Features 

- Interactive **REPL terminal** for rapid experimentation  
- Run `.whlst` scripts with `run <file.whlst>`  
- `say "text"` → prints text in REPL  
- Strong typing for variables declared with `let`  
- Conditional blocks with `if`, `elif`, `else`, and `end`  
- Time commands:  
  - `time.print` → prints current time  
  - `time.set "<FORMAT>" "<PREF>"` → change time format  

---

## Installation 

Clone the repository:

```bash
git clone https://github.com/CoolyDucks/WhistlerLang
cd WhistlerLang
./build.sh
```

---

## REPL Commands

- `quit` / `exit` → exit WhistlerLang  
- `run <file.whlst>` → run a WhistlerLang script  
- `say "text"` → print a string  
- `time.print` → show current time  
- `time.set "<FORMAT>" "<PREF>"` → change time format  
- `help` → show help menu  

---

## Example Usage 

```
say "Hello World from WhistlerLang!"
let user = "Alice"
say "Welcome " + user + " to Syntexly That Beautiful"
time.print
time.set "{date} {hou}:{min}:{sec}" "ms"
time.print
let score = 85
if score >= 90
    say "Grade: A"
elif score >= 80
    say "Grade: B"
else
    say "Grade: C or below"
end
```
# Licence 

```
BSD 3-Clause License

Copyright (c) 2026, CoolyDucks Del taco 

Redistribution and use in source and binary forms, with or without
modification, are permitted provided that the following conditions are met:

1. Redistributions of source code must retain the above copyright notice, this
   list of conditions and the following disclaimer.

2. Redistributions in binary form must reproduce the above copyright notice,
   this list of conditions and the following disclaimer in the documentation
   and/or other materials provided with the distribution.

3. Neither the name of the copyright holder nor the names of its
   contributors may be used to endorse or promote products derived from
   this software without specific prior written permission.

THIS SOFTWARE IS PROVIDED BY THE COPYRIGHT HOLDERS AND CONTRIBUTORS "AS IS"
AND ANY EXPRESS OR IMPLIED WARRANTIES, INCLUDING, BUT NOT LIMITED TO, THE
IMPLIED WARRANTIES OF MERCHANTABILITY AND FITNESS FOR A PARTICULAR PURPOSE ARE
DISCLAIMED. IN NO EVENT SHALL THE COPYRIGHT HOLDER OR CONTRIBUTORS BE LIABLE
FOR ANY DIRECT, INDIRECT, INCIDENTAL, SPECIAL, EXEMPLARY, OR CONSEQUENTIAL
DAMAGES (INCLUDING, BUT NOT LIMITED TO, PROCUREMENT OF SUBSTITUTE GOODS OR
SERVICES; LOSS OF USE, DATA, OR PROFITS; OR BUSINESS INTERRUPTION) HOWEVER
CAUSED AND ON ANY THEORY OF LIABILITY, WHETHER IN CONTRACT, STRICT LIABILITY,
OR TORT (INCLUDING NEGLIGENCE OR OTHERWISE) ARISING IN ANY WAY OUT OF THE USE
OF THIS SOFTWARE, EVEN IF ADVISED OF THE POSSIBILITY OF SUCH DAMAGE.
```
---



# Join Us in Our Community in Zulip

https://whistlerlang.zulipchat.com

## GitHub 

https://github.com/CoolyDucks/WhistlerLang
