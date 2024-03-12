.PHONY: check build test staging release clean

BIN=etime
VERSION="$(shell git describe --dirty --tags --always)"

check:
	golangci-lint run ./...

build: check
	go build -o $(BIN) -ldflags "-X main.version=${VERSION}"

test: build
	go test ./...

staging: clean
	goreleaser build --snapshot --config .goreleaser.staging.yml

release: clean
	goreleaser build --snapshot --config .goreleaser.release.yml

clean:
	rm -f $(BIN)
	rm -rf ~/.config/etime
	rm -f .gon.hcl
	rm -rf dist
