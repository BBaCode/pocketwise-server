package db

import (
	"context"
	"fmt"
	"log"
	"time"

	config "github.com/BBaCode/pocketwise-server/internal/app"
	"github.com/BBaCode/pocketwise-server/models"
	"github.com/jackc/pgx/v5"
	"github.com/joho/godotenv"
)

func InsertNewAccounts(account models.StoredAccount) error {
	// Load configuration (you can expand this later)
	cfg := config.LoadConfig()
	err := godotenv.Load("../../.env")
	if err != nil {
		log.Fatalf("Error loading .env file: %v", err)
	}
	// Connect to the database
	pool, err := Connect(DBConfig(cfg))
	if err != nil {
		log.Fatalf("Unable to connect to database: %v\n", err)
	}
	defer pool.Close() // Ensure the connection is closed when you're done

	query := `INSERT INTO public.accounts (id, user_id, name, account_type, currency, balance, available_balance, balance_date, org_name) VALUES ($1, $2, $3, $4, $5, $6, $7, $8, $9)`
	_, err = pool.Exec(context.Background(), query, account.ID, account.UserId, account.Name, account.AccountType, account.Currency, account.Balance, account.AvailableBalance, account.BalanceDate, account.Org.Name)
	if err != nil {
		log.Fatalf("Unable to insert new accounts into database: %v\n", err)
		return err
	}

	return nil
}

///////////////// TRANSACTIONS //////////////////////

func FetchExistingTransactions(accountId string) (map[string]models.Transaction, error) {

	logger := log.Default()
	// Load configuration (you can expand this later)
	cfg := config.LoadConfig()
	err := godotenv.Load("../../.env")
	if err != nil {
		log.Fatalf("Error loading .env file: %v", err)
	}

	// Connect to the database
	pool, err := Connect(DBConfig(cfg))
	if err != nil {
		log.Fatalf("Unable to connect to database: %v\n", err)
	}
	defer pool.Close() // Ensure the connection is closed when you're done

	query := `SELECT * FROM public.transactions WHERE account_id = $1`

	rows, err := pool.Query(context.Background(), query, accountId)
	if err != nil {
		log.Fatalf("Failed to get transactions: %s", err)
	}
	logger.Println(rows)
	transactions := make(map[string]models.Transaction)
	rowCount := 0

	// get all transactions and map them to the transaction map
	for rows.Next() {
		rowCount++
		var txn models.Transaction
		err := rows.Scan(&txn.ID, &txn.AccountID, &txn.Amount, &txn.Description, &txn.Payee, &txn.Memo, &txn.Category, &txn.TransactedAt, &txn.Posted)
		if err != nil {
			return nil, err
		}
		transactions[txn.ID] = txn // Use ID as a map key for easy lookups
	}
	log.Printf("Number of transactions fetched: %d", rowCount)

	return transactions, nil

}

func FetchMostRecentTransaction(accountId string) (int64, error) {

	// Load configuration (you can expand this later)
	cfg := config.LoadConfig()
	err := godotenv.Load("../../.env")
	if err != nil {
		log.Fatalf("Error loading .env file: %v", err)
	}
	// Connect to the database
	pool, err := Connect(DBConfig(cfg))
	if err != nil {
		log.Fatalf("Unable to connect to database: %v\n", err)
	}
	defer pool.Close() // Ensure the connection is closed when you're done
	var lastTransactionDate *int64
	err = pool.QueryRow(context.Background(), "SELECT MAX(transacted_at) FROM public.transactions WHERE account_id = $1", accountId).Scan(&lastTransactionDate)
	if err == pgx.ErrNoRows || lastTransactionDate == nil {
		fmt.Println("No transactions found for this account")
		return getLast30DaysTimestamp(), nil
	} else if err != nil {
		log.Fatalf("Failed to get most recent transaction: %s", err)
	}

	// Add 1 second buffer to avoid duplicates
	adjustedStartDate := *lastTransactionDate + 1
	return adjustedStartDate, nil
}

func getLast30DaysTimestamp() int64 {
	// Get the current time
	now := time.Now()

	// Subtract 30 days (in hours) from the current time
	last30Days := now.AddDate(0, 0, -30)

	// Return the Unix timestamp in seconds
	return last30Days.Unix()
}

func InsertNewTransactions(txns []models.Transaction) error {
	// Load configuration (you can expand this later)
	cfg := config.LoadConfig()
	err := godotenv.Load("../../.env")
	if err != nil {
		log.Fatalf("Error loading .env file: %v", err)
	}
	// Connect to the database
	pool, err := Connect(DBConfig(cfg))
	if err != nil {
		log.Fatalf("Unable to connect to database: %v\n", err)
	}
	defer pool.Close() // Ensure the connection is closed when you're done

	for _, txn := range txns {
		// Check if transaction ID already exists
		var exists bool
		err = pool.QueryRow(context.Background(), "SELECT EXISTS(SELECT 1 FROM public.transactions WHERE id = $1)", txn.ID).Scan(&exists)
		if err != nil {
			log.Fatalf("Failed to check existing transaction: %v\n", err)
		}
		if exists {
			continue // Skip this transaction if it already exists
		}

		query := `INSERT INTO public.transactions (id, account_id, posted, amount, description, payee, memo, transacted_at, category) VALUES ($1, $2, $3, $4, $5, $6, $7, $8, $9)`
		_, err = pool.Exec(context.Background(), query, txn.ID, txn.AccountID, txn.Posted, txn.Amount, txn.Description, txn.Payee, txn.Memo, txn.TransactedAt, txn.Category)
		if err != nil {
			log.Printf("Failed to insert transaction with ID: %s, AccountID: %s\n", txn.ID, txn.AccountID)
			log.Fatalf("Unable to insert new transactions into database: %v\n", err)
			return err
		}
	}
	return nil
}
