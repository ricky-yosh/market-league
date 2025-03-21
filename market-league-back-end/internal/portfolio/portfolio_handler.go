package portfolio

import (
	"encoding/json"
	"fmt"

	ws "github.com/market-league/api/websocket"
)

// PortfolioHandler Interface
type PortfolioHandlerInterface interface {
	GetPortfolioWithID(conn *ws.Connection, rawData json.RawMessage) error
	GetLeaguePortfolio(conn *ws.Connection, rawData json.RawMessage) error
	GetStocksValueChange(conn *ws.Connection, rawData json.RawMessage) error
	GetPortfolioPointsHistory(conn *ws.Connection, rawData json.RawMessage) error
	CreatePortfolio(conn *ws.Connection, rawData json.RawMessage) error
	AddStockToPortfolio(conn *ws.Connection, rawData json.RawMessage) error
	RemoveStockFromPortfolio(conn *ws.Connection, rawData json.RawMessage) error
}

// Compile-time check
var _ PortfolioHandlerInterface = (*PortfolioHandler)(nil)

// PortfolioHandler defines the HTTP handler for portfolio-related operations.
type PortfolioHandler struct {
	service *PortfolioService
}

// NewPortfolioHandler creates a new instance of PortfolioHandler.
func NewPortfolioHandler(service *PortfolioService) *PortfolioHandler {
	return &PortfolioHandler{service: service}
}

// * Implementation of Interface

// GetPortfolio handles fetching a portfolio by its ID.
func (h *PortfolioHandler) GetPortfolioWithID(conn *ws.Connection, rawData json.RawMessage) error {
	// Step 1: Parse the WebSocket message
	var request struct {
		PortfolioID uint `json:"portfolio_id" binding:"required"`
	}

	// Step 2: Parse data from WebSocket JSON payload
	if err := json.Unmarshal(rawData, &request); err != nil {
		ws.SendError(conn, ws.MessageType_Portfolio_PortfolioWithID, "Invalid input: "+err.Error())
		return fmt.Errorf("invalid input: %v", err)
	}

	// Step 3: Process business logic (reuse the service layer)
	portfolio, err := h.service.GetPortfolioWithID(request.PortfolioID)
	if err != nil {
		ws.SendError(conn, ws.MessageType_Portfolio_PortfolioWithID, err.Error())
		return fmt.Errorf("failed to retrieve portfolio with ID: %v", err)
	}

	// Step 4: Marshal the portfolio into JSON
	portfolioJSON, err := json.Marshal(portfolio)
	if err != nil {
		ws.SendError(conn, ws.MessageType_Portfolio_PortfolioWithID, "Failed to serialize portfolio")
		return fmt.Errorf("serialization error: %v", err)
	}

	// Step 5: Send success response back via WebSocket
	response := ws.WebsocketMessage{
		Type: ws.MessageType_Portfolio_PortfolioWithID,
		Data: json.RawMessage(portfolioJSON), // Use marshaled JSON bytes
	}
	if err := conn.Ws.WriteJSON(response); err != nil {
		return fmt.Errorf("failed to send response: %v", err)
	}

	return nil
}

// GetUserPortfolio handles fetching a user's portfolio in a specific league.
func (h *PortfolioHandler) GetLeaguePortfolio(conn *ws.Connection, rawData json.RawMessage) error {
	// Step 1: Parse the WebSocket message
	var request struct {
		UserID   uint `json:"user_id" binding:"required"`
		LeagueID uint `json:"league_id" binding:"required"`
	}

	// Step 2: Parse data from WebSocket JSON payload
	if err := json.Unmarshal(rawData, &request); err != nil {
		ws.SendError(conn, ws.MessageType_Portfolio_LeaguePortfolio, "Invalid input: "+err.Error())
		return fmt.Errorf("invalid input: %v", err)
	}

	// Step 3: Process business logic (reuse the service layer)
	portfolio, err := h.service.GetLeaguePortfolio(request.UserID, request.LeagueID)
	if err != nil {
		ws.SendError(conn, ws.MessageType_Portfolio_LeaguePortfolio, err.Error())
		return fmt.Errorf("failed to retrieve User's Portfolio from a specific league: %v", err)
	}

	// Step 4: Marshal the portfolio into JSON
	portfolioJSON, err := json.Marshal(portfolio)
	if err != nil {
		ws.SendError(conn, ws.MessageType_Portfolio_LeaguePortfolio, "Failed to serialize portfolio")
		return fmt.Errorf("serialization error: %v", err)
	}

	// Step 5: Send success response back via WebSocket
	response := ws.WebsocketMessage{
		Type: ws.MessageType_Portfolio_LeaguePortfolio,
		Data: json.RawMessage(portfolioJSON), // Use marshaled JSON bytes
	}
	if err := conn.Ws.WriteJSON(response); err != nil {
		return fmt.Errorf("failed to send response: %v", err)
	}

	return nil
}

// CreatePortfolio handles the creation of a new portfolio for a user in a league.
func (h *PortfolioHandler) CreatePortfolio(conn *ws.Connection, rawData json.RawMessage) error {
	// Step 1: Parse the WebSocket message
	var request struct {
		UserID   uint `json:"user_id"`
		LeagueID uint `json:"league_id"`
	}

	// Step 2: Parse data from WebSocket JSON payload
	if err := json.Unmarshal(rawData, &request); err != nil {
		ws.SendError(conn, ws.MessageType_Portfolio_CreatePortfolio, "Invalid input: "+err.Error())
		return fmt.Errorf("invalid input: %v", err)
	}

	// Step 3: Process business logic (reuse the service layer)
	portfolio, err := h.service.CreatePortfolio(request.UserID, request.LeagueID)
	if err != nil {
		ws.SendError(conn, ws.MessageType_Portfolio_CreatePortfolio, err.Error())
		return fmt.Errorf("failed to create portfolio: %v", err)
	}

	// Step 4: Marshal the portfolio into JSON
	portfolioJSON, err := json.Marshal(portfolio)
	if err != nil {
		ws.SendError(conn, ws.MessageType_Portfolio_CreatePortfolio, "Failed to serialize portfolio")
		return fmt.Errorf("serialization error: %v", err)
	}

	// Step 5: Send success response back via WebSocket
	response := ws.WebsocketMessage{
		Type: ws.MessageType_Portfolio_CreatePortfolio,
		Data: json.RawMessage(portfolioJSON), // Use marshaled JSON bytes
	}
	if err := conn.Ws.WriteJSON(response); err != nil {
		return fmt.Errorf("failed to send response: %v", err)
	}

	return nil
}

// AddStockToPortfolio handles adding a stock to a user's portfolio.
func (h *PortfolioHandler) AddStockToPortfolio(conn *ws.Connection, rawData json.RawMessage) error {
	// Step 1: Parse the WebSocket message
	var request struct {
		PortfolioID uint `json:"portfolio_id" binding:"required"`
		StockID     uint `json:"stock_id" binding:"required"`
	}

	// Step 2: Parse data from WebSocket JSON payload
	if err := json.Unmarshal(rawData, &request); err != nil {
		ws.SendError(conn, ws.MessageType_Portfolio_AddStock, "Invalid input: "+err.Error())
		return fmt.Errorf("invalid input: %v", err)
	}

	// Step 3: Process business logic (reuse the service layer)
	if err := h.service.AddStockToPortfolio(request.PortfolioID, request.StockID); err != nil {
		ws.SendError(conn, ws.MessageType_Portfolio_AddStock, err.Error())
		return fmt.Errorf("failed to add stock to portfolio: %v", err)
	}

	// Step 4: Send success response (no data, just confirmation)
	response := ws.WebsocketMessage{
		Type: ws.MessageType_Portfolio_AddStock,
		Data: json.RawMessage(`{"message": "Stock added successfully"}`), // Simple JSON message
	}
	if err := conn.Ws.WriteJSON(response); err != nil {
		return fmt.Errorf("failed to send response: %v", err)
	}

	return nil
}

// RemoveStockFromPortfolio handles removing a stock from a user's portfolio.
func (h *PortfolioHandler) RemoveStockFromPortfolio(conn *ws.Connection, rawData json.RawMessage) error {
	// Step 1: Parse the WebSocket message
	var request struct {
		PortfolioID uint `json:"portfolio_id" binding:"required"`
		StockID     uint `json:"stock_id" binding:"required"`
	}

	// Step 2: Parse data from WebSocket JSON payload
	if err := json.Unmarshal(rawData, &request); err != nil {
		ws.SendError(conn, ws.MessageType_Portfolio_RemoveStock, "Invalid input: "+err.Error())
		return fmt.Errorf("invalid input: %v", err)
	}

	// Step 3: Process business logic (reuse the service layer)
	if err := h.service.RemoveStockFromPortfolio(request.PortfolioID, request.StockID); err != nil {
		ws.SendError(conn, ws.MessageType_Portfolio_RemoveStock, err.Error())
		return fmt.Errorf("failed to remove a stock from a portfolio: %v", err)
	}

	// Step 4: Send success response (no data, just confirmation)
	response := ws.WebsocketMessage{
		Type: ws.MessageType_Portfolio_RemoveStock,
		Data: json.RawMessage(`{"message": "Stock removed successfully"}`), // Simple JSON message
	}
	if err := conn.Ws.WriteJSON(response); err != nil {
		return fmt.Errorf("failed to send response: %v", err)
	}

	return nil
}

// GetPortfolioPointsHistory gets the past history points the portfolio was in and sends it back as a list
func (h *PortfolioHandler) GetPortfolioPointsHistory(conn *ws.Connection, rawData json.RawMessage) error {
	// Step 1: Parse the WebSocket message
	var request struct {
		PortfolioID uint `json:"portfolio_id" binding:"required"`
	}

	// Step 2: Parse data from WebSocket JSON payload
	if err := json.Unmarshal(rawData, &request); err != nil {
		ws.SendError(conn, ws.MessageType_Portfolio_GetPortfolioPointsHistory, "Invalid input: "+err.Error())
		return fmt.Errorf("invalid input: %v", err)
	}

	// Step 3: Process business logic (reuse the service layer)
	portfolio, err := h.service.GetPortfolioPointsHistory(request.PortfolioID)
	if err != nil {
		ws.SendError(conn, ws.MessageType_Portfolio_GetPortfolioPointsHistory, err.Error())
		return fmt.Errorf("failed to retrieve portfolio with ID: %v", err)
	}

	// Step 4: Marshal the portfolio into JSON
	portfolioJSON, err := json.Marshal(portfolio)
	if err != nil {
		ws.SendError(conn, ws.MessageType_Portfolio_GetPortfolioPointsHistory, "Failed to serialize portfolio")
		return fmt.Errorf("serialization error: %v", err)
	}

	// Step 5: Send success response back via WebSocket
	response := ws.WebsocketMessage{
		Type: ws.MessageType_Portfolio_GetPortfolioPointsHistory,
		Data: json.RawMessage(portfolioJSON), // Use marshaled JSON bytes
	}
	if err := conn.Ws.WriteJSON(response); err != nil {
		return fmt.Errorf("failed to send response: %v", err)
	}

	return nil
}

// GetStockValueChange implements PortfolioHandlerInterface.
func (h *PortfolioHandler) GetStocksValueChange(conn *ws.Connection, rawData json.RawMessage) error {
	// Step 1: Parse the WebSocket message
	var request struct {
		PortfolioID uint `json:"portfolio_id" binding:"required"`
	}

	// Step 2: Parse data from WebSocket JSON payload
	if err := json.Unmarshal(rawData, &request); err != nil {
		ws.SendError(conn, ws.MessageType_Portfolio_GetStocksValueChange, "Invalid input: "+err.Error())
		return fmt.Errorf("invalid input: %v", err)
	}

	// Step 3: Process business logic (reuse the service layer)
	portfolio, err := h.service.GetStocksValueChange(request.PortfolioID)
	if err != nil {
		ws.SendError(conn, ws.MessageType_Portfolio_GetStocksValueChange, err.Error())
		return fmt.Errorf("failed to retrieve portfolio with ID: %v", err)
	}

	// Step 4: Marshal the portfolio into JSON
	portfolioJSON, err := json.Marshal(portfolio)
	if err != nil {
		ws.SendError(conn, ws.MessageType_Portfolio_GetStocksValueChange, "Failed to serialize portfolio")
		return fmt.Errorf("serialization error: %v", err)
	}

	// Step 5: Send success response back via WebSocket
	response := ws.WebsocketMessage{
		Type: ws.MessageType_Portfolio_GetStocksValueChange,
		Data: json.RawMessage(portfolioJSON), // Use marshaled JSON bytes
	}
	if err := conn.Ws.WriteJSON(response); err != nil {
		return fmt.Errorf("failed to send response: %v", err)
	}

	return nil
}
