# Define variables
BINARY_CLI=ezcd-cli
BINARY_SERVER=ezcd-server
DIST_DIR=dist

# Default target
all: install test build

# Install dependencies
install:
	go mod tidy

# Run tests
test:
	go test ./...

# Build CLI
build-cli:
	go build -o $(DIST_DIR)/$(BINARY_CLI) ./cmd/cli

# Build API
build-api:
	go build -o $(DIST_DIR)/$(BINARY_SERVER) ./cmd/api

# Build both CLI and API
build: build-cli build-api

# Clean up
clean:
	rm -rf $(DIST_DIR)

# PHONY targets
.PHONY: all install test build-cli build-api build clean