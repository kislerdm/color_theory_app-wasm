.DEFAULT_GOAL := help

help: ## Prints help message.
	@ grep -h -E '^[a-zA-Z_.-]+:.*?## .*$$' $(MAKEFILE_LIST) | awk 'BEGIN {FS = ":.*?## "}; {printf "\033[1m%-30s\033[0m %s\n", $$1, $$2}'

test:  ## Runs unit tests.
	@ GOARCH=wasm GOOS=js go test --tags=unittest -exec="$(GOROOT)/misc/wasm/go_js_wasm_exec" \
	    -v -coverpkg=./... -coverprofile=coverage.out ./...
	@ go tool cover -func coverage.out
	@ rm coverage.out

WASM_PORT := 9090
WASM_DIR := ./client/assets
CODE_DIR := ./cmd/wasm

web.setup: ## Sets WASM Go dependencies.
	@ cp "$(GOROOT)/misc/wasm/wasm_exec.js" $(WASM_DIR)

web.build: ## Compiles the app.
	@ GOOS=js GOARCH=wasm CGO_ENABLED=0 \
		go build \
		-o $(WASM_DIR)/logic.wasm cmd/wasm/*.go

web.server: ## Run a temp web server.
	@ PORT=$(WASM_PORT) DIR_WEB=./client/ go run cmd/server/main.go
