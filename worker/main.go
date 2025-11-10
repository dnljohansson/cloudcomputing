package main

import (
	"context"
	"fmt"
	"log"
	"math/rand"
	"net/http"
	"os"
	"strings"
	"sync"
	"time"

	"github.com/gin-gonic/gin"
	"go.mongodb.org/mongo-driver/v2/bson"
	"go.mongodb.org/mongo-driver/v2/mongo"
	"go.mongodb.org/mongo-driver/v2/mongo/options"
)

const port = "3001"

type Work struct {
	Segment []string `json:"segment"`
}

var ginMutex sync.Mutex

func main() {
	uri := os.Getenv("DB")

	//before we access the database, we want to have a job.

	timeout := 5
	maxRetries := 5
	var client *mongo.Client
	var coll *mongo.Collection
	var err error

	for i := 0; i < maxRetries; i++ {
		client, err = mongo.Connect(options.Client().ApplyURI(uri))

		coll = client.Database("spellwrong").Collection("misspellings")
		count, err := coll.CountDocuments(context.TODO(), bson.D{})
		if count != 0 && err == nil {
			fmt.Println("Connected to database")
			break
		}

		fmt.Println("Database not ready. Retrying...")

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

	var misspellings map[string][]string

	options := options.FindOne().SetProjection(bson.M{"_id": 0})
	err = coll.FindOne(context.TODO(), bson.D{}, options).Decode(&misspellings)
	if err != nil {
		panic(fmt.Sprintf("Failed to retrieve document: %v", err))
	}
	fmt.Printf("Retrieved misspellings document.")

	fmt.Printf("Misspellings on abandon key: %s\n", misspellings["abandon"])

	router := gin.Default()
	//Gin is concurrent by standard, which makes little sense for this assignment.
	//Therefore, we use a mutex to demonstrate horizontal scaling, i.e. the only "parallellism" achieved is by horizontal scaling

	router.POST("/", func(c *gin.Context) {
		ginMutex.Lock()

		defer ginMutex.Unlock()

		var work Work
		if err := c.BindJSON(&work); err != nil {
			log.Printf("Could not start working")
			c.JSON(http.StatusBadRequest, gin.H{"error": "bad request"})
			return
		}
		processedText := work.Segment
		fmt.Printf("Started work on: %s", processedText)

		time.Sleep(1 * time.Second) // Simulate work

		for attempts := 0; attempts < 5; attempts++ {
			index := rand.Intn(len(processedText))
			value, ok := misspellings[strings.ToLower(processedText[index])]
			if ok {
				var word string
				if len(value) > 1 {
					word = value[rand.Intn(len(value))]
				} else {
					word = value[0]
				}
				processedText[index] = word
				break
			}
		}

		//might have swapped the word. maybe not
		c.JSON(http.StatusOK, gin.H{"segment": processedText})
	})

	fmt.Printf("✅ Worker server starting on port %s\n", port)
	router.Run(":" + port)

}
