package main

import (
	"fmt"
	"log"
)

func main() {
	// Vulnerable code: hardcoded API key, which should not be in source code
	apiKey := "12345-SECRET-KEY" // Sensitive information hardcoded

	// Simulate using the key in an API call
	fmt.Printf("Using API Key: %s\n", apiKey)

	// Log sensitive information (bad practice)
	log.Printf("Logging sensitive key: %s", apiKey)
}
