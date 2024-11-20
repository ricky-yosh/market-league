package stock

import (
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

func (s *StockService) CreateStock(tickerSymbol string, companyName string, currentPrice float64) (*models.Stock, error) {
	stock := &models.Stock{
		TickerSymbol: tickerSymbol,
		CompanyName:  companyName,
		CurrentPrice: currentPrice,
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
