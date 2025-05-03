# Makefile

BINARY_NAME=time-server
# Default build target
.PHONY: all
all: build

# Build the application
.PHONY: build
build:
	go build -o $(BINARY_NAME) *.go

# Run the application
.PHONY: run
run: build
	./$(BINARY_NAME)

# Clean up generated files
.PHONY: clean
clean:
	go clean
	rm -f $(BINARY_NAME)

# Format code
.PHONY: fmt
fmt:
	go fmt ./...