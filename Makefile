.PHONY: build lint test test-coverage run clean migrate deps

build:
	go build -o bin/server ./cmd/server

lint:
	golangci-lint run ./...

test:
	go test -v -race ./...

test-coverage:
	go test -v -race -coverprofile=coverage.out ./...
	go tool cover -html=coverage.out -o coverage.html

run:
	go run ./cmd/server

clean:
	rm -rf bin/

migrate:
	go run ./cmd/migrate

deps:
	go mod download
	go mod tidy
