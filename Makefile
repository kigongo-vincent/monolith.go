.PHONY: dev build test cli

dev:
	go run .

build:
	go build -o bin/monolith .

cli:
	go build -o bin/monolith ./cmd/monolith

test:
	go test ./...
