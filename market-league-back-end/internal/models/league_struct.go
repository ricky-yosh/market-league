package models

import "time"

type League struct {
	ID            uint           `json:"id" gorm:"primaryKey;autoIncrement"`
	LeagueName    string         `json:"league_name"`
	StartDate     time.Time      `json:"start_date"`
	EndDate       time.Time      `json:"end_date"`
	LeagueState   LeagueState    `json:"league_state" gorm:"type:varchar(20);default:'pre_draft'"`
	Users         []User         `json:"users" gorm:"many2many:user_leagues;"` // Many-to-many Users <-> Leagues
	MaxPlayers    *int           `json:"max_players"`
	LeaguePlayers []LeaguePlayer `json:"players" gorm:"foreignKey:LeagueID"`
}
