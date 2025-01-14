package models

type StoredBudget struct {
	ID          string `json:"id"`
	UserId      string `json:"user_id"`
	Year        int    `json:"year"`
	Month       int    `json:"month"`
	Amount      int    `json:"amount"`
	CreatedAt   string `json:"created_at"`
	LastUpdated string `json:"last_updated"`
}

type BudgetRequest struct {
	UserId string `json:"user_id"`
	Year   int    `json:"year"`
	Month  int    `json:"month"`
	Amount int    `json:"amount"`
}

type MessageResponse struct {
	Message string `json:"message"`
}
