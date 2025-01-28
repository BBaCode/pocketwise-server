package db

import (
	"context"
	"errors"
	"fmt"
	"log"
	"time"

	"github.com/BBaCode/pocketwise-server/models"
	"github.com/google/uuid"
	"github.com/jackc/pgx/v5"
	"github.com/jackc/pgx/v5/pgxpool"
	"github.com/joho/godotenv"
)

func FetchExistingAccounts(userId uuid.UUID, pool *pgxpool.Pool) ([]models.StoredAccount, error) {

	logger := log.Default()
	// Load configuration (you can expand this later)
	err := godotenv.Load("../../.env")
	if err != nil {
		log.Fatalf("Error loading .env file: %v", err)
	}

	query := `SELECT * FROM public.accounts WHERE user_id = $1`
	rows, err := pool.Query(context.Background(), query, userId)
	if err != nil {
		log.Fatalf("Failed to get accounts: %s", err)
	}
	logger.Println(rows)
	accounts := []models.StoredAccount{}
	rowCount := 0

	// get all accounts and map them to the transaction map
	for rows.Next() {
		// storing these values as numeric in the database, even though they return from the simplefin as strings
		// this does a conversion to allow us to store them as strings again after getting back from the db
		// Maybe update the DB instead?
		var (
			balance          float64
			availableBalance float64
		)
		var acc models.StoredAccount
		err := rows.Scan(&acc.ID, &acc.UserId, &acc.Name, &acc.AccountType, &acc.Currency, &balance, &availableBalance, &acc.Org.Name, &acc.BalanceDate)
		if err != nil {
			return nil, err
		}
		acc.Balance = fmt.Sprintf("%.2f", balance)
		acc.AvailableBalance = fmt.Sprintf("%.2f", availableBalance)
		accounts = append(accounts, acc)
		rowCount++
	}
	log.Printf("Number of accounts fetched: %d", rowCount)

	return accounts, nil

}

func InsertNewAccounts(account models.StoredAccount, pool *pgxpool.Pool) error {

	err := godotenv.Load("../../.env")
	if err != nil {
		log.Fatalf("Error loading .env file: %v", err)
	}

	query := `INSERT INTO public.accounts (id, user_id, name, account_type, currency, balance, available_balance, org_name, balance_date) VALUES ($1, $2, $3, $4, $5, $6, $7, $8, $9)`
	result, err := pool.Exec(context.Background(), query, account.ID, account.UserId, account.Name, account.AccountType, account.Currency, account.Balance, account.AvailableBalance, account.Org.Name, account.BalanceDate)
	if err != nil {
		log.Printf("Unable to insert new account into database: %v\n", err)
		return err
	}
	rowsAffected := result.RowsAffected()
	log.Printf("Rows affected: %d\n", rowsAffected)

	if rowsAffected == 0 {
		log.Println("No rows were inserted. Check your query or data.")
	}

	return nil
}

func UpdateExistingAccounts(account models.UpdatedAccountData, pool *pgxpool.Pool) error {

	err := godotenv.Load("../../.env")
	if err != nil {
		log.Fatalf("Error loading .env file: %v", err)
	}
	query := `UPDATE public.accounts 
          SET balance = $1, available_balance = $2, balance_date = $3
          WHERE id = $4`
	_, err = pool.Exec(context.Background(), query, account.Balance, account.AvailableBalance, account.BalanceDate, account.ID)
	if err != nil {
		log.Printf("Failed to update category for transaction with ID: %s\n", account.ID)
		log.Fatalf("Unable to update transaction in database: %v\n", err)
		return err
	}

	return nil
}

///////////////// TRANSACTIONS //////////////////////

// Fetches every single transaction in the db. Will want to update this to fetch by userId
func FetchAllTransactions(pool *pgxpool.Pool) ([]models.Transaction, error) {

	logger := log.Default()
	// Load configuration (you can expand this later)
	err := godotenv.Load("../../.env")
	if err != nil {
		log.Fatalf("Error loading .env file: %v", err)
	}

	query := `SELECT * FROM public.transactions`
	rows, err := pool.Query(context.Background(), query)
	if err != nil {
		log.Fatalf("Failed to get transactions: %s", err)
	}
	logger.Println(rows)
	transactions := []models.Transaction{}
	rowCount := 0

	// get all transactions and map them to the transaction map
	for rows.Next() {
		rowCount++
		var txn models.Transaction
		err := rows.Scan(&txn.ID, &txn.AccountID, &txn.Amount, &txn.Description, &txn.Payee, &txn.Memo, &txn.Category, &txn.TransactedAt, &txn.Posted)
		if err != nil {
			return nil, err
		}
		transactions = append(transactions, txn) // Use ID as a map key for easy lookups
	}
	log.Printf("Number of transactions fetched: %d", rowCount)

	return transactions, nil

}

func FetchExistingTransactions(accountId string, pool *pgxpool.Pool) ([]models.Transaction, error) {

	logger := log.Default()
	// Load configuration (you can expand this later)
	err := godotenv.Load("../../.env")
	if err != nil {
		log.Fatalf("Error loading .env file: %v", err)
	}

	query := `SELECT * FROM public.transactions WHERE account_id = $1`
	rows, err := pool.Query(context.Background(), query, accountId)
	if err != nil {
		log.Fatalf("Failed to get transactions: %s", err)
	}
	logger.Println(rows)
	transactions := []models.Transaction{}
	rowCount := 0

	// get all transactions and map them to the transaction map
	for rows.Next() {
		rowCount++
		var txn models.Transaction
		err := rows.Scan(&txn.ID, &txn.AccountID, &txn.Amount, &txn.Description, &txn.Payee, &txn.Memo, &txn.Category, &txn.TransactedAt, &txn.Posted)
		if err != nil {
			return nil, err
		}
		transactions = append(transactions, txn) // Use ID as a map key for easy lookups
	}
	log.Printf("Number of transactions fetched: %d", rowCount)

	return transactions, nil

}

func FetchMostRecentTransactionForAllAccounts(pool *pgxpool.Pool) (int64, error) {
	err := godotenv.Load("../../.env")
	if err != nil {
		log.Fatalf("Error loading .env file: %v", err)
	}

	var lastTransactionDate *int64
	err = pool.QueryRow(context.Background(), "SELECT MAX(transacted_at) FROM public.transactions").Scan(&lastTransactionDate)
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

func FetchMostRecentTransactionForAnAccount(accountId string, pool *pgxpool.Pool) (int64, error) {
	err := godotenv.Load("../../.env")
	if err != nil {
		log.Fatalf("Error loading .env file: %v", err)
	}

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

func InsertNewTransactions(txns []models.Transaction, pool *pgxpool.Pool) error {
	err := godotenv.Load("../../.env")
	if err != nil {
		log.Fatalf("Error loading .env file: %v", err)
	}

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

func UpdateTransactionCategory(txns models.UpdatedTransactions, pool *pgxpool.Pool) error {
	err := godotenv.Load("../../.env")
	if err != nil {
		log.Fatalf("Error loading .env file: %v", err)
	}

	for _, txn := range txns.UpdatedTransactions {
		query := `UPDATE public.transactions 
          SET category = $1 
          WHERE id = $2`
		_, err = pool.Exec(context.Background(), query, txn.Category, txn.ID)
		if err != nil {
			log.Printf("Failed to update category for transaction with ID: %s\n", txn.ID)
			log.Fatalf("Unable to update transaction in database: %v\n", err)
			return err
		}

	}
	return nil
}

func FetchCategoryByPayee(txn models.Transaction, pool *pgxpool.Pool) (string, error) {
	// Load environment variables (optional if already loaded elsewhere)
	err := godotenv.Load("../../.env")
	if err != nil {
		return "", fmt.Errorf("error loading .env file: %w", err)
	}

	// Query to check if a similar transaction with the same payee exists
	var category string
	query := `SELECT category FROM public.transactions WHERE payee = $1 LIMIT 1`
	err = pool.QueryRow(context.Background(), query, txn.Payee).Scan(&category)
	if err != nil {
		if errors.Is(err, pgx.ErrNoRows) {
			// If no rows are found, return an empty category (not an error)
			return "", nil
		}
		// Return other database errors
		return "", fmt.Errorf("failed to fetch category: %w", err)
	}

	// Return the category if found
	return category, nil
}

////////////////////////// BUDGET /////////////////////////////////////

func FetchExistingBudget(budgetRequest models.BudgetRequest, pool *pgxpool.Pool) (models.StoredBudget, error) {
	err := godotenv.Load("../../.env")
	if err != nil {
		return models.StoredBudget{}, err
	}
	var budget models.StoredBudget
	query := `SELECT * FROM public.budgets WHERE user_id = $1 AND year = $2 AND month = $3`
	err = pool.QueryRow(context.Background(), query, budgetRequest.UserId, budgetRequest.Year, budgetRequest.Month).Scan(&budget.ID, &budget.UserId, &budget.Year, &budget.Month, &budget.Total, &budget.Food, &budget.Groceries, &budget.Transportation, &budget.Entertainment, &budget.Health, &budget.Shopping, &budget.Utilities, &budget.Housing, &budget.Travel, &budget.Education, &budget.Subscriptions, &budget.Gifts, &budget.Insurance, &budget.PersonalCare, &budget.Other, &budget.Unknown, &budget.CreatedAt, &budget.LastUpdated)
	if err != nil {
		return models.StoredBudget{}, err
	}
	return budget, nil
}

func FetchAllExistingBudgets(userId string, pool *pgxpool.Pool) ([]models.StoredBudget, error) {
	err := godotenv.Load("../../.env")
	if err != nil {
		return []models.StoredBudget{}, err
	}
	userUUID, err := uuid.Parse(userId)
	if err != nil {
		log.Printf("Invalid user ID: %v\n", err)
		return []models.StoredBudget{}, err
	}

	var budgets []models.StoredBudget
	query := `SELECT * FROM public.budgets WHERE user_id = $1`
	rows, err := pool.Query(context.Background(), query, userUUID)
	if err != nil {
		return []models.StoredBudget{}, err
	}

	for rows.Next() {
		var budget models.StoredBudget
		err := rows.Scan(&budget.ID, &budget.UserId, &budget.Year, &budget.Month, &budget.Total, &budget.Food, &budget.Groceries, &budget.Transportation, &budget.Entertainment, &budget.Health, &budget.Shopping, &budget.Utilities, &budget.Housing, &budget.Travel, &budget.Education, &budget.Subscriptions, &budget.Gifts, &budget.Insurance, &budget.PersonalCare, &budget.Other, &budget.Unknown, &budget.CreatedAt, &budget.LastUpdated)
		if err != nil {
			return nil, err
		}
		budgets = append(budgets, budget) // Use ID as a map key for easy lookups
	}
	return budgets, nil

}

func InsertNewBudget(budgetRequest models.BudgetRequest, pool *pgxpool.Pool) error {
	err := godotenv.Load("../../.env")
	if err != nil {
		return err
	}
	userUUID, err := uuid.Parse(budgetRequest.UserId)
	if err != nil {
		log.Printf("Invalid user ID: %v\n", err)
		return err
	}

	query := `INSERT INTO public.budgets (user_id, year, month, total, food, groceries, transportation, entertainment, health, shopping, utilities, housing, travel, education, subscriptions, gifts, insurance, personal_care, other, unknown) VALUES ($1, $2, $3, $4, $5, $6, $7, $8, $9, $10, $11, $12, $13, $14, $15, $16, $17, $18, $19, $20)`
	_, err = pool.Exec(context.Background(), query, userUUID, budgetRequest.Year, budgetRequest.Month, budgetRequest.Total, budgetRequest.Food, budgetRequest.Groceries, budgetRequest.Transportation, budgetRequest.Entertainment, budgetRequest.Health, budgetRequest.Shopping, budgetRequest.Utilities, budgetRequest.Housing, budgetRequest.Travel, budgetRequest.Education, budgetRequest.Subscriptions, budgetRequest.Gifts, budgetRequest.Insurance, budgetRequest.PersonalCare, budgetRequest.Other, budgetRequest.Unknown)
	if err != nil {
		return err
	}
	return nil
}

func DeleteBudget(budgetId string, pool *pgxpool.Pool) error {
	err := godotenv.Load("../../.env")
	if err != nil {
		return err
	}

	query := `DELETE FROM public.budgets WHERE id = $1 `
	_, err = pool.Exec(context.Background(), query, budgetId)
	if err != nil {
		return err
	}
	return nil
}

func UpdateExistingBudget(budgetId string, updateBudgetRequest models.UpdateBudgetRequest, pool *pgxpool.Pool) error {
	err := godotenv.Load("../../.env")
	if err != nil {
		return err
	}

	query := `UPDATE public.budgets 
          SET year = $1, month = $2, total = $3, food = $4, groceries = $5, transportation = $6, entertainment = $7, health = $8, shopping = $9, utilities = $10, housing = $11, travel = $12, education = $13, subscriptions = $14, gifts = $15, insurance = $16, personal_care = $17, other = $18, unknown = $19
          WHERE id = $20`
	_, err = pool.Exec(context.Background(), query, updateBudgetRequest.Year, updateBudgetRequest.Month, updateBudgetRequest.Total, updateBudgetRequest.Food, updateBudgetRequest.Groceries, updateBudgetRequest.Transportation, updateBudgetRequest.Entertainment, updateBudgetRequest.Health, updateBudgetRequest.Shopping, updateBudgetRequest.Utilities, updateBudgetRequest.Housing, updateBudgetRequest.Travel, updateBudgetRequest.Education, updateBudgetRequest.Subscriptions, updateBudgetRequest.Gifts, updateBudgetRequest.Insurance, updateBudgetRequest.PersonalCare, updateBudgetRequest.Other, updateBudgetRequest.Unknown, budgetId)
	if err != nil {
		return err
	}
	return nil
}
