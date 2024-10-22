package models

type Portfolio struct {
	ID         uint    `gorm:"primaryKey;autoIncrement"`             // Auto-incrementing primary key
	UserID     uint    `gorm:"not null"`                             // Foreign key to User
	User       User    `gorm:"foreignKey:UserID"`                    // Association with User
	TotalValue float64 `json:"total_value"`                          // Portfolio's total value
	Stocks     []Stock `gorm:"many2many:portfolio_stocks;"`          // Many-to-many relationship with Stocks through Portfolio_Stocks
	Trades     []Trade `json:"trades" gorm:"foreignKey:PortfolioID"` // One-to-Many relationship with Trades
}
