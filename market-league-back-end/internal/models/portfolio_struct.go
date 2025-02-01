package models

import "time"

// Portfolio represents a user's portfolio in a specific league.
type Portfolio struct {
	ID        uint      `gorm:"primaryKey;autoIncrement"`    // Auto-incrementing primary key
	UserID    uint      `json:"user_id"`                     // Foreign key to User
	User      User      `gorm:"foreignKey:UserID"`           // Association with User
	LeagueID  uint      `json:"league_id"`                   // Foreign key to League
	League    League    `gorm:"foreignKey:LeagueID"`         // Association with League
	Stocks    []Stock   `gorm:"many2many:portfolio_stocks;"` // Many-to-many relationship with Stocks
	Points    int       `json:"points" gorm:"default:0"`     // Points calculated based on stock performances
	CreatedAt time.Time `gorm:"autoCreateTime"`              // Timestamp of portfolio creation
}
