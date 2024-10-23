// internal/models/user_dto.go
package models

import "time"

// UserDTO represents the user data without sensitive information like password.
type SanitizedUser struct {
	ID        uint      `json:"id"`
	Username  string    `json:"username"`
	Email     string    `json:"email"`
	CreatedAt time.Time `json:"created_at"`
}
