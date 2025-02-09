// internal/models/user_dto.go
package models

import "time"

type SanitizedUser struct {
	ID        uint      `json:"id"`
	Username  string    `json:"username"`
	Email     string    `json:"email"`
	CreatedAt time.Time `json:"created_at"`
}
type SanitizedPortfolio struct {
	ID        uint             `json:"id"`
	UserID    uint             `json:"user_id"`
	User      SanitizedUser    `json:"user"`
	LeagueID  uint             `json:"league_id"`
	League    SanitizedLeague  `json:"league"`
	Stocks    []SanitizedStock `json:"stocks"`
	Points    int              `json:"points"`
	CreatedAt time.Time        `json:"created_at"`
}

type SanitizedLeague struct {
	ID         uint      `json:"id"`
	LeagueName string    `json:"league_name"`
	StartDate  time.Time `json:"start_date"`
	EndDate    time.Time `json:"end_date"`
}

type SanitizedStock struct {
	ID           uint    `json:"id"`
	TickerSymbol string  `json:"ticker_symbol"`
	CompanyName  string  `json:"company_name"`
	CurrentPrice float64 `json:"current_price"`
}

type SanitizedTrade struct {
	ID             uint          `json:"id"`
	LeagueID       uint          `json:"league_id"`
	User1          SanitizedUser `json:"user1"`
	User2          SanitizedUser `json:"user2"`
	Portfolio1ID   uint          `json:"portfolio1_id"`
	Portfolio2ID   uint          `json:"portfolio2_id"`
	Stocks1        []Stock       `json:"stocks1"`
	Stocks2        []Stock       `json:"stocks2"`
	User1Confirmed bool          `json:"user1_confirmed"`
	User2Confirmed bool          `json:"user2_confirmed"`
	Status         string        `json:"status"`
	CreatedAt      time.Time     `json:"created_at"`
	UpdatedAt      time.Time     `json:"updated_at"`
}
