all: check test

.PHONY: check test

check:
	go build ./roc
	golangci-lint run ./roc

test:
	go test -race ./roc

fmt:
	gofmt -s -w ./roc
