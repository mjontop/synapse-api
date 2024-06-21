package db

import (
	"context"
	"fmt"
	"log"
	"os"
	"time"

	"github.com/joho/godotenv"
	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/options"
)

var Client *mongo.Client

func connectToDB() (*mongo.Database, error) {

	mongoURI := os.Getenv("MONGODB_URI")
	databaseName := os.Getenv("DB_NAME")
	if mongoURI == "" {
		return nil, fmt.Errorf("MONGODB_URI environment variable is not set")
	}

	if databaseName == "" {
		return nil, fmt.Errorf("error: Database name is invalid")
	}

	clientOptions := options.Client().ApplyURI(mongoURI)
	ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
	defer cancel()

	client, err := mongo.Connect(ctx, clientOptions)
	if err != nil {
		return nil, err
	}

	err = client.Ping(ctx, nil)
	if err != nil {
		return nil, err
	}

	fmt.Println("Connected to MongoDB!")

	db := client.Database(databaseName)

	Client = client

	return db, nil

}

func ConnectDB() {

	if err := godotenv.Load(); err != nil {
		log.Fatalf("Error loading .env file: %v", err)
	}

	_, err := connectToDB()
	if err != nil {
		log.Fatal(err)
	}
}

func GetCollection(collectionName string) *mongo.Collection {
	return Client.Database("synapse").Collection(collectionName)
}
