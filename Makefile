# dnglab-gui Makefile
# Build Linux and macOS versions of the application

APP_NAME := dnglab-gui
VERSION := 0.1.0
BUILD_DIR := bin
CMD_PATH := ./cmd/dnglab-gui

# Build flags
LDFLAGS := -s -w

.PHONY: all clean build build-linux build-darwin build-darwin-amd64 build-darwin-arm64 test vet

# Default: build for the current platform
all: build

# Build for the current platform (native)
build:
	@mkdir -p $(BUILD_DIR)
	go build -ldflags "$(LDFLAGS)" -o $(BUILD_DIR)/$(APP_NAME) $(CMD_PATH)

# Linux amd64
# Note: Cross-compiling CGO apps for Linux from macOS requires a Linux cross-compiler.
# Run this target on Linux or use Docker:
#   docker run --rm -v "$(PWD)":/app -w /app golang:1.21 make build-linux
build-linux:
	@mkdir -p $(BUILD_DIR)
	CGO_ENABLED=1 GOOS=linux GOARCH=amd64 go build -ldflags "$(LDFLAGS)" -o $(BUILD_DIR)/$(APP_NAME)-linux-amd64 $(CMD_PATH)

# macOS (both architectures)
build-darwin: build-darwin-amd64 build-darwin-arm64

# macOS amd64 (Intel)
build-darwin-amd64:
	@mkdir -p $(BUILD_DIR)
	CGO_ENABLED=1 GOOS=darwin GOARCH=amd64 go build -ldflags "$(LDFLAGS)" -o $(BUILD_DIR)/$(APP_NAME)-darwin-amd64 $(CMD_PATH)

# macOS arm64 (Apple Silicon)
build-darwin-arm64:
	@mkdir -p $(BUILD_DIR)
	CGO_ENABLED=1 GOOS=darwin GOARCH=arm64 go build -ldflags "$(LDFLAGS)" -o $(BUILD_DIR)/$(APP_NAME)-darwin-arm64 $(CMD_PATH)

# Build Linux version using Docker (works on any host OS)
build-linux-docker:
	@mkdir -p $(BUILD_DIR)
	docker run --rm -v "$(PWD)":/app -w /app golang:1.21-bookworm \
		sh -c "apt-get update && apt-get install -y libgl1-mesa-dev xorg-dev && \
		CGO_ENABLED=1 GOOS=linux GOARCH=amd64 go build -ldflags '$(LDFLAGS)' -o $(BUILD_DIR)/$(APP_NAME)-linux-amd64 $(CMD_PATH)"

# Clean build artifacts
clean:
	rm -rf $(BUILD_DIR)
	rm -f $(APP_NAME)

# Run tests
test:
	go test ./...

# Run go vet
vet:
	go vet ./...
