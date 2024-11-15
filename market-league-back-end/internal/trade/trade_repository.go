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

// FetchTrades retrieves trades associated with specific filters and sanitizes the output.
func (r *TradeRepository) FetchTrades(filters map[string]interface{}) ([]models.SanitizedTrade, error) {
	var trades []models.Trade
	query := r.db.Model(&models.Trade{}).Preload("User1").Preload("User2").Preload("Stocks1").Preload("Stocks2")

	// Apply filters dynamically
	if leagueID, exists := filters["league_id"]; exists {
		query = query.Where("league_id = ?", leagueID)
	}
	if user1ID, exists := filters["user1_id"]; exists {
		query = query.Where("user1_id = ?", user1ID)
	}
	if user2ID, exists := filters["user2_id"]; exists {
		query = query.Where("user2_id = ?", user2ID)
	}

	// Execute query
	if err := query.Find(&trades).Error; err != nil {
		if err == gorm.ErrRecordNotFound {
			return nil, fmt.Errorf("no trades found with the provided filters")
		}
		return nil, fmt.Errorf("failed to fetch trades: %w", err)
	}

	// Sanitize trades
	var sanitizedTrades []models.SanitizedTrade
	for _, trade := range trades {
		sanitizedTrades = append(sanitizedTrades, models.SanitizedTrade{
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
		})
	}

	return sanitizedTrades, nil
}
