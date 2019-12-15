all: check test

.PHONY: gen check test

gen:
	c-for-go -ccincl ./roc.yml

check:
	go build ./roc
	golangci-lint run ./roc

test:
	go test ./roc

fmt:
	gofmt -s -w ./roc
