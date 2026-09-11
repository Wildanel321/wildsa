.PHONY: all build build-release build-amd64 build-arm64 build-armv7 build-web dist test clean install help

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
	GOOS=linux GOARCH=amd64 go build -ldflags "$(LDFLAGS)" -o $(BUILD_DIR)/sawit-health-linux-amd64 ./cmd/sawit-health
	GOOS=linux GOARCH=amd64 go build -ldflags "$(LDFLAGS)" -o $(BUILD_DIR)/sawit-update-linux-amd64 ./cmd/sawit-update

build-arm64:
	@echo "Cross-compiling SawitOS binaries for Linux ARM64 (Raspberry Pi 3/4/5 64-bit)..."
	@mkdir -p $(BUILD_DIR)
	GOOS=linux GOARCH=arm64 go build -ldflags "$(LDFLAGS)" -o $(BUILD_DIR)/sawitctl-linux-arm64 ./cmd/sawitctl
	GOOS=linux GOARCH=arm64 go build -ldflags "$(LDFLAGS)" -o $(BUILD_DIR)/sawitd-linux-arm64 ./cmd/sawitd
	GOOS=linux GOARCH=arm64 go build -ldflags "$(LDFLAGS)" -o $(BUILD_DIR)/sawit-agent-linux-arm64 ./cmd/sawit-agent
	GOOS=linux GOARCH=arm64 go build -ldflags "$(LDFLAGS)" -o $(BUILD_DIR)/sawit-health-linux-arm64 ./cmd/sawit-health
	GOOS=linux GOARCH=arm64 go build -ldflags "$(LDFLAGS)" -o $(BUILD_DIR)/sawit-update-linux-arm64 ./cmd/sawit-update

build-armv7:
	@echo "Cross-compiling SawitOS binaries for Linux ARMv7 (Raspberry Pi 3/4 32-bit)..."
	@mkdir -p $(BUILD_DIR)
	GOOS=linux GOARCH=arm GOARM=7 go build -ldflags "$(LDFLAGS)" -o $(BUILD_DIR)/sawitctl-linux-armv7 ./cmd/sawitctl
	GOOS=linux GOARCH=arm GOARM=7 go build -ldflags "$(LDFLAGS)" -o $(BUILD_DIR)/sawitd-linux-armv7 ./cmd/sawitd
	GOOS=linux GOARCH=arm GOARM=7 go build -ldflags "$(LDFLAGS)" -o $(BUILD_DIR)/sawit-agent-linux-armv7 ./cmd/sawit-agent
	GOOS=linux GOARCH=arm GOARM=7 go build -ldflags "$(LDFLAGS)" -o $(BUILD_DIR)/sawit-health-linux-armv7 ./cmd/sawit-health
	GOOS=linux GOARCH=arm GOARM=7 go build -ldflags "$(LDFLAGS)" -o $(BUILD_DIR)/sawit-update-linux-armv7 ./cmd/sawit-update

build-release: build-amd64 build-arm64 build-armv7
	@echo "All cross-platform release binaries compiled in $(BUILD_DIR)/"

dist: build-release build-web
	@echo "Packaging SawitOS release archives..."
	@mkdir -p $(BUILD_DIR)/release
	@tar -czf $(BUILD_DIR)/release/sawit-$(VERSION)-linux-amd64.tar.gz -C $(BUILD_DIR) sawitctl-linux-amd64 sawitd-linux-amd64 sawit-agent-linux-amd64 sawit-health-linux-amd64 sawit-update-linux-amd64
	@tar -czf $(BUILD_DIR)/release/sawit-$(VERSION)-linux-arm64.tar.gz -C $(BUILD_DIR) sawitctl-linux-arm64 sawitd-linux-arm64 sawit-agent-linux-arm64 sawit-health-linux-arm64 sawit-update-linux-arm64
	@tar -czf $(BUILD_DIR)/release/sawit-$(VERSION)-linux-armv7.tar.gz -C $(BUILD_DIR) sawitctl-linux-armv7 sawitd-linux-armv7 sawit-agent-linux-armv7 sawit-health-linux-armv7 sawit-update-linux-armv7
	@echo "Distribution packages created in $(BUILD_DIR)/release/"

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
	@echo "  make build        - Build native Go binaries into bin/"
	@echo "  make build-amd64  - Cross-compile release binaries for Linux x86_64"
	@echo "  make build-arm64  - Cross-compile release binaries for Linux ARM64 (Raspberry Pi 64-bit)"
	@echo "  make build-armv7  - Cross-compile release binaries for Linux ARMv7 (Raspberry Pi 3 32-bit)"
	@echo "  make build-release- Cross-compile release binaries for all architectures"
	@echo "  make dist         - Package tar.gz release bundles into bin/release/"
	@echo "  make build-web    - Build Sawit Web Next.js application"
	@echo "  make test         - Run Go test suite"
	@echo "  make clean        - Remove build directory"
	@echo "  make install      - Install binaries to system locations (Linux only)"

