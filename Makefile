build: fmt vet proto
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

proto:
	protoc proto/process_manager.proto --go_out=. --go_opt=paths=source_relative --go-grpc_out=. --go-grpc_opt=paths=source_relative
