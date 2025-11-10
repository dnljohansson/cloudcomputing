package main

import (
	"bytes"
	"encoding/json"
	"fmt"
	"log"
	"net/http"
	"os"
	"strings"
	"sync"

	"github.com/gin-contrib/cors"
	"github.com/gin-gonic/gin"
)

var client = &http.Client{
	Transport: &http.Transport{
		DisableKeepAlives: true,
	},
}

//allows the load balancer to do its work, as we are forcing this client to do a DNS lookup for every connection

const port = "3000"

// RequestData struct defines the expected structure of the JSON payload
type RequestData struct {
	Text      string `json:"text"`
	Mistakes  int    `json:"mistakes"`
	WordCount int    `json:"wordCount"`
}

type Work struct {
	Segment []string `json:"segment"`
}

func sendWork(work *Work, jobs *sync.WaitGroup, results [][]string, index int) {
	workerURL := os.Getenv("WORKER_URL")
	defer jobs.Done()

	log.Printf("▶️  Sending segment to worker: \"%.30s...\"", work.Segment)
	marshaled, err := json.Marshal(work)
	if err != nil {
		panic("Failed to marshal work.")
	}
	req, err := http.NewRequest("POST", workerURL, bytes.NewBuffer(marshaled))
	if err != nil {
		panic("Failed to create post request")
	}
	req.Header.Set("Content-Type", "application/json")

	resp, err := client.Do(req)
	if err != nil {
		panic("Failed to send request to worker")
	}

	defer resp.Body.Close()

	if resp.StatusCode != http.StatusOK {
		panic("Worker returned non-OK status code")
	}

	var respContent Work

	if err := json.NewDecoder(resp.Body).Decode(&respContent); err != nil {
		panic("Failed to decode worker response")
	}

	log.Printf("✅ Received processed segment from worker: \"%.30s...\"", respContent.Segment)

	results[index] = respContent.Segment

}

func main() {

	router := gin.Default()
	router.Use(cors.Default()) //permits communication

	// Define the POST endpoint
	router.POST("/", func(c *gin.Context) {
		var data RequestData

		if err := c.ShouldBindJSON(&data); err != nil {
			c.JSON(http.StatusBadRequest, gin.H{"error": "Error decoding JSON: " + err.Error()})
			return
		}

		log.Printf("Received data: Text='%.20s...', Mistakes=%d\n", data.Text, data.Mistakes)

		//here is where we split the data
		//var usertext = data.Text
		var segmentCount = data.Mistakes
		var jobs sync.WaitGroup
		words := strings.Fields(data.Text)
		wordCount := data.WordCount
		results := make([][]string, segmentCount)
		wordsPerSegment := wordCount / segmentCount
		fmt.Printf("Jobs for this request: %d", segmentCount)
		//so what we do here is to split the text into segments, which are then put in the queue
		for i := 0; i < segmentCount; i++ {
			start := i * wordsPerSegment
			end := (i + 1) * wordsPerSegment
			if i == segmentCount-1 {
				end = wordCount
			}
			segment := words[start:end]

			work := Work{Segment: segment}

			jobs.Add(1)

			//subroutine, allows async operations
			go sendWork(&work, &jobs, results, i)

		}
		//wait for all jobs to finish
		jobs.Wait()

		var finalWords []string
		for _, segment := range results {
			finalWords = append(finalWords, segment...)
		}
		finalText := strings.Join(finalWords, " ")

		c.JSON(http.StatusOK, gin.H{
			"message":      "All jobs completed successfully!",
			"receivedText": finalText,
		})
	})

	fmt.Printf("Server starting on port %s\n", port)
	// Start the server
	if err := router.Run(":" + port); err != nil {
		log.Fatalf("Error starting server: %v", err)
	}
}
