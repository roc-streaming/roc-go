GO111MODULE := on
export GO111MODULE

all: check test race

.PHONY: check
check:
	cd roc && go build .
	cd roc && go test . -run xxx
	cd roc && golangci-lint run .

.PHONY: test
test:
	cd roc && go test .

.PHONY: race
race:
	cd roc && go test -race .

.PHONY: fmt
fmt:
	cd roc && gofmt -s -w .
