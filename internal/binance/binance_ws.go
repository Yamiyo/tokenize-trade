package binance

import (
	"context"
	"fmt"
	binance "github.com/binance/binance-connector-go"
	"sync"
	"tokenize-trade/internal/utils/logger"
)

var (
	self *BinanceWs
	once sync.Once
)

func CreateWebSocketClient(ctx context.Context, cfg *BinanceCfg) BinanceWsInterface {
	logger.SysLog().Info(ctx, "CreateWebSocketClient")
	logger.SysLog().Info(ctx, fmt.Sprintf("cfg: %v", cfg))
	once.Do(func() {
		client := binance.NewWebsocketStreamClient(false, cfg.WsURL)
		self = &BinanceWs{
			cfg:    cfg,
			client: client,
		}
	})

	return self
}

func (ws *BinanceWs) DepthServe(ctx context.Context, symbol string, messageChan chan *binance.WsDepthEvent, errorChan chan error) error {
	defer func() {
		logger.SysLog().Info(ctx, "DepthServe done")
		close(messageChan)
		close(errorChan)
	}()

	_, stopCh, err := ws.client.WsDepthServe(symbol, func(resp *binance.WsDepthEvent) {
		logger.SysLog().Debug(ctx, fmt.Sprintf("DepthServe: %v", resp))
		messageChan <- resp
	}, func(err error) {
		logger.SysLog().Error(ctx, err.Error())
		errorChan <- err
	})
	if err != nil {
		return err
	}

	for {
		select {
		case err := <-errorChan:
			return err
		case <-ctx.Done():
			stopCh <- struct{}{}
			close(stopCh)
			return nil
		}
	}
}

func (ws *BinanceWs) TickerBookServe(ctx context.Context, symbol string, messageChan chan *binance.WsBookTickerEvent, errorChan chan error) error {
	defer func() {
		logger.SysLog().Info(ctx, "TickerBookServe done")
		close(messageChan)
		close(errorChan)
	}()

	_, stopCh, err := ws.client.WsBookTickerServe(symbol, func(resp *binance.WsBookTickerEvent) {
		logger.SysLog().Debug(ctx, fmt.Sprintf("TickerBookServe: %v", resp))
		messageChan <- resp
	}, func(err error) {
		logger.SysLog().Error(ctx, err.Error())
		errorChan <- err
	})
	if err != nil {
		logger.SysLog().Error(ctx, err.Error())
		return err
	}

	for {
		select {
		case err := <-errorChan:
			return err
		case <-ctx.Done():
			stopCh <- struct{}{}
			return nil
		}
	}
}
