package main

import (
	"log"
	"microservices/authentication/internal/api"
	"microservices/authentication/internal/database"
	"os"

	"github.com/joho/godotenv"
)

func main() {
	if err := godotenv.Load(); err != nil {
		log.Fatal("Couldn't load the .env file!")
	}

	port := os.Getenv("PORT")

	if port == "" {
		log.Fatal("Couldn't identify PORT variable!")
	}

	mongoDBUri := os.Getenv("MONGODB_URI")

	if mongoDBUri == "" {
		log.Fatal("Couldn't identify MONGODB_URI variable!")
	}

	database.SetupDBConnection(mongoDBUri)

	uri := ":" + port
	api.SetupAPI(uri)
}
