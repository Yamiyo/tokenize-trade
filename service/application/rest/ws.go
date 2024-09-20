package rest

import (
	"github.com/gin-gonic/gin"
)

func (s *restService) setWebsocketRoutes(parentRouteGroup *gin.RouterGroup) {
	ws := parentRouteGroup.Group("/ws")
	ws.GET("symbol-depth", s.WsSymbolDepthCtrl.Handle)
}
