package main

import (
	"fmt"
	"os"
	"payment_cli/internal"
)




func main() {
	// Load environment variables

	internal.Init()
	apiKey := os.Getenv("SWYPT_API_KEY")
	apiSecret := os.Getenv("SWYPT_API_SECRET")

	if apiKey == "" || apiSecret == "" {
		fmt.Println("Warning: SWYPT_API_KEY or SWYPT_API_SECRET environment variables are not set.")
	}

	fmt.Println("Welcome to the USDT Savings CLI!")
	fmt.Println("You can save an amount in USDT and protect yourself from local currency inflation.")

	// Next: retrieve assets from the API
}
