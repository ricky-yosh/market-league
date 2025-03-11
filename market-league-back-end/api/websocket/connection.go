package ws

import "github.com/gorilla/websocket"

// Connection wraps a raw WebSocket connection and tracks league subscriptions.
type Connection struct {
	Ws            *websocket.Conn
	Subscriptions map[uint]bool // key: leagueID, value: true if subscribed
}
