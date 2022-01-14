package db

import (
	"context"
	"log"
	"os"
	"time"

	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/options"
)

const databaseName = "chess"
const GamesCollectionName = "games"

func Connect() *mongo.Database {
	MONGODB_URI := os.Getenv("MONGODB_URI")

	if len(MONGODB_URI) == 0 {
		log.Fatal("Missing MONGODB_URI environment variable.")
	}

	ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
	defer cancel()

	clientOptions := options.Client().ApplyURI(MONGODB_URI)
	client, err := mongo.Connect(ctx, clientOptions)

	if err != nil {
		log.Fatal(err)
	}

	return client.Database(databaseName)
}

func Disconnect(client *mongo.Database) {
	ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
	defer cancel()

	if err := client.Client().Disconnect(ctx); err != nil {
		log.Fatal(err)
	}
}
