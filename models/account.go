package models

import (
	"encoding/json"
	"fmt"
	"time"
)

// CustomTime handles Unix timestamp parsing
type CustomTime time.Time

// UnmarshalJSON for CustomTime to handle Unix timestamps
func (ct *CustomTime) UnmarshalJSON(b []byte) error {
	// Parse Unix timestamp
	var timestamp int64
	if err := json.Unmarshal(b, &timestamp); err != nil {
		return fmt.Errorf("failed to parse timestamp: %v", err)
	}
	*ct = CustomTime(time.Unix(timestamp, 0))
	return nil
}

// MarshalJSON to ensure CustomTime formats correctly if re-encoded
func (ct CustomTime) MarshalJSON() ([]byte, error) {
	return json.Marshal(time.Time(ct).Unix())
}

type Transaction struct {
	ID           string     `json:"id"`
	AccountID    string     `json:"account_id"`
	Posted       CustomTime `json:"posted"`
	Amount       string     `json:"amount"`
	Description  string     `json:"description"`
	Payee        string     `json:"payee"`
	Memo         string     `json:"memo"`
	TransactedAt CustomTime `json:"transacted_at"`
	Category     string     `json:"category"`
}

// type CategorizedTransaction struct {
// 	Transaction Transaction `json:"transaction"`
// 	Category    string      `json:"category"`
// }

// Account represents the structure of account data
type Account struct {
	ID               string        `json:"id"`
	Name             string        `json:"name"`
	Balance          string        `json:"balance"`
	BalanceDate      CustomTime    `json:"balance-date"`
	Currency         string        `json:"currency"`
	AvailableBalance string        `json:"available-balance"`
	Org              Org           `json:"org"`
	Transactions     []Transaction `json:"transactions"`
}

// Response represents the structure of the accounts response
type AccountResponse struct {
	Accounts []Account `json:"accounts"`
}

type Org struct {
	Name string `json:"name"`
}

type StoredAccount struct {
	ID               string     `json:"id"`
	UserId           int        `json:"user_id"`
	Name             string     `json:"name"`
	AccountType      string     `json:"account_type"`
	Currency         string     `json:"currency"`
	Balance          string     `json:"balance"`
	BalanceDate      CustomTime `json:"balance-date"`
	AvailableBalance string     `json:"available-balance"`
	Org              Org        `json:"org"`
}
