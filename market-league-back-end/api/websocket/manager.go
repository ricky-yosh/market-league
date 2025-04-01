package ws

import (
	"sync"
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
