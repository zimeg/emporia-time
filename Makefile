.PHONY: build test

BIN=etime

build:
	go build -o $(BIN)

test: build
	go test ./...

clean:
	rm -f $(BIN)
	rm -rf ~/.config/etime
