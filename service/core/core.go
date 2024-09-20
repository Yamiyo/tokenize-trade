package core

import (
	"sync"
	"tokenize-trade/internal/binance"
	"tokenize-trade/internal/config"
)

var (
	self *core
	once sync.Once
)

type core struct {
	in CoreIn
	RestCoreOut
}

type CoreIn struct {
	Conf      config.ConfigSetup
	BinanceWs binance.BinanceWsInterface
}

type RestCoreOut struct {
	TickerBookCore TickerBookCoreInterface
}

func New(in CoreIn) RestCoreOut {
	once.Do(func() {
		self = &core{
			in: in,
			RestCoreOut: RestCoreOut{
				TickerBookCore: newTickerBookCore(in),
			},
		}
	})

	return self.RestCoreOut
}
