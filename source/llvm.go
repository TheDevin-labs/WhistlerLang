package main

import (
        "fmt"
        "os"
        "os/exec"
        "path/filepath"
        "runtime"
        "strings"
)

const termuxBase = "/data/data/com.termux/files"
const termuxUsr  = termuxBase + "/usr"
const termuxLib  = termuxUsr + "/lib"
const termuxInc  = termuxUsr + "/include"
const termuxClang = termuxUsr + "/bin/clang"

func isTermux() bool {
        _, err := os.Stat(termuxUsr)
        return err == nil
}

func isAndroid() bool {
        if isTermux() {
                return true
        }
        data, err := os.ReadFile("/proc/version")
        if err != nil {
                return false
        }
        return strings.Contains(strings.ToLower(string(data)), "android")
}

func findClang(goos string) (string, error) {
        candidates := []string{
                termuxClang,
                "/usr/bin/clang",
                "/usr/local/bin/clang",
                "/opt/homebrew/bin/clang",
                "/opt/homebrew/opt/llvm/bin/clang",
                "C:\\Program Files\\LLVM\\bin\\clang.exe",
                "C:\\Program Files (x86)\\LLVM\\bin\\clang.exe",
        }
        for _, p := range candidates {
                if info, err := os.Stat(p); err == nil && !info.IsDir() {
                        return p, nil
                }
        }
        for _, dir := range filepath.SplitList(os.Getenv("PATH")) {
                for _, name := range []string{"clang", "clang.exe"} {
                        full := filepath.Join(dir, name)
                        if info, err := os.Stat(full); err == nil && !info.IsDir() {
                                return full, nil
                        }
                }
        }
        return "", fmt.Errorf("clang not found")
}

func CompileToArm64(irFile, outputFile string) error {
        arch   := runtime.GOARCH
        goos   := runtime.GOOS
        termux := isTermux()
        android := isAndroid()

        if arch != "arm64" && arch != "amd64" {
                return fmt.Errorf("unsupported architecture: %s", arch)
        }

        clang, err := findClang(goos)
        if err != nil {
                msg := "clang not found. Please install LLVM:\n"
                switch {
                case termux:
                        msg += "  pkg install clang"
                case goos == "darwin":
                        msg += "  brew install llvm"
                case goos == "windows":
                        msg += "  Download from https://releases.llvm.org"
                default:
                        msg += "  sudo apt install clang"
                }
                return fmt.Errorf(msg)
        }

        platformName := goos
        if termux {
                platformName = "android/termux"
        } else if android {
                platformName = "android"
        }

        fmt.Printf("  Platform : %s/%s\n", platformName, arch)
        fmt.Printf("  Clang    : %s\n", clang)

        finalOutput := outputFile
        if goos == "windows" && !strings.HasSuffix(finalOutput, ".exe") {
                finalOutput += ".exe"
        }

        var args []string

        if termux {
                args = []string{
                        irFile,
                        "-o", finalOutput,
                        "-O2",
                        "--sysroot=" + termuxUsr,
                        "-I" + termuxInc,
                        "-L" + termuxLib,
                        "-Wl,-rpath=" + termuxLib,
                        "-lm", "-lc",
                        "-rtlib=compiler-rt",
                        "-unwindlib=none",
                        "-fPIE", "-pie",
                        "-Wl,-z,noexecstack",
                        "-Wl,-z,relro",
                        "-Wl,-z,now",
                }
        } else {
                var triple string
                switch {
                case android && arch == "arm64":
                        triple = "aarch64-linux-android"
                case android && arch == "amd64":
                        triple = "x86_64-linux-android"
                case goos == "linux" && arch == "arm64":
                        triple = "aarch64-unknown-linux-gnu"
                case goos == "linux" && arch == "amd64":
                        triple = "x86_64-unknown-linux-gnu"
                case goos == "darwin" && arch == "arm64":
                        triple = "aarch64-apple-macosx12.0"
                case goos == "darwin" && arch == "amd64":
                        triple = "x86_64-apple-macosx12.0"
                case goos == "windows" && arch == "arm64":
                        triple = "aarch64-pc-windows-msvc"
                case goos == "windows" && arch == "amd64":
                        triple = "x86_64-pc-windows-msvc"
                }
                fmt.Printf("  Target   : %s\n", triple)
                args = []string{
                        irFile,
                        "-o", finalOutput,
                        "-target", triple,
                        "-O2",
                        "-lm",
                }
                switch {
                case goos == "darwin":
                        args = append(args, "-lSystem")
                case goos == "windows":
                        args = append(args, "-lmsvcrt")
                }
        }

        cmd := exec.Command(clang, args...)
        cmd.Stdout = os.Stdout
        cmd.Stderr = os.Stderr

        if err := cmd.Run(); err != nil {
                return fmt.Errorf("clang failed: %w", err)
        }

        if termux || android {
                if err := os.Chmod(finalOutput, 0755); err != nil {
                        fmt.Printf("  Warning: could not chmod: %v\n", err)
                }
        }

        return nil
}
