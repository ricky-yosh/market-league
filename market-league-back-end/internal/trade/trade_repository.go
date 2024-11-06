package trade

import (
	"github.com/market-league/internal/models"
	"gorm.io/gorm"
)

// TradeRepository provides access to trade-related operations in the database.
type TradeRepository struct {
	db *gorm.DB
}

// NewTradeRepository creates a new instance of TradeRepository.
func NewTradeRepository(db *gorm.DB) *TradeRepository {
	return &TradeRepository{db: db}
}

// CreateTrade inserts a new trade into the database
func (r *TradeRepository) CreateTrade(trade *models.Trade) error {
	return r.db.Create(trade).Error
}
