package database

import (
	"context"
	"log"
	"os"

	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/options"
)

var MongoDb = os.Getenv("MONGODB")

func MongoInstance() *mongo.Client {
	// Create a new client and connect to the server
	client, err := mongo.Connect(context.Background(), options.Client().ApplyURI(MongoDb))
	if err != nil {
		log.Fatalln("Unable to establish connection:", err.Error())
	}
	log.Println("Connected to MongoDB!")
	return client
}

func OpenCollection(client *mongo.Client, name string) *mongo.Collection {
	var collection *mongo.Collection = client.Database("usersDB").Collection(name)
	return collection
}
