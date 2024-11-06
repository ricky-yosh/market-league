package models

import (
	"time"
)

type Trade struct {
	ID             uint      `json:"id" gorm:"primaryKey;autoIncrement"`
	LeagueID       uint      `json:"league_id" gorm:"not null"`
	User1ID        uint      `json:"user1_id" gorm:"not null"`                         // Initiating user
	User2ID        uint      `json:"user2_id" gorm:"not null"`                         // Counterparty user
	Portfolio1ID   uint      `json:"portfolio1_id" gorm:"not null"`                    // Portfolio of User 1
	Portfolio2ID   uint      `json:"portfolio2_id" gorm:"not null"`                    // Portfolio of User 2
	Stocks1        []Stock   `json:"stocks1" gorm:"many2many:trade_stocks1"`           // Stocks User 1 is offering
	Stocks2        []Stock   `json:"stocks2" gorm:"many2many:trade_stocks2"`           // Stocks User 2 is offering
	User1Confirmed bool      `json:"user1_confirmed" gorm:"default:false"`             // Confirmation status of User 1
	User2Confirmed bool      `json:"user2_confirmed" gorm:"default:false"`             // Confirmation status of User 2
	Status         string    `json:"status" gorm:"type:varchar(20);default:'pending'"` // Status of the trade
	CreatedAt      time.Time `json:"created_at" gorm:"autoCreateTime"`                 // Creation timestamp
	UpdatedAt      time.Time `json:"updated_at" gorm:"autoUpdateTime"`                 // Last update timestamp
}
