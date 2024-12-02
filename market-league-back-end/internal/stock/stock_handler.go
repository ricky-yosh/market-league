package stock

import (
	"net/http"
	"time"

	"github.com/gin-gonic/gin"
	"github.com/market-league/internal/models"
)

// StockHandler defines the HTTP handler for stock-related operations.
type StockHandler struct {
	StockService *StockService
}

// NewStockHandler creates a new instance of StockHandler.
func NewStockHandler(service *StockService) *StockHandler {
	return &StockHandler{StockService: service}
}

type CreateStockRequest struct {
	TickerSymbol string  `json:"ticker_symbol" binding:"required"`
	CompanyName  string  `json:"company_name" binding:"required"`
	CurrentPrice float64 `json:"current_price" binding:"required,gt=0"`
}

// CreateStockResponse represents the response after creating a stock
type CreateStockResponse struct {
	Success bool          `json:"success"`
	Message string        `json:"message"`
	Data    *models.Stock `json:"data,omitempty"`
}

type CreateMultipleStocksResponse struct {
	Success bool           `json:"success"`
	Message string         `json:"message"`
	Data    []models.Stock `json:"data,omitempty"`
}

// CreateStock handles the creation of a new stock
func (h *StockHandler) CreateStock(c *gin.Context) {
	var req CreateStockRequest

	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, CreateStockResponse{
			Success: false,
			Message: "Invalid request payload",
		})
		return
	}

	stock, err := h.StockService.CreateStock(req.TickerSymbol, req.CompanyName, req.CurrentPrice)
	if err != nil {
		// Handle unique constraint violation for TickerSymbol
		if isUniqueConstraintError(err, "ticker_symbol") {
			c.JSON(http.StatusConflict, CreateStockResponse{
				Success: false,
				Message: "Ticker symbol already exists",
			})
			return
		}

		// Handle other errors
		c.JSON(http.StatusInternalServerError, CreateStockResponse{
			Success: false,
			Message: "Failed to create stock",
		})
		return
	}

	c.JSON(http.StatusCreated, CreateStockResponse{
		Success: true,
		Message: "Stock created successfully",
		Data:    stock,
	})
}

// Helper function to detect unique constraint errors
func isUniqueConstraintError(err error, field string) bool {
	// This function needs to be implemented based on your database driver and error handling
	// For PostgreSQL with lib/pq, you can check for pq.Error and the specific constraint name
	// Here's a generic placeholder:
	return false
}

type CreateMultipleStocksRequest []CreateStockRequest

// CreateMultipleStocks handles the creation of multiple stocks
func (h *StockHandler) CreateMultipleStocks(c *gin.Context) {
	var req CreateMultipleStocksRequest

	// Bind JSON input to the array of CreateStockRequest structs
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, CreateMultipleStocksResponse{
			Success: false,
			Message: "Invalid input",
		})
		return
	}

	// Convert requests to models.Stock
	var stocks []*models.Stock
	for _, stockReq := range req {
		stock := &models.Stock{
			TickerSymbol: stockReq.TickerSymbol,
			CompanyName:  stockReq.CompanyName,
			CurrentPrice: stockReq.CurrentPrice,
		}
		stocks = append(stocks, stock)
	}

	// Call the service to create multiple stocks
	err := h.StockService.CreateMultipleStocks(stocks)
	if err != nil {
		// Handle specific errors if needed
		c.JSON(http.StatusInternalServerError, CreateMultipleStocksResponse{
			Success: false,
			Message: "Failed to create stocks",
		})
		return
	}

	// Return success response with created stocks
	c.JSON(http.StatusOK, CreateMultipleStocksResponse{
		Success: true,
		Message: "All stocks successfully created",
		Data:    extractStocksData(stocks),
	})
}

// Helper function to extract necessary data from stocks
func extractStocksData(stocks []*models.Stock) []models.Stock {
	var result []models.Stock
	for _, stock := range stocks {
		result = append(result, *stock)
	}
	return result
}

// UpdatePriceRequest represents the expected payload for updating stock price
type UpdatePriceRequest struct {
	StockID   uint       `json:"stock_id" binding:"required"`
	NewPrice  float64    `json:"new_price" binding:"required"`
	Timestamp *time.Time `json:"timestamp,omitempty"`
}

// UpdatePrice handles the request to update a stock's current price
func (h *StockHandler) UpdatePrice(c *gin.Context) {
	var req UpdatePriceRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{
			"success": false,
			"message": "Invalid request payload",
		})
		return
	}

	if err := h.StockService.UpdateStockPrice(req.StockID, req.NewPrice, req.Timestamp); err != nil {
		// You can customize error responses based on error types
		c.JSON(http.StatusInternalServerError, gin.H{
			"success": false,
			"message": err.Error(),
		})
		return
	}

	c.JSON(http.StatusOK, gin.H{
		"success": true,
		"message": "Stock price updated successfully",
	})
}

type GetStockInfoRequest struct {
	StockID uint `json:"stock_id" binding:"required"`
}

type GetStockInfoResponse struct {
	Success bool   `json:"success"`
	Message string `json:"message,omitempty"`
}

type StockInfo struct {
	ID             uint              `json:"id"`
	TickerSymbol   string            `json:"ticker_symbol"`
	CompanyName    string            `json:"company_name"`
	CurrentPrice   float64           `json:"current_price"`
	PriceHistories []PriceHistoryDTO `json:"price_histories"`
}

type PriceHistoryDTO struct {
	ID        uint      `json:"id"`
	Price     float64   `json:"price"`
	Timestamp time.Time `json:"timestamp"`
}

func (h *StockHandler) GetStockInfo(c *gin.Context) {
	var req GetStockInfoRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, GetStockInfoResponse{
			Success: false,
			Message: "Invalid request payload",
		})
		return
	}

	stock, err := h.StockService.GetStockInfo(req.StockID)
	if err != nil {
		c.JSON(http.StatusInternalServerError, GetStockInfoResponse{
			Success: false,
			Message: err.Error(),
		})
		return
	}

	// Map models.Stock to StockInfo DTO
	stockInfo := StockInfo{
		ID:           stock.ID,
		TickerSymbol: stock.TickerSymbol,
		CompanyName:  stock.CompanyName,
		CurrentPrice: stock.CurrentPrice,
		PriceHistories: func(histories []models.PriceHistory) []PriceHistoryDTO {
			dto := make([]PriceHistoryDTO, len(histories))
			for i, ph := range histories {
				dto[i] = PriceHistoryDTO{
					ID:        ph.ID,
					Price:     ph.Price,
					Timestamp: ph.Timestamp,
				}
			}
			return dto
		}(stock.PriceHistories),
	}

	c.JSON(http.StatusOK, stockInfo)
}
