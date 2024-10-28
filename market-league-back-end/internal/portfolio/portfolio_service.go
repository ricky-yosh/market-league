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
func (s *PortfolioService) GetPortfolioWithID(portfolioID uint) (*models.SanitizedPortfolio, error) {
	portfolio, err := s.repo.GetPortfolioWithID(portfolioID)
	if err != nil {
		return nil, err
	}

	// Convert to PortfolioDTO
	dto := models.SanitizedPortfolio{
		ID:        portfolio.ID,
		UserID:    portfolio.UserID,
		User:      models.SanitizedUser{ID: portfolio.User.ID, Username: portfolio.User.Username, Email: portfolio.User.Email, CreatedAt: portfolio.User.CreatedAt},
		LeagueID:  portfolio.LeagueID,
		League:    models.SanitizedLeague{ID: portfolio.League.ID, LeagueName: portfolio.League.LeagueName, StartDate: portfolio.League.StartDate, EndDate: portfolio.League.EndDate},
		CreatedAt: portfolio.CreatedAt,
	}

	// Map the stocks to the DTO
	for _, stock := range portfolio.Stocks {
		dto.Stocks = append(dto.Stocks, models.SanitizedStock{
			ID:           stock.ID,
			TickerSymbol: stock.TickerSymbol,
			CompanyName:  stock.CompanyName,
			CurrentPrice: stock.CurrentPrice,
		})
	}

	return &dto, nil
}

// GetUserPortfolio fetches a user's portfolio in a specific league.
func (s *PortfolioService) GetLeaguePortfolio(userID, leagueID uint) (*models.SanitizedPortfolio, error) {
	// Get the portfolio ID for the given user and league
	portfolioID, err := s.repo.GetPortfolioIDByUserAndLeague(userID, leagueID)
	if err != nil {
		return nil, err
	}

	// Use the portfolio ID to fetch the sanitized portfolio
	return s.GetPortfolioWithID(portfolioID)
}

// CreatePortfolio creates a new portfolio for a user in a league.
func (s *PortfolioService) CreatePortfolio(userID, leagueID uint) (*models.Portfolio, error) {
	// Step 1: Check if the user exists
	userExists, err := s.repo.UserExists(userID)
	if err != nil {
		return nil, fmt.Errorf("failed to check user existence: %v", err)
	}
	if !userExists {
		return nil, fmt.Errorf("user with ID %d not found", userID)
	}

	// Step 2: Check if the league exists
	leagueExists, err := s.repo.LeagueExists(leagueID)
	if err != nil {
		return nil, fmt.Errorf("failed to check league existence: %v", err)
	}
	if !leagueExists {
		return nil, fmt.Errorf("league with ID %d not found", leagueID)
	}

	// Step 3: Check if the user already has a portfolio in the league
	existingPortfolio, err := s.repo.GetPortfolioIDByUserAndLeague(userID, leagueID)
	if err == nil && existingPortfolio != 0 {
		return nil, fmt.Errorf("user already has a portfolio in this league")
	}

	// Step 4: Create a new portfolio
	portfolio := &models.Portfolio{
		UserID:   userID,
		LeagueID: leagueID,
	}

	// Step 5: Save the portfolio to the repository
	err = s.repo.CreatePortfolio(portfolio)
	if err != nil {
		return nil, fmt.Errorf("failed to create portfolio: %v", err)
	}

	return portfolio, nil
}

// AddStockToPortfolio adds a stock to the user's portfolio.
func (s *PortfolioService) AddStockToPortfolio(portfolioID, stockID uint) error {
	// Get the portfolio by ID
	portfolio, err := s.repo.GetPortfolioWithID(portfolioID)
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
	portfolio, err := s.repo.GetPortfolioWithID(portfolioID)
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

// CalculateTotalValue calculates the total value of the portfolio based on its stocks.
func (s *PortfolioService) CalculateTotalValue(portfolio *models.Portfolio) float64 {
	total := 0.0

	// Iterate over each stock in the portfolio and sum up their current prices.
	for _, stock := range portfolio.Stocks {
		total += stock.CurrentPrice
	}

	return total
}

// Helper Functions
func (r *PortfolioRepository) UserExists(userID uint) (bool, error) {
	var count int64
	if err := r.db.Model(&models.User{}).Where("id = ?", userID).Count(&count).Error; err != nil {
		return false, err
	}
	return count > 0, nil
}

func (r *PortfolioRepository) LeagueExists(leagueID uint) (bool, error) {
	var count int64
	if err := r.db.Model(&models.League{}).Where("id = ?", leagueID).Count(&count).Error; err != nil {
		return false, err
	}
	return count > 0, nil
}
