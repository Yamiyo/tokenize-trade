package application

import (
	"context"
	"errors"
	"sync"
	"tokenize-trade/internal/config"
	"tokenize-trade/internal/utils/logger"
	"tokenize-trade/service/application/rest"
)

var (
	templateApp   *TemplateApp
	serverSetOnce sync.Once
)

type TemplateApp struct {
	RestServer rest.ServiceInterface
}

func newTemplateServer(ctx context.Context, conf config.ConfigSetup) {
	serverSetOnce.Do(func() {
		templateApp = &TemplateApp{
			RestServer: rest.NewRestService(ctx, conf),
		}
	})
}

func Run(stop chan error) {
	if err := config.InitConfig(); err != nil {
		panic(err)
	}

	if err := logger.InitSysLog(
		config.GetLogConfig().Name,
		config.GetLogConfig().Level); err != nil {
		panic(err)
	}

	ctx := context.Background()

	if newTemplateServer(ctx, config.GetConfig()); templateApp == nil {
		panic(errors.New("templateApp is nil"))
	}

	go templateApp.RestServer.Run(ctx, stop)
}
