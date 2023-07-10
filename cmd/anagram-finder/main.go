package main

import (
	"log"
	"os"

	"github.com/tuanssm/anagram-finder/internal/api"
)

func main() {
	listenAddr := ":" + os.Getenv("ANAGRAM_FINDER_API_PORT")
	mongoURI := os.Getenv("MONGO_URI")

	apiServer := api.NewApiServer()
	log.Fatal(apiServer.Start(listenAddr, mongoURI))
}
