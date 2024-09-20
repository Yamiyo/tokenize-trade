package main

import (
	"context"
	"fmt"
	"runtime/debug"
	"tokenize-trade/internal/utils/logger"
	"tokenize-trade/service/application"
)

var stop = make(chan error, 1)

func main() {
	defer func() {
		logger.SysLog().Error(context.Background(), "Server shutdown...")
		if err := recover(); err != nil {
			logger.SysLog().Error(context.Background(), fmt.Sprintf("error: %v", err))
			logger.SysLog().Error(context.Background(), string(debug.Stack()))
		}
	}()

	application.Run(stop)

	<-stop
}
