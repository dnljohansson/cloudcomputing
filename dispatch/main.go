package main

import (
	"encoding/json"
	"fmt"
	"io"
	"log"
	"net/http"
)

const port = "3000"

// RequestData struct defines the expected structure of the JSON payload
type RequestData struct {
	Text string `json:"text"`
	Rate int    `json:"rate"`
}

// corsMiddleware adds the necessary CORS headers to each response.
func corsMiddleware(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		// Set headers
		w.Header().Set("Access-Control-Allow-Origin", "*") // Allow any origin
		w.Header().Set("Access-Control-Allow-Methods", "POST, GET, OPTIONS, PUT, DELETE")
		w.Header().Set("Access-Control-Allow-Headers", "Accept, Content-Type, Content-Length, Accept-Encoding, X-CSRF-Token, Authorization")

		// Handle pre-flight requests
		if r.Method == "OPTIONS" {
			w.WriteHeader(http.StatusOK)
			return
		}

		// Next
		next.ServeHTTP(w, r)
	})
}

// submitHandler processes requests to the /api/submit endpoint
func submitHandler(w http.ResponseWriter, r *http.Request) {
	log.Println("Received request on /") // Log that a request was received

	if r.Method != http.MethodPost {
		http.Error(w, "Invalid request method", http.StatusMethodNotAllowed)
		return
	}

	body, err := io.ReadAll(r.Body)
	if err != nil {
		http.Error(w, "Error reading request body", http.StatusInternalServerError)
		return
	}
	defer r.Body.Close()

	var data RequestData
	err = json.Unmarshal(body, &data)
	if err != nil {
		http.Error(w, "Error decoding JSON", http.StatusBadRequest)
		return
	}

	log.Printf("Received data: Text='%.20s...', Rate=%d\n", data.Text, data.Rate)

	// Respond to the client
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusOK)
	json.NewEncoder(w).Encode(map[string]string{"message": "Request received successfully!", "receivedText": data.Text})
}

func main() {
	// Create a new ServeMux
	mux := http.NewServeMux()
	mux.HandleFunc("/", submitHandler)

	// Wrap the mux with the CORS middleware
	handler := corsMiddleware(mux)

	fmt.Printf("✅ Go server starting on port %s\n", port)
	if err := http.ListenAndServe(":"+port, handler); err != nil {
		log.Fatalf("Error starting server: %v", err)
	}
}
