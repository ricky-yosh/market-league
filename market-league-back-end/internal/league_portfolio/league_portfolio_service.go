package leagueportfolio

import (
	"fmt"
	"time"

	"github.com/market-league/internal/draft"
	"github.com/market-league/internal/models"
	ownership_history "github.com/market-league/internal/ownership_history"
	"github.com/market-league/internal/portfolio"
	"github.com/market-league/internal/stock"
	"github.com/market-league/internal/utils"
)

type LeaguePortfolioService struct {
	repo                    *LeaguePortfolioRepository
	stockRepo               *stock.StockRepository
	portfolioRepo           *portfolio.PortfolioRepository
	ownershipHistoryService ownership_history.OwnershipHistoryServiceInterface
	draftProvider           draft.DraftChannelProvider
}

func NewLeaguePortfolioService(
	leaguePortfolioRepo *LeaguePortfolioRepository,
	stockRepo *stock.StockRepository,
	portfolioRepo *portfolio.PortfolioRepository,
	ownershipHistoryService ownership_history.OwnershipHistoryServiceInterface,
	draftProvider draft.DraftChannelProvider,
) *LeaguePortfolioService {
	return &LeaguePortfolioService{
		repo:                    leaguePortfolioRepo,
		stockRepo:               stockRepo,
		portfolioRepo:           portfolioRepo,
		ownershipHistoryService: ownershipHistoryService,
		draftProvider:           draftProvider,
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

func (s *LeaguePortfolioService) DraftStock(leagueID, userID, stockID uint) error {

	leaguePortfolioID, err := s.repo.GetLeaguePortfolioIDByLeagueID(leagueID)
	if err != nil {
		return fmt.Errorf("error fetching LeaguePortfolioID for LeagueID %d: %w", leagueID, err)
	}

	// Fetch the league portfolio
	leaguePortfolio, err := s.repo.GetLeaguePortfolioWithID(leaguePortfolioID)
	if err != nil {
		return fmt.Errorf("failed to fetch league portfolio: %v", err)
	}

	userPortfolioID, err := s.portfolioRepo.GetPortfolioIDByUserAndLeague(userID, leagueID)
	if err != nil {
		return fmt.Errorf("error fetching userPortfolioID for LeagueID %d and UserID %d: %w", leagueID, userID, err)
	}

	// Fetch the user portfolio
	userPortfolio, err := s.portfolioRepo.GetPortfolioWithID(userPortfolioID)
	if err != nil {
		return fmt.Errorf("failed to fetch user portfolio: %v", err)
	}

	// Check if the stock exists in the league portfolio
	var stockToDraft *models.Stock
	for _, stock := range leaguePortfolio.Stocks {
		if stock.ID == stockID {
			stockToDraft = &stock
			break
		}
	}

	if stockToDraft == nil {
		return fmt.Errorf("stock not found in league portfolio")
	}

	// Remove the stock from the league portfolio
	var updatedStocks []models.Stock
	for _, stock := range leaguePortfolio.Stocks {
		if stock.ID != stockID {
			updatedStocks = append(updatedStocks, stock)
		}
	}
	leaguePortfolio.Stocks = updatedStocks

	// Add the stock to the user's portfolio
	userPortfolio.Stocks = append(userPortfolio.Stocks, *stockToDraft)

	// Update both portfolios in the repository
	if err := s.repo.UpdateLeaguePortfolio(leaguePortfolio); err != nil {
		return fmt.Errorf("failed to update league portfolio: %v", err)
	}

	if err := s.portfolioRepo.UpdatePortfolio(userPortfolio); err != nil {
		return fmt.Errorf("failed to update user portfolio: %v", err)
	}

	// Update OwnershipHistory of stock
	portfolioID, err := s.portfolioRepo.GetPortfolioIDByUserAndLeague(userID, leagueID)
	if err != nil {
		return fmt.Errorf("failed to get user portfolioID: %v", err)
	}
	stockIDList := []uint{stockID} // Make list to have the correct input type
	stocks, err := s.stockRepo.GetStocksByIDs(stockIDList)
	if err != nil {
		return fmt.Errorf("unable to find stock: %v", err)
	}
	stock, err := utils.FirstStock(stocks)
	if err != nil {
		return fmt.Errorf("unable to access first stock: %v", err)
	}
	startingValue := stock.CurrentPrice
	startDate := time.Now()
	if err := s.ownershipHistoryService.CreateOwnershipHistory(portfolioID, stockID, startingValue, startDate); err != nil {
		return fmt.Errorf("failed to update user portfolio: %v", err)
	}

	return nil
}

// Helper function to get active drafts
func (s *LeaguePortfolioService) GetDraftSelectionChannel(leagueID uint) chan uint {
	return s.draftProvider.GetDraftSelectionChannel(leagueID)
}

// GetLeaguePortfolioInfo fetches the details of a LeaguePortfolio by ID.
func (s *LeaguePortfolioService) GetLeaguePortfolioInfo(leagueID uint) (*models.LeaguePortfolio, error) {
	// Fetch the LeaguePortfolioID from the LeagueID
	leaguePortfolioID, err := s.repo.GetLeaguePortfolioIDByLeagueID(leagueID)
	if err != nil {
		return nil, fmt.Errorf("error fetching LeaguePortfolioID for LeagueID %d: %w", leagueID, err)
	}

	// Fetch the LeaguePortfolio using the LeaguePortfolioID
	leaguePortfolio, err := s.repo.GetLeaguePortfolioWithID(leaguePortfolioID)
	if err != nil {
		return nil, fmt.Errorf("error fetching league portfolio: %w", err)
	}

	return leaguePortfolio, nil
}
