.PHONY: test
export GOBIN := $(PWD)/bin
export PATH := $(GOBIN):$(PATH)
export INSTALL_FLAG=
VERSION?=$(shell git describe --all --dirty | cut -d / -f2,3,4)

# Determine which OS.
OS?=$(shell uname -s | tr A-Z a-z)

default: build

init:
	@git config core.hooksPath .githooks	

dependencies:
	@go mod tidy

dep: dependencies

test:
	@go test ./... -race -coverprofile=coverage.txt -covermode=atomic ; go tool cover -html=coverage.txt -o coverage.html

build:
	CGO_ENABLED=0 GOOS=$(OS) go build $(INSTALL_FLAG) -ldflags="-X github.com/naemono/go-cloud-actions/cmd.version=$(VERSION)" -o $(GOBIN)/cloud
