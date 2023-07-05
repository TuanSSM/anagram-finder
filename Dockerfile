FROM golang:latest as builder
WORKDIR /app
COPY . .
RUN go mod download
RUN CGO_ENABLED=0 GOOS=linux make build

FROM scratch
WORKDIR /
COPY --from=builder /app/bin/anagram-finder .
COPY --from=builder /app/web ./

ENV ANAGRAM_FINDER_API_PORT=${ANAGRAM_FINDER_API_PORT:-3000}
CMD ["./anagram-finder"]