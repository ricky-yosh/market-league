package trade

import (
	"log"
	"time"

	"github.com/market-league/internal/models"
	"github.com/market-league/internal/portfolio"
	"github.com/market-league/internal/stock"
)

// TradeService handles the business logic for trades.
type TradeService struct {
	TradeRepo     *TradeRepository
	StockRepo     *stock.StockRepository
	PortfolioRepo *portfolio.PortfolioRepository
}

// NewTradeService creates a new instance of TradeService
func NewTradeService(tradeRepo *TradeRepository, stockRepo *stock.StockRepository, portfolioRepo *portfolio.PortfolioRepository) *TradeService {
	return &TradeService{
		TradeRepo:     tradeRepo,
		StockRepo:     stockRepo,
		PortfolioRepo: portfolioRepo,
	}
}

// CreateTrade initializes a new trade between two users.
func (s *TradeService) CreateTrade(leagueID, user1ID, user2ID uint, stocks1IDs, stocks2IDs []uint) (*models.Trade, error) {

	// Fetch stock details from the repository
	stocks1, err := s.StockRepo.GetStocksByIDs(stocks1IDs)
	if err != nil {
		return nil, err
	}
	stocks2, err := s.StockRepo.GetStocksByIDs(stocks2IDs)
	if err != nil {
		return nil, err
	}

	portfolio1ID, err := s.PortfolioRepo.GetPortfolioIDByUserAndLeague(user1ID, leagueID)
	if err != nil {
		// Handle the error appropriately (e.g., return it, log it, etc.)
		log.Printf("error fetching portfolio for user1: %v", err)
		return nil, err
	}
	log.Printf("Portfolio 1: %v", portfolio1ID)

	portfolio2ID, err := s.PortfolioRepo.GetPortfolioIDByUserAndLeague(user2ID, leagueID)
	if err != nil {
		// Handle the error appropriately (e.g., return it, log it, etc.)
		log.Printf("error fetching portfolio for user2: %v", err)
		return nil, err
	}
	log.Printf("Portfolio 2: %v", portfolio2ID)

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
