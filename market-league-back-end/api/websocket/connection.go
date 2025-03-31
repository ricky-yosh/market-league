package ws

import (
	"sync"

	"github.com/gorilla/websocket"
)

// First, update the Connection struct in connection.go
type Connection struct {
	Ws            *websocket.Conn
	Subscriptions map[uint]bool // key: leagueID, value: true if subscribed
	writeMutex    sync.Mutex    // Mutex to protect writes to this connection
}

// Then update the BroadcastToLeague method in manager.go
func (m *WebSocketManager) BroadcastToLeague(leagueID uint, message []byte) {
	m.mutex.Lock()
	// Create a copy of connections to avoid holding the lock while writing
	connections := make([]*Connection, 0, len(m.connections))
	for conn := range m.connections {
		if conn.Subscriptions[leagueID] {
			connections = append(connections, conn)
		}
	}
	m.mutex.Unlock()

	// Now write to each connection with its own lock
	for _, conn := range connections {
		conn.writeMutex.Lock()
		err := conn.Ws.WriteMessage(websocket.TextMessage, message)
		conn.writeMutex.Unlock()

		if err != nil {
			// Handle the error outside of the lock
			m.mutex.Lock()
			delete(m.connections, conn)
			m.mutex.Unlock()
			conn.Ws.Close()
		}
	}
}
