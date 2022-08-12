.DEFAULT_GOAL := help

help: ## Prints help message.
	@ grep -h -E '^[a-zA-Z_.-]+:.*?## .*$$' $(MAKEFILE_LIST) | awk 'BEGIN {FS = ":.*?## "}; {printf "\033[1m%-30s\033[0m %s\n", $$1, $$2}'

unittest: ## Runs unit tests.
	@ go test --tags=unittest -v -coverpkg=./... -coverprofile=.coverage_temp ./...
	@ go tool cover -func .coverage_temp

WASM_PORT := 9090
WASM_DIR := ./client/assets
WASMNAME_BIN := $(WASM_DIR)/logic.wasm
#GOROOT := `echo GOROOT`
CODE_DIR := ./cmd/wasm

web.setup: ## Sets WASM Go dependencies.
	@ cp "$(GOROOT)/misc/wasm/wasm_exec.js" $(WASM_DIR)

web.build: ## Compiles the app.
	@ GOOS=js GOARCH=wasm CGO_ENABLED=0 \
		go build -o $(WASMNAME_BIN) cmd/wasm/*.go

web.server: ## Run a temp web server.
	@ PORT=$(WASM_PORT) DIR_WEB=$(WASM_DIR) go run cmd/server/main.go
