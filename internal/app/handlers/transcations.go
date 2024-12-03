package handlers

import (
	"encoding/json"
	"fmt"
	"io"
	"log"
	"net/http"
	"os"

	"github.com/BBaCode/pocketwise-server/internal/app"
	"github.com/BBaCode/pocketwise-server/internal/db"
	"github.com/BBaCode/pocketwise-server/models"
	"github.com/jackc/pgx/v5/pgxpool"
	"github.com/joho/godotenv"
)

func HandleGetTransactions(w http.ResponseWriter, r *http.Request, pool *pgxpool.Pool) {
	w.Header().Set("Access-Control-Allow-Origin", "*")
	w.Header().Set("Access-Control-Allow-Methods", "POST, OPTIONS")
	w.Header().Set("Access-Control-Allow-Headers", "Content-Type, Authorization")
	fmt.Print(r.Method)

	// Extract user ID from request header (set by middleware)
	userID := r.Header.Get("X-User-ID")
	if userID == "" {
		http.Error(w, "Unauthorized", http.StatusUnauthorized)
		return
	}

	// Parse the JSON body
	var reqBody models.AccountRequest
	if err := json.NewDecoder(r.Body).Decode(&reqBody); err != nil {
		fmt.Printf("Bad request: %s", err)
		http.Error(w, "Invalid JSON body", http.StatusBadRequest)
		return
	}

	// Check if the account parameter is present
	account := reqBody.Account
	if account == "" {
		http.Error(w, "Missing 'account' parameter in request body", http.StatusBadRequest)
		return
	}

	startDate, err := db.FetchMostRecentTransaction(reqBody.Account, pool)
	if err != nil {
		http.Error(w, "Something went wrong. Please try again later.", http.StatusInternalServerError)
		log.Fatalf("Failed to get successful response from FetchMostRecentTransaction: %s", err)
	}

	req, err := http.NewRequest("GET", fmt.Sprintf("https://beta-bridge.simplefin.org/simplefin/accounts?start-date=%d&account=%s", startDate, reqBody.Account), nil)
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
		http.Error(w, "Failed to get account", http.StatusForbidden)
	}
	defer resp.Body.Close()

	if resp.StatusCode != http.StatusOK {
		body, _ := io.ReadAll(resp.Body)
		log.Fatalf("Failed to get account: %s", body)
	}

	// Read and parse the response
	var accountsResponse models.AccountResponse
	if err := json.NewDecoder(resp.Body).Decode(&accountsResponse); err != nil {
		log.Fatalf("Failed to decode accounts response: %v", err)
	}

	var accForTxns *models.Account
	for _, account := range accountsResponse.Accounts {
		if account.ID == reqBody.Account {
			accForTxns = &account
			break
		}
	}

	// Check if the account was found
	if accForTxns != nil {
		fmt.Printf("Found accounts: %+v\n", *accForTxns)
	} else {
		fmt.Println("Accounts not found")
	}

	var categorizedTxns []models.Transaction
	// categorize transactions and append them to a new array to send to database
	for _, txn := range accForTxns.Transactions {
		txn, err = app.CategorizeTransaction(&txn)
		txn.AccountID = reqBody.Account
		if err != nil {
			log.Fatalf("Failed to categorize transactions with error: %s", err)
		}
		categorizedTxns = append(categorizedTxns, txn)
	}

	// insert new transactions into the database
	err = db.InsertNewTransactions(categorizedTxns, pool)
	if err != nil {
		log.Fatalf("Failed to insert transactions with error: %s", err)
	}

	updatedTxns, err := db.FetchExistingTransactions(accForTxns.ID, pool)
	if err != nil {
		log.Fatalf("Failed to fetch transactions with error: %s", err)
	}

	// Send JSON response to the client
	w.Header().Set("Content-Type", "application/json")
	if err := json.NewEncoder(w).Encode(updatedTxns); err != nil {
		http.Error(w, "Failed to send accounts response", http.StatusInternalServerError)
	}
}
