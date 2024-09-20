package rest

import (
	"context"
	"sync"
	"tokenize-trade/internal/binance"
	"tokenize-trade/service/core"

	"tokenize-trade/internal/config"
	"tokenize-trade/internal/utils/logger"
	restctl "tokenize-trade/service/controller"

	"github.com/gin-gonic/gin"
)

var (
	self *restService
	once sync.Once
)

func NewRestService(ctx context.Context, conf config.ConfigSetup) ServiceInterface {
	once.Do(func() {
		cor := core.New(core.CoreIn{
			Conf:      conf,
			BinanceWs: binance.CreateWebSocketClient(ctx, &conf.BinanceConfig),
		})
		ctrl := restctl.New(restctl.RestCtrlIn{
			Conf:           conf,
			TickerBookCore: cor.TickerBookCore,
		})

		self = &restService{
			WsSymbolDepthCtrl: ctrl.WsSymbolDepthCtrl,
		}
	})

	return self
}

type ServiceInterface interface {
	Run(ctx context.Context, stop chan error)
}

type restService struct {
	WsSymbolDepthCtrl restctl.WsSymbolDepthCtrlInterface
}

func (s *restService) Run(ctx context.Context, stop chan error) {
	engine := s.newEngine()
	engine.LoadHTMLGlob("./web/dist/*.html")
	//設定靜態資源的讀取
	engine.Static("/assets", "./web/dist/assets")
	//engine.StaticFS("/assets", gin.Dir("./web/assets", true))

	s.setRoutes(engine)

	if err := engine.Run(config.GetGinConfig().Address); err != nil {
		logger.SysLog().Error(ctx, err.Error())
		stop <- err
	}
}

func (s *restService) newEngine() *gin.Engine {
	return gin.New()
}

func (s *restService) setRoutes(engine *gin.Engine) {
	// set middlewares
	engine.Use(
		gin.Logger(),
		gin.Recovery(),
	)

	// set router
	s.setPublicRoutes(engine)
}

func (s *restService) setPublicRoutes(engine *gin.Engine) {
	publicRouteGroup := engine.Group("")

	// setting router
	s.setViewRoutes(publicRouteGroup)
	s.setWebsocketRoutes(publicRouteGroup)
}
