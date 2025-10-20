# Makefile for TraceTrim project
# Provides portable build, test, and release targets

# Project configuration
PROJECT_NAME := tracetrim
MODULE_NAME := com.github/rethunk-tech/tracetrim
MAIN_PACKAGE := ./cmd
DIST_DIR := dist

# Go configuration
GO_VERSION := 1.24
GOOS := $(shell go env GOOS)
GOARCH := $(shell go env GOARCH)

# Build configuration
LDFLAGS := -s -w
BUILD_FLAGS := -v

# Version information (can be overridden)
VERSION ?= dev
COMMIT_HASH ?= $(shell git rev-parse HEAD 2>/dev/null || echo "unknown")
TAG ?= $(shell git describe --tags --abbrev=0 2>/dev/null || echo "dev")

# Detect if we're compiling on Windows
IS_WINDOWS := $(shell [ -n "$$WINDIR" ] && echo "true" || echo "false")

# Default target
.PHONY: help
help: ## Show this help message
	@echo "Available targets:"
	@grep -E '^[a-zA-Z_-]+:.*?## .*$$' $(MAKEFILE_LIST) | sort | awk 'BEGIN {FS = ":.*?## "}; {printf "  \033[36m%-20s\033[0m %s\n", $$1, $$2}'

# Development targets
.PHONY: deps
deps: ## Download and verify Go dependencies
	@echo "Downloading dependencies..."
	go mod download
	go mod verify

.PHONY: tidy
tidy: ## Tidy up Go modules
	@echo "Tidying Go modules..."
	go mod tidy

.PHONY: fmt
fmt: ## Format Go code
	@echo "Formatting Go code..."
	go fmt ./...

.PHONY: vet
vet: ## Run go vet
	@echo "Running go vet..."
	@packages=$$(go list ./...); \
	go vet $$packages

.PHONY: test
test: ## Run tests with race detection and coverage
	@echo "Running tests..."
	@packages=$$(go list ./...); \
	CGO_ENABLED=1 go test -race -coverprofile=coverage.out $$packages

.PHONY: test-no-race
test-no-race: ## Run tests without race detection (for environments without C compiler)
	@echo "Running tests without race detection..."
	@packages=$$(go list ./...); \
	go test -coverprofile=coverage.out $$packages

.PHONY: test-bench
test-bench: ## Run tests with benchmarks
	@echo "Running tests with benchmarks..."
	@packages=$$(go list ./...); \
	go test -bench=. -benchmem $$packages

.PHONY: test-verbose
test-verbose: ## Run tests with verbose output
	@echo "Running tests with verbose output..."
	go test -v ./...

.PHONY: coverage
coverage: test-no-race ## Generate and display test coverage
	@echo "Generating coverage report..."
	go tool cover -html=coverage.out -o coverage.html
	@echo "Coverage report generated: coverage.html"
	@echo "Coverage summary:"
	go tool cover -func=coverage.out | tail -1

# Linting targets
.PHONY: lint
lint: ## Run golangci-lint
	@echo "Running golangci-lint..."
	golangci-lint run --timeout=5m

.PHONY: lint-fix
lint-fix: ## Run golangci-lint with auto-fix
	@echo "Running golangci-lint with auto-fix..."
	golangci-lint run --timeout=5m --fix

# Security targets
.PHONY: security
security: ## Run security checks
	@echo "Running security checks..."
	go run golang.org/x/vuln/cmd/govulncheck@latest ./...

# Build targets
.PHONY: build
build: ## Build the application for current platform
	@echo "Building binary for $(GOOS)/$(GOARCH)..."
	@if [ "$(IS_WINDOWS)" = "true" ]; then \
		echo "Building TraceTrim.exe (compiled on Windows)"; \
		go build $(BUILD_FLAGS) -o TraceTrim.exe $(MAIN_PACKAGE); \
	else \
		echo "Building tracetrim (compiled on Unix-like system)"; \
		go build $(BUILD_FLAGS) -o tracetrim $(MAIN_PACKAGE); \
	fi

.PHONY: build-versioned
build-versioned: ## Build with version information
	@echo "Building binary with version $(VERSION)..."
	@if [ "$(IS_WINDOWS)" = "true" ]; then \
		go build $(BUILD_FLAGS) \
			-ldflags="$(LDFLAGS) -X main.version=$(VERSION) -X main.commitHash=$(COMMIT_HASH)" \
			-o TraceTrim.exe $(MAIN_PACKAGE); \
	else \
		go build $(BUILD_FLAGS) \
			-ldflags="$(LDFLAGS) -X main.version=$(VERSION) -X main.commitHash=$(COMMIT_HASH)" \
			-o tracetrim $(MAIN_PACKAGE); \
	fi

.PHONY: clean
clean: ## Clean build artifacts
	@echo "Cleaning build artifacts..."
	rm -f TraceTrim.exe tracetrim
	rm -rf $(DIST_DIR)
	rm -f coverage.out coverage.html

# Cross-compilation targets
.PHONY: cross-compile
cross-compile: ## Cross-compile for all supported platforms
	@echo "Cross-compiling for all platforms..."
	@mkdir -p $(DIST_DIR)
	@echo "Building for Windows AMD64..."
	GOOS=windows GOARCH=amd64 go build $(BUILD_FLAGS) \
		-ldflags="$(LDFLAGS) -X main.version=$(VERSION) -X main.commitHash=$(COMMIT_HASH)" \
		-o $(DIST_DIR)/TraceTrim-windows-amd64.exe $(MAIN_PACKAGE)
	@echo "Building for Windows ARM64..."
	GOOS=windows GOARCH=arm64 go build $(BUILD_FLAGS) \
		-ldflags="$(LDFLAGS) -X main.version=$(VERSION) -X main.commitHash=$(COMMIT_HASH)" \
		-o $(DIST_DIR)/TraceTrim-windows-arm64.exe $(MAIN_PACKAGE)
	@echo "Building for macOS AMD64..."
	GOOS=darwin GOARCH=amd64 go build $(BUILD_FLAGS) \
		-ldflags="$(LDFLAGS) -X main.version=$(VERSION) -X main.commitHash=$(COMMIT_HASH)" \
		-o $(DIST_DIR)/tracetrim-darwin-amd64 $(MAIN_PACKAGE)
	@echo "Building for macOS ARM64..."
	GOOS=darwin GOARCH=arm64 go build $(BUILD_FLAGS) \
		-ldflags="$(LDFLAGS) -X main.version=$(VERSION) -X main.commitHash=$(COMMIT_HASH)" \
		-o $(DIST_DIR)/tracetrim-darwin-arm64 $(MAIN_PACKAGE)
	@echo "Building for Linux AMD64..."
	GOOS=linux GOARCH=amd64 go build $(BUILD_FLAGS) \
		-ldflags="$(LDFLAGS) -X main.version=$(VERSION) -X main.commitHash=$(COMMIT_HASH)" \
		-o $(DIST_DIR)/tracetrim-linux-amd64 $(MAIN_PACKAGE)
	@echo "Building for Linux ARM64..."
	GOOS=linux GOARCH=arm64 go build $(BUILD_FLAGS) \
		-ldflags="$(LDFLAGS) -X main.version=$(VERSION) -X main.commitHash=$(COMMIT_HASH)" \
		-o $(DIST_DIR)/tracetrim-linux-arm64 $(MAIN_PACKAGE)

# Release targets
.PHONY: checksums
checksums: cross-compile ## Generate checksums for all binaries
	@echo "Generating checksums..."
	@cd $(DIST_DIR) && sha256sum * > checksums.txt
	@echo "Checksums generated: $(DIST_DIR)/checksums.txt"

.PHONY: verify-checksums
verify-checksums: checksums ## Verify checksums
	@echo "Verifying checksums..."
	@cd $(DIST_DIR) && sha256sum -c checksums.txt --strict
	@echo "✓ All checksums verified successfully"

.PHONY: release-prep
release-prep: clean test-no-race lint security cross-compile checksums verify-checksums ## Prepare release artifacts
	@echo "Release preparation complete!"
	@echo "Binaries available in $(DIST_DIR)/"
	@ls -la $(DIST_DIR)/

# Development workflow targets
.PHONY: dev-setup
dev-setup: deps ## Set up development environment
	@echo "Development environment setup complete!"

.PHONY: check
check: fmt vet lint test ## Run all checks (format, vet, lint, test)
	@echo "All checks passed!"

.PHONY: ci
ci: deps vet test-no-race test-bench ## Run complete CI suite
	@echo "CI suite completed!"

.PHONY: ci-lint
ci-lint: deps lint ## Run CI linting
	@echo "CI linting completed!"

.PHONY: ci-security
ci-security: deps security ## Run CI security checks
	@echo "CI security checks completed!"

.PHONY: ci-build
ci-build: deps build ## Run CI build
	@echo "CI build completed!"

.PHONY: ci-cross-compile
ci-cross-compile: deps cross-compile ## Run CI cross-compilation
	@echo "CI cross-compilation completed!"

# Utility targets
.PHONY: version
version: ## Show version information
	@echo "Version: $(VERSION)"
	@echo "Commit: $(COMMIT_HASH)"
	@echo "Tag: $(TAG)"
	@echo "Go version: $(shell go version)"
	@echo "Platform: $(GOOS)/$(GOARCH)"
	@echo "Compiling on Windows: $(IS_WINDOWS)"

.PHONY: install-tools
install-tools: ## Install development tools
	@echo "Installing development tools..."
	@go install github.com/golangci/golangci-lint/cmd/golangci-lint@latest
	@echo "Development tools installed!"

.PHONY: run
run: build ## Build and run the application
	@echo "Running application..."
	@if [ "$(IS_WINDOWS)" = "true" ]; then \
		./TraceTrim.exe; \
	else \
		./tracetrim; \
	fi

# System dependencies (for CI environments)
.PHONY: install-deps-linux
install-deps-linux: ## Install system dependencies for Linux
	@echo "Installing Linux system dependencies..."
	sudo apt-get update
	sudo apt-get install -y libx11-dev libxrandr-dev libxinerama-dev libxcursor-dev libxi-dev libgl1-mesa-dev libxext-dev xvfb

.PHONY: install-deps-macos
install-deps-macos: ## Install system dependencies for macOS
	@echo "Installing macOS system dependencies..."
	brew install pkg-config

.PHONY: setup-display-linux
setup-display-linux: ## Set up virtual display for Linux
	@echo "Setting up virtual display for Linux..."
	Xvfb :99 -screen 0 1024x768x24 > /dev/null 2>&1 &
	@echo "DISPLAY=:99.0" >> $$GITHUB_ENV

# Default target
.DEFAULT_GOAL := help
