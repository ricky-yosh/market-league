package ws

import (
	"sync"

	"github.com/gorilla/websocket"
)

type WebSocketManager struct {
	connections map[*Connection]bool
	mutex       sync.Mutex
}

var Manager = &WebSocketManager{
	connections: make(map[*Connection]bool),
}

func (m *WebSocketManager) Register(conn *Connection) {
	m.mutex.Lock()
	defer m.mutex.Unlock()
	m.connections[conn] = true
}

func (m *WebSocketManager) Unregister(conn *Connection) {
	m.mutex.Lock()
	defer m.mutex.Unlock()
	delete(m.connections, conn)
	conn.Ws.Close()
}

// BroadcastToLeague sends a message to only those connections subscribed to the given leagueID.
func (m *WebSocketManager) BroadcastToLeague(leagueID uint, message []byte) {
	m.mutex.Lock()
	defer m.mutex.Unlock()
	for conn := range m.connections {
		if conn.Subscriptions[leagueID] {
			if err := conn.Ws.WriteMessage(websocket.TextMessage, message); err != nil {
				// In case of error, remove the connection.
				delete(m.connections, conn)
				conn.Ws.Close()
			}
		}
	}
}
