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

# Build SERVER
build-server:
	go build -o $(DIST_DIR)/$(BINARY_SERVER) ./cmd/server

# Build both CLI and SERVER
build: build-cli build-server

# Clean up
clean:
	rm -rf $(DIST_DIR)

# PHONY targets
.PHONY: all install test build-cli build-server build clean