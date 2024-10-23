package models

// UserPortfolioValue represents a user's portfolio value in a league, used only for response purposes.
type UserPortfolioValue struct {
	UserID     uint    `json:"user_id"`
	Username   string  `json:"username"`
	TotalValue float64 `json:"total_value"`
}
