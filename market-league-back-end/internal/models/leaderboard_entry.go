package models

// LeaderboardEntry represents an entry in the league leaderboard.
type LeaderboardEntry struct {
	Username   string `json:"username"`
	TotalValue int    `json:"total_value"`
}
