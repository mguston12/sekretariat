SHELL := /bin/bash
GO := go
NAME := sekretariat # Replace with your application name
OS := $(shell uname)
MAIN_GO := cmd/http/main.go
GO_VERSION := $(shell $(GO) version | sed -e 's/^[^0-9.]*\([0-9.]*\).*/\1/')
PACKAGE_DIRS := $(shell $(GO) list ./... | grep -v /vendor/)
BUILDFLAGS := ''
CGO_ENABLED = 0
VENDOR_DIR = vendor

# Default target
all: build

# Check code formatting, build, and run tests
check: fmt build test

# Build the binary
.PHONY: build
build:
	mkdir -p bin
	CGO_ENABLED=$(CGO_ENABLED) $(GO) build -ldflags $(BUILDFLAGS) -o bin/$(NAME) $(MAIN_GO)

# Run tests
.PHONY: test
test:
	CGO_ENABLED=$(CGO_ENABLED) $(GO) test $(PACKAGE_DIRS) -v

# Build for Linux
.PHONY: linux
linux:
	mkdir -p bin
	CGO_ENABLED=$(CGO_ENABLED) GOOS=linux GOARCH=amd64 $(GO) build -ldflags $(BUILDFLAGS) -o bin/$(NAME) $(MAIN_GO)

# Format code
.PHONY: fmt
fmt:
	$(GO) fmt $(PACKAGE_DIRS)

# Clean build artifacts
.PHONY: clean
clean:
	rm -rf bin
