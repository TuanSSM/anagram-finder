package main

import (
	"log"

	"github.com/tuanssm/anagram-finder/internal/api"
)

func main() {
	//listenAddr := ":" + os.Getenv("ANAGRAM_FINDER_API_PORT")
	//mongoURI := os.Getenv("MONGO_URI")
	listenAddr := ":8080"
	mongoURI := "mongodb://mongodb:27017"

	apiServer := api.NewApiServer()
	log.Fatal(apiServer.Start(listenAddr, mongoURI))
}
