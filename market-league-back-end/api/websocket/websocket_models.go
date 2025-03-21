package ws

import (
	"encoding/json"
	"log"

	"github.com/gin-gonic/gin"
)

type WebsocketMessage struct {
	Type string          `json:"type"` // Message type for routing
	Data json.RawMessage `json:"data"` // Raw JSON data for flexibility
}

const (
	// Portfolio Routes
	MessageType_Portfolio_CreatePortfolio           = "MessageType_Portfolio_CreatePortfolio"
	MessageType_Portfolio_PortfolioWithID           = "MessageType_Portfolio_PortfolioWithID"
	MessageType_Portfolio_LeaguePortfolio           = "MessageType_Portfolio_LeaguePortfolio"
	MessageType_Portfolio_AddStock                  = "MessageType_Portfolio_AddStock"
	MessageType_Portfolio_RemoveStock               = "MessageType_Portfolio_RemoveStock"
	MessageType_Portfolio_GetPortfolioPointsHistory = "MessageType_Portfolio_GetPortfolioPointsHistory"
	MessageType_Portfolio_GetStocksValueChange      = "MessageType_Portfolio_GetStocksValueChange"

	// Stock Routes
	MessageType_Stock_CreateStock             = "MessageType_Stock_CreateStock"
	MessageType_Stock_CreateMultipleStocks    = "MessageType_Stock_CreateMultipleStocks"
	MessageType_Stock_GetStockInformation     = "MessageType_Stock_GetStockInformation"
	MessageType_Stock_UpdateCurrentStockPrice = "MessageType_Stock_UpdateCurrentStockPrice"
	MessageType_Stock_GetAllStocks            = "MessageType_Stock_GetAllStocks"

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
	MessageType_League_QueueUp         = "MessageType_League_QueueUp"
	MessageType_League_Portfolios      = "MessageType_League_Portfolios"
	MessageType_League_DraftUpdate     = "MessageType_League_DraftUpdate"
	MessageType_League_DraftPick       = "MessageType_League_DraftPick"
	MessageType_League_GetAllLeagues   = "MessageType_League_GetAllLeagues"

	// Websocket Subscription
	MessageType_SubscribeLeagues = "subscribe_leagues"

	// Error Message
	MessageType_Error = "MessageType_Error"
)

func SendError(conn *Connection, messageType string, errorMsg string) {
	// Construct the error response using gin.H
	errorData := gin.H{
		"type":    MessageType_Error, // Assume MessageType_Error is defined elsewhere
		"message": errorMsg,
	}

	// Marshal the error data into JSON bytes
	errorJSON, err := json.Marshal(errorData)
	if err != nil {
		log.Println("Failed to serialize error data:", err)
		return
	}

	// Create an error response using the WebsocketMessage struct
	errorResponse := WebsocketMessage{
		Type: messageType,
		Data: errorJSON,
	}

	// Write the error response back to the client using the underlying Ws connection.
	if err := conn.Ws.WriteJSON(errorResponse); err != nil {
		log.Println("Failed to send error message:", err)
	}
}
