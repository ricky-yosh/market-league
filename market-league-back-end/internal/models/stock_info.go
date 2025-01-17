// models/stock_info.go

package models

type StockInfo struct {
	ID             uint              `json:"id"`
	TickerSymbol   string            `json:"ticker_symbol"`
	CompanyName    string            `json:"company_name"`
	CurrentPrice   float64           `json:"current_price"`
	PriceHistories []PriceHistoryDTO `json:"price_histories"`
}
