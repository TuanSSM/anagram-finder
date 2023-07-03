FROM golang:latest as builder
WORKDIR /anagram-app
COPY go.mod go.sum Makefile ./
RUN go mod download
COPY . .
RUN CGO_ENABLED=0 GOOS=linux make build

FROM scratch
WORKDIR /
COPY --from=builder /anagram-app/bin/anagram-finder .

ENV ANAGRAM_FINDER_API_PORT=${ANAGRAM_FINDER_API_PORT:-3000}
CMD ["/anagram-finder"]