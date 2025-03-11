package api

import (
	"encoding/json"
	"fmt"
	"log"
	"net/http"

	"github.com/gin-gonic/gin"
	"github.com/gorilla/websocket"
	ws "github.com/market-league/api/websocket"
	"github.com/market-league/internal/league"
	league_portfolio "github.com/market-league/internal/league_portfolio"
	"github.com/market-league/internal/portfolio"
	"github.com/market-league/internal/stock"
	"github.com/market-league/internal/trade"
	"github.com/market-league/internal/user"
)

// Handle Dependencies
type WebSocketHandler struct {
	portfolioHandler       portfolio.PortfolioHandlerInterface
	stockHandler           stock.StockHandlerInterface
	userHandler            user.UserHandlerInterface
	tradeHandler           trade.TradeHandlerInterface
	leaguePortfolioHandler league_portfolio.LeaguePortfolioHandlerInterface
	leagueHandler          league.LeagueHandlerInterface
}

func NewWebSocketHandler(
	portfolioHandler portfolio.PortfolioHandlerInterface,
	stockHandler stock.StockHandlerInterface,
	userHandler user.UserHandlerInterface,
	tradeHandler trade.TradeHandlerInterface,
	leaguePortfolioHandler league_portfolio.LeaguePortfolioHandlerInterface,
	leagueHandler league.LeagueHandlerInterface,
) *WebSocketHandler {
	return &WebSocketHandler{
		portfolioHandler:       portfolioHandler,
		stockHandler:           stockHandler,
		userHandler:            userHandler,
		tradeHandler:           tradeHandler,
		leaguePortfolioHandler: leaguePortfolioHandler,
		leagueHandler:          leagueHandler,
	}
}

func (h *WebSocketHandler) routeTransmission(conn *ws.Connection, message ws.WebsocketMessage) error {
	// Route the message based on its type
	fmt.Printf("Message: %s", message)
	switch message.Type {

	// Portfolio Routes
	case ws.MessageType_Portfolio_CreatePortfolio:
		return h.portfolioHandler.CreatePortfolio(conn, message.Data)
	case ws.MessageType_Portfolio_PortfolioWithID:
		return h.portfolioHandler.GetPortfolioWithID(conn, message.Data)
	case ws.MessageType_Portfolio_LeaguePortfolio:
		return h.portfolioHandler.GetLeaguePortfolio(conn, message.Data)
	case ws.MessageType_Portfolio_AddStock:
		return h.portfolioHandler.AddStockToPortfolio(conn, message.Data)
	case ws.MessageType_Portfolio_RemoveStock:
		return h.portfolioHandler.RemoveStockFromPortfolio(conn, message.Data)
	case ws.MessageType_Portfolio_GetPortfolioPointsHistory:
		return h.portfolioHandler.GetPortfolioPointsHistory(conn, message.Data)
	case ws.MessageType_Portfolio_GetStocksValueChange:
		return h.portfolioHandler.GetStocksValueChange(conn, message.Data)

	// Stock Routes
	case ws.MessageType_Stock_CreateStock:
		return h.stockHandler.CreateStock(conn, message.Data)
	case ws.MessageType_Stock_CreateMultipleStocks:
		return h.stockHandler.CreateMultipleStocks(conn, message.Data)
	case ws.MessageType_Stock_UpdateCurrentStockPrice:
		return h.stockHandler.UpdatePrice(conn, message.Data)
	case ws.MessageType_Stock_GetStockInformation:
		return h.stockHandler.GetStockInfo(conn, message.Data)
	case ws.MessageType_Stock_GetAllStocks:
		return h.stockHandler.GetAllStocks(conn, message.Data)

	// User Routes
	case ws.MessageType_User_UserInfo:
		return h.userHandler.GetUserByID(conn, message.Data)
	case ws.MessageType_User_UserLeagues:
		return h.userHandler.GetUserLeagues(conn, message.Data)
	case ws.MessageType_User_UserTrades:
		return h.userHandler.GetUserTrades(conn, message.Data)
	case ws.MessageType_User_UserPortfolios:
		return h.userHandler.GetUserPortfolios(conn, message.Data)

	// Trade Routes
	case ws.MessageType_Trade_CreateTrade:
		return h.tradeHandler.CreateTrade(conn, message.Data)
	case ws.MessageType_Trade_ConfirmTrade:
		return h.tradeHandler.ConfirmTrade(conn, message.Data)
	case ws.MessageType_Trade_GetTrades:
		return h.tradeHandler.GetTrades(conn, message.Data)

	// League Portfolio Routes
	case ws.MessageType_LeaguePortfolio_DraftStock:
		return h.leaguePortfolioHandler.DraftStock(conn, message.Data)
	case ws.MessageType_LeaguePortfolio_GetLeaguePortfolioInfo:
		return h.leaguePortfolioHandler.GetLeaguePortfolioInfo(conn, message.Data)

	// League Routes
	case ws.MessageType_League_CreateLeague:
		return h.leagueHandler.CreateLeague(conn, message.Data)
	case ws.MessageType_League_RemoveLeague:
		return h.leagueHandler.RemoveLeague(conn, message.Data)
	case ws.MessageType_League_AddUserToLeague:
		return h.leagueHandler.AddUserToLeague(conn, message.Data)
	case ws.MessageType_League_GetDetails:
		return h.leagueHandler.GetLeagueDetails(conn, message.Data)
	case ws.MessageType_League_GetLeaderboard:
		return h.leagueHandler.GetLeaderboard(conn, message.Data)

	case ws.MessageType_SubscribeLeagues:
		var req struct {
			LeagueIDs []uint `json:"league_ids"`
		}
		if err := json.Unmarshal(message.Data, &req); err != nil {
			ws.SendError(conn, ws.MessageType_Error, "Invalid subscription format")
			return fmt.Errorf("invalid subscription format: %v", err)
		}
		// Update the connection's subscriptions.
		for _, leagueID := range req.LeagueIDs {
			conn.Subscriptions[leagueID] = true
		}
		// Optionally, send back an acknowledgment.
		return nil

	// Error or Unknown Message Type
	default:
		log.Println("Unknown message type:", message.Type)
		ws.SendError(conn, ws.MessageType_Error, "Unknown message type")
		return nil
	}
}

// WebSocket Upgrader - Upgrades HTTP to WebSocket
var upgrader = websocket.Upgrader{
	CheckOrigin: func(r *http.Request) bool {
		// TODO:
		// Allow all origins (modify for production)
		return true
	},
}

// HandleWebSocket - Handles incoming WebSocket connections
func (h *WebSocketHandler) HandleWebSocket(c *gin.Context) {
	// Upgrade HTTP request to WebSocket.
	rawConn, err := upgrader.Upgrade(c.Writer, c.Request, nil)
	if err != nil {
		log.Println("Failed to upgrade connection:", err)
		return
	}

	// Create our custom Connection with an empty subscriptions map.
	conn := &ws.Connection{
		Ws:            rawConn,
		Subscriptions: make(map[uint]bool),
	}
	ws.Manager.Register(conn)
	defer ws.Manager.Unregister(conn)

	log.Println("New WebSocket Connection Established!")

	for {
		// Read message from client.
		_, transmissionRaw, err := conn.Ws.ReadMessage()
		if err != nil {
			log.Println("Read error:", err)
			break
		}

		var message ws.WebsocketMessage
		if err = json.Unmarshal(transmissionRaw, &message); err != nil {
			// Pass the raw connection to SendError if that function expects *websocket.Conn.
			ws.SendError(conn, ws.MessageType_Error, "Invalid JSON format")
			continue
		}

		// Pass our custom Connection to routeTransmission.
		if err = h.routeTransmission(conn, message); err != nil {
			log.Println("Error with transmission:", message.Type, err.Error())
		}
	}
}
