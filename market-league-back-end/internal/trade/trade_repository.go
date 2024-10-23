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

// CreateTrade creates a new trade in the database.
func (r *TradeRepository) CreateTrade(trade *models.Trade) error {
	return r.db.Create(trade).Error
}

// GetTradesByUser fetches all trades made by a specific user.
func (r *TradeRepository) GetTradesByUser(userID uint) ([]models.Trade, error) {
	var trades []models.Trade
	err := r.db.Where("user_id = ?", userID).Find(&trades).Error
	return trades, err
}

// GetTradesByPortfolio fetches all trades related to a specific portfolio.
func (r *TradeRepository) GetTradesByPortfolio(portfolioID uint) ([]models.Trade, error) {
	var trades []models.Trade
	err := r.db.Where("portfolio_id = ?", portfolioID).Find(&trades).Error
	return trades, err
}
