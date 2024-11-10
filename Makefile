# Define variables
BINARY_CLI=ezcd-cli
BINARY_SERVER=ezcd-server
DIST_DIR=dist

# Default target
all: install test build

# Install dependencies
install:
	go mod tidy
	cd web && pnpm install

# Run tests
test:
	go test ./...

# Build CLI
build-cli:
	go build -o $(DIST_DIR)/$(BINARY_CLI) ./cmd/cli

build-web:
	cd web && pnpm run build

# Build SERVER
build-server: build-web
	go build -o $(DIST_DIR)/$(BINARY_SERVER) ./cmd/server

# Build both CLI and SERVER
build: build-cli build-server

dev-web:
	cd web && pnpm run dev

dev-server:
	go run --tags dev cmd/server/main.go

acceptance:
	cd acceptance && pnpm exec playwright test

playwright:
	cd acceptance && pnpm exec playwright test --ui --ui-host=0.0.0.0

# Clean up
clean:
	rm -rf $(DIST_DIR)

# PHONY targets
.PHONY: all install test build-cli build-server build clean dev-web dev-server