package ownershiphistory

import (
	"fmt"

	"github.com/market-league/internal/models"
	"gorm.io/gorm"
)

// OwnershipHistoryRepositoryInterface defines the interface for database operations
type OwnershipHistoryRepositoryInterface interface {
	Create(history *models.OwnershipHistory) error
	Update(history *models.OwnershipHistory) error
	FindByStockIDAndPortfolioID(stockID uint, portfolioID uint) (*models.OwnershipHistory, error)
	GetActiveHistories() ([]*models.OwnershipHistory, error)
}

// ownershipHistoryRepository implements OwnershipHistoryRepository
type ownershipHistoryRepository struct {
	db *gorm.DB
}

// NewOwnershipHistoryRepository creates a new repository
func NewOwnershipHistoryRepository(db *gorm.DB) OwnershipHistoryRepositoryInterface {
	return &ownershipHistoryRepository{db: db}
}

// Create adds a new OwnershipHistory record to the database
func (r *ownershipHistoryRepository) Create(history *models.OwnershipHistory) error {
	return r.db.Create(history).Error
}

// Update modifies an existing OwnershipHistory record in the database
func (r *ownershipHistoryRepository) Update(history *models.OwnershipHistory) error {
	return r.db.Save(history).Error
}

// FindByID fetches an OwnershipHistory record by its ID
func (r *ownershipHistoryRepository) FindByStockIDAndPortfolioID(stockID uint, portfolioID uint) (*models.OwnershipHistory, error) {
	var history models.OwnershipHistory
	err := r.db.Where("stock_id = ? AND portfolio_id = ?", stockID, portfolioID).First(&history).Error
	if err != nil {
		return nil, err
	}
	if history.EndDate != nil {
		return nil, fmt.Errorf("EndDate is not nil, which means that this history section is not mutable")
	}
	return &history, nil
}

// GetActiveHistories gets all the currently active histories
func (r *ownershipHistoryRepository) GetActiveHistories() ([]*models.OwnershipHistory, error) {
	var histories []*models.OwnershipHistory
	err := r.db.Where("end_date IS NULL").Find(histories).Error
	if err != nil {
		return nil, err
	}
	return histories, nil
}
