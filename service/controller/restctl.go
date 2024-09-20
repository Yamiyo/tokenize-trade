package controller

import (
	"github.com/gorilla/websocket"
	"net/http"
	"sync"
	"tokenize-trade/service/core"

	"tokenize-trade/internal/config"
)

var (
	self *restCtrl
	once sync.Once
)

type restCtrl struct {
	in RestCtrlIn
	RestCtrlOut
}

type RestCtrlIn struct {
	Conf           config.ConfigSetup
	TickerBookCore core.TickerBookCoreInterface
}

type RestCtrlOut struct {
	WsSymbolDepthCtrl WsSymbolDepthCtrlInterface
}

var upgrader = websocket.Upgrader{
	CheckOrigin: func(r *http.Request) bool {
		return true
	},
}

func New(in RestCtrlIn) RestCtrlOut {
	once.Do(func() {
		self = &restCtrl{
			in: in,
			RestCtrlOut: RestCtrlOut{
				WsSymbolDepthCtrl: newWsSymbolDepthCtrl(in),
			},
		}
	})

	return self.RestCtrlOut
}
