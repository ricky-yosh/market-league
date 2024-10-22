package models

import (
	"time"
)

type Trade struct {
	ID          uint      `gorm:"primaryKey;autoIncrement"`
	StockID     uint      `json:"stock_id"`               // Foreign key to Stock
	Stock       Stock     `gorm:"foreignKey:StockID"`     // Association with Stock
	UserID      uint      `json:"user_id"`                // Foreign key to User
	User        User      `gorm:"foreignKey:UserID"`      // Association with User
	PortfolioID uint      `json:"portfolio_id"`           // Foreign key to Portfolio
	Portfolio   Portfolio `gorm:"foreignKey:PortfolioID"` // Association with Portfolio
	LeagueID    uint      `json:"league_id"`              // Foreign key to League
	League      League    `gorm:"foreignKey:LeagueID"`    // Association with League
	TradeType   string    `json:"trade_type"`             // Buy or Sell
	TradePrice  float64   `json:"trade_price"`            // Price per stock during trade
	TradeDate   time.Time `json:"trade_date" gorm:"autoCreateTime"`
}
