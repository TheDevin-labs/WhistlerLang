SRC_DIR := source
BUILD_DIR := release
BIN := $(BUILD_DIR)/wl
YELLOW := \033[1;33m
RESET := \033[0m

.PHONY: help
help:
	@echo "WhistlerLang v1.2 Commands:"
	@echo "  make run           - Run a WhistlerLang script"
	@echo "  make install repo  - Clone WhistlerLang repo from GitHub"
	@echo "  make dash          - Check and fix source folder structure"
	@echo "  make gnu           - Prepare GNU tools"
	@echo "  make busybox       - Prepare BusyBox tools"
	@echo "  make etcinfo       - Display Whistler word in yellow ASCII art"

.PHONY: build
build:
	@mkdir -p $(BUILD_DIR)
	@echo "Building WhistlerLang..."
	@go build -o $(BIN) $(SRC_DIR)/main.go
	@echo "Build completed: $(BIN)"

.PHONY: run
run:
	@if [ -z "$(file)" ]; then \
		echo "Usage: make run file=<script.whlst>"; \
	else \
		$(BIN) run $(file); \
	fi

.PHONY: install
install:
	@git clone https://github.com/CoolyDucks/WhistlerLang.git
	@echo "Repository cloned successfully"

.PHONY: dash
dash:
	@mkdir -p $(SRC_DIR)/parser $(SRC_DIR)/evaluator $(SRC_DIR)/runtime $(SRC_DIR)/objecting $(SRC_DIR)/utils
	@echo "Folders verified/created"

.PHONY: gnu
gnu:
	@which gcc >/dev/null 2>&1 || { echo "GCC not found, install it first"; exit 1; }
	@which make >/dev/null 2>&1 || { echo "Make not found, install it first"; exit 1; }
	@echo "GNU tools ready"

.PHONY: busybox
busybox:
	@which busybox >/dev/null 2>&1 || { echo "BusyBox not installed"; exit 1; }
	@echo "BusyBox tools ready"

.PHONY: etcinfo
etcinfo:
	@echo -e "$(YELLOW)"
	@echo "██       ██ ███████ ███████ ███████ ███████ ███████"
	@echo "██       ██ ██      ██      ██      ██      ██     "
	@echo "██   █   ██ █████   █████   █████   █████   █████ "
	@echo "██  ███  ██ ██      ██      ██      ██      ██     "
	@echo " ███   ███  ███████ ███████ ███████ ███████ ███████"
	@echo -e "$(RESET)"
