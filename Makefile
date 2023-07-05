build:
	@go build -o ./bin/anagram-finder -v ./cmd/anagram-finder

run: build
	./bin/anagram-finder

test:
	go test -v ./...
