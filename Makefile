default: test

test:
	go test ./...

staticcheck:
	staticcheck ./...

build: plugin.so

plugin.so: go.mod go.sum $(shell find . ../riposo/ -name '*.go')
	go build -ldflags '-s -w' -buildmode=plugin -o $@ .
