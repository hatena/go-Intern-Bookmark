BIN=server

all: clean build

build: deps
	go generate ./...
	go build -o build/$(BIN) .

run: build
	build/$(BIN)

setup:
	command -v dep >/dev/null || go get -u github.com/golang/dep/cmd/dep
	command -v go-assets-builder >/dev/null || go get -u github.com/jessevdk/go-assets-builder
	command -v reflex >/dev/null || go get -u github.com/cespare/reflex

deps:
	dep ensure -vendor-only

test: build
	DATABASE_DSN=$$DATABASE_DSN_TEST go test -v -p 1 ./...

clean:
	rm -rf build */*-gen.go
	go clean

.PHONY: build setup deps test clean
