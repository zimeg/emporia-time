.PHONY: build test release clean

BIN=etime
VERSION="$(shell git describe --dirty --tags --always)"

build:
	go build -o $(BIN) -ldflags "-X main.version=${VERSION}"

test: build
	go test ./...

release: clean
	goreleaser build --snapshot

clean:
	rm -f $(BIN)
	rm -rf ~/.config/etime
	rm -rf dist
