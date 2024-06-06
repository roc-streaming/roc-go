GO111MODULE := on
export GO111MODULE

ifeq ($(shell uname -s),Linux)
TERM := xterm
export TERM
endif

ifneq ($(shell which gotest),)
gotest := gotest
else
gotest := go test
endif

all: gen build lint testall

gen:
	cd roc && go generate

build:
	cd roc && go build .
	cd roc && $(gotest) -run none .

lint:
	cd roc && golangci-lint run .

test:
	cd roc && $(gotest) -count=1 .

testall:
	cd roc && $(gotest) -count=1 .
	cd roc && $(gotest) -count=1 -race .
	cd roc && GOEXPERIMENT=cgocheck2 go build && $(gotest) -count=1 .

clean:
	cd roc && go clean -cache -testcache

tidy:
	cd roc && go mod tidy

fmt:
	cd roc && gofmt -s -w .
