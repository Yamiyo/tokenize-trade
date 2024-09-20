package binance

import (
	"context"
	"fmt"
	binance "github.com/binance/binance-connector-go"
	"github.com/stretchr/testify/suite"
	"testing"
	"time"
)

type TestBinance struct {
	ctx context.Context
	ws  *BinanceWs
	suite.Suite
}

func TestBinanceSuite(t *testing.T) {
	suite.Run(t, new(TestBinance))
}

func (s *TestBinance) SetupSuite() {
	s.ctx = context.Background()
	cfg := &BinanceCfg{
		ApiKey:    "",
		ApiSecret: "",
		WsURL:     "wss://stream.binance.com:9443",
	}

	client := binance.NewWebsocketStreamClient(false, cfg.WsURL)
	s.ws = &BinanceWs{
		cfg:    cfg,
		client: client,
	}
}

func (s *TestBinance) TestBinance_DepthServe() {
	respCh := make(chan *binance.WsDepthEvent)
	errCh := make(chan error)
	go s.ws.DepthServe(s.ctx, "ETHBTC", respCh, errCh)

	go func() {
		for resp := range respCh {
			fmt.Printf("resp: %+v\n", resp)
		}
	}()

	time.Sleep(5 * time.Minute)
}

func (s *TestBinance) TestBinance_TickerBookServe() {
	respCh := make(chan *binance.WsBookTickerEvent)
	errCh := make(chan error)

	go s.ws.TickerBookServe(s.ctx, "ETHBTC", respCh, errCh)

	go func() {
		for resp := range respCh {
			fmt.Printf("resp: %+v\n", resp)
		}
	}()

	time.Sleep(15 * time.Minute)
}
