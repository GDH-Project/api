build: download-deps tidy-deps compile lint

fmt:
	@echo "✓ Formatting all source code..."
	@gofmt -l -s -w .

download-deps:
	@echo "✓ Downloading all dependencies..."
	@go mod download

tidy-deps:
	@echo "✓ Tidying up all dependencies..."
	@go mod tidy

update-deps:
	@echo "✓ Upgrading all compile and test dependencies..."
	@go get -u -t ./...
	@echo "✓ Upgrading all dev dependencies..."
	@go get tool
	@echo "✓ Tidying up the go.mod file..."
	@go mod tidy

compile:
	@echo "✓ Compiling the source code..."
	@go build -o server cmd/main.go

lint:
	@echo "✓ Running goreleaser..."
	@go tool goreleaser check