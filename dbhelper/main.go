package main

import (
	"context"
	"encoding/json"
	"fmt"
	"os"
	"time"

	"go.mongodb.org/mongo-driver/v2/bson"
	"go.mongodb.org/mongo-driver/v2/mongo"
	"go.mongodb.org/mongo-driver/v2/mongo/options"
)

func main() {
	uri := os.Getenv("DB")
	maxRetries := 5
	var client *mongo.Client
	var err error

	for range maxRetries { //new way of doing for loops. cool.
		client, err = mongo.Connect(options.Client().ApplyURI(uri))
		if err == nil {
			fmt.Println("Connected to database")
			break
		}

		time.Sleep(time.Duration(10) * time.Second)

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
		fmt.Println("Database is empty. Populating with misspellings.json...")
		file, err := os.ReadFile("misspellings.json")
		if err != nil {
			panic("Failed to read misspellings.json")
		}

		// Unmarshal the JSON file into a stucture that fits out object
		var misspellingsData map[string]any
		if err := json.Unmarshal(file, &misspellingsData); err != nil {
			panic("Failed to unmarshal JSON")
		}

		// FIX: Insert the map directly, without wrapping it in another object.
		_, err = coll.InsertOne(context.TODO(), misspellingsData)
		if err != nil {
			panic("Failed to insert document into database")
		}
		fmt.Println("Created database with misspellings")
	}
	fmt.Println("Database ready.")

}
