BINARY_NAME=getctx
BINARY_PATH=bin/$(BINARY_NAME)
MAIN_PACKAGE_PATH=./cmd/getctx

build:
	@go build -o $(BINARY_PATH) $(MAIN_PACKAGE_PATH)

run: build
	@$(BINARY_PATH)

run-debug: build
	@$(BINARY_PATH) --debug

profile-cpu: build
	@$(BINARY_PATH) --cpuprofile=cpu.pprof

analyze-cpu:
	@go tool pprof $(BINARY_PATH) cpu.pprof