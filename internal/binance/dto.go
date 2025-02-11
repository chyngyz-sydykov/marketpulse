package binance

import "time"

type RecordDto struct {
	Symbol    string
	Timeframe string
	Timestamp time.Time
	Open      float64
	Close     float64
	High      float64
	Low       float64
	Volume    float64
}

type BinanceTickerResponse struct {
	Symbol    string `json:"symbol"`
	OpenPrice string `json:"openPrice"`
	HighPrice string `json:"highPrice"`
	LowPrice  string `json:"lowPrice"`
	LastPrice string `json:"lastPrice"`
	Volume    string `json:"volume"`
}
