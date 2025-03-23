package models

import (
	"time"
)

// User struct with auto-incrementing ID, many-to-many relationship, and timestamps
type PortfolioPointsHistory struct {
	ID          uint      `json:"id" gorm:"primaryKey;autoIncrement"` // Auto-incrementing primary key
	PortfolioID uint      `json:"portfolio_id"`                       // Foreign key to Portfolio
	Portfolio   Portfolio `gorm:"foreignKey:PortfolioID"`             // Relationship with Portfolio
	Points      int       `json:"points"`                             // Points value at a specific moment
	RecordedAt  time.Time `json:"recorded_at" gorm:"autoCreateTime"`  // Timestamp of when the points were recorded
}
