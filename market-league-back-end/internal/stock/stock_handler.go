package stock

import (
	"net/http"

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
