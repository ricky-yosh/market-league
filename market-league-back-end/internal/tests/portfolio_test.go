package tests

// "github.com/market-league/internal/testutils"

// func TestAddStockToPortfolio_Integration(t *testing.T) {
// 	// Step 1: Setup the test database
// 	db := testutils.SetupTestDB()

// 	// Step 2: Seed test data
// 	// db.Create(&models.Portfolio{ID: 1, Name: "Tech Portfolio"})
// 	// db.Create(&models.Stock{ID: 101, Symbol: "AAPL"})

// 	// Step 3: Setup repository, service, and handler
// 	portfolioRepo := portfolio.NewPortfolioRepository(db)
// 	portfolioService := portfolio.NewPortfolioService(portfolioRepo)
// 	handler := portfolio.NewPortfolioHandler(portfolioService)

// 	// Step 4: Setup Gin router and WebSocket endpoint
// 	router := http.NewServeMux()
// 	server := httptest.NewServer(router)
// 	defer server.Close()

// 	router.HandleFunc("/ws", func(w http.ResponseWriter, r *http.Request) {
// 		upgrader := websocket.Upgrader{}
// 		conn, err := upgrader.Upgrade(w, r, nil)
// 		assert.NoError(t, err)

// 		// Simulate WebSocket message
// 		rawMessage := ws.WebsocketMessage{
// 			Type: ws.MessageType_LeaguePortfolio_DraftStock,
// 			Data: json.RawMessage(`{"portfolio_id": 1, "stock_id": 101}`),
// 		}
// 		message, _ := json.Marshal(rawMessage)

// 		// Call the actual handler
// 		err = handler.AddStockToPortfolio(conn, json.RawMessage(message))
// 		assert.NoError(t, err)
// 	})

// 	// Step 5: Simulate WebSocket Client
// 	url := "ws" + server.URL[4:] + "/ws"
// 	wsClient, _, err := websocket.DefaultDialer.Dial(url, nil)
// 	assert.NoError(t, err)
// 	defer wsClient.Close()

// 	// Step 6: Send WebSocket Message
// 	request := ws.WebsocketMessage{
// 		Type: ws.MessageType_LeaguePortfolio_DraftStock,
// 		Data: json.RawMessage(`{"portfolio_id": 1, "stock_id": 101}`),
// 	}
// 	reqData, _ := json.Marshal(request)
// 	err = wsClient.WriteMessage(websocket.TextMessage, reqData)
// 	assert.NoError(t, err)

// 	// Step 7: Read Response
// 	_, resp, err := wsClient.ReadMessage()
// 	assert.NoError(t, err)

// 	var response ws.WebsocketMessage
// 	err = json.Unmarshal(resp, &response)
// 	assert.NoError(t, err)

// 	// Step 8: Assertions
// 	assert.Equal(t, ws.MessageType_LeaguePortfolio_DraftStock, response.Type)
// 	assert.JSONEq(t, `{"message": "Stock added successfully"}`, string(response.Data))

// 	// Verify that the stock was added in the database
// 	var portfolio models.Portfolio
// 	db.Preload("Stocks").First(&portfolio, 1)
// 	assert.Equal(t, 1, len(portfolio.Stocks))
// 	assert.Equal(t, uint(101), portfolio.Stocks[0].ID)
// }
