.PHONY: all build build-release build-arm64 build-amd64 build-web test clean install help

VERSION := $(shell cat VERSION 2>/dev/null || echo "0.1.0-dev")
BUILD_DIR := bin
LDFLAGS := -s -w -X main.Version=$(VERSION)

all: build build-web

build:
	@echo "Building SawitOS Go binaries..."
	@mkdir -p $(BUILD_DIR)
	go build -ldflags "$(LDFLAGS)" -o $(BUILD_DIR)/sawitctl ./cmd/sawitctl
	go build -ldflags "$(LDFLAGS)" -o $(BUILD_DIR)/sawitd ./cmd/sawitd
	go build -ldflags "$(LDFLAGS)" -o $(BUILD_DIR)/sawit-agent ./cmd/sawit-agent
	go build -ldflags "$(LDFLAGS)" -o $(BUILD_DIR)/sawit-health ./cmd/sawit-health
	go build -ldflags "$(LDFLAGS)" -o $(BUILD_DIR)/sawit-update ./cmd/sawit-update
	@echo "Build complete. Binaries output to $(BUILD_DIR)/"

build-amd64:
	@echo "Cross-compiling SawitOS binaries for Linux AMD64 (x86_64)..."
	@mkdir -p $(BUILD_DIR)
	GOOS=linux GOARCH=amd64 go build -ldflags "$(LDFLAGS)" -o $(BUILD_DIR)/sawitctl-linux-amd64 ./cmd/sawitctl
	GOOS=linux GOARCH=amd64 go build -ldflags "$(LDFLAGS)" -o $(BUILD_DIR)/sawitd-linux-amd64 ./cmd/sawitd
	GOOS=linux GOARCH=amd64 go build -ldflags "$(LDFLAGS)" -o $(BUILD_DIR)/sawit-agent-linux-amd64 ./cmd/sawit-agent

build-arm64:
	@echo "Cross-compiling SawitOS binaries for Linux ARM64 (Raspberry Pi 3/4/5 / ARM Server)..."
	@mkdir -p $(BUILD_DIR)
	GOOS=linux GOARCH=arm64 go build -ldflags "$(LDFLAGS)" -o $(BUILD_DIR)/sawitctl-linux-arm64 ./cmd/sawitctl
	GOOS=linux GOARCH=arm64 go build -ldflags "$(LDFLAGS)" -o $(BUILD_DIR)/sawitd-linux-arm64 ./cmd/sawitd
	GOOS=linux GOARCH=arm64 go build -ldflags "$(LDFLAGS)" -o $(BUILD_DIR)/sawit-agent-linux-arm64 ./cmd/sawit-agent

build-web:
	@echo "Building Sawit Web Next.js frontend..."
	@cd web && npm run build

test:
	@echo "Running SawitOS unit and integration tests..."
	go test -v ./...

clean:
	@echo "Cleaning build artifacts..."
	rm -rf $(BUILD_DIR)
	rm -rf web/.next

install: build
	@echo "Installing SawitOS binaries to system paths..."
	install -d /usr/bin /usr/sbin /usr/libexec /etc/sawit
	install -m 0755 $(BUILD_DIR)/sawitctl /usr/bin/sawitctl
	install -m 0755 $(BUILD_DIR)/sawitd /usr/sbin/sawitd
	install -m 0755 $(BUILD_DIR)/sawit-agent /usr/libexec/sawit-agent
	install -m 0755 $(BUILD_DIR)/sawit-health /usr/bin/sawit-health
	install -m 0755 $(BUILD_DIR)/sawit-update /usr/bin/sawit-update
	@if [ ! -f /etc/sawit/sawitd.yaml ]; then install -m 0644 configs/sawitd.yaml /etc/sawit/sawitd.yaml; fi

help:
	@echo "SawitOS Build System"
	@echo "  make build        - Build all Go binaries into bin/"
	@echo "  make build-amd64  - Cross-compile release binaries for Linux x86_64"
	@echo "  make build-arm64  - Cross-compile release binaries for Linux ARM64 (Raspberry Pi)"
	@echo "  make build-web    - Build Sawit Web Next.js application"
	@echo "  make test         - Run Go test suite"
	@echo "  make clean        - Remove build directory"
	@echo "  make install      - Install binaries to system locations (Linux only)"
