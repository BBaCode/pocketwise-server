package app

import (
	"context"
	"fmt"
	"log"
	"os"

	"github.com/BBaCode/pocketwise-server/models"
	"github.com/joho/godotenv"
	openai "github.com/sashabaranov/go-openai"
)

func CategorizeTransaction(transaction *models.Transaction) (models.Transaction, error) {

	err := godotenv.Load("../../.env")
	if err != nil {
		log.Fatalf("Error loading .env file: %v", err)
	}

	// Create the prompt that helps categorize the transaction
	// prompt := fmt.Sprintf("Categorize the following transaction based on the description: '%s' with an amount of $%.2f", transaction.Description, transaction.Amount)

	client := openai.NewClient(os.Getenv("OPENAI_API_KEY"))

	var categories = []string{
		"Food & Dining", "Groceries", "Transportation", "Entertainment",
		"Health & Wellness", "Shopping", "Utilities", "Rent", "Travel",
		"Education", "Subscriptions", "Gifts & Donations", "Insurance",
		"Personal Care", "Other",
	}

	prompt := fmt.Sprintf(
		"You are a transaction categorizer. Classify each transaction into only one of these categories: %v. If it's unclear, categorize it as 'Other'. Respond with only the category name, without any extra words or punctuation.",
		categories,
	)

	userPrompt := fmt.Sprintf(
		"Transaction: '%s' Amount: $%s", transaction.Description, transaction.Amount,
	)

	// The system level role set is telling the chatgpt bot what to do / what its job is
	// the user level role is the actual prompt that will be acted upon.
	resp, err := client.CreateChatCompletion(
		context.Background(),
		openai.ChatCompletionRequest{
			Model: openai.GPT3Dot5Turbo,
			Messages: []openai.ChatCompletionMessage{
				{
					Role:    openai.ChatMessageRoleSystem,
					Content: prompt,
				},
				{
					Role:    openai.ChatMessageRoleUser,
					Content: userPrompt,
				},
			},
		},
	)

	if err != nil {
		return models.Transaction{}, err
	}
	fmt.Println(resp.Choices[0].Message.Content)

	return models.Transaction{
		ID:           transaction.ID,
		Posted:       transaction.Posted,
		Amount:       transaction.Amount,
		Description:  transaction.Description,
		Payee:        transaction.Payee,
		Memo:         transaction.Memo,
		TransactedAt: transaction.TransactedAt,
		Category:     resp.Choices[0].Message.Content,
	}, nil
}
