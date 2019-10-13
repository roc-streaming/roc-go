all: build test

.PHONY: build test

build:
	go build ./roc

test:
	go test ./roc/...
