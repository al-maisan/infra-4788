.PHONY: all build

BIN_DIR := ./bin
version := $(shell git rev-parse --short=12 HEAD)
timestamp := $(shell date -u +"%Y-%m-%dT%H:%M:%SZ")

ROOT_DIR:=$(shell dirname $(realpath $(firstword $(MAKEFILE_LIST))))

all: build

clean:
	rm -f $(BIN_DIR)/pgen

build: lint
	rm -f $(BIN_DIR)/pgen
	go build -o $(BIN_DIR)/pgen -v -ldflags \
		"-X main.rev=$(version) -X main.bts=$(timestamp)" cmd/pgen/*.go

lint:
	golangci-lint run

test: lint
	go test ./...
