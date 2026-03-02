.PHONY: build test run clean

BINARY=sudoku-tui
BUILD_DIR=./bin

build:
	go build -o $(BUILD_DIR)/$(BINARY) .

test:
	go test -v -race ./...

run:
	go run .

clean:
	rm -rf $(BUILD_DIR)

deps:
	go mod tidy
