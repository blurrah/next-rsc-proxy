.PHONY: build run test

build:
	go build -o rsc-proxy ./cmd/server/main.go

run:
	go run ./cmd/server/main.go

test:
	go test -v -race ./...
