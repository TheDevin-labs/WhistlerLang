#!/data/data/com.termux/files/usr/bin/bash
set -e
OUT_DIR=build
mkdir -p $OUT_DIR

# Linux and Unix ELF maker (for linuuuuux users)
GOOS=linux GOARCH=amd64 go build -o $OUT_DIR/WhistlerLang-linux-amd64 ./source
GOOS=linux GOARCH=386 go build -o $OUT_DIR/WhistlerLang-linux-x86 ./source
GOOS=linux GOARCH=arm64 go build -o $OUT_DIR/WhistlerLang-linux-arm64 ./source
GOOS=linux GOARCH=arm go build -o $OUT_DIR/WhistlerLang-linux-arm ./source

# Windows (hmmmm...fake dev?)
GOOS=windows GOARCH=amd64 go build -o $OUT_DIR/WhistlerLang-win-x64.exe ./source
GOOS=windows GOARCH=386 go build -o $OUT_DIR/WhistlerLang-win-x86.exe ./source

# MacOS (coming for homebrew repo's)
GOOS=darwin GOARCH=amd64 go build -o $OUT_DIR/WhistlerLang-mac-x64 ./source
GOOS=darwin GOARCH=arm64 go build -o $OUT_DIR/WhistlerLang-mac-arm ./source

# BSD
GOOS=freebsd GOARCH=amd64 go build -o $OUT_DIR/WhistlerLang-freebsd-amd64 ./source
GOOS=openbsd GOARCH=386 go build -o $OUT_DIR/WhistlerLang-openbsd-x86 ./source
