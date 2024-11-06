package trade

import (
	"time"

	"github.com/market-league/internal/models"
	"github.com/market-league/internal/stock"
)

// TradeService handles the business logic for trades.
type TradeService struct {
	TradeRepo *TradeRepository
	StockRepo *stock.StockRepository
}

// NewTradeService creates a new instance of TradeService
func NewTradeService(tradeRepo *TradeRepository, stockRepo *stock.StockRepository) *TradeService {
	return &TradeService{
		TradeRepo: tradeRepo,
		StockRepo: stockRepo,
	}
}

// CreateTrade initializes a new trade between two users.
func (s *TradeService) CreateTrade(leagueID, user1ID, user2ID, portfolio1ID, portfolio2ID uint, stocks1IDs, stocks2IDs []uint) (*models.Trade, error) {
	// Fetch stock details from the repository
	stocks1, err := s.StockRepo.GetStocksByIDs(stocks1IDs)
	if err != nil {
		return nil, err
	}
	stocks2, err := s.StockRepo.GetStocksByIDs(stocks2IDs)
	if err != nil {
		return nil, err
	}

	trade := &models.Trade{
		LeagueID:     leagueID,
		User1ID:      user1ID,
		User2ID:      user2ID,
		Portfolio1ID: portfolio1ID,
		Portfolio2ID: portfolio2ID,
		Stocks1:      stocks1,
		Stocks2:      stocks2,
		Status:       "pending",
		CreatedAt:    time.Now(),
		UpdatedAt:    time.Now(),
	}

	if err := s.TradeRepo.CreateTrade(trade); err != nil {
		return nil, err
	}

	return trade, nil
}
