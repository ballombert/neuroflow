.PHONY: build test run clean build-windows build-linux build-all

export PATH := /opt/homebrew/bin:$(PATH)

BINARY_NAME=neuroflow
WINDOWS_BINARY=$(BINARY_NAME).exe

build:
	go build -ldflags="-s -w" -o $(BINARY_NAME) ./cmd/neuroflow

build-windows:
	GOOS=windows GOARCH=amd64 go build -ldflags="-s -w" -o $(WINDOWS_BINARY) ./cmd/neuroflow

build-linux:
	GOOS=linux GOARCH=amd64 go build -ldflags="-s -w" -o $(BINARY_NAME)-linux ./cmd/neuroflow

build-all: build build-windows build-linux

test:
	go test -v ./...

run: build
	./$(BINARY_NAME)

clean:
	rm -f $(BINARY_NAME) $(WINDOWS_BINARY) $(BINARY_NAME)-linux
