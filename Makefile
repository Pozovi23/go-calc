P ?= 10
ARGS ?=

.PHONY: build run ui test lint clean

build:
	go build -o bin/calc ./cmd/calc

run:
	go run ./cmd/calc -p $(P) $(ARGS)

ui:
	go run ./cmd/calc -p $(P) --ui

test:
	go test -v ./...

lint:
	golangci-lint run

clean:
	rm -rf bin/ coverage.out
