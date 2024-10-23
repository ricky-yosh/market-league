package stock

import (
	"encoding/json"
	"fmt"

	"github.com/market-league/internal/models"
	"gorm.io/gorm"
)

// StockService handles business logic related to stocks.
type StockService struct {
	db *gorm.DB
}

// NewStockService creates a new instance of StockService.
func NewStockService(repo *StockRepository) *StockService {
	return &StockService{db: repo.db}
}

// GetPrice retrieves the current price of a stock by its ID.
func (s *StockService) GetPrice(stockID uint) (float64, error) {
	// Fetch the stock by ID
	stock := &models.Stock{}
	if err := s.db.First(stock, stockID).Error; err != nil {
		return 0, fmt.Errorf("failed to find stock: %v", err)
	}

	return stock.CurrentPrice, nil
}

// GetPriceHistory retrieves the price history of a stock by its ID.
func (s *StockService) GetPriceHistory(stockID uint) ([]float64, error) {
	// Fetch the stock by ID
	stock := &models.Stock{}
	if err := s.db.First(stock, stockID).Error; err != nil {
		return nil, fmt.Errorf("failed to find stock: %v", err)
	}

	// Decode the price history from JSON
	var history []float64
	if stock.PriceHistory != "" {
		if err := json.Unmarshal([]byte(stock.PriceHistory), &history); err != nil {
			return nil, fmt.Errorf("failed to decode price history: %v", err)
		}
	}

	return history, nil
}

// CreateStock creates a new stock in the database.
func (s *StockService) CreateStock(stock *models.Stock) error {
	// Save the stock to the database
	return s.db.Create(stock).Error
}

// UpdateCurrentPrice updates the current price of a stock and records it in the price history.
func (s *StockService) UpdateCurrentPrice(stockID uint, newPrice float64) error {
	// Fetch the stock by ID
	stock := &models.Stock{}
	if err := s.db.First(stock, stockID).Error; err != nil {
		return fmt.Errorf("failed to find stock: %v", err)
	}

	// Update the stock's current price
	stock.CurrentPrice = newPrice

	// Record the price change in the stock's price history
	if err := s.RecordPriceHistory(stock, newPrice); err != nil {
		return err
	}

	// Save the updated stock back to the database
	return s.db.Save(stock).Error
}

// RecordPriceHistory appends a new price to the stock's price history.
func (s *StockService) RecordPriceHistory(stock *models.Stock, newPrice float64) error {
	// Decode the existing price history from JSON
	var history []float64
	if stock.PriceHistory != "" {
		if err := json.Unmarshal([]byte(stock.PriceHistory), &history); err != nil {
			return fmt.Errorf("failed to decode price history: %v", err)
		}
	}

	// Add the new price point to the history
	history = append(history, newPrice)

	// Encode the updated history back to JSON
	historyJSON, err := json.Marshal(history)
	if err != nil {
		return fmt.Errorf("failed to encode price history: %v", err)
	}

	stock.PriceHistory = string(historyJSON)
	return nil
}
