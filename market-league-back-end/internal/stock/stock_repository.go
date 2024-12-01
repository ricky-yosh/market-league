package stock

import (
	"fmt"
	"time"

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

// GetStocksByIDs fetches multiple stocks by their IDs from the database.
func (r *StockRepository) GetStocksByIDs(stockIDs []uint) ([]models.Stock, error) {
	var stocks []models.Stock
	if err := r.db.Where("id IN ?", stockIDs).Find(&stocks).Error; err != nil {
		return nil, fmt.Errorf("failed to find stocks with IDs %v: %w", stockIDs, err)
	}
	return stocks, nil
}

// CreateStock creates a new stock and its initial price history within a transaction
func (r *StockRepository) CreateStock(stock *models.Stock) error {
	return r.db.Transaction(func(tx *gorm.DB) error {
		// Create the stock
		if err := tx.Create(stock).Error; err != nil {
			return err
		}

		// Create the initial price history entry
		priceHistory := models.PriceHistory{
			StockID:   stock.ID,
			Price:     stock.CurrentPrice,
			Timestamp: time.Now(), // Assuming CreatedAt is set
		}

		if err := tx.Create(&priceHistory).Error; err != nil {
			return err
		}

		// Optionally, associate the price history with the stock
		stock.PriceHistories = append(stock.PriceHistories, priceHistory)

		return nil
	})
}

// New CreateMultipleStocks method
func (r *StockRepository) CreateMultipleStocks(stocks []*models.Stock) error {
	return r.db.Transaction(func(tx *gorm.DB) error {
		for _, stock := range stocks {
			if err := tx.Create(stock).Error; err != nil {
				return err
			}

			priceHistory := models.PriceHistory{
				StockID:   stock.ID,
				Price:     stock.CurrentPrice,
				Timestamp: time.Now(), // Ensure current time is used
			}

			if err := tx.Create(&priceHistory).Error; err != nil {
				return err
			}

			stock.PriceHistories = append(stock.PriceHistories, priceHistory)
		}
		return nil
	})
}

func (r *StockRepository) UpdateCurrentPrice(stockID uint, newPrice float64, timestamp *time.Time) error {
	// Start a transaction
	return r.db.Transaction(func(tx *gorm.DB) error {
		var stock models.Stock
		if err := tx.First(&stock, stockID).Error; err != nil {
			return err
		}

		// Update the current price
		stock.CurrentPrice = newPrice
		if err := tx.Save(&stock).Error; err != nil {
			return err
		}

		// Create a new PriceHistory entry
		priceHistory := models.PriceHistory{
			StockID: stock.ID,
			Price:   newPrice,
		}

		if timestamp != nil {
			priceHistory.Timestamp = *timestamp
		}

		if err := tx.Create(&priceHistory).Error; err != nil {
			return err
		}

		return nil
	})
}

func (r *StockRepository) GetStockWithHistory(stockID uint) (models.Stock, error) {
	var stock models.Stock
	if err := r.db.Preload("PriceHistories").First(&stock, stockID).Error; err != nil {
		return models.Stock{}, err
	}
	return stock, nil
}

// GetAllStocks retrieves all stocks from the database.
func (r *StockRepository) GetAllStocks() ([]models.Stock, error) {
	var stocks []models.Stock
	if err := r.db.Find(&stocks).Error; err != nil {
		return nil, err
	}
	return stocks, nil
}
