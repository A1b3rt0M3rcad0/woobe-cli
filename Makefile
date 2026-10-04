.PHONY: build test check
build:
	go build -trimpath -o bin/woobe ./cmd/woobe
test:
	go test -race ./...
check:
	test -z "$$(gofmt -l cmd internal)"
	go mod verify
	go vet ./...
	go test -race ./...
