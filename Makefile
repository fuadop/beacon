# Variables
BINARY_NAME=beacon
BUILD_DIR=bin
MAIN_PATH=cmd/main.go

.PHONY: all build run clean tidy help

# Default target when running just 'make'
all: tidy build

## build: Compiles the application binary into the bin/ directory
build:
	@echo "🔨 Building binary..."
	@mkdir -p $(BUILD_DIR)
	@go build -o $(BUILD_DIR)/$(BINARY_NAME) $(MAIN_PATH)
	@echo "✔ Binary built successfully at: $(BUILD_DIR)/$(BINARY_NAME)"

## run: Builds the binary and runs it immediately
run: build
	@echo "🚀 Executing app..."
	@./$(BUILD_DIR)/$(BINARY_NAME)

## tidy: Cleans up module dependencies and downloads missing items
tidy:
	@echo "📦 Tidying Go modules..."
	@go mod tidy

## clean: Removes local compilation output artifacts
clean:
	@echo "🧹 Cleaning up build artifacts..."
	@rm -rf $(BUILD_DIR)
	@echo "✔ Clean complete."

## help: Shows all available Makefile command targets
help:
	@echo "Available commands:"
	@sed -n 's/^##//p' $(MAKEFILE_LIST) | column -t -s ':' |  sed -e 's/^/ /'
