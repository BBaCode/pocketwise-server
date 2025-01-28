package models

type StoredBudget struct {
	ID             string  `json:"id"`
	UserId         string  `json:"user_id"`
	Year           int     `json:"year"`
	Month          int     `json:"month"`
	Total          float64 `json:"total"`
	CreatedAt      string  `json:"created_at"`
	LastUpdated    string  `json:"last_updated"`
	Food           float64 `json:"food"`
	Groceries      float64 `json:"groceries"`
	Transportation float64 `json:"transportation"`
	Entertainment  float64 `json:"entertainment"`
	Health         float64 `json:"health"`
	Shopping       float64 `json:"shopping"`
	Utilities      float64 `json:"utilities"`
	Housing        float64 `json:"housing"`
	Travel         float64 `json:"travel"`
	Education      float64 `json:"education"`
	Subscriptions  float64 `json:"subscriptions"`
	Gifts          float64 `json:"gifts"`
	Insurance      float64 `json:"insurance"`
	PersonalCare   float64 `json:"personal_care"`
	Other          float64 `json:"other"`
	Unknown        float64 `json:"unknown"`
}

// This is for new budgets, not for updating existing budgets
type BudgetRequest struct {
	UserId         string  `json:"user_id"`
	Year           int     `json:"year"`
	Month          int     `json:"month"`
	Total          float64 `json:"total"`
	Food           float64 `json:"food"`
	Groceries      float64 `json:"groceries"`
	Transportation float64 `json:"transportation"`
	Entertainment  float64 `json:"entertainment"`
	Health         float64 `json:"health"`
	Shopping       float64 `json:"shopping"`
	Utilities      float64 `json:"utilities"`
	Housing        float64 `json:"housing"`
	Travel         float64 `json:"travel"`
	Education      float64 `json:"education"`
	Subscriptions  float64 `json:"subscriptions"`
	Gifts          float64 `json:"gifts"`
	Insurance      float64 `json:"insurance"`
	PersonalCare   float64 `json:"personal_care"`
	Other          float64 `json:"other"`
	Unknown        float64 `json:"unknown"`
}

type UpdateBudgetRequest struct {
	Id             string  `json:"id"`
	Year           int     `json:"year"`
	Month          int     `json:"month"`
	Total          float64 `json:"total"`
	Food           float64 `json:"food"`
	Groceries      float64 `json:"groceries"`
	Transportation float64 `json:"transportation"`
	Entertainment  float64 `json:"entertainment"`
	Health         float64 `json:"health"`
	Shopping       float64 `json:"shopping"`
	Utilities      float64 `json:"utilities"`
	Housing        float64 `json:"housing"`
	Travel         float64 `json:"travel"`
	Education      float64 `json:"education"`
	Subscriptions  float64 `json:"subscriptions"`
	Gifts          float64 `json:"gifts"`
	Insurance      float64 `json:"insurance"`
	PersonalCare   float64 `json:"personal_care"`
	Other          float64 `json:"other"`
	Unknown        float64 `json:"unknown"`
}

type MessageResponse struct {
	Message string `json:"message"`
}
