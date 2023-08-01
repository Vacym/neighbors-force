.PHONY: build run test coverage

BINARY_NAME=proxyserver.exe
BUILD_DIR=build
COVERAGE_FILE=coverage.out

build:
	go build -o $(BUILD_DIR)/$(BINARY_NAME) ./cmd/proxyserver

run:
	go run ./cmd/proxyserver/main.go

test:
	go test -v ./...

coverage:
	go test -coverprofile=$(COVERAGE_FILE) ./...
	go tool cover -html=$(COVERAGE_FILE)

.DEFAULT_GOAL := run
