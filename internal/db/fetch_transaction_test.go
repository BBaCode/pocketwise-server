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
	categorizedTransactions, err := FetchExistingTransactions("ACT-17dbc9ca-ce58-4d16-b4f1-f8edc1dd7364")
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
	mostRecentTransaction, err := FetchMostRecentTransaction("ACT-17dbc9ca-ce58-4d16-b4f1-f8edc1dd7364")
	if err != nil {
		t.Errorf("Expected no error, got %v", err)
	}
	t.Logf("Most Recent Transaction: %+v\n", mostRecentTransaction)
}
