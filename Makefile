.DEFAULT_GOAL := help

.PHONY: help test build

help: ## Prints help message.
	@ grep -h -E '^[a-zA-Z_.-]+:.*?## .*$$' $(MAKEFILE_LIST) | awk 'BEGIN {FS = ":.*?## "}; {printf "\033[1m%-30s\033[0m %s\n", $$1, $$2}'

test:  ## Runs unit tests.
	@ GOARCH=wasm GOOS=js go test --tags=unittest -exec="$(GOROOT)/misc/wasm/go_js_wasm_exec" \
	    -v -coverpkg=./... -coverprofile=coverage.out ./...
	@ go tool cover -func coverage.out
	@ rm coverage.out

train:
	@ cd ./internal/colortype/train && \
	    docker-compose run train "python main.py"

generate.model: train ## Re-trains and generates the model object.
	@ cd ./internal/colortype/train && go run main.go

WASM_PORT := 9090
WASM_DIR := ./client/assets
CODE_DIR := ./cmd/wasm

build:
	@ docker run --rm \
        -v $(PWD):/src \
        -w="/src" \
        tinygo/tinygo:0.24.0 tinygo build -o client/assets/logic.wasm -target=wasm cmd/wasm/*.go

web.setup: ## Sets WASM Go dependencies.
	@ cp "$(GOROOT)/misc/wasm/wasm_exec.js" $(WASM_DIR)

web.build: ## Compiles the app.
	@ GOOS=js GOARCH=wasm CGO_ENABLED=0 \
		go build \
		-o $(WASM_DIR)/logic.wasm cmd/wasm/*.go

web.server: ## Run a temp web server.
	@ PORT=$(WASM_PORT) DIR_WEB=./client/ go run cmd/server/main.go
