# gjfs Makefile

.PHONY: build test test-verbose test-coverage fmt vet lint clean install run-example help

# Variables
BINARY_NAME := gjfs
BUILD_DIR := ./bin
MAIN_PACKAGE := ./cmd/gjfs
VERSION := $(shell git describe --tags --always --dirty 2>/dev/null || echo "dev")
COMMIT := $(shell git rev-parse --short HEAD 2>/dev/null || echo "none")
BUILD_DATE := $(shell date -u +"%Y-%m-%dT%H:%M:%SZ")
LDFLAGS := -ldflags="-s -w -X main.version=$(VERSION) -X main.commit=$(COMMIT) -X main.buildDate=$(BUILD_DATE)"

# Default target
all: build

## Build the binary
build:
	@echo "Building $(BINARY_NAME)..."
	@mkdir -p $(BUILD_DIR)
	@go build $(LDFLAGS) -o $(BUILD_DIR)/$(BINARY_NAME) $(MAIN_PACKAGE)
	@echo "Built $(BUILD_DIR)/$(BINARY_NAME)"

## Build for multiple platforms
build-all:
	@echo "Building for multiple platforms..."
	@mkdir -p $(BUILD_DIR)
	@GOOS=linux GOARCH=amd64 go build $(LDFLAGS) -o $(BUILD_DIR)/$(BINARY_NAME)-linux-amd64 $(MAIN_PACKAGE)
	@GOOS=linux GOARCH=arm64 go build $(LDFLAGS) -o $(BUILD_DIR)/$(BINARY_NAME)-linux-arm64 $(MAIN_PACKAGE)
	@GOOS=darwin GOARCH=amd64 go build $(LDFLAGS) -o $(BUILD_DIR)/$(BINARY_NAME)-darwin-amd64 $(MAIN_PACKAGE)
	@GOOS=darwin GOARCH=arm64 go build $(LDFLAGS) -o $(BUILD_DIR)/$(BINARY_NAME)-darwin-arm64 $(MAIN_PACKAGE)
	@GOOS=windows GOARCH=amd64 go build $(LDFLAGS) -o $(BUILD_DIR)/$(BINARY_NAME)-windows-amd64.exe $(MAIN_PACKAGE)
	@echo "Built all platforms in $(BUILD_DIR)/"

## Run tests
test:
	@echo "Running tests..."
	@go test ./... -v

## Run tests with coverage
test-coverage:
	@echo "Running tests with coverage..."
	@go test ./... -coverprofile=coverage.out -covermode=atomic
	@go tool cover -html=coverage.out -o coverage.html
	@echo "Coverage report: coverage.html"

## Run benchmarks
bench:
	@go test ./... -bench=. -benchmem

## Format code
fmt:
	@echo "Formatting code..."
	@go fmt ./...

## Vet code
vet:
	@echo "Vetting code..."
	@go vet ./...

## Run linter (requires golangci-lint)
lint:
	@which golangci-lint > /dev/null || (echo "golangci-lint not installed" && exit 1)
	@golangci-lint run ./...

## Tidy dependencies
tidy:
	@echo "Tidying dependencies..."
	@go mod tidy

## Install binary to GOPATH/bin
install:
	@echo "Installing $(BINARY_NAME)..."
	@go install $(LDFLAGS) $(MAIN_PACKAGE)

## Clean build artifacts
clean:
	@echo "Cleaning..."
	@rm -rf $(BUILD_DIR) coverage.out coverage.html

## Run example
run-example: build
	@echo "Running example..."
	@$(BUILD_DIR)/$(BINARY_NAME) -schema ./examples/user-schema.json

## Run example with strict mode
run-example-strict: build
	@$(BUILD_DIR)/$(BINARY_NAME) -schema ./examples/user-schema.json -strict

## Run example with seed
run-example-seed: build
	@$(BUILD_DIR)/$(BINARY_NAME) -schema ./examples/user-schema.json -seed 42

## Run example output to file
run-example-file: build
	@$(BUILD_DIR)/$(BINARY_NAME) -schema ./examples/user-schema.json -output /tmp/example.json
	@echo "Output written to /tmp/example.json"
	@cat /tmp/example.json

## Validate example
validate-example: build
	@echo "Validating example..."
	@$(BUILD_DIR)/$(BINARY_NAME) -schema ./examples/user-schema.json -validate /tmp/example.json

## Generate JSON from schema string
run-string-example: build
	@$(BUILD_DIR)/$(BINARY_NAME) -schema-string '{"type": "object", "properties": {"name": {"type": "string"}, "age": {"type": "integer"}}, "required": ["name"]}'

## Download dependencies
deps:
	@go mod download

## Update dependencies
update-deps:
	@go get -u ./...
	@go mod tidy

## Show help
help:
	@echo "gjfs Makefile"
	@echo ""
	@echo "Usage: make [target]"
	@echo ""
	@echo "Targets:"
	@grep -E '^##' Makefile | sed 's/## /  /'

## Check for common issues
check: fmt vet test
	@echo "All checks passed!"

## Run all CI checks
ci: tidy fmt vet test-coverage
	@echo "CI checks completed!"
