package models

import (
	"time"
)

type Trade struct {
	ID                uint      `gorm:"primaryKey;autoIncrement"`         // Auto-incrementing primary key
	SellerPortfolioID uint      `json:"seller_portfolio_id"`              // Foreign key to Seller's Portfolio
	SellerPortfolio   Portfolio `gorm:"foreignKey:SellerPortfolioID"`     // Association with Seller's Portfolio
	BuyerPortfolioID  uint      `json:"buyer_portfolio_id"`               // Foreign key to Buyer's Portfolio
	BuyerPortfolio    Portfolio `gorm:"foreignKey:BuyerPortfolioID"`      // Association with Buyer's Portfolio
	LeagueID          uint      `json:"league_id"`                        // Foreign key to League
	League            League    `gorm:"foreignKey:LeagueID"`              // Association with League
	StocksSeller      []Stock   `gorm:"many2many:trade_stocks_seller;"`   // Stocks being sold by the seller
	StocksBuyer       []Stock   `gorm:"many2many:trade_stocks_buyer;"`    // Stocks being given by the buyer
	TradeDate         time.Time `json:"trade_date" gorm:"autoCreateTime"` // Timestamp of the trade
}
