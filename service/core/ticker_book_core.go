package core

import (
	"context"
	binance_connector "github.com/binance/binance-connector-go"
	"github.com/google/uuid"
	"github.com/shopspring/decimal"
	"math/rand"
	"strconv"
	"time"
	"tokenize-trade/internal/utils/logger"
	"tokenize-trade/service/model"
)

type TickerBookCoreInterface interface {
	Subscript(ctx context.Context, symbol string) (chan *model.WsSymbolDepthResponse, error)
	UnSubscript(ctx context.Context)
}

type TickerBookCore struct {
	in     CoreIn
	msgCh  chan *model.WsSymbolDepthResponse
	errCh  chan error
	cancel context.CancelFunc
}

func newTickerBookCore(in CoreIn) TickerBookCoreInterface {
	return &TickerBookCore{
		in:    in,
		msgCh: make(chan *model.WsSymbolDepthResponse),
		errCh: make(chan error),
	}
}

func (c *TickerBookCore) Subscript(ctx context.Context, symbol string) (chan *model.WsSymbolDepthResponse, error) {
	cc, cancel := context.WithCancel(ctx)
	c.cancel = cancel
	go func(ctx context.Context) {
		msgCh := make(chan *binance_connector.WsBookTickerEvent, 10)
		go func() {
			for msg := range msgCh {
				data := c.convertResponse(msg)
				c.msgCh <- data
			}
		}()

		if err := c.in.BinanceWs.TickerBookServe(ctx, symbol, msgCh, c.errCh); err != nil {
			logger.SysLog().Error(ctx, err.Error())
			c.UnSubscript(ctx)
		}
	}(cc)

	return c.msgCh, nil
}

func (c *TickerBookCore) UnSubscript(ctx context.Context) {
	c.cancel()
}

func randomBetween(min, max float64) float64 {
	rand.Seed(time.Now().UnixNano())
	num := rand.Float64()
	if num < 0 {
		num = -num
	}
	return min + num*(max-min)
}

func (c *TickerBookCore) convertResponse(msg *binance_connector.WsBookTickerEvent) *model.WsSymbolDepthResponse {
	bidPrice, _ := strconv.ParseFloat(msg.BestBidPrice, 64)
	bidSize, _ := strconv.ParseFloat(msg.BestBidQty, 64)
	askPrice, _ := strconv.ParseFloat(msg.BestAskPrice, 64)
	askSize, _ := strconv.ParseFloat(msg.BestAskQty, 64)

	bids := []*model.SymbolTicker{
		{Price: decimal.NewFromFloat(bidPrice), Size: decimal.NewFromFloat(bidSize)},
	}
	asks := []*model.SymbolTicker{
		{Price: decimal.NewFromFloat(askPrice), Size: decimal.NewFromFloat(askSize)},
	}
	bidSum := decimal.NewFromFloat(bidPrice).Mul(decimal.NewFromFloat(bidSize))
	askSum := decimal.NewFromFloat(askSize)
	bidTotal := bidPrice * bidSize
	askTotal := askSize

	// Generate Bid List where sum(size * price) < 5
	for bidTotal < 5 {
		bidPrice = bidPrice - randomBetween(0.0001, 0.001)
		bidSize = randomBetween(1, 50)
		bidValue := bidPrice * bidSize

		// Check if adding this bid will exceed the total
		if (bidTotal + bidValue) < 5 {
			bidTotal += bidValue
			bids = append(bids, &model.SymbolTicker{
				Price: decimal.NewFromFloat(bidPrice),
				Size:  decimal.NewFromFloat(bidSize),
			})
			bidSum = bidSum.Add(decimal.NewFromFloat(bidPrice).Mul(decimal.NewFromFloat(bidSize)))
		} else {
			break
		}
	}

	// Generate Ask List where sum(size) < 150
	for askTotal < 150 {
		askPrice = askPrice + randomBetween(0.0001, 0.001)
		askSize = randomBetween(1, 50)

		// Check if adding this ask will exceed the total
		if (askTotal + askSize) < 150 {
			askTotal += askSize
			asks = append(asks, &model.SymbolTicker{
				Price: decimal.NewFromFloat(askPrice),
				Size:  decimal.NewFromFloat(askSize),
			})
			askSum = askSum.Add(decimal.NewFromFloat(askSize))
		} else {
			break
		}
	}

	return &model.WsSymbolDepthResponse{
		ID:      uuid.New().String(),
		Symbol:  msg.Symbol,
		Bids:    bids,
		Asks:    asks,
		BidsSum: bidSum,
		AsksSum: askSum,
	}
}
