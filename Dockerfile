FROM golang:latest as builder
WORKDIR /anagram-app
COPY go.mod go.sum Makefile ./
RUN go mod download
COPY . .
RUN CGO_ENABLED=0 GOOS=linux make build

FROM scratch
WORKDIR /
COPY --from=builder /anagram-app/bin/anagram-finder .
EXPOSE 8080
CMD ["/anagram-finder"]