.PHONY: all vendor release

all: build

vendor:
	@go mod tidy
	@go mod vendor
	@go mod download

build: build-report

build-report:
	CGO_ENABLED=0 GOOS=linux GOARCH=amd64 go build -a -o ~/bin/sc-marker-by-pr-hash cmd/sc-marker-by-pr-hash/*.go

fmt:
	go fmt ./...
	go vet ./...

run:
	go run cmd/sc-marker-by-pr-hash/main.go

play:
	go run cmd/sc-marker-by-pr-hash/main.go

