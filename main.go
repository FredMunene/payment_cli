package main

import (
	"bytes"
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

	// Only show LISK network and Tether LISK (USDT) as default
	fmt.Print("Supported Asset (Default):")
	if crypto, ok := assetsResp["crypto"].(map[string]any); ok {
		if liskAssets, ok := crypto["lisk"].([]any); ok {
			for _, asset := range liskAssets {
				if assetMap, ok := asset.(map[string]any); ok {
					if assetMap["symbol"] == "USDT" {
						fmt.Printf("Lisk -")
						fmt.Printf(" %s (%s)\n", assetMap["name"], assetMap["symbol"])
					}
				}
			}
		} else {
			fmt.Println("No LISK assets found.")
		}
	} else {
		fmt.Println("No assets found or unexpected response format.")
	}
	// Next: ask user for the amount in KES they want to save
	var amount int
	for {
		fmt.Print("Enter the amount in KES you want to save (minimum 50): ")
		_, err := fmt.Scan(&amount)
		if err != nil {
			fmt.Println("Invalid input. Please enter a valid integer amount.")
			// Clear input buffer
			var discard string
			fmt.Scanln(&discard)
			continue
		}
		if amount < 50 {
			fmt.Println("Amount must be at least 50 KES. Please try again.")
			continue
		}
		break
	}
	fmt.Printf("You have chosen to save %d KES.\n", amount)
	// Next: call the API endpoint with the value provided
	quotePayload := map[string]any{
		"type": "onramp",
		"amount": fmt.Sprintf("%d", amount),
		"fiatCurrency": "KES",
		"cryptoCurrency": "USDT",
		"network": "lisk",
	}
	payloadBytes, err := json.Marshal(quotePayload)
	if err != nil {
		fmt.Println("Error encoding quote payload:", err)
		return
	}

	quoteReq, err := http.NewRequest("POST", "https://pool.swypt.io/api/swypt-quotes", io.NopCloser(bytes.NewReader(payloadBytes)))
	if err != nil {
		fmt.Println("Error creating quote request:", err)
		return
	}
	quoteReq.Header.Add("x-api-key", apiKey)
	quoteReq.Header.Add("x-api-secret", apiSecret)
	quoteReq.Header.Add("Content-Type", "application/json")

	quoteResp, err := client.Do(quoteReq)
	if err != nil {
		fmt.Println("Error making quote request:", err)
		return
	}
	defer quoteResp.Body.Close()

	quoteBody, err := io.ReadAll(quoteResp.Body)
	if err != nil {
		fmt.Println("Error reading quote response:", err)
		return
	}

	var quoteRespData map[string]any
	err = json.Unmarshal(quoteBody, &quoteRespData)
	if err != nil {
		fmt.Println("Error parsing quote response:", err)
		return
	}

	if data, ok := quoteRespData["data"].(map[string]any); ok {
		fmt.Printf("You will receive approximately %.6f USDT.\n", data["outputAmount"])
	} else {
		fmt.Println("Unexpected quote response format:", string(quoteBody))
	}
	// Next: ask user to enter wallet address to agree, advise the network to use is LISK
}
