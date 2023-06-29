build:
	@go build -o bin/anagram-finder


run: build
	./bin/anagram-finder

test:
	go test -v ./...
