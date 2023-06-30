build:
	@go build -o ./bin/anagram-finder -v ./app/

run: build
	./bin/anagram-finder

test:
	go test -v ./...
