package models

type Stock struct {
	ID           uint        `gorm:"primaryKey;autoIncrement"`             // Auto-incrementing primary key
	TickerSymbol string      `json:"ticker_symbol" gorm:"unique;not null"` // Unique ticker symbol
	CompanyName  string      `json:"company_name"`                         // Company name
	CurrentPrice float64     `json:"current_price"`                        // Current price of the stock
	PriceHistory string      `json:"price_history" gorm:"type:jsonb"`      // Price history stored as JSONB
	Portfolios   []Portfolio `gorm:"many2many:portfolio_stocks;"`          // Many-to-many relation through Portfolio_Stocks
	Trades       []Trade     `json:"trades" gorm:"foreignKey:StockID"`     // One-to-many relationship with Trades
}
