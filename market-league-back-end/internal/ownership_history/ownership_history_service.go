package ownershiphistory

import (
	"fmt"
	"time"

	"github.com/market-league/internal/models"
	"github.com/market-league/internal/stock"
	"github.com/market-league/internal/utils"
)

// OwnershipHistoryServiceInterface defines the interface for business logic
type OwnershipHistoryServiceInterface interface {
	CreateOwnershipHistory(portfolioID uint, stockID uint, startingValue float64, startDate time.Time) error
	UpdateOwnershipHistory(portfolioID uint, stockID uint, currentValue float64, endDate *time.Time) error
	UpdateActiveOwnershipHistoryCurrentPrices() error
}

// ownershipHistoryService implements OwnershipHistoryService
type ownershipHistoryService struct {
	repo      OwnershipHistoryRepositoryInterface
	stockRepo *stock.StockRepository
}

// NewOwnershipHistoryService creates a new service
func NewOwnershipHistoryService(
	repo OwnershipHistoryRepositoryInterface,
	stockRepo *stock.StockRepository,
) OwnershipHistoryServiceInterface {
	return &ownershipHistoryService{
		repo:      repo,
		stockRepo: stockRepo,
	}
}

// CreateOwnershipHistory creates a new ownership record
func (s *ownershipHistoryService) CreateOwnershipHistory(portfolioID uint, stockID uint, startingValue float64, startDate time.Time) error {
	history := &models.OwnershipHistory{
		PortfolioID:   portfolioID,
		StockID:       stockID,
		StartingValue: startingValue,
		CurrentValue:  startingValue, // Initial current value is the same as starting value
		StartDate:     startDate,
	}
	err := s.repo.Create(history)
	if err != nil {
		return err
	}
	return nil
}

// UpdateOwnershipHistory updates an existing ownership record
func (s *ownershipHistoryService) UpdateOwnershipHistory(portfolioID uint, stockID uint, currentValue float64, endDate *time.Time) error {
	history, err := s.repo.FindActiveByStockIDAndPortfolioID(stockID, portfolioID)
	if err != nil {
		return err
	}
	history.CurrentValue = currentValue
	history.EndDate = endDate
	err = s.repo.Update(history)
	if err != nil {
		return err
	}
	return nil
}

// UpdateOwnershipHistory updates all active ownership records
func (s *ownershipHistoryService) UpdateActiveOwnershipHistoryCurrentPrices() error {
	histories, err := s.repo.GetActiveHistories()
	if err != nil {
		return fmt.Errorf("unable to retrieve active ownership histories: %v", err)
	}

	// Iterate over each history and update its current value
	for i := range histories {
		stockIDList := []uint{histories[i].StockID}
		// Fetch the current stock price
		stocks, err := s.stockRepo.GetStocksByIDs(stockIDList)
		if err != nil {
			return fmt.Errorf("unable to fetch current price for stock ID %d: %v", histories[i].StockID, err)
		}
		updatedStock, err := utils.FirstStock(stocks)
		if err != nil {
			return fmt.Errorf("unable to get first stock: %v", err)
		}

		// Update the history entry
		histories[i].CurrentValue = updatedStock.CurrentPrice

		// Save the updated history record
		err = s.repo.Update(histories[i]) // Pass by reference to update in-place
		if err != nil {
			return fmt.Errorf("unable to update ownership history ID %d: %v", histories[i].ID, err)
		}
	}
	fmt.Println("updated activeOwnershipHistory prices")
	return nil
}
