package stock

import (
	"net/http"

	"github.com/gin-gonic/gin"
	"github.com/market-league/internal/models"
)

// StockHandler defines the HTTP handler for stock-related operations.
type StockHandler struct {
	service *StockService
}

// NewStockHandler creates a new instance of StockHandler.
func NewStockHandler(service *StockService) *StockHandler {
	return &StockHandler{service: service}
}

// GetPrice fetches the current price of a stock by its ID.
func (h *StockHandler) GetPrice(c *gin.Context) {
	var request struct {
		StockID uint `json:"stock_id" binding:"required"`
	}

	// Bind the request data to the struct
	if err := c.ShouldBindJSON(&request); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	// Call the service to get the stock price
	price, err := h.service.GetPrice(request.StockID)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	c.JSON(http.StatusOK, gin.H{"price": price})
}

// GetPriceHistory fetches the price history of a stock by its ID.
func (h *StockHandler) GetPriceHistory(c *gin.Context) {
	var request struct {
		StockID uint `json:"stock_id" binding:"required"`
	}

	// Bind the request data to the struct
	if err := c.ShouldBindJSON(&request); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	// Call the service to get the stock price history
	priceHistory, err := h.service.GetPriceHistory(request.StockID)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	c.JSON(http.StatusOK, gin.H{"price_history": priceHistory})
}

// UpdatePriceHistory updates the price history of a stock.
func (h *StockHandler) UpdatePriceHistory(c *gin.Context) {
	var request struct {
		StockID      uint      `json:"stock_id" binding:"required"`
		PriceHistory []float64 `json:"price_history" binding:"required"`
	}

	// Bind the request data to the struct
	if err := c.ShouldBindJSON(&request); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	// Call the service to update the price history
	err := h.service.UpdatePriceHistory(request.StockID, request.PriceHistory)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	c.JSON(http.StatusOK, gin.H{"message": "Price history updated successfully"})
}

// CreateStock handles creating a new stock (if needed for administrative purposes).
func (h *StockHandler) CreateStock(c *gin.Context) {
	var stock models.Stock

	// Bind JSON data to the stock model
	if err := c.ShouldBindJSON(&stock); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	// Create the stock using the service
	if err := h.service.CreateStock(&stock); err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to create stock"})
		return
	}

	c.JSON(http.StatusCreated, stock)
}

func (h *StockHandler) CreateMultipleStocks(c *gin.Context) {
	var stocks []models.Stock

	// Bind JSON input to the array of stock structs
	if err := c.ShouldBindJSON(&stocks); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid input"})
		return
	}

	// Iterate over the stocks and create each one
	for _, stock := range stocks {
		if err := h.service.CreateStock(&stock); err != nil {
			c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to create stock"})
			return
		}
	}

	// Return success response
	c.JSON(http.StatusOK, gin.H{"message": "All stocks successfully created"})
}

// UpdateStockPrice updates the current price of a stock by its ID.
func (h *StockHandler) UpdateStockPrice(c *gin.Context) {
	var request struct {
		StockID  uint    `json:"stock_id" binding:"required"`
		NewPrice float64 `json:"new_price" binding:"required"`
	}

	// Bind the request data to the struct
	if err := c.ShouldBindJSON(&request); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	// Call the service to update the stock price
	err := h.service.UpdateStockPrice(request.StockID, request.NewPrice)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	c.JSON(http.StatusOK, gin.H{"message": "Stock price updated successfully"})
}
