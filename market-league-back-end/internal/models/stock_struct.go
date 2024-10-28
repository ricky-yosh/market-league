package models

type Stock struct {
	ID           uint        `gorm:"primaryKey;autoIncrement"`
	TickerSymbol string      `json:"ticker_symbol" gorm:"unique;not null"`
	CompanyName  string      `json:"company_name"`
	CurrentPrice float64     `json:"current_price"`
	PriceHistory string      `json:"price_history" gorm:"type:jsonb"`
	Portfolios   []Portfolio `gorm:"many2many:portfolio_stocks;"` // Many-to-many with portfolios
}
