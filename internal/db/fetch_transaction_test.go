package db

import (
	"testing"

	"github.com/joho/godotenv"
)

func TestFetchExistingTransaction(t *testing.T) {
	// Load environment variables (optional if already done globally in main)
	err := godotenv.Load("../../.env")
	if err != nil {
		t.Fatalf("Error loading .env file: %v", err)
	}

	// Call FetchExistingTransaction
	categorizedTransactions, err := FetchExistingTransactions()
	if err != nil {
		t.Errorf("Expected no error, got %v", err)
	}

	if len(categorizedTransactions) == 0 {
		t.Log("0 Transactions found")
	} else {
		t.Logf("Map of Transactions: %+v\n", categorizedTransactions)
	}
}

func TestFetchMostRecentTransaction(t *testing.T) {
	// Load environment variables (optional if already done globally in main)
	err := godotenv.Load("../../.env")
	if err != nil {
		t.Fatalf("Error loading .env file: %v", err)
	}

	// Call FetchMostRecentTransaction
	mostRecentTransaction, err := FetchMostRecentTransaction()
	if err != nil {
		t.Errorf("Expected no error, got %v", err)
	}
	t.Logf("Most Recent Transaction: %+v\n", mostRecentTransaction)
}
