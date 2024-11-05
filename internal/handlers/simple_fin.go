package handlers

import (
	"encoding/json"
	"fmt"
	"io"
	"log"
	"net/http"
	"os"

	"github.com/joho/godotenv"
)

// Account represents the structure of account data
type Account struct {
	ID          string `json:"id"`
	Name        string `json:"name"`
	Balance     string `json:"balance"`
	BalanceDate int    `json:"balance-date"`
}

// Response represents the structure of the accounts response
type AccountResponse struct {
	Accounts []Account `json:"accounts"`
}

func HandleGetAccounts(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Access-Control-Allow-Origin", "*")
	w.Header().Set("Access-Control-Allow-Methods", "GET, POST, OPTIONS")
	w.Header().Set("Access-Control-Allow-Headers", "Content-Type")

	req, err := http.NewRequest("GET", fmt.Sprintf("%s/accounts", "https://beta-bridge.simplefin.org/simplefin"), nil)
	if err != nil {
		log.Fatalf("Failed to create request: %v", err)
	}

	err = godotenv.Load("../../.env")
	if err != nil {
		log.Fatalf("Error loading .env file: %v", err)
	}

	req.SetBasicAuth(os.Getenv("SIMPLE_FIN_USERNAME"), os.Getenv("SIMPLE_FIN_PASSWORD")) // Split username and password

	client := &http.Client{}
	resp, err := client.Do(req)
	if err != nil {
		http.Error(w, "Failed to get accounts", http.StatusForbidden)
	}
	defer resp.Body.Close()

	if resp.StatusCode != http.StatusOK {
		body, _ := io.ReadAll(resp.Body)
		log.Fatalf("Failed to get accounts: %s", body)
	}

	// Read and parse the response
	var accountsResponse AccountResponse
	if err := json.NewDecoder(resp.Body).Decode(&accountsResponse); err != nil {
		log.Fatalf("Failed to decode accounts response: %v", err)
	}

	// Print account data
	for _, account := range accountsResponse.Accounts {
		fmt.Printf("Account ID: %s, Name: %s, Balance: %s", account.ID, account.Name, account.Balance)
	}

	// Send JSON response to the client
	w.Header().Set("Content-Type", "application/json")
	if err := json.NewEncoder(w).Encode(accountsResponse); err != nil {
		http.Error(w, "Failed to send accounts response", http.StatusInternalServerError)
	}
}
