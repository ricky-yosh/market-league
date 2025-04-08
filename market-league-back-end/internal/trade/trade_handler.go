package trade

import (
	"encoding/json"
	"fmt"

	ws "github.com/market-league/api/websocket"
)

// UserHandler Interface
type TradeHandlerInterface interface {
	CreateTrade(conn *ws.Connection, rawData json.RawMessage) error
	ConfirmTrade(conn *ws.Connection, rawData json.RawMessage) error
	RefuseTrade(conn *ws.Connection, rawData json.RawMessage) error
	GetTrades(conn *ws.Connection, rawData json.RawMessage) error
}

// Compile-time check
var _ TradeHandlerInterface = (*TradeHandler)(nil)

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

// * Implementation of Interface

// CreateTradeHandler handles the creation of a new trade
func (h *TradeHandler) CreateTrade(conn *ws.Connection, rawData json.RawMessage) error {
	// Step 1: Parse the WebSocket message
	var request struct {
		LeagueID   uint   `json:"league_id"`
		User1ID    uint   `json:"user1_id"`
		User2ID    uint   `json:"user2_id"`
		Stocks1IDs []uint `json:"stocks1_ids"`
		Stocks2IDs []uint `json:"stocks2_ids"`
	}

	// Step 2: Parse data from WebSocket JSON payload
	if err := json.Unmarshal(rawData, &request); err != nil {
		ws.SendError(conn, ws.MessageType_Trade_CreateTrade, "Invalid input: "+err.Error())
		return fmt.Errorf("invalid input: %v", err)
	}

	// Step 3: Process business logic (reuse the service layer)
	trade, err := h.TradeService.CreateTrade(request.LeagueID, request.User1ID, request.User2ID, request.Stocks1IDs, request.Stocks2IDs)
	if err != nil {
		ws.SendError(conn, ws.MessageType_Trade_CreateTrade, err.Error())
		return fmt.Errorf("failed to retrieve portfolio with ID: %v", err)
	}

	// Step 4: Marshal the portfolio into JSON
	portfolioJSON, err := json.Marshal(trade)
	if err != nil {
		ws.SendError(conn, ws.MessageType_Trade_CreateTrade, "Failed to serialize portfolio")
		return fmt.Errorf("serialization error: %v", err)
	}

	// Step 5: Send success response back via WebSocket
	response := ws.WebsocketMessage{
		Type: ws.MessageType_Trade_CreateTrade,
		Data: json.RawMessage(portfolioJSON), // Use marshaled JSON bytes
	}
	if err := conn.Ws.WriteJSON(response); err != nil {
		return fmt.Errorf("failed to send response: %v", err)
	}

	return nil
}

// GetTrades handles the retrieval of all trades for a given League with the option of specifying a user
func (h *TradeHandler) GetTrades(conn *ws.Connection, rawData json.RawMessage) error {
	// Step 1: Parse the WebSocket message
	var request struct {
		UserID         *uint `json:"user_id"`         // Optional User ID
		LeagueID       uint  `json:"league_id"`       // Required League ID
		ReceivingTrade *bool `json:"receiving_trade"` // Optional: Filter for receiving trades
		SendingTrade   *bool `json:"sending_trade"`   // Optional: Filter for sending trades
	}

	// Step 2: Parse data from WebSocket JSON payload
	if err := json.Unmarshal(rawData, &request); err != nil {
		ws.SendError(conn, ws.MessageType_Trade_GetTrades, "Invalid input: "+err.Error())
		return fmt.Errorf("invalid input: %v", err)
	}

	// Step 3: Process business logic (reuse the service layer)
	trades, err := h.TradeService.GetTrades(request.LeagueID, request.UserID, request.ReceivingTrade, request.SendingTrade)
	if err != nil {
		ws.SendError(conn, ws.MessageType_Trade_GetTrades, err.Error())
		return fmt.Errorf("failed to retrieve portfolio with ID: %v", err)
	}

	// Step 4: Marshal the portfolio into JSON
	portfolioJSON, err := json.Marshal(trades)
	if err != nil {
		ws.SendError(conn, ws.MessageType_Trade_GetTrades, "Failed to serialize portfolio")
		return fmt.Errorf("serialization error: %v", err)
	}

	// Step 5: Send success response back via WebSocket
	response := ws.WebsocketMessage{
		Type: ws.MessageType_Trade_GetTrades,
		Data: json.RawMessage(portfolioJSON), // Use marshaled JSON bytes
	}
	if err := conn.Ws.WriteJSON(response); err != nil {
		return fmt.Errorf("failed to send response: %v", err)
	}

	return nil

}

// ConfirmTrade handles the confirmation of a trade
func (h *TradeHandler) ConfirmTrade(conn *ws.Connection, rawData json.RawMessage) error {
	// Step 1: Parse the WebSocket message
	var request struct {
		TradeID uint `json:"trade_id" binding:"required"`
		UserID  uint `json:"user_id" binding:"required"`
	}

	// Step 2: Parse data from WebSocket JSON payload
	if err := json.Unmarshal(rawData, &request); err != nil {
		ws.SendError(conn, ws.MessageType_Trade_ConfirmTrade, "Invalid input: "+err.Error())
		return fmt.Errorf("invalid input: %v", err)
	}

	// Step 3: Process business logic (reuse the service layer)
	if err := h.TradeService.ConfirmTrade(request.TradeID, request.UserID); err != nil {
		ws.SendError(conn, ws.MessageType_Trade_ConfirmTrade, err.Error())
		return fmt.Errorf("failed to confirm trade: %v", err)
	}

	// Step 4: Send success response (no data, just confirmation)
	response := ws.WebsocketMessage{
		Type: ws.MessageType_Trade_ConfirmTrade,
		Data: json.RawMessage(`{"message": "Trade confirmed successfully"}`), // Simple JSON message
	}
	if err := conn.Ws.WriteJSON(response); err != nil {
		return fmt.Errorf("failed to send response: %v", err)
	}

	return nil
}
// RefuseTrade handles the confirmation of a trade
func (h *TradeHandler) RefuseTrade(conn *ws.Connection, rawData json.RawMessage) error {
	// Step 1: Parse the WebSocket message
	var request struct {
		TradeID uint `json:"trade_id" binding:"required"`
		UserID  uint `json:"user_id" binding:"required"`
	}

	// Step 2: Parse data from WebSocket JSON payload
	if err := json.Unmarshal(rawData, &request); err != nil {
		ws.SendError(conn, ws.MessageType_Trade_RefuseTrade, "Invalid input: "+err.Error())
		return fmt.Errorf("invalid input: %v", err)
	}

	// Step 3: Process business logic (reuse the service layer)
	if err := h.TradeService.RefuseTrade(request.TradeID, request.UserID); err != nil {
		ws.SendError(conn, ws.MessageType_Trade_RefuseTrade, err.Error())
		return fmt.Errorf("failed to refuse trade: %v", err)
	}

	// Step 4: Send success response (no data, just confirmation)
	response := ws.WebsocketMessage{
		Type: ws.MessageType_Trade_RefuseTrade,
		Data: json.RawMessage(`{"message": "Trade refused successfully"}`), // Simple JSON message
	}
	if err := conn.Ws.WriteJSON(response); err != nil {
		return fmt.Errorf("failed to send response: %v", err)
	}

	return nil
}
