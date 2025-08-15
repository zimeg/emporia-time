.PHONY: check build test staging release clean

BIN=etime
VERSION="$(shell git describe --dirty --tags --always)"

check:
	golangci-lint run ./...

build: check
	go build -o $(BIN) -ldflags "-X main.version=${VERSION}"

test: build
	go test -v ./...

coverage: build
	mkdir -p coverage
	go test -v 2>&1 -coverprofile=coverage/coverage.txt ./... | tee coverage/results.out
	go tool cover -html coverage/coverage.txt -o coverage/coverage.html
	go-junit-report -in coverage/results.out -set-exit-code > coverage/coverage.xml

staging: clean
	goreleaser build --snapshot --config .goreleaser.staging.yaml

release: clean
	goreleaser build --snapshot --config .goreleaser.release.yaml

clean:
	rm -f $(BIN)
	rm -rf ~/.config/etime
	rm -f .gon.hcl
	rm -rf coverage
	rm -rf dist
	rm -f result
