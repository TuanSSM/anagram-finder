package main

import (
	"log"
	"os"
)

func main() {
	svc := NewAnagramService()
	svc = NewLoggingService(svc)

	apiServer := NewApiServer(svc)
	var listenAddr string = ":" + os.Getenv("ANAGRAM_FINDER_API_PORT")
	log.Fatal(apiServer.Start(listenAddr))
}
