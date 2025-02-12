package binance

import "time"

type RecordDto struct {
	Id        *int
	Symbol    string
	Timeframe string
	Timestamp time.Time
	Open      float64
	Close     float64
	High      float64
	Low       float64
	Volume    float64
}

type IndicatorDto struct {
	Id             *int
	SMA            float64
	EMA            float64
	StdDev         float64
	LowerBollinger float64
	UpperBollinger float64
	RSI            float64
	Volatility     float64
	MACD           float64
	MACDSignal     float64
}

type BinanceTickerResponse struct {
	Symbol    string `json:"symbol"`
	OpenPrice string `json:"openPrice"`
	HighPrice string `json:"highPrice"`
	LowPrice  string `json:"lowPrice"`
	LastPrice string `json:"lastPrice"`
	Volume    string `json:"volume"`
}
