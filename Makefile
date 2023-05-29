GO111MODULE := on
export GO111MODULE

all: gen build lint test

gen:
	cd roc && go generate

build:
	cd roc && go build .
	cd roc && go test . -run xxx

lint:
	cd roc && golangci-lint run .

test:
	cd roc && go test .
	cd roc && GODEBUG=cgocheck=2 go test -count=1 .
	cd roc && go test -race .

clean:
	cd roc && go clean -cache -testcache

tidy:
	cd roc && go mod tidy

fmt:
	cd roc && gofmt -s -w .
