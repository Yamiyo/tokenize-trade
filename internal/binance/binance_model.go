package binance

import (
	"context"
	binance "github.com/binance/binance-connector-go"
)

type BinanceCfg struct {
	ApiKey    string `mapstructure:"api_key"`
	ApiSecret string `mapstructure:"api_secret"`
	WsURL     string `mapstructure:"ws_url"`
}

type BinanceWs struct {
	cfg    *BinanceCfg
	client *binance.WebsocketStreamClient
}

type BinanceWsInterface interface {
	DepthServe(ctx context.Context, symbol string, messageChan chan *binance.WsDepthEvent, errorChan chan error) error
	TickerBookServe(ctx context.Context, symbol string, messageChan chan *binance.WsBookTickerEvent, errorChan chan error) error
}
