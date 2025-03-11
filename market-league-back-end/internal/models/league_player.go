package models

import "gorm.io/gorm"

type LeaguePlayer struct {
	gorm.Model
	LeagueID    uint        `json:"league_id" gorm:"index;constraint:OnDelete:CASCADE;"` // Explicit foreign key
	PlayerID    uint        `json:"player_id"`
	DraftStatus DraftStatus `gorm:"type:varchar(20);default:'not_ready'" json:"draft_status"`
}
