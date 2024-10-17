package user

import (
	"time"
)

type User struct {
	ID        uint      `gorm:"primaryKey;autoIncrement"` // ID will auto-increment
	Name      string    `gorm:"size:100;not null"`        // User's name, with a limit of 100 characters, can't be null
	Email     string    `gorm:"size:100;unique;not null"` // Unique email, can't be null
	Password  string    `gorm:"not null"`                 // Password (hash should be stored, not plaintext)
	CreatedAt time.Time `gorm:"autoCreateTime"`           // Created timestamp, auto-generated
	UpdatedAt time.Time `gorm:"autoUpdateTime"`           // Updated timestamp, auto-generated
}
