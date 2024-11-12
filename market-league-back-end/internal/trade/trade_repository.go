package trade

import (
	"fmt"

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

// FetchTradesByUserAndLeague retrieves trades associated with a user and league from the database.
func (r *TradeRepository) FetchTradesByUserAndLeague(userID, leagueID uint) ([]models.Trade, error) {
	var trades []models.Trade
	err := r.db.
		Where("(user1_id = ? OR user2_id = ?) AND league_id = ?", userID, userID, leagueID).
		Preload("Stocks1").
		Preload("Stocks2").
		Find(&trades).Error

	if err != nil {
		if err == gorm.ErrRecordNotFound {
			return nil, fmt.Errorf("no trades found for user ID %d in league ID %d", userID, leagueID)
		}
		return nil, fmt.Errorf("failed to fetch trades: %w", err)
	}
	return trades, nil
}
