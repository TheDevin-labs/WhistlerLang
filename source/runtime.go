package main

import (
	"context"
	"errors"
	"os"
	"os/exec"
	"path/filepath"
	"runtime"
	"sync"
	"time"
)

type ValKind int

const (
	ValNone ValKind = iota
	ValNumber
	ValString
)

type Value struct {
	Kind ValKind
	Num  float64
	Str  string
}

type RuntimeEnv struct {
	WorkDir   string
	Variables map[string]Value
	mu        sync.RWMutex
}

var GlobalRuntime = &RuntimeEnv{
	WorkDir:   ".",
	Variables: map[string]Value{},
}

func (r *RuntimeEnv) SetWorkDir(dir string) error {
	if dir == "" {
		return errors.New("empty")
	}
	r.mu.Lock()
	r.WorkDir = dir
	r.mu.Unlock()
	return nil
}

func (r *RuntimeEnv) RunScript(path string) error {
	abs := path
	if !filepath.IsAbs(path) {
		abs = filepath.Join(r.WorkDir, path)
	}
	nodes, err := ParseFileToNodes(abs)
	if err != nil {
		return err
	}
	return EvaluateNodes(nodes, r)
}

func (r *RuntimeEnv) BuildScript(path, out, arch string) (string, error) {
	abs := path
	if !filepath.IsAbs(path) {
		abs = filepath.Join(r.WorkDir, path)
	}
	return BuildScriptToObject(abs, out, arch)
}

func (r *RuntimeEnv) ExecShell(cmdline string, safe bool) (string, error) {
	if cmdline == "" {
		return "", errors.New("empty")
	}
	if safe {
		parts := splitArgs(cmdline)
		if len(parts) == 0 {
			return "", errors.New("empty")
		}
		allowed := map[string]bool{"echo": true, "ls": true, "date": true}
		if !allowed[parts[0]] {
			return "", errors.New("not allowed")
		}
	}
	ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
	defer cancel()
	var shell, flag string
	if runtime.GOOS == "windows" {
		shell = "cmd"
		flag = "/C"
	} else {
		shell = "sh"
		flag = "-c"
	}
	c := exec.CommandContext(ctx, shell, flag, cmdline)
	c.Env = os.Environ()
	c.Dir = r.WorkDir
	out, err := c.CombinedOutput()
	return string(out), err
}
