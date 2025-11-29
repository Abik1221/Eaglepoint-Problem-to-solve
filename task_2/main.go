package main

import (
	"encoding/json"
	"fmt"
	"math/rand"
	"strings"
	"time"
)

// this struct represents the structure of our API response
type Response struct {
	Success bool   `json:"success"`
	URL     string `json:"url"`
	Data    Data   `json:"data"`
}

// data holds the actual payload from the API
type Data struct {
	Message   string `json:"message"`
	Timestamp string `json:"timestamp"`
	Value     int    `json:"value"`
}

// callAPI simulates a network request that might fail
func callAPI(url string) (*Response, error) {
	// Simulate network latency (500ms)
	time.Sleep(500 * time.Millisecond)

	// randomly succeed or fail to test our retry logic
	// 60% chance of failure, 40% chance of success
	if rand.Float64() > 0.6 {
		return &Response{
			Success: true,
			URL:     url,
			Data: Data{
				Message:   "Got the data!",
				Timestamp: time.Now().Format(time.RFC3339),
				Value:     rand.Intn(1000),
			},
		}, nil
	}

	return nil, fmt.Errorf("couldn't fetch from %s", url)
}

// fetch tries to get data from the URL, retrying if it fails
func fetch(url string, retries int) (*Response, error) {
	var err error

	// loop through our allowed number of attempts
	for i := 1; i <= retries; i++ {
		fmt.Printf("Try %d/%d...\n", i, retries)

		// Make the actual call
		resp, e := callAPI(url)

		// If it worked, we're done! Return the data immediately.
		if e == nil {
			fmt.Printf("Success!\n")
			return resp, nil
		}

		// If it failed, save the error and log it
		err = e
		fmt.Printf("Failed: %s\n", e.Error())

		// If we still have retries left, wait a bit before trying again.
		// This "backoff" is crucial for not hammering a struggling server.
		if i < retries {
			fmt.Println("Retrying in 1 sec...\n")
			time.Sleep(time.Second)
		}
	}

	// If we're here, we used up all our retries and still failed.
	return nil, fmt.Errorf("gave up after %d tries: %s", retries, err.Error())
}

func main() {
	// Seed the random number generator so we get different results each run
	rand.Seed(time.Now().UnixNano())

	url := "https://api.example.com/data"

	fmt.Printf("Fetching: %s\n\n", url)
	fmt.Println(strings.Repeat("=", 50) + "\n")

	// Try to fetch the data with up to 3 retries
	data, err := fetch(url, 3)

	fmt.Println("\n" + strings.Repeat("=", 50))

	if err != nil {
		// It failed completely
		fmt.Println("Error:")
		fmt.Println(err.Error())
	} else {
		// It worked! Show the result.
		fmt.Println("Result:")
		output, _ := json.MarshalIndent(data, "", "  ")
		fmt.Println(string(output))
	}
}
