package models

import (
	"time"
)

// User struct with auto-incrementing ID, many-to-many relationship, and timestamps
type OwnershipHistory struct {
	ID            uint       `gorm:"primaryKey;autoIncrement"`                // Auto-incrementing primary key
	PortfolioID   uint       `json:"portfolio_id"`                            // Foreign key to Portfolio
	Portfolio     Portfolio  `json:"portfolio" gorm:"foreignKey:PortfolioID"` // Association with Portfolio
	StockID       uint       `json:"stock_id"`                                // Foreign key to Stock
	Stock         Stock      `json:"stock" gorm:"foreignKey:StockID"`         // Association with Stock
	StartingValue float64    `json:"starting_value"`                          // Value of the stock when acquired
	CurrentValue  float64    `json:"current_value"`                           // Current or ending value of the stock
	StartDate     time.Time  `json:"start_date"`                              // Timestamp when the stock was acquired
	EndDate       *time.Time `json:"end_date"`                                // Nullable timestamp for when the stock was sold
}
