package models

type Transaction struct {
	ID           string `json:"id"`
	Posted       int    `json:"posted"`
	Amount       string `json:"amount"`
	Description  string `json:"description"`
	Payee        string `json:"payee"`
	Memo         string `json:"memo"`
	TransactedAt int    `json:"transacted_at"`
	Category     string `json:"category"`
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
	BalanceDate      int           `json:"balance-date"`
	AvailableBalance string        `json:"available-balance"`
	Transactions     []Transaction `json:"transactions"`
}

// Response represents the structure of the accounts response
type AccountResponse struct {
	Accounts []Account `json:"accounts"`
}
