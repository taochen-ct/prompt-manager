run:
	@CGO_ENABLED=1 go run ./cmd/server

wire:
	@cd ./internal/app/ && wire

build:
	@echo "Static Building..."

	@CGO_ENABLED=1 GOOS=linux GOARCH=amd64 go build -o build/linux/main ./cmd/server
	@CGO_ENABLED=1 GOOS=darwin GOARCH=arm64 go build -o build/darwin/main ./cmd/server

	@echo "Static Build Completed"