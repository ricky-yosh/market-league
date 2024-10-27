package stock

import (
	"fmt"

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
	if err := r.db.First(&stock, stockID).Error; err != nil {
		return nil, fmt.Errorf("failed to find stock with ID %d: %w", stockID, err)
	}
	return &stock, nil
}

// GetStocksByIDs fetches multiple stocks by their IDs from the database.
func (r *StockRepository) GetStocksByIDs(stockIDs []uint) ([]models.Stock, error) {
	var stocks []models.Stock
	if err := r.db.Where("id IN ?", stockIDs).Find(&stocks).Error; err != nil {
		return nil, fmt.Errorf("failed to find stocks with IDs %v: %w", stockIDs, err)
	}
	return stocks, nil
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
	// Insert the new stock into the database
	if err := r.db.Create(stock).Error; err != nil {
		return fmt.Errorf("failed to create stock: %w", err)
	}
	return nil
}

// UpdateStockPriceHistory updates the price history of a stock.
func (r *StockRepository) UpdateStockPriceHistory(stock *models.Stock) error {
	return r.db.Model(stock).Update("price_history", stock.PriceHistory).Error
}

// UpdateStock updates an existing stock in the database.
func (r *StockRepository) UpdateStock(stock *models.Stock) error {
	return r.db.Save(stock).Error
}

// DeleteStock deletes a stock by its ID from the database.
func (r *StockRepository) DeleteStock(stockID uint) error {
	return r.db.Delete(&models.Stock{}, stockID).Error
}

// UpdateStockPrice updates the current price of a stock in the database.
func (r *StockRepository) UpdateStockPrice(stock *models.Stock) error {
	return r.db.Model(stock).Update("current_price", stock.CurrentPrice).Error
}
