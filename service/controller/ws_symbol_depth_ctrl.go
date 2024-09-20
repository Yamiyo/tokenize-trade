package controller

import (
	"context"
	"encoding/json"
	"github.com/gin-gonic/gin"
	"github.com/gorilla/websocket"
	"log"
	"time"
	ws_manager "tokenize-trade/internal/ws-manager"
)

type WsSymbolDepthCtrlInterface interface {
	Handle(ctx *gin.Context)
}

type WsSymbolDepthCtrl struct {
	in      RestCtrlIn
	manager *ws_manager.WebSocketManager
}

func newWsSymbolDepthCtrl(in RestCtrlIn) WsSymbolDepthCtrlInterface {
	manage := ws_manager.NewWebSocketManager()
	go manage.Start()

	msgCh, err := in.TickerBookCore.Subscript(context.Background(), "ETHBTC")
	if err != nil {
		log.Println("Failed to subscript:", err)
		return nil
	}

	go func() {
		for {
			select {
			case msg := <-msgCh:
				data, err := json.Marshal(msg)
				if err != nil {
					log.Println("Failed to marshal:", err)
					continue
				}
				manage.Broadcast <- data
			}
		}
	}()

	return &WsSymbolDepthCtrl{
		in:      in,
		manager: manage,
	}
}

func (ctrl *WsSymbolDepthCtrl) Handle(ctx *gin.Context) {
	// 升級為 WebSocket 連接
	conn, err := upgrader.Upgrade(ctx.Writer, ctx.Request, nil)
	if err != nil {
		log.Println("Failed to upgrade:", err)
		return
	}
	// 註冊新的客戶端
	ctrl.manager.Register <- conn

	// 設置讀取超時，等待客戶端發送 ping 消息
	conn.SetReadDeadline(time.Now().Add(15 * time.Second))
	// 設置 PingHandler 來回應 pong 消息
	conn.SetPingHandler(func(appData string) error {
		// 收到 ping 時重置讀取超時
		log.Println("Ping received from client")
		conn.SetReadDeadline(time.Now().Add(15 * time.Second)) // 收到 ping 後重置讀取超時
		// 回應 pong
		err := conn.WriteMessage(websocket.PongMessage, nil)
		if err != nil {
			log.Println("Error sending pong:", err)
			ctrl.manager.Unregister <- conn
			return err
		}
		return nil
	})

	for {
		_, message, err := conn.ReadMessage()
		if err != nil {
			log.Println("Read error:", err)
			ctrl.manager.Unregister <- conn
			break
		}

		if string(message) == "ping" {
			conn.SetReadDeadline(time.Now().Add(15 * time.Second))
		}
	}
}
