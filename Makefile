.DEFAULT_GOAL := help

.PHONY: help test build

help: ## Prints help message.
	@ grep -h -E '^[a-zA-Z_.-]+:.*?## .*$$' $(MAKEFILE_LIST) | awk 'BEGIN {FS = ":.*?## "}; {printf "\033[1m%-30s\033[0m %s\n", $$1, $$2}'

test:  ## Runs unit tests.
	@ go test --tags=unittest -v -coverpkg=./... -coverprofile=coverage.out ./...
	@ go tool cover -func coverage.out
	@ rm coverage.out

train:
	@ cd ./internal/colortype/train && \
	    docker-compose run train "python main.py"

generate.model: train ## Re-trains and generates the model object.
	@ cd ./internal/colortype/train && go run --tags=gen main.go

generate.colors: ## Generate Go struct with colors for colorname feature.
	@ cd ./internal/colorname/data && go run --tags=gen main.go

WASM_PORT := 9090
WASM_DIR := ./client/assets
CODE_DIR := ./cmd/wasm

build:
	@ docker run --rm \
        -v $(PWD):/src \
        -w="/src" \
        tinygo/tinygo:0.24.0 tinygo build \
                -target=wasm \
            -o client/assets/logic.wasm cmd/wasm/main.go

setup: ## Sets WASM Go dependencies.
	@ docker run --rm \
              -v $(PWD):/src \
              -w="/src" \
              tinygo/tinygo:0.24.0 \
              /bin/bash -c "cp /usr/local/tinygo/targets/wasm_exec.js client/assets"

web.build:
	@ GOOS=js GOARCH=wasm CGO_ENABLED=0 \
		go build \
		-o $(WASM_DIR)/logic.wasm cmd/wasm/*.go

web.setup:
	@ cp $(GOROOT)/misc/wasm/wasm_exec.js ./client/assets/

web.server: ## Run a temp web server.
	@ PORT=$(WASM_PORT) DIR_WEB=./client/ go run cmd/server/main.go
