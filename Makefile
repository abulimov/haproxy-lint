.PHONY: all test clean build linux release

GOFLAGS ?= $(GOFLAGS:)
GOOS ?= linux
GOARCH ?= amd64
CGO_ENABLED ?= 0
REPO = github.com/abulimov/haproxy-lint

all: test

get:
	@go get $(GOFLAGS) -t ./...

release:
	mkdir -p release
	GOOS=linux GOARCH=amd64 go build -o haproxy-lint $(REPO)
	zip -0 release/haproxy-lint.linux-amd64.zip  haproxy-lint
	GOOS=linux GOARCH=arm go build -o haproxy-lint $(REPO)
	zip -0 release/haproxy-lint.linux-arm.zip  haproxy-lint
	GOOS=darwin GOARCH=amd64 go build -o haproxy-lint $(REPO)
	zip -0 release/haproxy-lint.darwin-amd64.zip  haproxy-lint

build: get
	go build $(GOFLAGS) $(REPO)

test: get
	@go test -v $(GOFLAGS) ./...

clean:
	@go clean $(GOFLAGS) -i $(REPO)
