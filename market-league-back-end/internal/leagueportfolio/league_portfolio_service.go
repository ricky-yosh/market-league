package leagueportfolio

import (
	"time"

	"github.com/market-league/internal/models"
	"github.com/market-league/internal/stock"
)

type LeaguePortfolioService struct {
	repo      LeaguePortfolioRepository
	stockRepo stock.StockRepository
}

func NewLeaguePortfolioService(leaguePortfolioRepo LeaguePortfolioRepository, stockRepo stock.StockRepository) *LeaguePortfolioService {
	return &LeaguePortfolioService{
		repo:      leaguePortfolioRepo,
		stockRepo: stockRepo,
	}
}

func (s *LeaguePortfolioService) CreateLeaguePortfolio(leagueID uint) (*models.LeaguePortfolio, error) {
	// Fetch league details
	league, err := s.repo.GetLeagueDetails(leagueID)
	if err != nil {
		return nil, err
	}

	// Initialize League Portfolio
	leaguePortfolio := &models.LeaguePortfolio{
		LeagueID:  league.ID,
		Name:      "Remaining League Stocks",
		CreatedAt: time.Now(),
	}

	// Create the League Portfolio
	createdLeaguePortfolio, err := s.repo.CreateLeaguePortfolio(leaguePortfolio)
	if err != nil {
		return nil, err
	}

	// Initialize stock pool (example: add initial stocks)
	initialStocks, err := s.stockRepo.GetAllStocks()
	if err != nil {
		return nil, err
	}

	// Assign stocks to the League Portfolio
	if err := s.repo.AddStocksToLeaguePortfolio(createdLeaguePortfolio.ID, initialStocks); err != nil {
		return nil, err
	}

	return createdLeaguePortfolio, nil
}
