package trade

import (
	"github.com/market-league/internal/models"
	"gorm.io/gorm"
)

// TradeService handles business logic related to trades.
type TradeService struct {
	db *gorm.DB
}

// NewTradeService creates a new instance of TradeService.
func NewTradeService(repo *TradeRepository) *TradeService {
	return &TradeService{db: repo.db}
}

// CreateTrade creates a new trade in the database.
func (s *TradeService) CreateTrade(trade *models.Trade) error {
	return s.db.Create(trade).Error
}

// GetTradesByUser fetches all trades made by a specific user.
func (s *TradeService) GetTradesByUser(userID uint) ([]models.Trade, error) {
	var trades []models.Trade
	err := s.db.Where("user_id = ?", userID).Find(&trades).Error
	return trades, err
}

// GetTradesByPortfolio fetches all trades related to a specific portfolio.
func (s *TradeService) GetTradesByPortfolio(portfolioID uint) ([]models.Trade, error) {
	var trades []models.Trade
	err := s.db.Where("portfolio_id = ?", portfolioID).Find(&trades).Error
	return trades, err
}
