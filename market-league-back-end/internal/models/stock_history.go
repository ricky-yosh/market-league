// models/price_history.go

package models

import (
	"time"
)

type PriceHistory struct {
	ID        uint      `gorm:"primaryKey;autoIncrement" json:"id"`
	StockID   uint      `gorm:"not null;index" json:"stock_id"`
	Stock     Stock     `gorm:"foreignKey:StockID;constraint:OnUpdate:CASCADE,OnDelete:CASCADE" json:"stock"`
	Price     float64   `gorm:"not null" json:"price"`
	Timestamp time.Time `gorm:"not null;default:CURRENT_TIMESTAMP" json:"timestamp"`
}
