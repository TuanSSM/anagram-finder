package main

import (
	"log"

	"github.com/tuanssm/anagram-finder/internal/api"
)

func main() {
	listenAddr := ":8080"                         //os.Getenv("ANAGRAM_FINDER_API_PORT")
	mongoURI := "mongodb://mongodb-service:27017" //os.Getenv("MONGO_URI")

	log.Printf("MONGO URI: %v", mongoURI)

	apiServer := api.NewApiServer()
	log.Fatal(apiServer.Start(listenAddr, mongoURI))
}
