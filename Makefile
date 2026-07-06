.PHONY: run build lint format test clean

BINARY := bin/shell

run:
	go run ./app

build:
	go build -o $(BINARY) ./app

lint:
	go vet ./...
	@unformatted=$$(gofmt -l app); \
	if [ -n "$$unformatted" ]; then \
		echo "gofmt needed on:"; echo "$$unformatted"; exit 1; \
	fi

format:
	gofmt -w app

test:
	go test ./...

clean:
	rm -rf bin
