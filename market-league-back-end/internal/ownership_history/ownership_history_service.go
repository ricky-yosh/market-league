package ownershiphistory

import (
	"fmt"
	"time"

	"github.com/market-league/internal/models"
)

// OwnershipHistoryServiceInterface defines the interface for business logic
type OwnershipHistoryServiceInterface interface {
	CreateOwnershipHistory(portfolioID uint, stockID uint, startingValue float64, startDate time.Time) error
	UpdateOwnershipHistory(portfolioID uint, stockID uint, currentValue float64, endDate *time.Time) error
	UpdateActiveOwnershipHistoryCurrentPrices() error
}

// ownershipHistoryService implements OwnershipHistoryService
type ownershipHistoryService struct {
	repo OwnershipHistoryRepositoryInterface
}

// NewOwnershipHistoryService creates a new service
func NewOwnershipHistoryService(repo OwnershipHistoryRepositoryInterface) OwnershipHistoryServiceInterface {
	return &ownershipHistoryService{repo: repo}
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
		err = s.repo.Update(histories[i]) // Pass by reference
		if err != nil {
			return fmt.Errorf("unable to update current value of ownershipHistory ID %d: %v", histories[i].ID, err)
		}
	}

	return nil
}
