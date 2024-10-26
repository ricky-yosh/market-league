package stock

import (
	"encoding/json"
	"fmt"

	"github.com/market-league/internal/models"
)

// StockService handles business logic related to stocks.
type StockService struct {
	repo *StockRepository // Reference to the repository layer
}

// NewStockService creates a new instance of StockService.
func NewStockService(repo *StockRepository) *StockService {
	return &StockService{repo: repo}
}

// GetPrice retrieves the current price of a stock by its ID.
func (s *StockService) GetPrice(stockID uint) (float64, error) {
	// Fetch the stock by ID using the repository
	stock, err := s.repo.GetStockByID(stockID)
	if err != nil {
		return 0, fmt.Errorf("failed to find stock: %v", err)
	}

	return stock.CurrentPrice, nil
}

// GetPriceHistory retrieves the price history of a stock by its ID.
func (s *StockService) GetPriceHistory(stockID uint) ([]float64, error) {
	// Fetch the stock by ID using the repository
	stock, err := s.repo.GetStockByID(stockID)
	if err != nil {
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
	// Call the repository to create the stock
	return s.repo.CreateStock(stock)
}

// UpdateCurrentPrice updates the current price of a stock and records it in the price history.
func (s *StockService) UpdatePriceHistory(stockID uint, priceHistory []float64) error {
	// Fetch the stock by ID using the repository
	stock, err := s.repo.GetStockByID(stockID)
	if err != nil {
		return fmt.Errorf("failed to find stock: %v", err)
	}

	// Convert the price history array to JSON
	historyJSON, err := json.Marshal(priceHistory)
	if err != nil {
		return fmt.Errorf("failed to encode price history: %v", err)
	}

	// Update the stock's price history
	stock.PriceHistory = string(historyJSON)

	// Save the updated price history using the repository
	return s.repo.UpdateStockPriceHistory(stock)
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

// UpdateStockPrice updates the current price of a stock by its ID.
func (s *StockService) UpdateStockPrice(stockID uint, newPrice float64) error {
	// Fetch the stock by ID using the repository
	stock, err := s.repo.GetStockByID(stockID)
	if err != nil {
		return fmt.Errorf("failed to find stock: %v", err)
	}

	// Update the stock's current price
	stock.CurrentPrice = newPrice

	// Save the updated stock back to the database using the repository
	return s.repo.UpdateStockPrice(stock)
}
