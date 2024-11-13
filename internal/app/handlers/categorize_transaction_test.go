package handlers

import (
	"fmt"
	"testing"
	"time"

	"github.com/BBaCode/pocketwise-server/models"
	"github.com/joho/godotenv"
)

func TestCategorizeTransaction(t *testing.T) {
	// Load environment variables (optional if already done globally in main)
	err := godotenv.Load("../../.env")
	if err != nil {
		t.Fatalf("Error loading .env file: %v", err)
	}

	// Define a sample transaction
	transaction := models.Transaction{
		ID:           "txn_123456",
		Posted:       int(time.Now().Unix()),
		Amount:       "45.67",
		Description:  "Dinner at Olive Garden",
		Payee:        "Olive Garden",
		Memo:         "Family dinner",
		TransactedAt: int(time.Now().Add(-48 * time.Hour).Unix()), // 2 days ago
		Category:     "",
	}

	// Call CategorizeTransaction
	categorizedTransaction, err := CategorizeTransaction(transaction)
	if err != nil {
		t.Errorf("Expected no error, got %v", err)
	}

	// Print or verify the result for debugging
	fmt.Printf("Categorized Transaction: %+v\n", categorizedTransaction)

	// Add assertion to check if categorization makes sense
	if categorizedTransaction.Category == "" {
		t.Errorf("Expected category, got empty string")
	}
}
