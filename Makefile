.PHONY: build test

BIN=etime

build:
	go build -o $(BIN)

run: build
	./$(BIN) sleep 4

test: build
	go test

clean:
	rm -f $(BIN)
