package models

import (
	"time"
)

type Trade struct {
	ID                 uint       `gorm:"primaryKey;autoIncrement"`
	LeagueID           uint       `json:"league_id"`                                             // ID of the league in which the trade is taking place
	Player1ID          uint       `json:"player1_id"`                                            // ID of the first player (initiator)
	Player2ID          uint       `json:"player2_id"`                                            // ID of the second player (recipient)
	Player1PortfolioID uint       `json:"player1_portfolio_id"`                                  // ID of the first player's portfolio
	Player2PortfolioID uint       `json:"player2_portfolio_id"`                                  // ID of the second player's portfolio
	Player1Stocks      []Stock    `json:"player1_stocks" gorm:"many2many:trade_player1_stocks;"` // Stocks offered by player 1
	Player2Stocks      []Stock    `json:"player2_stocks" gorm:"many2many:trade_player2_stocks;"` // Stocks offered by player 2
	Player1Confirmed   bool       `json:"player1_confirmed"`                                     // Whether Player 1 has confirmed the trade
	Player2Confirmed   bool       `json:"player2_confirmed"`                                     // Whether Player 2 has confirmed the trade
	CreatedAt          time.Time  `json:"created_at" gorm:"autoCreateTime"`                      // Timestamp of when the trade was created
	ConfirmedAt        *time.Time `json:"confirmed_at"`                                          // Timestamp of when the trade was confirmed (nullable)
}
