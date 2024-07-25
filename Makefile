.PHONY: build run test clean

# Build the main executable
build:
	@go build -o bin/gpu-cloudsim cmd/main.go

# Run the simulation
run: build
	@./bin/gpu-cloudsim

test:
	@go test ./...

# Clean up build artifacts
clean:
	@go clean