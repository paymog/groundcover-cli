VERSION ?= dev
LDFLAGS := -s -w -X github.com/paymog/groundcover-cli/internal/cli.version=$(VERSION)

.PHONY: build test generate

build:
	go build -ldflags "$(LDFLAGS)" ./cmd/groundcover

test:
	go test ./...

generate:
	go run ./scripts/generate-commands.go
