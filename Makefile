.PHONY: dev build test

dev:
	go run .

build:
	go build -o bin/monolith .

test:
	go test ./...
