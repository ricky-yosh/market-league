package user

import (
	"encoding/json"
	"fmt"

	ws "github.com/market-league/api/websocket"
)

// UserHandler Interface
type UserHandlerInterface interface {
	GetUserByID(conn *ws.Connection, rawData json.RawMessage) error
	GetUserLeagues(conn *ws.Connection, rawData json.RawMessage) error
	GetUserTrades(conn *ws.Connection, rawData json.RawMessage) error
	GetUserPortfolios(conn *ws.Connection, rawData json.RawMessage) error
}

// Compile-time check
var _ UserHandlerInterface = (*UserHandler)(nil)

// UserHandler defines the HTTP handler for user-related operations.
type UserHandler struct {
	UserService *UserService
}

// NewUserHandler creates a new instance of UserHandler.
func NewUserHandler(service *UserService) *UserHandler {
	return &UserHandler{UserService: service}
}

// * Implementation of Interface

// GetUserByID fetches user information based on filter criteria.
func (h *UserHandler) GetUserByID(conn *ws.Connection, rawData json.RawMessage) error {
	// Step 1: Parse the WebSocket message
	var request struct {
		UserID uint `json:"user_id" binding:"required"` // User ID to fetch information for
	}

	// Step 2: Parse data from WebSocket JSON payload
	if err := json.Unmarshal(rawData, &request); err != nil {
		ws.SendError(conn, ws.MessageType_User_UserInfo, "Invalid input: "+err.Error())
		return fmt.Errorf("invalid input: %v", err)
	}

	// Step 3: Process business logic (reuse the service layer)
	userInfo, err := h.UserService.GetUserByID(request.UserID)
	if err != nil {
		ws.SendError(conn, ws.MessageType_User_UserInfo, err.Error())
		return fmt.Errorf("failed to retrieve portfolio with ID: %v", err)
	}

	// Step 4: Marshal the portfolio into JSON
	userInfoJSON, err := json.Marshal(userInfo)
	if err != nil {
		ws.SendError(conn, ws.MessageType_User_UserInfo, "Failed to serialize portfolio")
		return fmt.Errorf("serialization error: %v", err)
	}

	// Step 5: Send success response back via WebSocket
	response := ws.WebsocketMessage{
		Type: ws.MessageType_User_UserInfo,
		Data: json.RawMessage(userInfoJSON), // Use marshaled JSON bytes
	}
	if err := conn.Ws.WriteJSON(response); err != nil {
		return fmt.Errorf("failed to send response: %v", err)
	}

	return nil
}

// GetUserLeagues handles requests to retrieve leagues that a user is a member of.
func (h *UserHandler) GetUserLeagues(conn *ws.Connection, rawData json.RawMessage) error {
	// Step 1: Parse the WebSocket message
	var request struct {
		UserID uint `json:"user_id" binding:"required"`
	}

	// Step 2: Parse data from WebSocket JSON payload
	if err := json.Unmarshal(rawData, &request); err != nil {
		ws.SendError(conn, ws.MessageType_User_UserLeagues, "Invalid input: "+err.Error())
		return fmt.Errorf("invalid input: %v", err)
	}

	// Step 3: Process business logic (reuse the service layer)
	leagues, err := h.UserService.GetUserLeagues(request.UserID)
	if err != nil {
		ws.SendError(conn, ws.MessageType_User_UserLeagues, err.Error())
		return fmt.Errorf("failed to retrieve portfolio with ID: %v", err)
	}

	// Step 4: Marshal the portfolio into JSON
	leaguesJSON, err := json.Marshal(leagues)
	if err != nil {
		ws.SendError(conn, ws.MessageType_User_UserLeagues, "Failed to serialize portfolio")
		return fmt.Errorf("serialization error: %v", err)
	}

	// Step 5: Send success response back via WebSocket
	response := ws.WebsocketMessage{
		Type: ws.MessageType_User_UserLeagues,
		Data: json.RawMessage(leaguesJSON), // Use marshaled JSON bytes
	}
	if err := conn.Ws.WriteJSON(response); err != nil {
		return fmt.Errorf("failed to send response: %v", err)
	}

	return nil
}

// GetUserTrades handles requests to retrieve trades that a user is involved in within a specific league.
func (h *UserHandler) GetUserTrades(conn *ws.Connection, rawData json.RawMessage) error {
	// Step 1: Parse the WebSocket message
	var request struct {
		UserID   uint `json:"user_id" binding:"required"`
		LeagueID uint `json:"league_id" binding:"required"`
	}

	// Step 2: Parse data from WebSocket JSON payload
	if err := json.Unmarshal(rawData, &request); err != nil {
		ws.SendError(conn, ws.MessageType_User_UserTrades, "Invalid input: "+err.Error())
		return fmt.Errorf("invalid input: %v", err)
	}

	// Step 3: Process business logic (reuse the service layer)
	trades, err := h.UserService.GetUserTrades(request.UserID, request.LeagueID)
	if err != nil {
		ws.SendError(conn, ws.MessageType_User_UserTrades, err.Error())
		return fmt.Errorf("failed to retrieve portfolio with ID: %v", err)
	}

	// Step 4: Marshal the portfolio into JSON
	portfolioJSON, err := json.Marshal(trades)
	if err != nil {
		ws.SendError(conn, ws.MessageType_User_UserTrades, "Failed to serialize portfolio")
		return fmt.Errorf("serialization error: %v", err)
	}

	// Step 5: Send success response back via WebSocket
	response := ws.WebsocketMessage{
		Type: ws.MessageType_User_UserTrades,
		Data: json.RawMessage(portfolioJSON), // Use marshaled JSON bytes
	}
	if err := conn.Ws.WriteJSON(response); err != nil {
		return fmt.Errorf("failed to send response: %v", err)
	}

	return nil
}

// GetUserPortfolios handles requests to retrieve portfolios that a user is associated with.
func (h *UserHandler) GetUserPortfolios(conn *ws.Connection, rawData json.RawMessage) error {
	// Step 1: Parse the WebSocket message
	var request struct {
		UserID uint `json:"user_id" binding:"required"`
	}

	// Step 2: Parse data from WebSocket JSON payload
	if err := json.Unmarshal(rawData, &request); err != nil {
		ws.SendError(conn, ws.MessageType_User_UserPortfolios, "Invalid input: "+err.Error())
		return fmt.Errorf("invalid input: %v", err)
	}

	// Step 3: Process business logic (reuse the service layer)
	portfolios, err := h.UserService.GetUserPortfolios(request.UserID)
	if err != nil {
		ws.SendError(conn, ws.MessageType_User_UserPortfolios, err.Error())
		return fmt.Errorf("failed to retrieve portfolio with ID: %v", err)
	}

	// Step 4: Marshal the portfolio into JSON
	portfolioJSON, err := json.Marshal(portfolios)
	if err != nil {
		ws.SendError(conn, ws.MessageType_User_UserPortfolios, "Failed to serialize portfolio")
		return fmt.Errorf("serialization error: %v", err)
	}

	// Step 5: Send success response back via WebSocket
	response := ws.WebsocketMessage{
		Type: ws.MessageType_User_UserPortfolios,
		Data: json.RawMessage(portfolioJSON), // Use marshaled JSON bytes
	}
	if err := conn.Ws.WriteJSON(response); err != nil {
		return fmt.Errorf("failed to send response: %v", err)
	}

	return nil
}
