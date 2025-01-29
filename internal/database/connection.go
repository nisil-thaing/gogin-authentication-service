package database

import (
	"context"
	"fmt"
	"log"
	"time"

	"go.mongodb.org/mongo-driver/v2/mongo"
	"go.mongodb.org/mongo-driver/v2/mongo/options"
)

const DB_TIMEOUT = 10 * time.Second

var dbClient *mongo.Client

func SetupDBConnection(connectionUri string) {
	if connectionUri == "" {
		log.Fatal("Couldn't identify the database's connection string!")
	}

	ctx, cancel := context.WithTimeout(context.Background(), DB_TIMEOUT)
	defer cancel()

	clientOptions := options.Client().ApplyURI(connectionUri)

	client, err := mongo.Connect(clientOptions)
	if err != nil {
		log.Fatal(err)
	}

	if err = client.Ping(ctx, nil); err != nil {
		log.Fatalf("❌ Couldn't ping to the MongoDB! %v", err)
	}

	fmt.Println("✅ The MongoDB connection has been successfully established!")

	dbClient = client
}

func GetDBInstance() *mongo.Client {
	return dbClient
}

func OpenCollection(client *mongo.Client, databaseName string, collectionName string) *mongo.Collection {
	database := client.Database(databaseName)
	collection := database.Collection(collectionName)

	return collection
}
