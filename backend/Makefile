.PHONY: build test lint run clean

# Build the application
build:
	go build -o bin/glossary-server main.go

# Run tests
test:
	go test -v ./...

# Run tests with coverage
test-coverage:
	go test -v -coverprofile=coverage.out ./...
	go tool cover -html=coverage.out -o coverage.html

# Run linter
lint:
	golangci-lint run

# Run the application
run:
	go run main.go

# Clean build artifacts
clean:
	rm -rf bin/
	rm -f glossary.db
	rm -f coverage.out coverage.html

# Install dependencies
deps:
	go mod download
	go mod tidy

# Format code
fmt:
	go fmt ./...
	goimports -w .

