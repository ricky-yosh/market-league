package trade

import (
	"fmt"
	"time"

	"github.com/market-league/internal/models"
	"github.com/market-league/internal/stock"
)

// TradeService handles business logic related to trades.
type TradeService struct {
	tradeRepo *TradeRepository       // Reference to the repository layer
	stockRepo *stock.StockRepository // Reference to the StockRepository
}

// NewTradeService creates a new instance of TradeService.
func NewTradeService(tradeRepo *TradeRepository, stockRepo *stock.StockRepository) *TradeService {
	return &TradeService{
		tradeRepo: tradeRepo,
		stockRepo: stockRepo,
	}
}

// CreateTrade creates a new trade between two players within a specific league and portfolios.
func (s *TradeService) CreateTrade(
	leagueID, player1ID, player2ID, player1PortfolioID, player2PortfolioID uint,
	player1StockIDs, player2StockIDs []uint) error {

	// Fetch stocks for Player 1 using the StockRepository's GetStocksByIDs method
	player1Stocks, err := s.stockRepo.GetStocksByIDs(player1StockIDs)
	if err != nil {
		return fmt.Errorf("failed to fetch stocks for Player 1: %v", err)
	}

	// Fetch stocks for Player 2 using the StockRepository's GetStocksByIDs method
	player2Stocks, err := s.stockRepo.GetStocksByIDs(player2StockIDs)
	if err != nil {
		return fmt.Errorf("failed to fetch stocks for Player 2: %v", err)
	}

	// Create a new trade instance
	trade := &models.Trade{
		LeagueID:           leagueID,
		Player1ID:          player1ID,
		Player2ID:          player2ID,
		Player1PortfolioID: player1PortfolioID,
		Player2PortfolioID: player2PortfolioID,
		Player1Stocks:      player1Stocks,
		Player2Stocks:      player2Stocks,
		Player1Confirmed:   false,
		Player2Confirmed:   false,
	}

	// Save the trade using the TradeRepository
	return s.tradeRepo.CreateTrade(trade)
}

// CreateTrade creates a new trade in the database.
func (s *TradeService) ConfirmTrade(tradeID, playerID uint) error {
	// Fetch the trade by ID
	trade, err := s.tradeRepo.GetTradeByID(tradeID)
	if err != nil {
		return fmt.Errorf("failed to fetch trade: %v", err)
	}

	// Mark the trade as confirmed by the player
	if playerID == trade.Player1ID {
		trade.Player1Confirmed = true
	} else if playerID == trade.Player2ID {
		trade.Player2Confirmed = true
	} else {
		return fmt.Errorf("player ID %d is not part of this trade", playerID)
	}

	// Check if both players have confirmed the trade
	if trade.Player1Confirmed && trade.Player2Confirmed {
		// Both players have confirmed; execute the stock swap
		if err := s.tradeRepo.SwapStocks(trade); err != nil {
			return fmt.Errorf("failed to swap stocks: %v", err)
		}

		// Set the confirmed timestamp
		now := time.Now()
		trade.ConfirmedAt = &now
	}

	// Save the updated trade
	if err := s.tradeRepo.SaveTrade(trade); err != nil {
		return fmt.Errorf("failed to save trade confirmation: %v", err)
	}

	return nil
}

func (s *TradeService) GetTrades(portfolioID, leagueID uint, filterByPortfolio, filterByLeague bool) ([]models.Trade, error) {
	// Call the repository to fetch trades based on the filter criteria
	return s.tradeRepo.GetTrades(portfolioID, leagueID, filterByPortfolio, filterByLeague)
}
