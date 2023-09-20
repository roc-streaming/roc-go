GO111MODULE := on
export GO111MODULE

all: gen build lint testall

gen:
	cd roc && go generate

build:
	cd roc && go build .
	cd roc && go test . -run xxx

lint:
	cd roc && golangci-lint run .

test:
	cd roc && go test . -count=1 .

testall:
	cd roc && go test . -count=1 .
	cd roc && go test -count=1 -race .
	cd roc && GOEXPERIMENT=cgocheck2 go build && go test -count=1 .

clean:
	cd roc && go clean -cache -testcache

tidy:
	cd roc && go mod tidy

fmt:
	cd roc && gofmt -s -w .
