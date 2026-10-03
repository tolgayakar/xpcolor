BINARY  := xpcolor
VERSION ?= $(shell git describe --tags --always --dirty 2>/dev/null || echo dev)
LDFLAGS := -s -w -X main.Version=$(VERSION)

.PHONY: all build test lint fmt tidy clean

all: lint test build

build:
	go build -ldflags "$(LDFLAGS)" -o bin/$(BINARY) ./cmd/xpcolor

test:
	go test -race -cover ./...

lint:
	golangci-lint run ./...

fmt:
	gofmt -s -w .

tidy:
	go mod tidy

clean:
	rm -rf bin dist
