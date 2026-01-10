.PHONY: help test build run clean lint fmt coverage windows

# Default target
help:
	@echo "Hexy and I Know It - Makefile Commands"
	@echo ""
	@echo "  make test      - Run all tests"
	@echo "  make build     - Build the game binary"
	@echo "  make run       - Build and run the game"
	@echo "  make windows   - Cross-compile for Windows (creates game.exe)"
	@echo "  make clean     - Remove build artifacts"
	@echo "  make lint      - Run linter"
	@echo "  make fmt       - Format code"
	@echo "  make coverage  - Generate test coverage report"
	@echo ""

# Run tests
test:
	@echo "Running tests..."
	go test -v -race ./...

# Run tests with coverage
coverage:
	@echo "Running tests with coverage..."
	go test -v -race -coverprofile=coverage.out ./...
	go tool cover -html=coverage.out -o coverage.html
	@echo "Coverage report generated: coverage.html"

# Build the game
build:
	@echo "Building game..."
	go build -o hexy cmd/game/main.go

# Cross-compile for Windows
windows:
	@echo "Cross-compiling for Windows..."
	@echo "Note: Requires gcc-mingw-w64-x86-64 (install with: sudo apt-get install gcc-mingw-w64-x86-64)"
	GOOS=windows GOARCH=amd64 CGO_ENABLED=1 CC=x86_64-w64-mingw32-gcc go build -o /mnt/c/Temp/hexy.exe cmd/game/main.go
	@echo "Windows executable created: C:\\Temp\\hexy.exe"
	@echo "Run from Windows: C:\\Temp\\hexy.exe"

# Build and run the game
run: build
	@echo "Running game..."
	./hexy

# Clean build artifacts
clean:
	@echo "Cleaning..."
	rm -f hexy game*.exe coverage.out coverage.html

# Run linter
lint:
	@echo "Running linter..."
	golangci-lint run

# Format code
fmt:
	@echo "Formatting code..."
	gofmt -s -w .
	goimports -w .

# Install development dependencies
deps:
	@echo "Installing development dependencies..."
	go install github.com/golangci/golangci-lint/cmd/golangci-lint@latest
	go install golang.org/x/tools/cmd/goimports@latest
	@echo "Installing game dependencies..."
	go mod download
