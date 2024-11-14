package app

import (
	"fmt"
	"testing"

	"github.com/BBaCode/pocketwise-server/models"
	"github.com/joho/godotenv"
)

// TODO: Fix the time thing
func TestCategorizeTransaction(t *testing.T) {
	// Load environment variables (optional if already done globally in main)
	err := godotenv.Load("../../.env")
	if err != nil {
		t.Fatalf("Error loading .env file: %v", err)
	}

	// Define a sample transaction
	transaction := models.Transaction{
		ID:           "txn_123456",
		Posted:       123451251,
		Amount:       "45.67",
		Description:  "Dinner at Olive Garden",
		Payee:        "Olive Garden",
		Memo:         "Family dinner",
		TransactedAt: 123451251, // 2 days ago
		Category:     "",
	}

	// Call CategorizeTransaction
	categorizedTransaction, err := CategorizeTransaction(&transaction)
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
