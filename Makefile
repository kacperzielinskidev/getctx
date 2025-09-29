BINARY_NAME=getctx
MAIN_PACKAGE_PATH=./cmd/getctx

BIN_DIR=bin
DIST_DIR=dist

VERSION ?= $(shell git describe --tags --abbrev=0 2>/dev/null || echo "dev")

BINARY_UNIX_PATH=$(BIN_DIR)/$(BINARY_NAME)
BINARY_WINDOWS_PATH=$(BIN_DIR)/$(BINARY_NAME).exe



ifeq ($(OS),Windows_NT)
	BINARY_PATH=$(BINARY_WINDOWS_PATH)
else
	BINARY_PATH=$(BINARY_UNIX_PATH)
endif

build:
	@go build -o $(BINARY_PATH) $(MAIN_PACKAGE_PATH)

build-windows:
	@go build -o $(BINARY_WINDOWS_PATH) $(MAIN_PACKAGE_PATH)

run: build
	@$(BINARY_PATH)

run-debug: build
	@$(BINARY_PATH) --debug

profile-cpu: build
	@$(BINARY_PATH) --cpuprofile=cpu.pprof

analyze-cpu:
	@go tool pprof $(BINARY_PATH) cpu.pprof

package-linux-amd64:
	@echo "--> Packaging for Linux AMD64 (wersja $(VERSION))..."
	@mkdir -p $(DIST_DIR)
	@GOOS=linux GOARCH=amd64 go build -o $(DIST_DIR)/$(BINARY_NAME) $(MAIN_PACKAGE_PATH)
	@cd $(DIST_DIR) && tar -czf $(BINARY_NAME)_$(VERSION)_linux.tar.gz $(BINARY_NAME)
	@rm $(DIST_DIR)/$(BINARY_NAME)

package-windows-amd64:
	@echo "--> Packaging for Windows AMD64 (wersja $(VERSION))..."
	@mkdir -p $(DIST_DIR)
	@GOOS=windows GOARCH=amd64 go build -o $(DIST_DIR)/$(BINARY_NAME).exe $(MAIN_PACKAGE_PATH)
	@cd $(DIST_DIR) && zip $(BINARY_NAME)_$(VERSION)_windows.zip $(BINARY_NAME).exe
	@rm $(DIST_DIR)/$(BINARY_NAME).exe

package-darwin-amd64:
	@echo "--> Packaging for macOS AMD64 (wersja $(VERSION))..."
	@mkdir -p $(DIST_DIR)
	@GOOS=darwin GOARCH=amd64 go build -o $(DIST_DIR)/$(BINARY_NAME) $(MAIN_PACKAGE_PATH)
	@cd $(DIST_DIR) && tar -czf $(BINARY_NAME)_$(VERSION)_darwin_amd64.tar.gz $(BINARY_NAME)
	@rm $(DIST_DIR)/$(BINARY_NAME)

package-darwin-arm64:
	@echo "--> Packaging for macOS ARM64 (wersja $(VERSION))..."
	@mkdir -p $(DIST_DIR)
	@GOOS=darwin GOARCH=arm64 go build -o $(DIST_DIR)/$(BINARY_NAME) $(MAIN_PACKAGE_PATH)
	@cd $(DIST_DIR) && tar -czf $(BINARY_NAME)_$(VERSION)_darwin_arm64.tar.gz $(BINARY_NAME)
	@rm $(DIST_DIR)/$(BINARY_NAME)


clean:
	@rm -rf $(BIN_DIR) $(DIST_DIR) *.pprof *.log