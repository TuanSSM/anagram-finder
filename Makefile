build:
	@go build -o /anagram-finder

run: build
	/anagram-finder

test:
	go test -v ./...
