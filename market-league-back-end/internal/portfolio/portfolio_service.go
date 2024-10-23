package portfolio

import (
	"fmt"

	"github.com/market-league/internal/models"
)

// PortfolioService handles business logic related to portfolios.
type PortfolioService struct {
	repo *PortfolioRepository
}

// NewPortfolioService creates a new instance of PortfolioService.
func NewPortfolioService(repo *PortfolioRepository) *PortfolioService {
	return &PortfolioService{repo: repo}
}

// GetPortfolio fetches a portfolio by its ID.
func (s *PortfolioService) GetPortfolio(portfolioID uint) (*models.Portfolio, error) {
	return s.repo.GetPortfolioByID(portfolioID)
}

// GetUserPortfolio fetches a user's portfolio in a specific league.
func (s *PortfolioService) GetUserPortfolio(userID, leagueID uint) (*models.Portfolio, error) {
	return s.repo.GetUserPortfolioInLeague(userID, leagueID)
}

// CreatePortfolio creates a new portfolio for a user in a league.
func (s *PortfolioService) CreatePortfolio(userID, leagueID uint) (*models.Portfolio, error) {
	// Check if the user already has a portfolio in the league
	existingPortfolio, err := s.repo.GetUserPortfolioInLeague(userID, leagueID)
	if err == nil && existingPortfolio != nil {
		return nil, fmt.Errorf("user already has a portfolio in this league")
	}

	// Create a new portfolio
	portfolio := &models.Portfolio{
		UserID:   userID,
		LeagueID: leagueID,
	}

	// Save the portfolio to the repository
	err = s.repo.CreatePortfolio(portfolio)
	if err != nil {
		return nil, fmt.Errorf("failed to create portfolio: %v", err)
	}

	return portfolio, nil
}

// AddStockToPortfolio adds a stock to the user's portfolio.
func (s *PortfolioService) AddStockToPortfolio(portfolioID, stockID uint) error {
	// Get the portfolio by ID
	portfolio, err := s.repo.GetPortfolioByID(portfolioID)
	if err != nil {
		return fmt.Errorf("failed to fetch portfolio: %v", err)
	}

	// Get the stock by ID (assuming a StockRepository exists)
	stock := &models.Stock{}
	err = s.repo.db.First(stock, stockID).Error
	if err != nil {
		return fmt.Errorf("failed to fetch stock: %v", err)
	}

	// Add the stock to the portfolio
	portfolio.Stocks = append(portfolio.Stocks, *stock)

	// Update the portfolio in the repository
	return s.repo.UpdatePortfolio(portfolio)
}

// RemoveStockFromPortfolio removes a stock from the user's portfolio.
func (s *PortfolioService) RemoveStockFromPortfolio(portfolioID, stockID uint) error {
	// Get the portfolio by ID
	portfolio, err := s.repo.GetPortfolioByID(portfolioID)
	if err != nil {
		return fmt.Errorf("failed to fetch portfolio: %v", err)
	}

	// Remove the stock from the portfolio
	var updatedStocks []models.Stock
	for _, s := range portfolio.Stocks {
		if s.ID != stockID {
			updatedStocks = append(updatedStocks, s)
		}
	}
	portfolio.Stocks = updatedStocks

	// Update the portfolio in the repository
	return s.repo.UpdatePortfolio(portfolio)
}
