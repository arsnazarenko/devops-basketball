#!/usr/bin/make -f

BIN_DIR = ./bin
BIN_NAME = devops-basketball
SRC_DIR = ./cmd/devops-basketball

.PHONY: all build generate-api test clean

all: build

# Build habit service
build:
	@mkdir -p $(BIN_DIR)
	CGO_ENABLED=0 go build -v -o $(BIN_DIR)/$(BIN_NAME) $(SRC_DIR)

# Generate api from openapi spec
generate-api:
	go generate ./api/...

# Run tests
test:
	go test -v -cover ./...

clean:
	rm -rf ./bin
