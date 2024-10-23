package stock

import (
	"net/http"
	"strconv"

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

// GetPrice handles fetching the current price of a stock by its ID.
func (h *StockHandler) GetPrice(c *gin.Context) {
	// Parse the stock ID from the URL parameter
	stockID, err := strconv.ParseUint(c.Param("stockID"), 10, 64)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid stock ID"})
		return
	}

	// Get the stock price using the service
	price, err := h.service.GetPrice(uint(stockID))
	if err != nil {
		c.JSON(http.StatusNotFound, gin.H{"error": "Stock not found"})
		return
	}

	c.JSON(http.StatusOK, gin.H{"stock_id": stockID, "price": price})
}

// GetPriceHistory handles fetching the price history of a stock by its ID.
func (h *StockHandler) GetPriceHistory(c *gin.Context) {
	// Parse the stock ID from the URL parameter
	stockID, err := strconv.ParseUint(c.Param("stockID"), 10, 64)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid stock ID"})
		return
	}

	// Get the stock's price history using the service
	history, err := h.service.GetPriceHistory(uint(stockID))
	if err != nil {
		c.JSON(http.StatusNotFound, gin.H{"error": "Stock not found or failed to retrieve price history"})
		return
	}

	c.JSON(http.StatusOK, gin.H{"stock_id": stockID, "price_history": history})
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

// UpdateStockPrice handles updating the current price of a stock.
func (h *StockHandler) UpdateStockPrice(c *gin.Context) {
	// Parse the stock ID from the URL parameter
	stockID, err := strconv.ParseUint(c.Param("stockID"), 10, 64)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid stock ID"})
		return
	}

	// Bind JSON data to get the new price
	var request struct {
		NewPrice float64 `json:"new_price" binding:"required"`
	}
	if err := c.ShouldBindJSON(&request); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	// Update the stock price using the service
	err = h.service.UpdateCurrentPrice(uint(stockID), request.NewPrice)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to update stock price"})
		return
	}

	c.JSON(http.StatusOK, gin.H{"message": "Stock price updated successfully"})
}
