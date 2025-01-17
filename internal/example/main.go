package main

import (
	"fmt"
	"go-shell/internal/inputReader"
)

func main() {
	fmt.Println("Type something and use arrow keys to navigate. Press Enter to finish:")

	// Initialize the InputReader
	input := inputReader.NewInputReader()

	// Call the Read method to start input capture and handle navigation
	finalInput, err := input.Read()
	if err != nil {
		fmt.Println("Error reading input:", err)
		return
	}

	// Output the final input captured
	fmt.Printf("Final input: %s\n", finalInput)
}
