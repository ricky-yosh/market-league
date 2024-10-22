package models

import "time"

// League represents the Leagues table in the database.
type League struct {
	ID         int       `json:"id" gorm:"primaryKey;autoIncrement"` // Auto-incrementing primary key
	LeagueName string    `json:"league_name"`
	StartDate  time.Time `json:"start_date"`
	EndDate    time.Time `json:"end_date"`
	Users      []User    `json:"users" gorm:"many2many:user_leagues;"` // Many-to-many relation with Users
}
