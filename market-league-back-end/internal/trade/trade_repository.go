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
func (r *TradeRepository) FetchTradesByUserAndLeague(userID, leagueID uint) ([]models.SanitizedTrade, error) {
	var trades []models.Trade
	err := r.db.
		Where("(user1_id = ? OR user2_id = ?) AND league_id = ?", userID, userID, leagueID).
		Preload("User1").
		Preload("User2").
		Preload("Stocks1").
		Preload("Stocks2").
		Find(&trades).Error

	if err != nil {
		if err == gorm.ErrRecordNotFound {
			return nil, fmt.Errorf("no trades found for user ID %d in league ID %d", userID, leagueID)
		}
		return nil, fmt.Errorf("failed to fetch trades: %w", err)
	}

	// Map to sanitized trades
	var sanitizedTrades []models.SanitizedTrade
	for _, trade := range trades {
		sanitizedTrade := models.SanitizedTrade{
			ID:       trade.ID,
			LeagueID: trade.LeagueID,
			User1: models.SanitizedUser{
				ID:        trade.User1.ID,
				Username:  trade.User1.Username,
				Email:     trade.User1.Email,
				CreatedAt: trade.User1.CreatedAt,
			},
			User2: models.SanitizedUser{
				ID:        trade.User2.ID,
				Username:  trade.User2.Username,
				Email:     trade.User2.Email,
				CreatedAt: trade.User2.CreatedAt,
			},
			Portfolio1ID:   trade.Portfolio1ID,
			Portfolio2ID:   trade.Portfolio2ID,
			Stocks1:        trade.Stocks1,
			Stocks2:        trade.Stocks2,
			User1Confirmed: trade.User1Confirmed,
			User2Confirmed: trade.User2Confirmed,
			Status:         trade.Status,
			CreatedAt:      trade.CreatedAt,
			UpdatedAt:      trade.UpdatedAt,
		}
		sanitizedTrades = append(sanitizedTrades, sanitizedTrade)
	}
	return sanitizedTrades, nil
}
