package main

import (
	"errors"
	"os"
	"strings"
	"sync"
	"time"
)

func EnsureDirs() error {
	return os.MkdirAll("release", 0755)
}

func FileExists(path string) bool {
	_, err := os.Stat(path)
	return err == nil
}

func EnsureWhlstExtension(path string) error {
	if !strings.HasSuffix(path, ".whlst") {
		return errors.New("file must have .whlst extension")
	}
	return nil
}

func EnsureObjectExtension(path string) error {
	if !strings.HasSuffix(path, ".o") {
		return errors.New("file must have .o extension")
	}
	return nil
}

func SafeWriteFile(path string, data []byte) error {
	dir := ""
	for i := len(path) - 1; i >= 0; i-- {
		if path[i] == '/' || path[i] == '\\' {
			dir = path[:i]
			break
		}
	}
	if dir != "" {
		if err := os.MkdirAll(dir, 0755); err != nil {
			return err
		}
	}
	return os.WriteFile(path, data, 0644)
}

func ReadFile(path string) (string, error) {
	b, err := os.ReadFile(path)
	if err != nil {
		return "", err
	}
	return string(b), nil
}

// --- Time formatting ---

var (
	globalTimeMu     sync.RWMutex
	globalTimeFormat = "{date} {hou}:{min}:{sec}"
	globalTimePref   = "" // e.g. "ms" for milliseconds
)

// SetGlobalTime updates the global time format and preference.
func SetGlobalTime(format, pref string) {
	globalTimeMu.Lock()
	defer globalTimeMu.Unlock()
	if format != "" {
		globalTimeFormat = format
	}
	globalTimePref = pref
}

// PrintTime returns the current time rendered with the global format string.
// Supported placeholders: {date}, {hou}, {min}, {sec}, {ms}, {year}, {mon}, {day}.
func PrintTime() string {
	globalTimeMu.RLock()
	format := globalTimeFormat
	pref := globalTimePref
	globalTimeMu.RUnlock()

	now := time.Now()
	if pref == "ms" {
		// millisecond precision already handled via {ms} placeholder
		_ = pref
	}

	r := format
	r = strings.ReplaceAll(r, "{year}", now.Format("2006"))
	r = strings.ReplaceAll(r, "{mon}", now.Format("01"))
	r = strings.ReplaceAll(r, "{day}", now.Format("02"))
	r = strings.ReplaceAll(r, "{date}", now.Format("2006-01-02"))
	r = strings.ReplaceAll(r, "{hou}", now.Format("15"))
	r = strings.ReplaceAll(r, "{min}", now.Format("04"))
	r = strings.ReplaceAll(r, "{sec}", now.Format("05"))
	r = strings.ReplaceAll(r, "{ms}", now.Format("000"))
	return r
}
