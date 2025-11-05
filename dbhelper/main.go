package main

import (
	"context"
	"fmt"
	"os"
	"time"

	"go.mongodb.org/mongo-driver/v2/bson"
	"go.mongodb.org/mongo-driver/v2/mongo"
	"go.mongodb.org/mongo-driver/v2/mongo/options"
)

func main() {
	uri := os.Getenv("TEXTSTORE_HOST")
	timeout := 10
	maxRetries := 5
	var client *mongo.Client
	var err error

	for i := 0; i < maxRetries; i++ {
		client, err = mongo.Connect(options.Client().ApplyURI(uri))
		if err == nil {
			fmt.Println("Connected to database")
			break
		}

		time.Sleep(time.Duration(timeout) * time.Second)

	}
	if err != nil {
		panic("Could not connect to database")
	}

	defer func() {
		if err := client.Disconnect(context.TODO()); err != nil {
			panic(err)
		}
	}()

	// we now check if the collection of misspelled words exists

	coll := client.Database("spellwrong").Collection("misspellings")

	count, err := coll.CountDocuments(context.TODO(), bson.D{})
	if err != nil || count == 0 {
		//collection does not exist. Create it
		fmt.Println("Created database with misspellings")
	}
	fmt.Println("Database ready.")

}
