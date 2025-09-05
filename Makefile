build:
	@go build -o bin/getctx

run: build
	@./bin/getctx