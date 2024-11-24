# Define variables
BINARY_CLI=ezcd-cli
BINARY_SERVER=ezcd-server
DIST_DIR=dist
COVERAGE_FILE=coverage.out
COVERAGE_HTML_DIR=coverage

export EZCD_DATABASE_URL=postgres://${PGUSER}:${PGPASSWORD}@${PGHOST}:5432/ezcd?sslmode=disable

# Default target
all: install test build

# Install dependencies
install:
	go mod tidy
	cd web && pnpm install

format:
	go fmt ./...

# Run tests
test:
	go test ./...

# coverage: test
# 	mkdir -p $(COVERAGE_HTML_DIR)
# 	go tool cover -html=$(COVERAGE_FILE) -o $(COVERAGE_HTML_DIR)/index.html

serve-coverage: coverage
	pnpm dlx serve $(COVERAGE_HTML_DIR)

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

# we run the server on 3924 for dev mode and dev-web runs on 3923 with a 
# proxy to 3924
dev-server:
	EZCD_PORT=3924 EZCD_DATABASE_URL=${EZCD_DATABASE_URL} go run --tags dev cmd/server/main.go

test-server:
	EZCD_DATABASE_URL=${EZCD_DATABASE_URL} ./dist/ezcd-server

acceptance:
	cd acceptance && pnpm exec playwright test

playwright:
	cd acceptance && pnpm exec playwright test --ui --ui-host=0.0.0.0

env:
	@echo 'export EZCD_DATABASE_URL='${EZCD_DATABASE_URL}

db-reset:
	psql -c "DROP DATABASE IF EXISTS ezcd;"
	psql -c "CREATE DATABASE ezcd;"

db-migrate:
	pg-schema-diff apply --dsn "postgres://${PGUSER}:${PGPASSWORD}@${PGHOST}:5432/ezcd" --schema-dir schema --allow-hazards DELETES_DATA

# Clean up
clean:
	rm -rf $(DIST_DIR)

# PHONY targets
.PHONY: all install test build-cli build-server build clean dev-web dev-server acceptance db-reset db-migrate db-set-envn coverage serve-coverage