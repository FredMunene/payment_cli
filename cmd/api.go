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

type QuoteRequest struct {
	Amount int    `json:"amount"`
	Wallet string `json:"wallet"`
}

type QuoteResponse struct {
	OutputAmount float64 `json:"outputAmount"`
	Message      string  `json:"message"`
}

func quoteHandler(w http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodPost {
		w.WriteHeader(http.StatusMethodNotAllowed)
		fmt.Fprint(w, "Method not allowed")
		return
	}

	var req QuoteRequest
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		w.WriteHeader(http.StatusBadRequest)
		fmt.Fprint(w, "Invalid request body")
		return
	}
	if req.Amount < 50 || req.Wallet == "" {
		w.WriteHeader(http.StatusBadRequest)
		fmt.Fprint(w, "Amount must be at least 50 and wallet must not be empty")
		return
	}

	apiKey := os.Getenv("SWYPT_API_KEY")
	apiSecret := os.Getenv("SWYPT_API_SECRET")
	if apiKey == "" || apiSecret == "" {
		w.WriteHeader(http.StatusInternalServerError)
		fmt.Fprint(w, "API credentials not set")
		return
	}

	client := &http.Client{}
	quotePayload := map[string]any{
		"type":           "onramp",
		"amount":         fmt.Sprintf("%d", req.Amount),
		"fiatCurrency":   "KES",
		"cryptoCurrency": "USDT",
		"network":        "lisk",
	}
	payloadBytes, err := json.Marshal(quotePayload)
	if err != nil {
		w.WriteHeader(http.StatusInternalServerError)
		fmt.Fprint(w, "Error encoding payload")
		return
	}

	quoteReq, err := http.NewRequest("POST", "https://pool.swypt.io/api/swypt-quotes", io.NopCloser(bytes.NewReader(payloadBytes)))
	if err != nil {
		w.WriteHeader(http.StatusInternalServerError)
		fmt.Fprint(w, "Error creating request")
		return
	}
	quoteReq.Header.Add("x-api-key", apiKey)
	quoteReq.Header.Add("x-api-secret", apiSecret)
	quoteReq.Header.Add("Content-Type", "application/json")

	quoteResp, err := client.Do(quoteReq)
	if err != nil {
		w.WriteHeader(http.StatusInternalServerError)
		fmt.Fprint(w, "Error making quote request")
		return
	}
	defer quoteResp.Body.Close()

	quoteBody, err := io.ReadAll(quoteResp.Body)
	if err != nil {
		w.WriteHeader(http.StatusInternalServerError)
		fmt.Fprint(w, "Error reading quote response")
		return
	}

	var quoteRespData map[string]any
	err = json.Unmarshal(quoteBody, &quoteRespData)
	if err != nil {
		w.WriteHeader(http.StatusInternalServerError)
		fmt.Fprint(w, "Error parsing quote response")
		return
	}

	if data, ok := quoteRespData["data"].(map[string]any); ok {
		outputAmount, _ := data["outputAmount"].(float64)
		resp := QuoteResponse{
			OutputAmount: outputAmount,
			Message:      "Quote retrieved successfully",
		}
		w.Header().Set("Content-Type", "application/json")
		json.NewEncoder(w).Encode(resp)
	} else {
		w.WriteHeader(http.StatusInternalServerError)
		fmt.Fprintf(w, "Unexpected quote response format: %s", string(quoteBody))
	}
}

func main() {
	internal.Init()
	port := os.Getenv("PORT")
	if port == "" {
		port = "8080"
	}
	http.HandleFunc("/api/quote", quoteHandler)
	log.Printf("API server running on :%s...", port)
	log.Fatal(http.ListenAndServe(":"+port, nil))
}
