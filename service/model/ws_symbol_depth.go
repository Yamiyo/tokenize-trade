package model

import "github.com/shopspring/decimal"

type WsSymbolDepthResponse struct {
	ID      string          `json:"id"`
	Symbol  string          `json:"symbol"`
	Bids    []*SymbolTicker `json:"bids"`
	Asks    []*SymbolTicker `json:"asks"`
	BidsSum decimal.Decimal `json:"bids_sum"`
	AsksSum decimal.Decimal `json:"asks_sum"`
}

type SymbolTicker struct {
	Price decimal.Decimal `json:"price"`
	Size  decimal.Decimal `json:"size"`
}
