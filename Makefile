.DEFAULT_GOAL := help

help: ## Prints help message.
	@ grep -h -E '^[a-zA-Z_.-]+:.*?## .*$$' $(MAKEFILE_LIST) | awk 'BEGIN {FS = ":.*?## "}; {printf "\033[1m%-30s\033[0m %s\n", $$1, $$2}'

unittest: ## Runs unit tests.
	@ GOOS=js GOARCH=wasm go test --tags=unittest -exec="$(GOROOT)/misc/wasm/go_js_wasm_exec" \
	    -v -coverpkg=./... -coverprofile=coverage.out ./...
	@ go tool cover -func coverage.out
	@ rm coverage.out

WASM_PORT := 9090
WASMNAME_BIN := ./client/assets/logic.wasm
CODE_DIR := ./cmd/wasm

web.setup: ## Sets WASM Go dependencies.
	@ cp "$(GOROOT)/misc/wasm/wasm_exec.js" $(WASM_DIR)

web.build: ## Compiles the app.
	@ GOOS=js GOARCH=wasm CGO_ENABLED=0 \
		go build -o $(WASMNAME_BIN) cmd/wasm/*.go

web.server: ## Run a temp web server.
	@ PORT=$(WASM_PORT) DIR_WEB=./client/ go run cmd/server/main.go
