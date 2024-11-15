package handlers

import (
	"encoding/json"
	"io"
	"log"
	"net/http"
	"os"

	"github.com/BBaCode/pocketwise-server/internal/db"
	"github.com/BBaCode/pocketwise-server/models"
	"github.com/joho/godotenv"
)

func HandleGetAccounts(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Access-Control-Allow-Origin", "*")
	w.Header().Set("Access-Control-Allow-Methods", "GET, POST, OPTIONS")
	w.Header().Set("Access-Control-Allow-Headers", "Content-Type")

	req, err := http.NewRequest("GET", "https://beta-bridge.simplefin.org/simplefin/accounts", nil)
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
	var accountsResponse models.AccountResponse
	if err := json.NewDecoder(resp.Body).Decode(&accountsResponse); err != nil {
		log.Fatalf("Failed to decode accounts response: %v", err)
	}

	// Send JSON response to the client
	w.Header().Set("Content-Type", "application/json")
	if err := json.NewEncoder(w).Encode(accountsResponse); err != nil {
		http.Error(w, "Failed to send accounts response", http.StatusInternalServerError)
	}
}

func HandleAddAccounts(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Access-Control-Allow-Origin", "*")
	w.Header().Set("Access-Control-Allow-Methods", "GET, POST, OPTIONS")
	w.Header().Set("Access-Control-Allow-Headers", "Content-Type")

	req, err := http.NewRequest("GET", "https://beta-bridge.simplefin.org/simplefin/accounts", nil)
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
	var accountsResponse models.AccountResponse
	if err := json.NewDecoder(resp.Body).Decode(&accountsResponse); err != nil {
		log.Fatalf("Failed to decode accounts response: %v", err)
	}

	for _, account := range accountsResponse.Accounts {
		// Assume each account has a list of transactions (you may need to retrieve this separately if not included)
		var storedAccount models.StoredAccount
		storedAccount.AccountType = "General" // hardcoded but can probably create some function to identify it
		storedAccount.AvailableBalance = account.AvailableBalance
		storedAccount.Balance = account.Balance
		storedAccount.BalanceDate = account.BalanceDate
		storedAccount.Currency = account.Currency
		storedAccount.ID = account.ID
		storedAccount.Org.Name = account.Org.Name
		storedAccount.Name = account.Name
		storedAccount.UserId = 39 // hardcoded but can probably take from user login information
		db.InsertNewAccounts(storedAccount)
	}

	// Send JSON response to the client
	w.Header().Set("Content-Type", "application/json")
	if err := json.NewEncoder(w).Encode(accountsResponse); err != nil {
		http.Error(w, "Failed to send accounts response", http.StatusInternalServerError)
	}
}
