package main

import (
	"log"
)

func main() {
	svc := NewAnagramService()
	svc = NewLoggingService(svc)

	apiServer := NewApiServer(svc)
	log.Fatal(apiServer.Start(":8080"))
}
