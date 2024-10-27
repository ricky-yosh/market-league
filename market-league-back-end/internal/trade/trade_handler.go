package trade

import (
	"net/http"

	"github.com/gin-gonic/gin"
)

// TradeHandler defines the HTTP handler for trade-related operations.
type TradeHandler struct {
	service *TradeService
}

// NewTradeHandler creates a new instance of TradeHandler.
func NewTradeHandler(service *TradeService) *TradeHandler {
	return &TradeHandler{service: service}
}

// CreateTrade creates a new trade between two players within a league.
func (h *TradeHandler) CreateTrade(c *gin.Context) {
	var request struct {
		LeagueID           uint   `json:"league_id" binding:"required"`
		Player1ID          uint   `json:"player1_id" binding:"required"`
		Player2ID          uint   `json:"player2_id" binding:"required"`
		Player1PortfolioID uint   `json:"player1_portfolio_id" binding:"required"`
		Player2PortfolioID uint   `json:"player2_portfolio_id" binding:"required"`
		Player1Stocks      []uint `json:"player1_stocks" binding:"required"` // List of stock IDs offered by Player 1
		Player2Stocks      []uint `json:"player2_stocks" binding:"required"` // List of stock IDs offered by Player 2
	}

	// Bind the request data to the struct
	if err := c.ShouldBindJSON(&request); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	// Call the service to create the trade
	err := h.service.CreateTrade(
		request.LeagueID,
		request.Player1ID,
		request.Player2ID,
		request.Player1PortfolioID,
		request.Player2PortfolioID,
		request.Player1Stocks,
		request.Player2Stocks,
	)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	c.JSON(http.StatusCreated, gin.H{"message": "Trade created successfully"})
}

// ConfirmTrade confirms a trade for a player.
func (h *TradeHandler) ConfirmTrade(c *gin.Context) {
	var request struct {
		TradeID  uint `json:"trade_id" binding:"required"`
		PlayerID uint `json:"player_id" binding:"required"`
	}

	// Bind the request data to the struct
	if err := c.ShouldBindJSON(&request); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	// Call the service to confirm the trade
	err := h.service.ConfirmTrade(request.TradeID, request.PlayerID)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	c.JSON(http.StatusOK, gin.H{"message": "Trade confirmed successfully"})
}

// GetTrades fetches trades based on filter criteria.
func (h *TradeHandler) GetTrades(c *gin.Context) {
	var request struct {
		PortfolioID       uint `json:"portfolio_id"`        // Portfolio ID to filter by, if applicable
		LeagueID          uint `json:"league_id"`           // League ID to filter by, if applicable
		FilterByPortfolio bool `json:"filter_by_portfolio"` // Whether to filter by portfolio
		FilterByLeague    bool `json:"filter_by_league"`    // Whether to filter by league
	}

	// Bind the JSON request data to the struct
	if err := c.ShouldBindJSON(&request); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	// Call the service to get the trades based on the filter criteria
	trades, err := h.service.GetTrades(request.PortfolioID, request.LeagueID, request.FilterByPortfolio, request.FilterByLeague)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to fetch trades"})
		return
	}

	c.JSON(http.StatusOK, trades)
}
