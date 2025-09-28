BINARY_NAME=getctx
MAIN_PACKAGE_PATH=./cmd/getctx

BINARY_UNIX_PATH=bin/$(BINARY_NAME)
BINARY_WINDOWS_PATH=bin/$(BINARY_NAME).exe

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