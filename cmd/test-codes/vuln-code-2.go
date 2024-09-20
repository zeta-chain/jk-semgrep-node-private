package main

import (
	"context"
	"fmt"
	"os/exec"
	"time"
)

func main() {
	// Simulate user input (e.g., from a web request, CLI argument, etc.)
	userInput := "ls; rm -rf /" // Potentially dangerous input

	// Create a context with a timeout to avoid hanging processes
	ctx, cancel := context.WithTimeout(context.Background(), 2*time.Second)
	defer cancel()

	// Vulnerable code: passing unsanitized user input to exec.CommandContext
	cmd := exec.CommandContext(ctx, "sh", "-c", userInput)

	// Execute the command and capture output
	output, err := cmd.CombinedOutput()
	if err != nil {
		fmt.Printf("Error executing command: %v\n", err)
	}

	// Print the command output
	fmt.Printf("Command output: %s\n", output)
}

