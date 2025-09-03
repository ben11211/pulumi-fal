.PHONY: build install clean test lint fmt vet check-deps schema example

# Variables
PROVIDER_NAME := pulumi-provider-fal
BINARY := $(PROVIDER_NAME)

# Default target
all: build

# Build the provider binary
build:
	@echo "Building $(PROVIDER_NAME)..."
	go build -ldflags "$(LDFLAGS)" -o $(BINARY) .

# Clean build artifacts
clean:
	@echo "Cleaning..."
	rm -f $(BINARY)
	go clean

# Generate schema (if needed)
schema: build
	@echo "Generating schema..."
	./$(BINARY) gen-schema --out schema.json

# Help
help:
	@echo "Available targets:"
	@echo "  build      - Build the provider binary"
	@echo "  clean      - Clean build artifacts"
	@echo "  schema     - Generate schema file"
	@echo "  help       - Show this help"