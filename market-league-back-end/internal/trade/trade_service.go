package trade

import (
	"log"
	"time"

	"github.com/market-league/internal/models"
	"github.com/market-league/internal/portfolio"
	"github.com/market-league/internal/stock"
	"github.com/market-league/internal/user"
)

// TradeService handles the business logic for trades.
type TradeService struct {
	TradeRepo     *TradeRepository
	StockRepo     *stock.StockRepository
	PortfolioRepo *portfolio.PortfolioRepository
	UserRepo      *user.UserRepository
}

// NewTradeService creates a new instance of TradeService
func NewTradeService(tradeRepo *TradeRepository, stockRepo *stock.StockRepository, portfolioRepo *portfolio.PortfolioRepository, userRepo *user.UserRepository) *TradeService {
	return &TradeService{
		TradeRepo:     tradeRepo,
		StockRepo:     stockRepo,
		PortfolioRepo: portfolioRepo,
		UserRepo:      userRepo,
	}
}

// CreateTrade initializes a new trade between two users.
func (s *TradeService) CreateTrade(leagueID, user1ID, user2ID uint, stocks1IDs, stocks2IDs []uint) (*models.SanitizedTrade, error) {

	// Fetch stock details from the repository
	stocks1, err := s.StockRepo.GetStocksByIDs(stocks1IDs)
	if err != nil {
		return nil, err
	}
	stocks2, err := s.StockRepo.GetStocksByIDs(stocks2IDs)
	if err != nil {
		return nil, err
	}

	// Fetch user details from the repository
	user1, err := s.UserRepo.GetUserByID(user1ID)
	if err != nil {
		return nil, err
	}
	user2, err := s.UserRepo.GetUserByID(user2ID)
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
		User1:        user1,
		User1ID:      user1ID,
		User2:        user2,
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
	// Convert to a sanitized trade
	sanitizedTrade := &models.SanitizedTrade{
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

	return sanitizedTrade, nil
}

func (s *TradeService) GetTrades(leagueID uint, userID *uint, receivingTrade *bool, sendingTrade *bool) ([]models.SanitizedTrade, error) {
	// Build filters based on input
	filters := map[string]interface{}{
		"league_id": leagueID,
	}

	if userID != nil {
		if receivingTrade != nil && *receivingTrade {
			filters["user2_id"] = *userID
		}
		if sendingTrade != nil && *sendingTrade {
			filters["user1_id"] = *userID
		}
	}

	// Call the repository to fetch filtered trades
	return s.TradeRepo.FetchTrades(filters)
}
