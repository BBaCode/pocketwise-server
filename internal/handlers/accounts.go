package handlers

import (
	"encoding/json"
	"fmt"
	"io"
	"log"
	"net/http"
	"os"
	"time"

	"github.com/joho/godotenv"
)

type Transaction struct {
	ID           string `json:"id"`
	Posted       int    `json:"posted"`
	Amount       string `json:"amount"`
	Description  string `json:"description"`
	Payee        string `json:"payee"`
	Memo         string `json:"memo"`
	TransactedAt int    `json:"transacted_at"`
}

// Account represents the structure of account data
type Account struct {
	ID               string        `json:"id"`
	Name             string        `json:"name"`
	Balance          string        `json:"balance"`
	BalanceDate      int           `json:"balance-date"`
	AvailableBalance string        `json:"available-balance"`
	Transactions     []Transaction `json:"transactions"`
}

// Response represents the structure of the accounts response
type AccountResponse struct {
	Accounts []Account `json:"accounts"`
}

func getLast30DaysTimestamp() int64 {
	// Get the current time
	now := time.Now()

	// Subtract 30 days (in hours) from the current time
	last30Days := now.AddDate(0, 0, -30)

	// Return the Unix timestamp in seconds
	return last30Days.Unix()
}

func HandleGetAccounts(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Access-Control-Allow-Origin", "*")
	w.Header().Set("Access-Control-Allow-Methods", "GET, POST, OPTIONS")
	w.Header().Set("Access-Control-Allow-Headers", "Content-Type")

	startDate := getLast30DaysTimestamp()

	req, err := http.NewRequest("GET", fmt.Sprintf("https://beta-bridge.simplefin.org/simplefin/accounts?start-date=%d", startDate), nil)
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
	// for _, account := range accountsResponse.Accounts {
	// 	fmt.Printf("Account ID: %s, Name: %s, Balance: %s", account.ID, account.Name, account.Balance)
	// }

	// Send JSON response to the client
	w.Header().Set("Content-Type", "application/json")
	if err := json.NewEncoder(w).Encode(accountsResponse); err != nil {
		http.Error(w, "Failed to send accounts response", http.StatusInternalServerError)
	}
}
