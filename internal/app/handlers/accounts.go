package handlers

import (
	"encoding/json"
	"fmt"
	"io"
	"log"
	"net/http"
	"os"
	"time"

	"github.com/BBaCode/pocketwise-server/internal/app"
	"github.com/BBaCode/pocketwise-server/internal/app/constants"
	"github.com/BBaCode/pocketwise-server/internal/db"
	"github.com/BBaCode/pocketwise-server/models"
	"github.com/google/uuid"
	"github.com/jackc/pgx/v5/pgxpool"
)

// TODO: I suspect that this function is no longer needed because everything flows
// through transactions. The "HandleAddAccounts" gets accounts into DB and thats
// then used for all of the transactions
func HandleGetAccounts(w http.ResponseWriter, r *http.Request, pool *pgxpool.Pool) {
	allowedOrigins := constants.AllowedOrigins

	origin := r.Header.Get("Origin")
	if allowedOrigins[origin] {
		w.Header().Set("Access-Control-Allow-Origin", origin)
	}
	w.Header().Set("Access-Control-Allow-Methods", "GET, POST, OPTIONS")
	w.Header().Set("Access-Control-Allow-Headers", "Content-Type, Authorization")

	// Handle preflight OPTIONS request
	if r.Method == "OPTIONS" {
		w.WriteHeader(http.StatusOK)
		return
	}

	userID := r.Header.Get("X-User-ID")
	if userID == "" {
		http.Error(w, "Unauthorized", http.StatusUnauthorized)
		return
	}

	userUUID, err := uuid.Parse(userID)
	if err != nil {
		log.Printf("Invalid user ID: %v\n", err)
		http.Error(w, "Invalid user ID", http.StatusBadRequest)
		return
	}

	// CONTINUE HERE
	var fetchedAccounts []models.StoredAccount
	fetchedAccounts, err = db.FetchExistingAccounts(userUUID, pool)
	if err != nil {
		fmt.Printf("Unexpected error fetching existing accounts: %v", err)
	}

	// Send JSON response to the client
	w.Header().Set("Content-Type", "application/json")
	if err := json.NewEncoder(w).Encode(fetchedAccounts); err != nil {
		http.Error(w, "Failed to send accounts response", http.StatusInternalServerError)
	}
}

// This is really only used as an initial setup (need to add accounts to simplefin first)
// If a user needs to add new accounts, they need to go to simplefin first I believe,
// but maybe there is an API that can do it programmaticly
func HandleAddAccounts(w http.ResponseWriter, r *http.Request, pool *pgxpool.Pool) {
	allowedOrigins := constants.AllowedOrigins

	origin := r.Header.Get("Origin")
	if allowedOrigins[origin] {
		w.Header().Set("Access-Control-Allow-Origin", origin)
	}
	w.Header().Set("Access-Control-Allow-Methods", "GET, POST, OPTIONS")
	w.Header().Set("Access-Control-Allow-Headers", "Content-Type, Authorization")

	// Extract user ID from request header (set by middleware)
	userID := r.Header.Get("X-User-ID")
	if userID == "" {
		http.Error(w, "Unauthorized", http.StatusUnauthorized)
		return
	}

	userUUID, err := uuid.Parse(userID)
	if err != nil {
		log.Printf("Invalid user ID: %v\n", err)
		http.Error(w, "Invalid user ID", http.StatusBadRequest)
		return
	}

	req, err := http.NewRequest("GET", "https://beta-bridge.simplefin.org/simplefin/accounts", nil)
	if err != nil {
		log.Fatalf("Failed to create request: %v", err)
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

	// account response parsed from the body
	var accountsResponse models.AccountResponse
	if err := json.NewDecoder(resp.Body).Decode(&accountsResponse); err != nil {
		log.Fatalf("Failed to decode accounts response: %v", err)
	}

	for _, account := range accountsResponse.Accounts {
		// Assume each account has a list of transactions (you may need to retrieve this separately if not included)
		var storedAccount models.StoredAccount
		storedAccount.AccountType = "General" // hardcoded but can probably create some function to identify it (checking/savings/investment)
		storedAccount.AvailableBalance = account.AvailableBalance
		storedAccount.Balance = account.Balance
		storedAccount.BalanceDate = account.BalanceDate
		storedAccount.Currency = account.Currency
		storedAccount.ID = account.ID
		storedAccount.Org.Name = account.Org.Name
		storedAccount.Name = account.Name
		storedAccount.UserId = userUUID
		db.InsertNewAccounts(storedAccount, pool)
	}

	// Send JSON response to the client
	w.Header().Set("Content-Type", "application/json")
	if err := json.NewEncoder(w).Encode(accountsResponse); err != nil {
		http.Error(w, "Failed to send accounts response", http.StatusInternalServerError)
	}
}

func UpdateAccountData(pool *pgxpool.Pool) error {
	fmt.Println("Updating account data for all accounts and transactions at", time.Now())

	startDate, err := db.FetchMostRecentTransactionForAllAccounts(pool)
	if err != nil {
		return fmt.Errorf("failed to get most recent transaction: %w", err)
	}

	req, err := http.NewRequest("GET", fmt.Sprintf("https://beta-bridge.simplefin.org/simplefin/accounts?start-date=%d", startDate), nil)
	if err != nil {
		return fmt.Errorf("failed to create request: %w", err)
	}

	req.SetBasicAuth(os.Getenv("SIMPLE_FIN_USERNAME"), os.Getenv("SIMPLE_FIN_PASSWORD"))

	client := &http.Client{}
	resp, err := client.Do(req)
	if err != nil {
		return fmt.Errorf("failed to get accounts: %w", err)
	}
	defer resp.Body.Close()

	if resp.StatusCode != http.StatusOK {
		body, _ := io.ReadAll(resp.Body)
		return fmt.Errorf("failed to get accounts: %s", body)
	}

	var accountsResponse models.AccountResponse
	if err := json.NewDecoder(resp.Body).Decode(&accountsResponse); err != nil {
		return fmt.Errorf("failed to decode accounts response: %w", err)
	}

	for _, account := range accountsResponse.Accounts {
		updatedAccountData := models.UpdatedAccountData{
			ID:               account.ID,
			AvailableBalance: account.AvailableBalance,
			Balance:          account.Balance,
			BalanceDate:      account.BalanceDate,
		}

		db.UpdateExistingAccounts(updatedAccountData, pool)

		var categorizedTxns []models.Transaction
		for _, txn := range account.Transactions {
			category, err := db.FetchCategoryByPayee(txn, pool)
			if err != nil {
				return fmt.Errorf("failed to fetch category: %w", err)
			}
			if len(category) > 0 {
				txn.Category = category
			} else {
				txn, err = app.CategorizeTransaction(&txn)
				if err != nil {
					return fmt.Errorf("failed to categorize transaction: %w", err)
				}
			}
			txn.AccountID = account.ID
			categorizedTxns = append(categorizedTxns, txn)
		}

		if err := db.InsertNewTransactions(categorizedTxns, pool); err != nil {
			return fmt.Errorf("failed to insert transactions: %w", err)
		}
	}

	return nil
}

// HTTP Handler Function
func HandleGetUpdatedAccountData(w http.ResponseWriter, r *http.Request, pool *pgxpool.Pool) {
	if err := UpdateAccountData(pool); err != nil {
		http.Error(w, "Failed to update account data", http.StatusInternalServerError)
		log.Println(err)
		return
	}

	w.WriteHeader(http.StatusOK)
	w.Write([]byte("Account data updated successfully"))
}
