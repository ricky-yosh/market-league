package ws

import (
	"encoding/json"
	"log"

	"github.com/gin-gonic/gin"
	"github.com/gorilla/websocket"
)

type WebsocketMessage struct {
	Type string          `json:"type"` // Message type for routing
	Data json.RawMessage `json:"data"` // Raw JSON data for flexibility
}

const (
	// Portfolio Routes
	MessageType_Portfolio_CreatePortfolio = "MessageType_Portfolio_CreatePortfolio"
	MessageType_Portfolio_PortfolioWithID = "MessageType_Portfolio_PortfolioWithID"
	MessageType_Portfolio_LeaguePortfolio = "MessageType_Portfolio_LeaguePortfolio"
	MessageType_Portfolio_AddStock        = "MessageType_Portfolio_AddStock"
	MessageType_Portfolio_RemoveStock     = "MessageType_Portfolio_RemoveStock"

	// Stock Routes
	MessageType_Stock_CreateStock             = "MessageType_Stock_CreateStock"
	MessageType_Stock_CreateMultipleStocks    = "MessageType_Stock_CreateMultipleStocks"
	MessageType_Stock_GetStockInformation     = "MessageType_Stock_GetStockInformation"
	MessageType_Stock_UpdateCurrentStockPrice = "MessageType_Stock_UpdateCurrentStockPrice"
	MessageType_Stock_GetAllStocks			  = "MessageType_Stock_GetAllStocks"

	// User Routes
	MessageType_User_UserInfo       = "MessageType_User_UserInfo"
	MessageType_User_UserLeagues    = "MessageType_User_UserLeagues"
	MessageType_User_UserTrades     = "MessageType_User_UserTrades"
	MessageType_User_UserPortfolios = "MessageType_User_UserPortfolios"

	// Trade Routes
	MessageType_Trade_CreateTrade  = "MessageType_Trade_CreateTrade"
	MessageType_Trade_ConfirmTrade = "MessageType_Trade_ConfirmTrade"
	MessageType_Trade_GetTrades    = "MessageType_Trade_GetTrades"

	// League Portfolio Routes
	MessageType_LeaguePortfolio_DraftStock             = "MessageType_LeaguePortfolio_DraftStock"
	MessageType_LeaguePortfolio_GetLeaguePortfolioInfo = "MessageType_LeaguePortfolio_GetLeaguePortfolioInfo"

	// League Routes
	MessageType_League_CreateLeague    = "MessageType_League_CreateLeague"
	MessageType_League_RemoveLeague    = "MessageType_League_RemoveLeague"
	MessageType_League_AddUserToLeague = "MessageType_League_AddUserToLeague"
	MessageType_League_GetDetails      = "MessageType_League_GetDetails"
	MessageType_League_GetLeaderboard  = "MessageType_League_GetLeaderboard"

	// Error Message
	MessageType_Error = "MessageType_Error"
)

func SendError(conn *websocket.Conn, messageType string, errorMsg string) {
	// Construct the error response using gin.H
	errorData := gin.H{
		"type":    MessageType_Error, // Include the error type
		"message": errorMsg,          // Include the error message
	}

	// Marshal the error data into JSON bytes
	errorJSON, err := json.Marshal(errorData)
	if err != nil {
		log.Println("Failed to serialize error data:", err)
		return
	}
	// Create an error response using the Transmission struct
	errorResponse := WebsocketMessage{
		Type: messageType, // Message type is 'error'
		Data: errorJSON,   // Encodes the error message as JSON
	}

	// Write the error response back to the client
	if err := conn.WriteJSON(errorResponse); err != nil {
		log.Println("Failed to send error message:", err)
	}
}
