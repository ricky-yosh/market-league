// models/stock.go

package models

type Stock struct {
	ID             uint           `gorm:"primaryKey;autoIncrement" json:"id"`
	TickerSymbol   string         `json:"ticker_symbol" gorm:"unique;not null"`
	CompanyName    string         `json:"company_name"`
	CurrentPrice   float64        `json:"current_price"`
	PriceHistories []PriceHistory `json:"price_histories" gorm:"foreignKey:StockID;constraint:OnUpdate:CASCADE,OnDelete:CASCADE"`
}
