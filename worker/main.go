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

	//before we access the database, we want to have a job.

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

	var result bson.M
	coll := client.Database("spellwrong").Collection("misspellings")
	err = coll.FindOne(context.TODO(), bson.D{}).Decode(&result)
	if err != nil {
		panic(fmt.Sprintf("Failed to retrieve document: %v", err))
	}
	fmt.Printf("Retrieved document: %+v\n", result)
}
