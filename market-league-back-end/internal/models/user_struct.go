package models

import (
	"time"
)

// User struct with auto-incrementing ID, many-to-many relationship, and timestamps
type User struct {
	ID        uint      `gorm:"primaryKey;autoIncrement"`               // Auto-incrementing primary key
	Username  string    `json:"username"`                               // Username
	Email     string    `json:"email"`                                  // User email
	Password  string    `gorm:"not null"`                               // Store hashed password (not plaintext)
	Leagues   []League  `json:"leagues" gorm:"many2many:user_leagues;"` // Many-to-many relation with Leagues
	CreatedAt time.Time `gorm:"autoCreateTime"`                         // Auto-create timestamp
}
