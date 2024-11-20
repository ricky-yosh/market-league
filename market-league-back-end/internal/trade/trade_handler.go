package trade

import (
	"errors"
	"net/http"

	"github.com/gin-gonic/gin"
	"gorm.io/gorm"
)

// TradeHandler handles HTTP requests for trades
type TradeHandler struct {
	TradeService *TradeService
}

// NewTradeHandler creates a new instance of TradeHandler
func NewTradeHandler(tradeService *TradeService) *TradeHandler {
	return &TradeHandler{
		TradeService: tradeService,
	}
}

// CreateTradeHandler handles the creation of a new trade
func (h *TradeHandler) CreateTrade(c *gin.Context) {
	var request struct {
		LeagueID   uint   `json:"league_id"`
		User1ID    uint   `json:"user1_id"`
		User2ID    uint   `json:"user2_id"`
		Stocks1IDs []uint `json:"stocks1_ids"`
		Stocks2IDs []uint `json:"stocks2_ids"`
	}

	if err := c.ShouldBindJSON(&request); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	trade, err := h.TradeService.CreateTrade(request.LeagueID, request.User1ID, request.User2ID, request.Stocks1IDs, request.Stocks2IDs)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	c.JSON(http.StatusOK, trade)
}

func (h *TradeHandler) GetTrades(c *gin.Context) {
	var request struct {
		UserID         *uint `json:"user_id"`         // Optional User ID
		LeagueID       uint  `json:"league_id"`       // Required League ID
		ReceivingTrade *bool `json:"receiving_trade"` // Optional: Filter for receiving trades
		SendingTrade   *bool `json:"sending_trade"`   // Optional: Filter for sending trades
	}

	// Parse the JSON request
	if err := c.ShouldBindJSON(&request); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid request payload"})
		return
	}

	// Ensure LeagueID is always provided
	if request.LeagueID == 0 {
		c.JSON(http.StatusBadRequest, gin.H{"error": "league_id is required"})
		return
	}

	// Call the service to fetch trades
	trades, err := h.TradeService.GetTrades(request.LeagueID, request.UserID, request.ReceivingTrade, request.SendingTrade)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	c.JSON(http.StatusOK, trades)
}

// ConfirmTrade handles the confirmation of a trade
func (h *TradeHandler) ConfirmTrade(c *gin.Context) {
	var req struct {
		TradeID uint `json:"trade_id" binding:"required"`
		UserID  uint `json:"user_id" binding:"required"`
	}

	// ConfirmTradeResponse represents the response after confirming a trade
	type ConfirmTradeResponse struct {
		Success bool   `json:"success"`
		Message string `json:"message"`
	}

	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, ConfirmTradeResponse{
			Success: false,
			Message: "Invalid request payload",
		})
		return
	}

	// Call the service to confirm the trade
	if err := h.TradeService.ConfirmTrade(req.TradeID, req.UserID); err != nil {
		// Determine the appropriate status code based on the error
		var statusCode int
		switch {
		case errors.Is(err, gorm.ErrRecordNotFound) || err.Error() == "trade not found":
			statusCode = http.StatusNotFound
		case err.Error() == "trade is already confirmed" ||
			err.Error() == "user1 has already confirmed this trade" ||
			err.Error() == "user2 has already confirmed this trade" ||
			err.Error() == "user is not part of this trade":
			statusCode = http.StatusBadRequest
		default:
			statusCode = http.StatusInternalServerError
		}

		c.JSON(statusCode, ConfirmTradeResponse{
			Success: false,
			Message: err.Error(),
		})
		return
	}

	// Success response
	c.JSON(http.StatusOK, ConfirmTradeResponse{
		Success: true,
		Message: "Trade confirmed successfully",
	})
}
