package models

type LeaguePlayer struct {
	LeagueID    uint        `json:"league_id" gorm:"primaryKey;constraint:OnDelete:CASCADE;"`
	PlayerID    uint        `json:"player_id" gorm:"primaryKey"`
	DraftStatus DraftStatus `gorm:"type:varchar(20);default:'not_ready'" json:"draft_status"`
}
