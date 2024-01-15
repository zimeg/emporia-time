.PHONY: build test release clean

BIN=etime

build:
	go build -o $(BIN)

test: build
	go test ./...

release: clean
	goreleaser build --snapshot

clean:
	rm -f $(BIN)
	rm -rf ~/.config/etime
	rm -rf dist
