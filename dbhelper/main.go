package main

import (
	"context"
	"encoding/json"
	"fmt"
	"log"
	"os"
	"time"

	"go.mongodb.org/mongo-driver/v2/bson"
	"go.mongodb.org/mongo-driver/v2/mongo"
	"go.mongodb.org/mongo-driver/v2/mongo/options"
)

func main() {
	uri := os.Getenv("DB")
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
	if err != nil {
		panic("Could not read from database")
	}
	if count == 0 {
		//collection does not exist. Create it
		log.Println("Database is empty. Populating with misspellings.json...")
		file, err := os.ReadFile("misspellings.json")
		if err != nil {
			log.Fatalf("Failed to read misspellings.json: %v", err)
		}

		// Unmarshal the JSON into a generic map
		var misspellingsData map[string]interface{}
		if err := json.Unmarshal(file, &misspellingsData); err != nil {
			log.Fatalf("Failed to unmarshal JSON: %v", err)
		}

		// FIX: Insert the map directly, without wrapping it in another object.
		_, err = coll.InsertOne(context.TODO(), misspellingsData)
		if err != nil {
			log.Fatalf("Failed to insert document: %v", err)
		}
		fmt.Println("Created database with misspellings")
	}
	fmt.Println("Database ready.")

}
