package ws_manager

import (
	"github.com/gorilla/websocket"
	"log"
)

// WebSocket 管理器
type WebSocketManager struct {
	Clients    map[*websocket.Conn]bool // 所有連接的客戶端
	Broadcast  chan []byte              // 廣播消息的 channel
	Register   chan *websocket.Conn     // 註冊新的客戶端
	Unregister chan *websocket.Conn     // 注銷客戶端
}

// 初始化 WebSocket 管理器
func NewWebSocketManager() *WebSocketManager {
	return &WebSocketManager{
		Clients:    make(map[*websocket.Conn]bool),
		Broadcast:  make(chan []byte),
		Register:   make(chan *websocket.Conn),
		Unregister: make(chan *websocket.Conn),
	}
}

// 啟動 WebSocket 管理器
func (manager *WebSocketManager) Start() {
	for {
		select {
		case conn := <-manager.Register:
			// 註冊新客戶端
			manager.Clients[conn] = true
			log.Println("New client connected")

		case conn := <-manager.Unregister:
			// 注銷客戶端
			if _, ok := manager.Clients[conn]; ok {
				delete(manager.Clients, conn)
				conn.Close()
				log.Println("Client disconnected")
			}

		case message := <-manager.Broadcast:
			// 廣播消息給所有客戶端
			for conn := range manager.Clients {
				err := conn.WriteMessage(websocket.TextMessage, message)
				if err != nil {
					log.Printf("Broadcast error: %v", err)
					conn.Close()
					delete(manager.Clients, conn)
				}
			}
		}
	}
}
