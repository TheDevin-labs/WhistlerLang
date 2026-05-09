package main

import (
	"bufio"
	"fmt"
	"os"
	"strings"
	"time"
)

var banner = `
 __        __  _     _ _     _ _           
 \ \      / / (_)   | (_)   (_) |          
  \ \ /\ / /__ _ ___| |_  __ _| |__  _   _ 
   \ V  V / _ \ / __| | |/ _` + "`" + ` | '_ \| | | |
    \_/\_/_/ \_\__ \_|_|\__,_|_.__/ \__, |
                                     __/ |
                                    |___/ 
WHISTLER
WhistlerLang v1.2 — Syntexly That Beautiful
`

var longHelp = `
WhistlerLang v1.2 — Syntexly That Beautiful
Stable release: recommended for teaching, prototyping, and experimental research.

Why use this stable version?
- Stability: core parser, evaluator, and runtime are tested and organized.
- Predictability: Strong typing rules (except 'say' sugar) reduce runtime surprises.
- Education-first design: syntax and examples are tailored for classrooms.
- Interactivity: REPL for rapid experimentation (like Lisp-style REPL).
- Portability: single-file build and simple object exporter for experiments.

Quickstart:
  run <file.whlst>               Run a script file
  build <file.whlst> <out.o>     Create an object file (simple format)
  info <file.o>                  Show object file metadata
  time.print                     Print current time (formatable)
  time.set "<FORMAT>" "<PREF>"   Change time format (e.g. "{date} {hou}:{min}:{sec}")
  say "Hello World"              Print a string
  math() 1 + 2 * (3 + 4)         One-line math
  math;                          Start math block
    1 + 2
    3 * 4
  end                            End block
  exec "ls -la"                  Execute shell command
  exec-safe "ls"                 Execute allowed safe shell command
  if sum > 10                   Conditional blocks are supported

Type 'help' to see this message again. Type 'quit' or 'exit' to leave the REPL.
`

func printBanner() {
	fmt.Print(banner)
	fmt.Println("Type 'help' for detailed information about this stable release.")
}

func main() {
	wd, _ := os.Getwd()
	_ = GlobalRuntime.SetWorkDir(wd)
	_ = EnsureDirs()
	if len(os.Args) > 1 {
		handleArgs(os.Args[1:])
		return
	}
	printBanner()
	reader := bufio.NewReader(os.Stdin)
	for {
		fmt.Print(">>> ")
		line, err := reader.ReadString('\n')
		if err != nil {
			fmt.Println("input error:", err)
			return
		}
		line = strings.TrimSpace(line)
		if line == "" {
			continue
		}
		if line == "quit" || line == "exit" {
			return
		}
		if line == "help" {
			fmt.Println(longHelp)
			continue
		}
		toks := splitArgs(line)
		if len(toks) == 0 {
			continue
		}
		handleArgs(toks)
	}
}

func handleArgs(args []string) {
	switch args[0] {
	case "run":
		if len(args) < 2 {
			fmt.Println("usage: run <file.whlst>")
			return
		}
		if err := GlobalRuntime.RunScript(args[1]); err != nil {
			fmt.Println("run error:", err)
		}
	case "build":
		if len(args) < 3 {
			fmt.Println("usage: build <file.whlst> <out.o>")
			return
		}
		out, err := GlobalRuntime.BuildScript(args[1], args[2], "")
		if err != nil {
			fmt.Println("build error:", err)
		} else {
			fmt.Println("Built:", out)
		}
	case "info":
		if len(args) < 2 {
			fmt.Println("usage: info <file.o>")
			return
		}
		info, err := GetObjectInfo(args[1])
		if err != nil {
			fmt.Println("info error:", err)
			return
		}
		fmt.Printf("Magic: %s\nVersion: 0x%04x\nArch: %s\nTimestamp: %s\nPayload: %d bytes\n",
			info.Magic, info.Version, info.Arch, time.Unix(info.Timestamp, 0).Format(time.RFC3339), len(info.Payload))
	default:
		line := strings.Join(args, " ")
		nodes, _ := ParseLine(line)
		_ = EvaluateNodes(nodes, GlobalRuntime)
	}
}
