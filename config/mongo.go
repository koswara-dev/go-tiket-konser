package config

import (
	"context"
	"log"
	"os"
	"time"

	"go.mongodb.org/mongo-driver/v2/mongo"
	"go.mongodb.org/mongo-driver/v2/mongo/options"
)

var MongoClient *mongo.Client
var MongoDatabase *mongo.Database

func InitMongo() {
	mongoURI := os.Getenv("MONGO_URI")
	mongoDBName := os.Getenv("MONGO_DB")

	if mongoURI == "" {
		mongoURI = "mongodb://localhost:27017"
	}
	if mongoDBName == "" {
		mongoDBName = "eticketdb"
	}

	ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
	defer cancel()

	client, err := mongo.Connect(options.Client().ApplyURI(mongoURI))
	if err != nil {
		log.Fatalf("failed to create MongoDB client: %v", err)
	}

	err = client.Ping(ctx, nil)
	if err != nil {
		log.Printf("Warning: Failed to ping MongoDB: %v\n", err)
	} else {
		log.Println("MongoDB connected successfully")
	}

	MongoClient = client
	MongoDatabase = client.Database(mongoDBName)
}
