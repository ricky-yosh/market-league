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

// CreateTrade inserts a new trade record into the database.
func (r *TradeRepository) CreateTrade(trade *models.Trade) error {
	// Use a database transaction to ensure data consistency.
	return r.db.Transaction(func(tx *gorm.DB) error {
		// Create the trade record
		if err := tx.Create(trade).Error; err != nil {
			return fmt.Errorf("failed to create trade record: %w", err)
		}

		// Associate Player 1's stocks with the trade
		if len(trade.Player1Stocks) > 0 {
			if err := tx.Model(trade).Association("Player1Stocks").Replace(trade.Player1Stocks); err != nil {
				return fmt.Errorf("failed to associate Player 1's stocks with the trade: %w", err)
			}
		}

		// Associate Player 2's stocks with the trade
		if len(trade.Player2Stocks) > 0 {
			if err := tx.Model(trade).Association("Player2Stocks").Replace(trade.Player2Stocks); err != nil {
				return fmt.Errorf("failed to associate Player 2's stocks with the trade: %w", err)
			}
		}

		return nil
	})
}

func (r *TradeRepository) GetTradeByID(tradeID uint) (*models.Trade, error) {
	var trade models.Trade
	if err := r.db.Preload("Player1Stocks").Preload("Player2Stocks").First(&trade, tradeID).Error; err != nil {
		return nil, fmt.Errorf("failed to find trade with ID %d: %w", tradeID, err)
	}
	return &trade, nil
}

// SaveTrade updates a trade record in the database.
func (r *TradeRepository) SaveTrade(trade *models.Trade) error {
	return r.db.Save(trade).Error
}

// SwapStocks swaps the stocks between two portfolios.
func (r *TradeRepository) SwapStocks(trade *models.Trade) error {
	// Transfer stocks from Player 1 to Player 2
	for _, stock := range trade.Player1Stocks {
		// Remove stock from Player 1's portfolio and add to Player 2's portfolio
		if err := r.db.Model(&trade.Player1PortfolioID).Association("Stocks").Delete(&stock); err != nil {
			return fmt.Errorf("failed to remove stock from Player 1's portfolio: %w", err)
		}
		if err := r.db.Model(&trade.Player2PortfolioID).Association("Stocks").Append(&stock); err != nil {
			return fmt.Errorf("failed to add stock to Player 2's portfolio: %w", err)
		}
	}

	// Transfer stocks from Player 2 to Player 1
	for _, stock := range trade.Player2Stocks {
		// Remove stock from Player 2's portfolio and add to Player 1's portfolio
		if err := r.db.Model(&trade.Player2PortfolioID).Association("Stocks").Delete(&stock); err != nil {
			return fmt.Errorf("failed to remove stock from Player 2's portfolio: %w", err)
		}
		if err := r.db.Model(&trade.Player1PortfolioID).Association("Stocks").Append(&stock); err != nil {
			return fmt.Errorf("failed to add stock to Player 1's portfolio: %w", err)
		}
	}

	return nil
}

// GetTrades fetches trades based on the provided filter criteria.
func (r *TradeRepository) GetTrades(portfolioID, leagueID uint, filterByPortfolio, filterByLeague bool) ([]models.Trade, error) {
	var trades []models.Trade
	query := r.db

	// Apply filters if requested
	if filterByPortfolio {
		query = query.Where("player1_portfolio_id = ? OR player2_portfolio_id = ?", portfolioID, portfolioID)
	}

	if filterByLeague {
		query = query.Where("league_id = ?", leagueID)
	}

	// Fetch the filtered trades from the database
	if err := query.Preload("Player1Stocks").Preload("Player2Stocks").Find(&trades).Error; err != nil {
		return nil, fmt.Errorf("failed to fetch trades: %w", err)
	}

	return trades, nil
}
