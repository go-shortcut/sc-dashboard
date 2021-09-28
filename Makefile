.PHONY: all vendor release

all: build

vendor:
	@go mod tidy
	@go mod vendor
	@go mod download

build:
	CGO_ENABLED=0 GOOS=linux GOARCH=amd64 go build -a -o ~/bin/shortcut-story-marker cmd/app/*.go

fmt:
	go fmt ./...
	go vet ./...

run:
	go run cmd/app/*.go


