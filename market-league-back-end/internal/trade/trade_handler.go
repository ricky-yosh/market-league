package trade

import (
	"net/http"
	"strconv"

	"github.com/gin-gonic/gin"
	"github.com/market-league/internal/models"
)

// TradeHandler defines the HTTP handler for trade-related operations.
type TradeHandler struct {
	service *TradeService
}

// NewTradeHandler creates a new instance of TradeHandler.
func NewTradeHandler(service *TradeService) *TradeHandler {
	return &TradeHandler{service: service}
}

// CreateTrade handles creating a new trade.
func (h *TradeHandler) CreateTrade(c *gin.Context) {
	var trade models.Trade

	// Bind JSON data to the trade model
	if err := c.ShouldBindJSON(&trade); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	// Create the trade using the service
	if err := h.service.CreateTrade(&trade); err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to create trade"})
		return
	}

	c.JSON(http.StatusCreated, trade)
}

// GetTradesByUser handles fetching all trades made by a specific user.
func (h *TradeHandler) GetTradesByUser(c *gin.Context) {
	// Parse the user ID from the URL parameter
	userID, err := strconv.ParseUint(c.Param("userID"), 10, 64)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid user ID"})
		return
	}

	// Fetch trades by the user using the service
	trades, err := h.service.GetTradesByUser(uint(userID))
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to fetch trades for user"})
		return
	}

	c.JSON(http.StatusOK, trades)
}

// GetTradesByPortfolio handles fetching all trades related to a specific portfolio.
func (h *TradeHandler) GetTradesByPortfolio(c *gin.Context) {
	// Parse the portfolio ID from the URL parameter
	portfolioID, err := strconv.ParseUint(c.Param("portfolioID"), 10, 64)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid portfolio ID"})
		return
	}

	// Fetch trades by portfolio using the service
	trades, err := h.service.GetTradesByPortfolio(uint(portfolioID))
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to fetch trades for portfolio"})
		return
	}

	c.JSON(http.StatusOK, trades)
}
