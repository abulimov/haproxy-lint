.PHONY: all test clean build linux

GOFLAGS ?= $(GOFLAGS:)
GOOS ?= linux
GOARCH ?= amd64
CGO_ENABLED ?= 0
REPO = github.com/abulimov/haproxy-lint

all: test

get:
	@go get $(GOFLAGS) -t ./...

linux: get
	@GOOS=$(GOOS) GOARCH=$(GOARCH) CGO_ENABLED=$(CGO_ENABLED) go build $(GOFLAGS) $(REPO)

build: get
	go build $(GOFLAGS) $(REPO)

test: get
	@go test -v $(GOFLAGS) ./...

clean:
	@go clean $(GOFLAGS) -i $(REPO)
