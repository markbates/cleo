.PHONY: test lint vet fmt check build clean ci help
GO_VERSION := 1.24

# Default target
.DEFAULT_GOAL := help

help: ## Show this help message
	@echo 'Usage: make <target>'
	@echo ''
	@echo 'Targets:'
	@awk 'BEGIN {FS = ":.*?## "} /^[a-zA-Z_-]+:.*?## / {printf "  %-15s %s\n", $$1, $$2}' $(MAKEFILE_LIST)

test: ## Run tests with race detection and coverage
	go test -race -cover -count=1 ./...

test-verbose: ## Run tests with verbose output
	go test -race -cover -count=1 -v ./...

lint: ## Run golangci-lint
	@which golangci-lint > /dev/null || (echo "golangci-lint not found, installing..." && go install github.com/golangci/golangci-lint/cmd/golangci-lint@latest)
	golangci-lint run ./...

vet: ## Run go vet
	go vet ./...

fmt: ## Format code
	go fmt ./...

check: fmt vet test ## Run all checks (fmt, vet, test)

build: ## Build the project
	go build ./...

clean: ## Clean build artifacts
	go clean ./...
	go mod tidy

bench: ## Run benchmarks
	go test -bench=. -benchmem ./...

cov: ## Generate coverage report
	go test -short -coverprofile cover.out ./...
	go tool cover -html cover.out
	go mod tidy -v

ci: check lint build ## Run CI pipeline (check, lint, build)

deps: ## Download dependencies
	go mod download
	go mod tidy

update-deps: ## Update dependencies
	go get -u ./...
	go mod tidy

version: ## Show Go version
	@echo "Go version: $(shell go version)"
	@echo "Required: $(GO_VERSION)"
