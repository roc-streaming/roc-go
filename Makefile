all: check test race

.PHONY: check
check:
	go build ./roc
	go test ./roc -run xxx
	golangci-lint run ./roc

.PHONY: test
test:
	go test ./roc

.PHONY: race
race:
	go test -race ./roc

.PHONY: fmt
fmt:
	gofmt -s -w ./roc
