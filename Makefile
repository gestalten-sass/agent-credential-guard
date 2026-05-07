BINARY=guard

.PHONY: build install test fmt

build:
	go build -o $(BINARY) ./cmd/guard

install:
	go install ./cmd/guard

test:
	go test ./...

fmt:
	gofmt -w cmd internal
