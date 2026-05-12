package main

import (
	"fmt"
	"os"
	"os/exec"
	"runtime"
)

func CompileToArm64(irFile, outputFile string) error {
	arch := runtime.GOARCH
	goos := runtime.GOOS

	if arch != "arm64" {
		return fmt.Errorf(
			"WhistlerLang only compiles on arm64 (detected: %s/%s)\n"+
				"Please run on an arm64 Linux or macOS machine.", goos, arch)
	}
	if goos != "linux" && goos != "darwin" {
		return fmt.Errorf("WhistlerLang only supports Linux and macOS (detected: %s)", goos)
	}

	clang, err := exec.LookPath("clang")
	if err != nil {
		return fmt.Errorf(
			"clang not found. Please install LLVM:\n" +
				"  macOS: brew install llvm\n" +
				"  Linux: sudo apt install clang")
	}

	var triple string
	switch goos {
	case "linux":  triple = "aarch64-unknown-linux-gnu"
	case "darwin": triple = "aarch64-apple-macosx12.0"
	}

	fmt.Printf("  Platform : %s/%s\n", goos, arch)
	fmt.Printf("  Target   : %s\n", triple)
	fmt.Printf("  Clang    : %s\n", clang)

	args := []string{
		irFile,
		"-o", outputFile,
		"-target", triple,
		"-O2",
		"-lm",
	}
	if goos == "darwin" {
		args = append(args, "-lSystem")
	}

	cmd := exec.Command(clang, args...)
	cmd.Stdout = os.Stdout
	cmd.Stderr = os.Stderr

	if err := cmd.Run(); err != nil {
		return fmt.Errorf("clang failed: %w", err)
	}
	return nil
}

