package models

import "time"

type LeaguePortfolio struct {
	ID        uint      `gorm:"primaryKey;autoIncrement" json:"id"`
	LeagueID  uint      `json:"league_id"`
	League    League    `gorm:"foreignKey:LeagueID" json:"league"`
	Name      string    `json:"name"`
	Stocks    []Stock   `gorm:"many2many:league_portfolio_stocks;" json:"stocks"`
	CreatedAt time.Time `gorm:"autoCreateTime" json:"created_at"`
}
