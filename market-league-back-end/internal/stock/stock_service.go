package stock

import (
	"errors"
	"time"

	"github.com/market-league/internal/models"
)

// StockService handles business logic related to stocks.
type StockService struct {
	StockRepo *StockRepository // Reference to the repository layer
}

// NewStockService creates a new instance of StockService.
func NewStockService(repo *StockRepository) *StockService {
	return &StockService{StockRepo: repo}
}

func (s *StockService) CreateStock(tickerSymbol string, companyName string) (*models.Stock, error) {
	stock := &models.Stock{
		TickerSymbol: tickerSymbol,
		CompanyName:  companyName,
	}

	err := s.StockRepo.CreateStock(stock)
	if err != nil {
		return nil, err
	}

	return stock, nil
}

func (s *StockService) CreateMultipleStocks(stocks []*models.Stock) error {
	return s.StockRepo.CreateMultipleStocks(stocks)
}

func (s *StockService) UpdateStockPrice(stockID uint, newPrice float64, timestamp *time.Time) error {
	if newPrice < 0 {
		return errors.New("new price must be non-negative")
	}

	if timestamp != nil {
		now := time.Now().UTC()
		providedTime := timestamp.UTC()

		if providedTime.After(now) {
			return errors.New("timestamp cannot be in the future")
		}

		// Optionally, add more validations (e.g., not too far in the past)
	}

	// Add more business logic here if needed

	return s.StockRepo.UpdateCurrentPrice(stockID, newPrice, timestamp)
}

func (s *StockService) GetStockInfo(stockID uint) (models.Stock, error) {
	if stockID == 0 {
		return models.Stock{}, errors.New("invalid stock ID")
	}

	stock, err := s.StockRepo.GetStockWithHistory(stockID)
	if err != nil {
		return models.Stock{}, err
	}

	return stock, nil
}

func (s *StockService) GetAllStocks() ([]*models.Stock, error) {
	// Fetch all stocks from the repository
	stocks, err := s.StockRepo.GetAllStocks()
	if err != nil {
		return nil, err
	}

	// Convert []models.Stock to []*models.Stock
	var stocksPointers []*models.Stock
	for i := range stocks {
		stocksPointers = append(stocksPointers, &stocks[i])
	}

	return stocksPointers, nil
}