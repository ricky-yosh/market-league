package stock

import (
	"github.com/market-league/internal/models"
	"gorm.io/gorm"
)

// StockRepository provides access to stock-related operations in the database.
type StockRepository struct {
	db *gorm.DB
}

// NewStockRepository creates a new instance of StockRepository.
func NewStockRepository(db *gorm.DB) *StockRepository {
	return &StockRepository{db: db}
}

// GetStockByID fetches a stock by its ID.
func (r *StockRepository) GetStockByID(stockID uint) (*models.Stock, error) {
	var stock models.Stock
	err := r.db.First(&stock, stockID).Error
	if err != nil {
		return nil, err
	}
	return &stock, nil
}

// GetAllStocks fetches all stocks from the database.
func (r *StockRepository) GetAllStocks() ([]models.Stock, error) {
	var stocks []models.Stock
	err := r.db.Find(&stocks).Error
	if err != nil {
		return nil, err
	}
	return stocks, nil
}

// CreateStock creates a new stock in the database.
func (r *StockRepository) CreateStock(stock *models.Stock) error {
	return r.db.Create(stock).Error
}

// UpdateStock updates an existing stock in the database.
func (r *StockRepository) UpdateStock(stock *models.Stock) error {
	return r.db.Save(stock).Error
}

// DeleteStock deletes a stock by its ID from the database.
func (r *StockRepository) DeleteStock(stockID uint) error {
	return r.db.Delete(&models.Stock{}, stockID).Error
}
