package models

import "github.com/google/uuid"

type Transaction struct {
	ID           string `json:"id"`
	AccountID    string `json:"account_id"`
	Posted       int64  `json:"posted"`
	Amount       string `json:"amount"`
	Description  string `json:"description"`
	Payee        string `json:"payee"`
	Memo         string `json:"memo"`
	TransactedAt int64  `json:"transacted_at"`
	Category     string `json:"category"`
}

type TransactionCategoryRequest struct {
	ID       string `json:"id"`
	Category string `json:"category"`
}

type UpdatedTransactions struct {
	UpdatedTransactions []TransactionCategoryRequest `json:"updatedTxns"`
}

// Account represents the structure of account data
type Account struct {
	ID               string        `json:"id"`
	Name             string        `json:"name"`
	Balance          string        `json:"balance"`
	BalanceDate      int64         `json:"balance-date"`
	Currency         string        `json:"currency"`
	AvailableBalance string        `json:"available-balance"`
	Org              Org           `json:"org"`
	Transactions     []Transaction `json:"transactions"`
}

// Request represents the structure of the account request (when calling for transactions)
type AccountRequest struct {
	Account string `json:"account"`
}

// Response represents the structure of the accounts response
type AccountResponse struct {
	Accounts []Account `json:"accounts"`
}

type StoredAccount struct {
	ID               string    `json:"id"`
	UserId           uuid.UUID `json:"user_id"`
	Name             string    `json:"name"`
	AccountType      string    `json:"account_type"`
	Currency         string    `json:"currency"`
	Balance          string    `json:"balance"`
	AvailableBalance string    `json:"available-balance"`
	Org              Org       `json:"org"`
	BalanceDate      int64     `json:"balance-date"`
}

type UpdatedAccountData struct {
	ID               string `json:"id"`
	Balance          string `json:"balance"`
	AvailableBalance string `json:"available-balance"`
	BalanceDate      int64  `json:"balance-date"`
}

type Org struct {
	Name string `json:"name"`
}
