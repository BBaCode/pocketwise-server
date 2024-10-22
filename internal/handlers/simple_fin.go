package handlers

import (
	"encoding/json"
	"fmt"
	"io"
	"log"
	"net/http"
	"os"
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
	req, err := http.NewRequest("GET", fmt.Sprintf("%s/accounts", "https://beta-bridge.simplefin.org/simplefin"), nil)
	if err != nil {
		log.Fatalf("Failed to create request: %v", err)
	}

	req.SetBasicAuth(os.Getenv("SIMPLE_FIN_USERNAME"), os.Getenv("SIMPLE_FIN_PASSWORD")) // Split username and password

	client := &http.Client{}
	resp, err := client.Do(req)
	if err != nil {
		log.Fatalf("Failed to get accounts: %v", err)
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
}
