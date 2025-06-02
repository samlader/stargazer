.PHONY: build test clean run lint

BINARY_NAME=stargazer
BUILD_DIR=build

GO=go
GOFMT=gofmt
GOLINT=golangci-lint

build:
	@echo "Building..."
	@mkdir -p $(BUILD_DIR)
	@$(GO) build -o $(BUILD_DIR)/$(BINARY_NAME) .

test:
	@echo "Running tests..."
	@$(GO) test -v ./...

run:
	@echo "Running..."
	@$(GO) run .

fmt:
	@echo "Formatting code..."
	@$(GOFMT) -w .

lint:
	@echo "Running linter..."
	@$(GOLINT) run

deps:
	@echo "Installing dependencies..."
	@$(GO) mod download
	@$(GO) install github.com/golangci/golangci-lint/cmd/golangci-lint@latest
