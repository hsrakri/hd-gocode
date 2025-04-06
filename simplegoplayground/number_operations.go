package main

import (
	"fmt"
)

func main() {
	// Create a slice to store 5 numbers
	numbers := make([]int, 5)

	// Get input from user
	fmt.Println("Please enter 5 numbers:")
	for i := 0; i < 5; i++ {
		fmt.Printf("Enter number %d: ", i+1)
		fmt.Scan(&numbers[i])
	}

	// Calculate sum
	sum := 0
	for _, num := range numbers {
		sum += num
	}

	// Find largest number
	largest := numbers[0]
	for _, num := range numbers {
		if num > largest {
			largest = num
		}
	}

	// Print results
	fmt.Println("\nResults:")
	fmt.Println("Numbers:", numbers)
	fmt.Printf("Sum: %d\n", sum)
	fmt.Printf("Largest number: %d\n", largest)
}
