build: fmt vet
	go build -race ./...

all: build test linter

fmt:
	go fmt ./...

vet:
	go vet ./...

test:
	go test -race ./...

# https://golangci-lint.run/usage/install/#linux-and-windows
linter:
	golangci-lint run ./...
