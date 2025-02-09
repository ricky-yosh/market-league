package portfolio

import (
	"fmt"
	"math"

	"github.com/market-league/internal/models"
	ownership_history "github.com/market-league/internal/ownership_history"
)

// PortfolioService handles business logic related to portfolios.
type PortfolioService struct {
	repo                 *PortfolioRepository
	ownershipHistoryRepo ownership_history.OwnershipHistoryRepositoryInterface
}

// NewPortfolioService creates a new instance of PortfolioService.
func NewPortfolioService(
	repo *PortfolioRepository,
	ownershipHistoryRepo ownership_history.OwnershipHistoryRepositoryInterface,
) *PortfolioService {
	return &PortfolioService{
		repo:                 repo,
		ownershipHistoryRepo: ownershipHistoryRepo,
	}
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
		Points:    portfolio.Points,
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
	fmt.Printf("Creating Portfolio... Disregard error above. Error is a check to see if user already has a portfolio.\n")

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

	// Check if the stock is already in the portfolio
	for _, s := range portfolio.Stocks {
		if s.ID == stockID {
			return fmt.Errorf("stock with ID %d is already in the portfolio", stockID)
		}
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
	if err := s.repo.UpdatePortfolio(portfolio); err != nil {
		return fmt.Errorf("failed to update portfolio: %v", err)
	}

	return nil
}

// RemoveStockFromPortfolio removes a stock from the user's portfolio.
func (s *PortfolioService) RemoveStockFromPortfolio(portfolioID, stockID uint) error {
	// Get the portfolio by ID
	portfolio, err := s.repo.GetPortfolioWithID(portfolioID)
	if err != nil {
		return fmt.Errorf("failed to fetch portfolio: %v", err)
	}

	// Check if the stock exists in the portfolio
	var stockFound bool
	var updatedStocks []models.Stock
	for _, s := range portfolio.Stocks {
		if s.ID == stockID {
			stockFound = true
			continue // Skip adding this stock to the updated list
		}
		updatedStocks = append(updatedStocks, s)
	}

	if !stockFound {
		return fmt.Errorf("stock with ID %d is not in the portfolio", stockID)
	}

	// Update the portfolio's stocks
	portfolio.Stocks = updatedStocks

	// Update the portfolio in the repository
	if err := s.repo.UpdatePortfolio(portfolio); err != nil {
		return fmt.Errorf("failed to update portfolio: %v", err)
	}

	return nil
}

// CalculateAllPortfolioTotalValues calculates the value of every portfolio
func (s *PortfolioService) CalculateAllPortfolioTotalValues() error {
	// Get all portfolios to update
	allPortfolios, err := s.repo.GetAllPortfolios()
	if err != nil {
		return fmt.Errorf("unable to load all portfolios: %v", err)
	}
	for index := range allPortfolios {
		// update each portfolio's value
		portfolio := allPortfolios[index]
		err := s.CalculatePortfolioTotalValue(&portfolio)
		if err != nil {
			return fmt.Errorf("unable to calculate portfolio total value %v", err)
		}
	}

	return nil
}

// CalculateTotalValue calculates the total value of the portfolio based on its stocks.
func (s *PortfolioService) CalculatePortfolioTotalValue(portfolio *models.Portfolio) error {
	totalPercentChangeForPortfolio := 0.0
	// Get percent change of each ownership_history
	for index := range portfolio.Stocks {
		stock := portfolio.Stocks[index]
		ownershipHistoryList, err := s.ownershipHistoryRepo.GetAllStockHistoryByStockIDAndPortfolioID(stock.ID, portfolio.ID)
		if err != nil {
			return fmt.Errorf("unable to retrieve ownership history with stockID and portfolioID: %v", err)
		}
		totalPercentageChangeForStock := 0.0
		for index := range ownershipHistoryList {
			ownershipHistoryItem := ownershipHistoryList[index]
			currentVal := ownershipHistoryItem.CurrentValue

			previousVal := ownershipHistoryItem.StartingValue
			// Check for 0 and replace with 1 to avoid infinity
			if previousVal == 0 {
				previousVal = 1
			}
			// Calculate the percentage change and use that in the point scoring system
			percentChangeForItem := ((currentVal - previousVal) / math.Abs(previousVal)) * 100
			totalPercentageChangeForStock = totalPercentageChangeForStock + percentChangeForItem
		}
		totalPercentChangeForPortfolio = totalPercentChangeForPortfolio + totalPercentageChangeForStock
	}
	// Add all of them up and update the score
	portfolioValue := int(math.Round(totalPercentChangeForPortfolio))
	err := s.repo.UpdatePortfolioPoints(portfolio.ID, portfolioValue)
	if err != nil {
		return fmt.Errorf("unable to update portfolio points: %v", err)
	}
	err = s.repo.LogPortfolioPointsChange(portfolio.ID, portfolioValue)
	if err != nil {
		return fmt.Errorf("unable to log portfolio points change %v", err)
	}
	return nil
}

// GetPortfolioPointsHistory
func (s *PortfolioService) GetPortfolioPointsHistory(portfolioID uint) ([]models.PortfolioPointsHistory, error) {
	// Get the portfolio by ID
	portfolioPointsHistoryEntries, err := s.repo.GetPortfolioPointsHistoryEntry(portfolioID)
	if err != nil {
		return nil, fmt.Errorf("failed to fetch portfolioPointsHistoryEntries: %v", err)
	}

	return portfolioPointsHistoryEntries, nil
}

// GetStocksValueChange
func (s *PortfolioService) GetStocksValueChange(portfolioID uint) ([]*models.OwnershipHistory, error) {
	var stocksAndHistory []*models.OwnershipHistory
	// Get the portfolio by ID
	portfolio, err := s.repo.GetPortfolioWithID(portfolioID)
	if err != nil {
		return nil, fmt.Errorf("failed to fetch portfolio: %v", err)
	}
	for index := range portfolio.Stocks {
		stock := portfolio.Stocks[index]
		historyItem, err := s.ownershipHistoryRepo.FindActiveByStockIDAndPortfolioID(stock.ID, portfolioID)
		if err != nil {
			return nil, fmt.Errorf("unable to retrieve ownershipHistoryItem with stockID and portfolioID: %v", err)
		}
		stocksAndHistory = append(stocksAndHistory, historyItem)
	}

	return stocksAndHistory, nil
}
