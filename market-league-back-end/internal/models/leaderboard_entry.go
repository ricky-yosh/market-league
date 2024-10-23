package models

// LeaderboardEntry represents an entry in the league leaderboard.
type LeaderboardEntry struct {
	Username   string  `json:"username"`
	TotalValue float64 `json:"total_value"`
}
