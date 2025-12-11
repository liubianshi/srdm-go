.PHONY: all build test clean run

BINARY_NAME=srdm
BUILD_DIR=bin
MAIN_PATH=./cmd/srdm

all: build

build:
	@echo "Building..."
	@mkdir -p $(BUILD_DIR)
	@go build -o $(BUILD_DIR)/$(BINARY_NAME) $(MAIN_PATH)

test:
	@echo "Running tests..."
	@go test -v ./...

clean:
	@echo "Cleaning..."
	@rm -rf $(BUILD_DIR)
	@go clean

run: build
	@$(BUILD_DIR)/$(BINARY_NAME)
