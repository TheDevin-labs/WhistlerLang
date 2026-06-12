package main

import (
	"bufio"
	"fmt"
	"os"
	"path/filepath"
	"strings"
	"runtime"
)

const VERSION = "26.0"
const AUTHOR = "TheDevin-labs"

func main() {
	fmt.Printf("WhistlerLang %s, The Language by %s\n", VERSION, AUTHOR)
	fmt.Println("Commands: help | llvm | llvm --strict | quit/exit")
	scanner := bufio.NewScanner(os.Stdin)
	for {
		fmt.Print(">>> ")
		if !scanner.Scan() {
			break
		}
		input := strings.TrimSpace(scanner.Text())
		switch strings.ToLower(input) {
		case "":
			continue
		case "help":
			printHelp()
		case "llvm":
			runLLVM(false)
		case "llvm --strict":
			runLLVM(true)
		case "quit", "exit":
			fmt.Println("Goodbye!")
			os.Exit(0)
		default:
			fmt.Printf("Unknown command: \"%s\". Type \"help\" for available commands.\n", input)
		}
	}
}

func printHelp() {
	fmt.Println()
	fmt.Println("╔══════════════════════════════════════════════════════╗")
	fmt.Println("║         WhistlerLang 26.0 - Help                    ║")
	fmt.Println("╠══════════════════════════════════════════════════════╣")
	fmt.Println("║  COMMANDS                                            ║")
	fmt.Println("║    help       Show this help screen                  ║")
	fmt.Println("║    llvm           Compile .wh file                   ║")
	fmt.Println("║    llvm --strict  Compile with strict mode checks    ║")
	fmt.Println("║    quit/exit  Exit WhistlerLang                      ║")
	fmt.Println("╠══════════════════════════════════════════════════════╣")
	fmt.Println("║  SYNTAX EXAMPLES                                     ║")
	fmt.Println("║                                                      ║")
	fmt.Println("║  -- This is a comment                                ║")
	fmt.Println("║                                                      ║")
	fmt.Println("║  Variables:                                          ║")
	fmt.Println("║    let x = 10                                        ║")
	fmt.Println("║    let pi = 3.14159                                  ║")
	fmt.Println("║    let name = \"WhistlerLang\"                         ║")
	fmt.Println("║                                                      ║")
	fmt.Println("║  Print:                                              ║")
	fmt.Println("║    say \"Hello World\"                                 ║")
	fmt.Println("║    say x                                             ║")
	fmt.Println("║                                                      ║")
	fmt.Println("║  Functions:                                          ║")
	fmt.Println("║    fn add(a, b) -> {                                 ║")
	fmt.Println("║        a + b                                         ║")
	fmt.Println("║    }                                                 ║")
	fmt.Println("║                                                      ║")
	fmt.Println("║  Conditionals:                                       ║")
	fmt.Println("║    if x > 5 { } elif x == 3 { } else { }            ║")
	fmt.Println("║                                                      ║")
	fmt.Println("║  Loops:                                              ║")
	fmt.Println("║    for i in range(10) { }                            ║")
	fmt.Println("║                                                      ║")
	fmt.Println("║  Arrays & Matrices:                                  ║")
	fmt.Println("║    let nums = [1, 2, 3]                              ║")
	fmt.Println("║    let mat  = [[1, 2], [3, 4]]                       ║")
	fmt.Println("║                                                      ║")
	fmt.Println("║  Built-in:                                           ║")
	fmt.Println("║    Math:   sin cos sqrt log exp abs pow              ║")
	fmt.Println("║    Stats:  mean std variance median sum min max      ║")
	fmt.Println("║    Linalg: dot cross inverse transpose det norm      ║")
	fmt.Println("╚══════════════════════════════════════════════════════╝")
	fmt.Println()
}

func runLLVM(strict bool) {
	files, err := filepath.Glob("*.wh")
	if err != nil || len(files) == 0 {
		fmt.Println("Error: No .wh file found in the current directory.")
		fmt.Println("Place your WhistlerLang source file (.wh) in the same folder and try again.")
		return
	}
	if len(files) > 1 {
		fmt.Println("Multiple .wh files found. Please keep only one:")
		for _, f := range files {
			fmt.Println("  -", f)
		}
		return
	}
	sourceFile := files[0]
	fmt.Printf("Found: %s\n", sourceFile)
	fmt.Println("Compiling...")

	src, err := os.ReadFile(sourceFile)
	if err != nil {
		fmt.Println("Error reading file:", err)
		return
	}
	lexer := NewLexer(string(src))
	tokens := lexer.Tokenize()
	parser := NewParser(tokens)
	program, err := parser.Parse()
	if err != nil {
		fmt.Println("Parse error:", err)
		return
	}
	triple := resolveTriple(runtime.GOOS, runtime.GOARCH, isTermux(), isAndroid())
	if strict {
		checker := NewStrictChecker()
		errs := checker.Check(program)
		if len(errs) > 0 {
			fmt.Print(FormatStrictErrors(errs))
			return
		}
		fmt.Println("  Strict mode: all checks passed ✓")
	}
	cg := NewCodeGen(triple)
	ir, err := cg.Generate(program)
	if err != nil {
		fmt.Println("Codegen error:", err)
		return
	}
	irFile := strings.TrimSuffix(sourceFile, ".wh") + ".ll"
	if err := os.WriteFile(irFile, []byte(ir), 0644); err != nil {
		fmt.Println("Error writing IR:", err)
		return
	}
	outputFile := strings.TrimSuffix(sourceFile, ".wh")
	if err := CompileToArm64(irFile, outputFile); err != nil {
		fmt.Println("Compilation error:", err)
		return
	}
	fmt.Printf("✓ Compiled: ./%s\n", outputFile)
}
