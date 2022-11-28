.PHONY: build test

BIN=etime

build:
	go build -o $(BIN)

test: build
	./$(BIN) sleep 4

clean:
	rm -f $(BIN)
