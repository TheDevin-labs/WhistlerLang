BINARY = wl
SRC = $(wildcard *.go)

all: menu

menu:
	@echo "=================================="
	@echo "  WhistlerLang 26.0 Build System  "
	@echo "  by TheDevin-labs                "
	@echo "=================================="
	@echo ""
	@echo "Select target:"
	@echo "  make linux_arm64"
	@echo "  make linux_amd64"
	@echo "  make android"
	@echo "  make macos_arm64"
	@echo "  make macos_amd64"
	@echo "  make windows_amd64"
	@echo "  make clean"
	@echo ""

arch: build

linux_arm64:
	@echo "Building for linux/arm64..."
	GOOS=linux GOARCH=arm64 go build -o $(BINARY) .
	@echo "Done: ./$(BINARY)"

linux_amd64:
	@echo "Building for linux/amd64..."
	GOOS=linux GOARCH=amd64 go build -o $(BINARY) .
	@echo "Done: ./$(BINARY)"

android:
	@echo "Building for android/arm64 (Termux)..."
	GOOS=linux GOARCH=arm64 go build -o $(BINARY) .
	@echo "Done: ./$(BINARY)"

macos_arm64:
	@echo "Building for macOS/arm64..."
	GOOS=darwin GOARCH=arm64 go build -o $(BINARY) .
	@echo "Done: ./$(BINARY)"

macos_amd64:
	@echo "Building for macOS/amd64..."
	GOOS=darwin GOARCH=amd64 go build -o $(BINARY) .
	@echo "Done: ./$(BINARY)"

windows_amd64:
	@echo "Building for Windows/amd64..."
	GOOS=windows GOARCH=amd64 go build -o $(BINARY).exe .
	@echo "Done: ./$(BINARY).exe"

build:
	go build -o $(BINARY) .

clean:
	rm -f $(BINARY) $(BINARY).exe *.ll
	@echo "Cleaned."
