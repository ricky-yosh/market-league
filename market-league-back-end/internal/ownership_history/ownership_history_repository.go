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
	FindActiveByStockIDAndPortfolioID(stockID uint, portfolioID uint) (*models.OwnershipHistory, error)
	GetAllStockHistoryByStockIDAndPortfolioID(stockID uint, portfolioID uint) ([]models.OwnershipHistory, error)
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
func (r *ownershipHistoryRepository) FindActiveByStockIDAndPortfolioID(stockID uint, portfolioID uint) (*models.OwnershipHistory, error) {
	var history models.OwnershipHistory
	err := r.db.
		Preload("Stock").
		Where("stock_id = ? AND portfolio_id = ? AND end_date IS NULL", stockID, portfolioID).
		First(&history).Error
	if err != nil {
		return nil, fmt.Errorf("unable to retrieve active ownership history with stockID and portfolioID: %v", err)
	}
	return &history, nil
}

// GetAllStockHistoryByStockIDAndPortfolioID gets all the history of a single stock so that we can calculate portfolio points
func (r *ownershipHistoryRepository) GetAllStockHistoryByStockIDAndPortfolioID(stockID uint, portfolioID uint) ([]models.OwnershipHistory, error) {
	var history []models.OwnershipHistory
	err := r.db.Where("stock_id = ? AND portfolio_id = ?", stockID, portfolioID).Find(&history).Error
	if err != nil {
		return nil, fmt.Errorf("unable to retrieve ownership history with stockID %d and portfolioID %d: %v", stockID, portfolioID, err)
	}
	return history, nil
}

// GetActiveHistories gets all the currently active histories
func (r *ownershipHistoryRepository) GetActiveHistories() ([]*models.OwnershipHistory, error) {
	var histories []*models.OwnershipHistory
	err := r.db.Where("end_date IS NULL").Find(&histories).Error
	if err != nil {
		return nil, err
	}
	return histories, nil
}
