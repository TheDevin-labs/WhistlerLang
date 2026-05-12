SHELL := /bin/bash
REPO_URL := https://github.com/TheDevin-labs/WhistlerLang
BINARY := wl
BUILD_DIR := release

.PHONY: all
all: menu

menu:
	@clear
	@echo "======================================"
	@echo "       WHISTLER BUILD MANAGER         "
	@echo "======================================"
	@echo "1) Build for Linux"
	@echo "2) Build for Darwin (MacOS)"
	@echo "3) Build for BSD"
	@echo "4) Pull Latest Updates (The Risky Button)"
	@echo "5) Exit"
	@echo "======================================"
	@read -p "Select an option [1-5]: " opt; \
	if [ "$$opt" == "1" ]; then $(MAKE) arch OS=linux; \
	elif [ "$$opt" == "2" ]; then $(MAKE) arch OS=darwin; \
	elif [ "$$opt" == "3" ]; then $(MAKE) arch OS=freebsd; \
	elif [ "$$opt" == "4" ]; then $(MAKE) update; \
	elif [ "$$opt" == "5" ]; then exit 0; \
	else echo "Invalid option"; sleep 1; $(MAKE) menu; fi

arch:
	@echo "--------------------------------------"
	@echo " Target OS: $(OS)"
	@echo "--------------------------------------"
	@echo "1) arm64"
	@echo "2) amd64"
	@echo "3) i386"
	@echo "4) Back to Main Menu"
	@echo "--------------------------------------"
	@read -p "Select Architecture [1-4]: " aopt; \
	if [ "$$aopt" == "1" ]; then $(MAKE) build OS=$(OS) ARCH=arm64; \
	elif [ "$$aopt" == "2" ]; then $(MAKE) build OS=$(OS) ARCH=amd64; \
	elif [ "$$aopt" == "3" ]; then $(MAKE) build OS=$(OS) ARCH=386; \
	elif [ "$$aopt" == "4" ]; then $(MAKE) menu; \
	else echo "Invalid option"; sleep 1; $(MAKE) arch OS=$(OS); fi

build:
	@echo "Preparing $(OS)_$(ARCH)..."
	@mkdir -p $(BUILD_DIR)
	@GOOS=$(OS) GOARCH=$(ARCH) go build -o $(BUILD_DIR)/$(BINARY) ./source
	@cd $(BUILD_DIR) && tar -cJf WhistlerLang_$(OS)_$(ARCH).tar.xz $(BINARY)
	@rm $(BUILD_DIR)/$(BINARY)
	@echo "Done: $(BUILD_DIR)/WhistlerLang_$(OS)_$(ARCH).tar.xz"

update:
	@echo "Warning: Fetching latest updates from GitHub..."
	@git remote add origin $(REPO_URL) 2>/dev/null || true
	@git fetch --all
	@git reset --hard origin/main
	@echo "Update complete."
	@sleep 2
	@$(MAKE) menu
