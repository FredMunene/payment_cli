package main

import (
	"encoding/json"
	"fmt"
	"io"
	"log"
	"net/http"
	"os"

	"payment_cli/internal"
)

func main() {
	internal.Init()

	apiKey := os.Getenv("SWYPT_API_KEY")
	apiSecret := os.Getenv("SWYPT_API_SECRET")
	// apiUrl := os.Getenv("SWYPT_API_BASE_URL")

	if apiKey == "" || apiSecret == "" {
		log.Println("Warning: SWYPT_API_KEY or SWYPT_API_SECRET environment variables are not set.")
	}

	fmt.Println("Welcome to the USDT Savings CLI!")
	fmt.Println("You can save an amount in USDT and protect yourself from local currency inflation.")

	// Retrieve assets from the API
	client := &http.Client{}
	req, err := http.NewRequest("GET", "https://pool.swypt.io/api/swypt-supported-assets", nil)
	if err != nil {
		fmt.Println("Error creating request:", err)
		return
	}
	req.Header.Add("x-api-key", apiKey)
	req.Header.Add("x-api-secret", apiSecret)

	resp, err := client.Do(req)
	if err != nil {
		fmt.Println("Error making request:", err)
		return
	}
	defer resp.Body.Close()

	body, err := io.ReadAll(resp.Body)
	if err != nil {
		fmt.Println("Error reading response:", err)
		return
	}

	var assetsResp map[string]any
	err = json.Unmarshal(body, &assetsResp)
	if err != nil {
		fmt.Println("Error parsing response:", err)
		return
	}

	fmt.Println("Supported Assets:")
	if crypto, ok := assetsResp["crypto"].(map[string]any); ok {
		for network, assets := range crypto {
			fmt.Printf("Network: %s\n", network)
			if assetList, ok := assets.([]any); ok {
				for i, asset := range assetList {
					if assetMap, ok := asset.(map[string]any); ok {
						fmt.Printf("  %d. %s (%s)\n", i+1, assetMap["name"], assetMap["symbol"])
					}
				}
			}
		}
	} else {
		fmt.Println("No assets found or unexpected response format.")
	}

	
	// Next: show user the assets and ask to select one

}
